package shipmentprice

import (
	"context"
	"math"
	"sort"

	"o.o/api/main/location"
	"o.o/api/main/shipmentpricing/pricelist"
	"o.o/api/main/shipmentpricing/shipmentprice"
	"o.o/api/main/shipmentpricing/shopshipmentpricelist"
	"o.o/api/top/types/etc/route_type"
	"o.o/api/top/types/etc/shipping_fee_type"
	"o.o/api/top/types/etc/status3"
	com "o.o/backend/com/main"
	locationutil "o.o/backend/com/main/location/util"
	"o.o/backend/com/main/shipmentpricing/shipmentprice/sqlstore"
	"o.o/backend/com/main/shipmentpricing/util"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/redis"
	"o.o/capi/dot"
)

var _ shipmentprice.QueryService = &QueryService{}

type QueryService struct {
	redisStore         redis.Store
	shipmentPriceStore sqlstore.ShipmentPriceStoreFactory
	locationQS         location.QueryBus
	priceListQS        pricelist.QueryBus
	shopPriceListQS    shopshipmentpricelist.QueryBus
}

func NewQueryService(db com.MainDB, redisStore redis.Store, locationQS location.QueryBus, priceListQS pricelist.QueryBus, shopPriceListQS shopshipmentpricelist.QueryBus) *QueryService {
	return &QueryService{
		redisStore:         redisStore,
		shipmentPriceStore: sqlstore.NewShipmentPriceStore(db),
		locationQS:         locationQS,
		priceListQS:        priceListQS,
		shopPriceListQS:    shopPriceListQS,
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

type getActiveShipmentPricesArgs struct {
	AccountID           dot.ID
	ConnectionID        dot.ID
	ShipmentServiceID   dot.ID
	ShipmentPriceListID dot.ID
}

/*
	func: GetActiveShipmentPrices

	- Field `shipmentPriceListID`: dùng để tính giá cho một bảng giá cụ thể (sử dụng trong trường hợp admin kiểm tra giá trước khi public bảng giá)
	- Khi shop sử dụng & `shipmentPriceListID` = 0:
		- Tìm bảng giá của shop đó trong bảng `shop_shipment_price_list`, nếu có, sử dụng bảng giá đó
		- Nếu shop ko có bảng giá riêng, sử dụng bảng giá mặc định (GetActiveShipmentPriceListQuery)
*/

func (q *QueryService) GetActiveShipmentPrices(ctx context.Context, args getActiveShipmentPricesArgs) ([]*shipmentprice.ShipmentPrice, error) {
	shipmentServiceID, shipmentPriceListID := args.ShipmentServiceID, args.ShipmentPriceListID
	var res []*shipmentprice.ShipmentPrice

	if shipmentPriceListID == 0 && args.AccountID != 0 {
		queryShopPriceList := &shopshipmentpricelist.GetShopShipmentPriceListQuery{
			ShopID:       args.AccountID,
			ConnectionID: args.ConnectionID,
		}
		if err := q.shopPriceListQS.Dispatch(ctx, queryShopPriceList); err == nil {
			shipmentPriceListID = queryShopPriceList.Result.ShipmentPriceListID
		}
	}
	key := getActiveShipmentPricesRedisKey(ctx, shipmentPriceListID)

	err := q.redisStore.Get(key, &res)
	if err != nil {
		priceListIDs := []dot.ID{shipmentPriceListID}
		if shipmentPriceListID == 0 {
			priceListIDs, err = q.listActivePriceLists(ctx)
			if err != nil {
				return nil, err
			}
		}
		res, err = q.shipmentPriceStore(ctx).ShipmentPriceListIDs(priceListIDs...).Status(status3.P).ListShipmentPrices()
		if err != nil {
			return nil, err
		}
		_ = q.redisStore.SetWithTTL(key, res, util.DefaultTTL)
	}

	res = filterShipmentPricesByShipmentServiceID(res, shipmentServiceID)
	return res, nil
}

// getActivePriceList
//
// Lấy tất cả các bảng giá mặc định của các NVC (đang có is_default = true)
func (q *QueryService) listActivePriceLists(ctx context.Context) ([]dot.ID, error) {
	query := &pricelist.ListShipmentPriceListsQuery{
		IsDefault: dot.Bool(true),
	}
	if err := q.priceListQS.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	res := make([]dot.ID, len(query.Result))
	for i, priceList := range query.Result {
		res[i] = priceList.ID
	}
	return res, nil
}

func filterShipmentPricesByShipmentServiceID(shipmentPrices []*shipmentprice.ShipmentPrice, shipmentServiceID dot.ID) (res []*shipmentprice.ShipmentPrice) {
	for _, sp := range shipmentPrices {
		if sp.ShipmentServiceID == shipmentServiceID {
			res = append(res, sp)
		}
	}
	return res
}

func (q *QueryService) GetMatchingPricings(ctx context.Context, pricings []*shipmentprice.ShipmentPrice, fromProvince, toProvince *location.Province, toDistrict *location.District) (res []*shipmentprice.ShipmentPrice, err error) {
	queryFrom := &location.ListCustomRegionsByCodeQuery{
		ProvinceCode: fromProvince.Code,
	}
	queryTo := &location.ListCustomRegionsByCodeQuery{
		ProvinceCode: toProvince.Code,
	}
	if err := q.locationQS.DispatchAll(ctx, queryFrom, queryTo); err != nil && cm.ErrorCode(err) != cm.NotFound {
		return nil, err
	}
	fromCustomRegions, toCustomRegions := []dot.ID{}, []dot.ID{}
	for _, fromRegion := range queryFrom.Result {
		fromCustomRegions = append(fromCustomRegions, fromRegion.ID)
	}
	for _, toRegion := range queryTo.Result {
		toCustomRegions = append(toCustomRegions, toRegion.ID)
	}

	for _, pricing := range pricings {
		if checkMatchingPricing(pricing, fromProvince, toProvince, toDistrict, fromCustomRegions, toCustomRegions) {
			res = append(res, pricing)
		}
	}
	return res, nil
}

func checkMatchingPricing(pricing *shipmentprice.ShipmentPrice, fromProvince, toProvince *location.Province, toDistrict *location.District, fromCustomRegions, toCustomRegions []dot.ID) bool {
	if pricing.Status != status3.P {
		return false
	}
	var _fromCustomRegions, _toCustomRegions []dot.ID
	if len(pricing.CustomRegionIDs) > 0 {
		// yêu cầu cả địa điểm gửi và lấy hàng đều phải nằm trong danh sách CustomRegionIDs này
		for _, fromRegionID := range fromCustomRegions {
			if cm.IDsContain(pricing.CustomRegionIDs, fromRegionID) {
				_fromCustomRegions = append(_fromCustomRegions, fromRegionID)
			}
		}
		for _, toRegionID := range toCustomRegions {
			if cm.IDsContain(pricing.CustomRegionIDs, toRegionID) {
				_toCustomRegions = append(_toCustomRegions, toRegionID)
			}
		}
		if len(_fromCustomRegions) == 0 || len(_toCustomRegions) == 0 {
			return false
		}
	}

	// check CustomRegionRouteType
	// CustomRegionRouteType và CustomRegionIDs luôn đi chung với nhau
	if len(pricing.CustomRegionTypes) > 0 && len(_fromCustomRegions) > 0 {
		checkValidRegion := false
		for _, fromRegion := range _fromCustomRegions {
			for _, toRegion := range _toCustomRegions {
				customRegionRouteType := locationutil.GetCustomRegionRouteType(fromRegion, toRegion)
				if containCustomRegionRouteType(pricing.CustomRegionTypes, customRegionRouteType) {
					checkValidRegion = true
					break
				}
			}
			if checkValidRegion {
				break
			}
		}
		if !checkValidRegion {
			return false
		}
	}

	// check RegionRouteType
	if len(pricing.RegionTypes) > 0 {
		regionRouteTypes := locationutil.GetRegionRouteTypes(fromProvince, toProvince)
		matchRegion := false
		for _, regionType := range regionRouteTypes {
			if containRegionRouteType(pricing.RegionTypes, regionType) {
				matchRegion = true
				break
			}
		}
		if !matchRegion {
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

// GetPriceByPricingDetail
//
// Main fee
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
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "The weight is out of range")
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
		if ov.MaxWeight == shipmentprice.MaximumValue {
			return true
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
	if weight < ov.MinWeight {
		return 0
	}
	if weight == ov.MinWeight {
		return ov.PriceStep
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

func (q *QueryService) CalculateShippingFees(ctx context.Context, args *shipmentprice.CalculateShippingFeesArgs) (*shipmentprice.CalculateShippingFeesResponse, error) {
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

	calcShippingFeeArgs := calculateShippingFeesArgs{
		From:                from,
		To:                  to,
		AccountID:           args.AccountID,
		ConnectionID:        args.ConnectionID,
		ShipmentServiceID:   args.ShipmentServiceID,
		ShipmentPriceListID: args.ShipmentPriceListID,
		Weight:              args.Weight,
		BasketValue:         args.BasketValue,
		CODAmount:           args.CODAmount,
		AdditionalFeeTypes:  args.AdditionalFeeTypes,
	}
	// Trường hợp có bảng giá khuyến mãi (PromotionPriceListID)
	// Gán ShipmentPriceListID = PromotionPriceListID
	// Nếu tính được giá => trả về kết quả
	// Nếu không tính được giá => tính như bình thường
	if args.PromotionPriceListID != 0 {
		calcShippingFeeArgs.ShipmentPriceListID = args.PromotionPriceListID
		if resp, err := q.calcShippingFees(ctx, calcShippingFeeArgs); err == nil {
			return resp, nil
		}
		// Gán lại ShipmentPriceListID để tính lại giá
		calcShippingFeeArgs.ShipmentPriceListID = args.ShipmentPriceListID
	}

	return q.calcShippingFees(ctx, calcShippingFeeArgs)
}

type calculateShippingFeesArgs struct {
	From                *location.LocationQueryResult
	To                  *location.LocationQueryResult
	AccountID           dot.ID
	ConnectionID        dot.ID
	ShipmentServiceID   dot.ID
	ShipmentPriceListID dot.ID
	Weight              int
	BasketValue         int
	CODAmount           int
	AdditionalFeeTypes  []shipping_fee_type.ShippingFeeType
}

func (q *QueryService) calcShippingFees(ctx context.Context, args calculateShippingFeesArgs) (*shipmentprice.CalculateShippingFeesResponse, error) {
	getActiveShipmentPricesArgs := getActiveShipmentPricesArgs{
		AccountID:           args.AccountID,
		ConnectionID:        args.ConnectionID,
		ShipmentServiceID:   args.ShipmentServiceID,
		ShipmentPriceListID: args.ShipmentPriceListID,
	}
	pricings, err := q.GetActiveShipmentPrices(ctx, getActiveShipmentPricesArgs)
	if err != nil {
		return nil, err
	}

	if len(pricings) == 0 {
		return nil, cm.Errorf(cm.NotFound, nil, "Chưa cấu hình giá cho connection này (conn_id = %v)", args.ConnectionID)
	}

	matchingPricings, err := q.GetMatchingPricings(ctx, pricings, args.From.Province, args.To.Province, args.To.District)
	if len(matchingPricings) == 0 || err != nil {
		// Không thay đổi mã lỗi (yêu cầu: ko sử dụng mã not_found)
		// not_found chỉ sử dụng khi không có cấu hình giá nào trong bảng giá
		return nil, cm.Errorf(cm.FailedPrecondition, err, "Không có cấu hình giá phù hợp")
	}

	// calculate main fee
	pricing := getPricingByPriorityPoint(matchingPricings)
	pRuleDetail, err := GetPriceRuleDetail(args.Weight, pricing.Details)
	if err != nil {
		return nil, err
	}
	mainFee, err := GetPriceByPricingDetail(args.Weight, pRuleDetail)
	if err != nil {
		return nil, err
	}

	// calculate additional fee
	calcAdditionalFeeArgs := CalcAdditionalFeeArgs{
		BasketValue:        args.BasketValue,
		CODAmount:          args.CODAmount,
		MainFee:            mainFee,
		AdditionalFeeTypes: args.AdditionalFeeTypes,
	}
	feeLines, err := calcAdditionalFees(calcAdditionalFeeArgs, pricing.AdditionalFees)
	if err != nil {
		return nil, cm.Errorf(cm.FailedPrecondition, err, "Không thể tính phí từ cấu hình giá shipment_price = %v", pricing.ID)
	}
	feeLines = append(feeLines, &shipmentprice.ShippingFee{
		FeeType: shipping_fee_type.Main,
		Price:   mainFee,
	})

	totalFee := 0
	for _, line := range feeLines {
		totalFee += line.Price
	}

	return &shipmentprice.CalculateShippingFeesResponse{
		ShipmentPriceID:     pricing.ID,
		ShipmentPriceListID: pricing.ShipmentPriceListID,
		TotalFee:            totalFee,
		FeeLines:            feeLines,
	}, nil
}
