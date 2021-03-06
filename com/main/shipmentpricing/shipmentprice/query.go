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

	- Field `shipmentPriceListID`: d??ng ????? t??nh gi?? cho m???t b???ng gi?? c??? th??? (s??? d???ng trong tr?????ng h???p admin ki???m tra gi?? tr?????c khi public b???ng gi??)
	- Khi shop s??? d???ng & `shipmentPriceListID` = 0:
		- T??m b???ng gi?? c???a shop ???? trong b???ng `shop_shipment_price_list`, n???u c??, s??? d???ng b???ng gi?? ????
		- N???u shop ko c?? b???ng gi?? ri??ng, s??? d???ng b???ng gi?? m???c ?????nh (GetActiveShipmentPriceListQuery)
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
// L???y t???t c??? c??c b???ng gi?? m???c ?????nh c???a c??c NVC (??ang c?? is_default = true)
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
		// y??u c???u c??? ?????a ??i???m g???i v?? l???y h??ng ?????u ph???i n???m trong danh s??ch CustomRegionIDs n??y
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
	// CustomRegionRouteType v?? CustomRegionIDs lu??n ??i chung v???i nhau
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
// Ki???m tra ??i???u ki???n v?????t c??n c?? h???p l??? hay kh??ng
// Kh???i l?????ng ph???i n???m trong kho???ng c??n ???????c setting
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
	// Tr?????ng h???p c?? b???ng gi?? khuy???n m??i (PromotionPriceListID)
	// G??n ShipmentPriceListID = PromotionPriceListID
	// N???u t??nh ???????c gi?? => tr??? v??? k???t qu???
	// N???u kh??ng t??nh ???????c gi?? => t??nh nh?? b??nh th?????ng
	if args.PromotionPriceListID != 0 {
		calcShippingFeeArgs.ShipmentPriceListID = args.PromotionPriceListID
		if resp, err := q.calcShippingFees(ctx, calcShippingFeeArgs); err == nil {
			return resp, nil
		}
		// G??n l???i ShipmentPriceListID ????? t??nh l???i gi??
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
		return nil, cm.Errorf(cm.NotFound, nil, "Ch??a c???u h??nh gi?? cho connection n??y (conn_id = %v)", args.ConnectionID)
	}

	matchingPricings, err := q.GetMatchingPricings(ctx, pricings, args.From.Province, args.To.Province, args.To.District)
	if len(matchingPricings) == 0 || err != nil {
		// Kh??ng thay ?????i m?? l???i (y??u c???u: ko s??? d???ng m?? not_found)
		// not_found ch??? s??? d???ng khi kh??ng c?? c???u h??nh gi?? n??o trong b???ng gi??
		return nil, cm.Errorf(cm.FailedPrecondition, err, "Kh??ng c?? c???u h??nh gi?? ph?? h???p")
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
		return nil, cm.Errorf(cm.FailedPrecondition, err, "Kh??ng th??? t??nh ph?? t??? c???u h??nh gi?? shipment_price = %v", pricing.ID)
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
