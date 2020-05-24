package shipmentprice

import (
	"context"
	"math"
	"sort"

	"o.o/api/main/location"
	"o.o/api/main/shipmentpricing/pricelist"
	"o.o/api/main/shipmentpricing/shipmentprice"
	"o.o/api/top/types/etc/route_type"
	"o.o/api/top/types/etc/status3"
	locationutil "o.o/backend/com/main/location/util"
	"o.o/backend/com/main/shipmentpricing"
	"o.o/backend/com/main/shipmentpricing/shipmentprice/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

var _ shipmentprice.QueryService = &QueryService{}

type QueryService struct {
	redisStore         redis.Store
	shipmentPriceStore sqlstore.ShipmentPriceStoreFactory
	locationQS         location.QueryBus
	priceListQS        pricelist.QueryBus
}

func NewQueryService(db *cmsql.Database, redisStore redis.Store, locationQS location.QueryBus, priceListQS pricelist.QueryBus) *QueryService {
	return &QueryService{
		redisStore:         redisStore,
		shipmentPriceStore: sqlstore.NewShipmentPriceStore(db),
		locationQS:         locationQS,
		priceListQS:        priceListQS,
	}
}

func QueryServiceMessageBus(q *QueryService) shipmentprice.QueryBus {
	b := bus.New()
	return shipmentprice.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *QueryService) ListShipmentPrices(ctx context.Context, args *shipmentprice.ListShipmentPricesArgs) ([]*shipmentprice.ShipmentPrice, error) {
	return q.shipmentPriceStore(ctx).OptionalShipmentServiceID(args.ShipmentServiceID).OptionalShipmentPriceListID(args.ShipmentPriceListID).ListShipmentPrices()
}

func (q *QueryService) GetShipmentPrice(ctx context.Context, ID dot.ID) (*shipmentprice.ShipmentPrice, error) {
	return q.shipmentPriceStore(ctx).ID(ID).GetShipmentPrice()
}

/*
	func: GetActiveShipmentPrices

	- shipmentPriceListID: dùng để tính giá cho một bảng giá cụ thể (sử dụng trong trường hợp admin kiểm tra giá trước khi public bảng giá)
	- Khi shop sử dụng, shipmentPriceListID = 0 vì bảng giá sử dụng để tính toán mặc định lấy từ database ra (GetActiveShipmentPriceListQuery)
*/

func (q *QueryService) GetActiveShipmentPrices(ctx context.Context, shipmentServiceID, shipmentPriceListID dot.ID) ([]*shipmentprice.ShipmentPrice, error) {
	var res []*shipmentprice.ShipmentPrice
	key := getActiveShipmentPricesRedisKey(ctx)

	if shipmentPriceListID != 0 {
		// no cache in here
		res, err := q.shipmentPriceStore(ctx).ShipmentPriceListID(shipmentPriceListID).Status(status3.P).ListShipmentPrices()
		if err != nil {
			return nil, err
		}
		return filterShipmentPricesByShipmentServiceID(res, shipmentServiceID), nil
	}

	err := q.redisStore.Get(key, &res)
	if err != nil {
		query := &pricelist.GetActiveShipmentPriceListQuery{}
		if err := q.priceListQS.Dispatch(ctx, query); err != nil {
			return nil, err
		}
		res, err = q.shipmentPriceStore(ctx).ShipmentPriceListID(query.Result.ID).Status(status3.P).ListShipmentPrices()
		if err == nil {
			_ = q.redisStore.SetWithTTL(key, res, shipmentpricing.DefaultTTL)
		}
	}

	res = filterShipmentPricesByShipmentServiceID(res, shipmentServiceID)
	return res, err
}

func filterShipmentPricesByShipmentServiceID(shipmentPrices []*shipmentprice.ShipmentPrice, shipmentServiceID dot.ID) (res []*shipmentprice.ShipmentPrice) {
	for _, sp := range shipmentPrices {
		if sp.ShipmentServiceID == shipmentServiceID {
			res = append(res, sp)
		}
	}
	return res
}

func (q *QueryService) CalculatePrice(ctx context.Context, args *shipmentprice.CalculatePriceArgs) (*shipmentprice.CalculatePriceResult, error) {
	if args.Weight == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing weight")
	}
	fromQuery := &location.FindOrGetLocationQuery{
		Province:     args.FromProvince,
		District:     args.FromDistrict,
		ProvinceCode: args.FromProvinceCode,
		DistrictCode: args.FromDistrictCode,
	}
	toQuery := &location.FindOrGetLocationQuery{
		Province:     args.ToProvince,
		District:     args.ToDistrict,
		ProvinceCode: args.ToProvinceCode,
		DistrictCode: args.ToDistrictCode,
	}
	if err := q.locationQS.DispatchAll(ctx, fromQuery, toQuery); err != nil {
		return nil, err
	}
	from, to := fromQuery.Result, toQuery.Result

	pricings, err := q.GetActiveShipmentPrices(ctx, args.ShipmentServiceID, args.ShipmentPriceListID)
	if err != nil {
		return nil, err
	}
	matchingPricings, err := q.GetMatchingPricings(ctx, pricings, from.Province, to.Province, to.District)
	if len(matchingPricings) == 0 || err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Không có gói phù hợp")
	}

	pricing := getPricingByPriorityPoint(matchingPricings)
	pRuleDetail, err := GetPriceRuleDetail(args.Weight, pricing.Details)
	if err != nil {
		return nil, err
	}
	price, err := GetPriceByPricingDetail(args.Weight, pRuleDetail)
	if err != nil {
		return nil, err
	}
	return &shipmentprice.CalculatePriceResult{
		ShipmentPriceID: pricing.ID,
		Price:           price,
	}, nil
}

func (q *QueryService) GetMatchingPricings(ctx context.Context, pricings []*shipmentprice.ShipmentPrice, fromProvince, toProvince *location.Province, toDistrict *location.District) (res []*shipmentprice.ShipmentPrice, err error) {
	queryFrom := &location.GetCustomRegionByCodeQuery{
		ProvinceCode: fromProvince.Code,
	}
	queryTo := &location.GetCustomRegionByCodeQuery{
		ProvinceCode: toProvince.Code,
	}
	if err := q.locationQS.DispatchAll(ctx, queryFrom, queryTo); err != nil && cm.ErrorCode(err) != cm.NotFound {
		return nil, err
	}
	fromCustomRegion, toCustomRegion := dot.ID(0), dot.ID(0)
	if queryFrom.Result != nil {
		fromCustomRegion = queryFrom.Result.ID
	}
	if queryTo.Result != nil {
		toCustomRegion = queryTo.Result.ID
	}

	for _, pricing := range pricings {
		if checkMatchingPricing(pricing, fromProvince, toProvince, toDistrict, fromCustomRegion, toCustomRegion) {
			res = append(res, pricing)
		}
	}
	return res, nil
}

func checkMatchingPricing(pricing *shipmentprice.ShipmentPrice, fromProvince, toProvince *location.Province, toDistrict *location.District, fromCustomRegion, toCustomRegion dot.ID) bool {
	if pricing.Status != status3.P {
		return false
	}
	// check CustomRegionRouteType
	if len(pricing.CustomRegionTypes) > 0 {
		customRegionRouteType := locationutil.GetCustomRegionRouteType(fromCustomRegion, toCustomRegion)
		if !containCustomRegionRouteType(pricing.CustomRegionTypes, customRegionRouteType) {
			return false
		}
	}
	if len(pricing.CustomRegionIDs) > 0 {
		if !cm.IDsContain(pricing.CustomRegionIDs, fromCustomRegion) || !cm.IDsContain(pricing.CustomRegionIDs, toCustomRegion) {
			return false
		}
	}

	// check RegionRouteType
	if len(pricing.RegionTypes) > 0 {
		regionRouteType := locationutil.GetRegionRouteType(fromProvince, toProvince)
		if !containRegionRouteType(pricing.RegionTypes, regionRouteType) {
			return false
		}
	}

	// check ProvinceRouteType
	if len(pricing.ProvinceTypes) > 0 {
		provinceRouteType := locationutil.GetProvinceRouteType(fromProvince, toProvince)
		if !containProvinceRouteType(pricing.ProvinceTypes, provinceRouteType) {
			return false
		}
	}

	// check UrbanType
	if len(pricing.UrbanTypes) > 0 {
		check := false
		dType := locationutil.GetShippingDistrictType(toDistrict)
		for _, districtType := range pricing.UrbanTypes {
			if districtType == dType {
				check = true
				break
			}
		}
		if !check {
			return false
		}
	}
	return true
}

func GetPriceRuleDetail(weight int, priceRuleDetails []*shipmentprice.PricingDetail) (*shipmentprice.PricingDetail, error) {
	if priceRuleDetails == nil || len(priceRuleDetails) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing priceDetails")
	}

	var weightIndex []int
	var priceRuleDetailsMap = make(map[int]*shipmentprice.PricingDetail)
	for _, detail := range priceRuleDetails {
		weightIndex = append(weightIndex, detail.Weight)
		priceRuleDetailsMap[detail.Weight] = detail
	}
	sort.Ints(weightIndex) // increase
	var index int
	for _, wIndex := range weightIndex {
		index = wIndex
		if weight <= wIndex {
			break
		}
	}

	return priceRuleDetailsMap[index], nil
}

func GetPriceByPricingDetail(weight int, pRuleDetail *shipmentprice.PricingDetail) (int, error) {
	if (pRuleDetail.Overweight == nil || len(pRuleDetail.Overweight) == 0) && weight > pRuleDetail.Weight {
		// can not apply this rule
		return 0, cm.Error(cm.InvalidArgument, "Can not apply to this rule", nil)
	}

	addWeight := weight - pRuleDetail.Weight
	if addWeight <= 0 {
		return pRuleDetail.Price, nil
	}

	price := pRuleDetail.Price
	if !checkValidOverweight(weight, pRuleDetail.Overweight) {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "The weight is out of range")
	}
	for _, ov := range pRuleDetail.Overweight {
		ovPrice := GetOverweightPrice(weight, ov)
		price += ovPrice
	}
	return price, nil
}

// checkValidOverweight
//
// Kiểm tra điều kiện vượt cân có hợp lệ hay không
// Khối lượng phải nằm trong khoảng cân được setting
func checkValidOverweight(weight int, overWeights []*shipmentprice.PricingDetailOverweight) bool {
	maxWeight := 0
	for _, ov := range overWeights {
		if ov.MaxWeight == -1 {
			return true
			break
		}
		if maxWeight < ov.MaxWeight {
			maxWeight = ov.MaxWeight
		}
	}
	if weight <= maxWeight {
		return true
	}
	return false
}

func GetOverweightPrice(weight int, ov *shipmentprice.PricingDetailOverweight) int {
	if weight <= ov.MinWeight {
		return 0
	}
	if ov.MaxWeight != -1 && weight > ov.MaxWeight {
		weight = ov.MaxWeight
	}
	additionalWeight := weight - ov.MinWeight
	s := float64(additionalWeight) / float64(ov.WeightStep)
	step := int(math.Ceil(s))
	return step * ov.PriceStep
}

func containRegionRouteType(types []route_type.RegionRouteType, routeType route_type.RegionRouteType) bool {
	for _, rt := range types {
		if routeType == rt {
			return true
		}
	}
	return false
}

func containProvinceRouteType(types []route_type.ProvinceRouteType, routeType route_type.ProvinceRouteType) bool {
	for _, rt := range types {
		if routeType == rt {
			return true
		}
	}
	return false
}

func containCustomRegionRouteType(types []route_type.CustomRegionRouteType, routeType route_type.CustomRegionRouteType) bool {
	for _, rt := range types {
		if routeType == rt {
			return true
		}
	}
	return false
}

func getPricingByPriorityPoint(pricings []*shipmentprice.ShipmentPrice) *shipmentprice.ShipmentPrice {
	sort.Slice(pricings, func(i, j int) bool {
		return pricings[i].PriorityPoint > pricings[j].PriorityPoint
	})
	return pricings[0]
}
