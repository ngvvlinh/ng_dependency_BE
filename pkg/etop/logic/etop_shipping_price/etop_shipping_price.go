package etop_shipping_price

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"

	"etop.vn/api/main/location"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/gencode"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/etop/model"
)

type ESPriceRule struct {
	Carrier  model.ShippingProvider
	Pricings []*ESPricing
}

type ESPricing struct {
	// Type: Nhanh | Chuẩn
	Type          string
	FromProvince  *FromProvinceDetail
	RouteType     RouteTypeDetail
	DistrictTypes []model.ShippingDistrictType
	Details       map[int]*ESPricingDetail
}

type FromProvinceDetail struct {
	IncludeCode []string
	ExcludeCode []string
}

type RouteTypeDetail struct {
	Include model.ShippingRouteType
	Exclude []model.ShippingRouteType
}

type ESPricingDetail struct {
	ID     int
	Weight int
	Price  int
	// WeightStep: trọng lượng tăng thêm (g)
	WeightStep int
	// PriceStep: Mức giá của phần trọng lượng tăng thêm
	// VD: Step: 500, PriceStep: 4500
	// Cứ 500g tăng thêm 4500đ
	PriceStep int
}

var (
	ll             = l.New()
	priceRuleIndex = make(map[model.ShippingProvider][]*ESPricing)
)

func init() {
	model.GetShippingServiceRegistry().RegisterNameFunc(model.TypeShippingETOP, DecodeShippingServiceName)
	for _, rule := range ESPriceRules {
		key := rule.Carrier
		if priceRuleIndex[key] != nil {
			ll.Warn("Etop Shipping rule: Duplicate rule", l.String("key", string(key)))
			continue
		}
		priceRuleIndex[key] = rule.Pricings
	}
}

type GetEtopShippingServicesArgs struct {
	ArbitraryID  int64
	Carrier      model.ShippingProvider
	FromProvince *location.Province
	ToProvince   *location.Province
	ToDistrict   *location.District
	Weight       int
}

func GetEtopShippingServices(args *GetEtopShippingServicesArgs) []*model.AvailableShippingService {
	var res []*model.AvailableShippingService
	pricings := priceRuleIndex[args.Carrier]
	pricingsMatch := GetESPricingsMatch(pricings, args.FromProvince, args.ToProvince, args.ToDistrict)

	generator := newServiceIDGenerator(args.ArbitraryID, args.Carrier)
	for _, price := range pricingsMatch {
		if service, err := price.ToService(generator, args.Weight, args.Carrier); err == nil {
			res = append(res, service)
		}
	}
	return res
}

func GetESPricingsMatch(pricings []*ESPricing, fromProvince *location.Province, toProvince *location.Province, toDistrict *location.District) []*ESPricing {
	res := make([]*ESPricing, 0, len(pricings))
	for _, price := range pricings {
		if price.CheckESPriceMatch(fromProvince, toProvince, toDistrict) {
			res = append(res, price)
		}
	}
	return res
}

func (pricing *ESPricing) CheckESPriceMatch(fromProvince *location.Province, toProvince *location.Province, toDistrict *location.District) bool {
	// check from province
	if pricing.FromProvince != nil {
		provincesApply := pricing.FromProvince.IncludeCode
		provincesNotApply := pricing.FromProvince.ExcludeCode
		if len(provincesApply) > 0 && !cm.StringsContain(provincesApply, fromProvince.Code) {
			return false
		}
		if len(provincesNotApply) > 0 && cm.StringsContain(provincesNotApply, fromProvince.Code) {
			return false
		}
	}
	// check Route type
	routes := GetShippingRouteTypes(fromProvince, toProvince)
	foundRouteInclude, foundRouteExclude := false, false
	for _, route := range routes {
		if route == pricing.RouteType.Include {
			foundRouteInclude = true
		}
		if ContainRouteType(pricing.RouteType.Exclude, route) {
			foundRouteExclude = true
			break
		}
	}
	if !foundRouteInclude || foundRouteExclude {
		return false
	}

	// check District Type
	if pricing.DistrictTypes != nil {
		check := false
		dType := GetShippingDistrictType(toDistrict)
		for _, districtType := range pricing.DistrictTypes {
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

func ContainRouteType(types []model.ShippingRouteType, routeType model.ShippingRouteType) bool {
	for _, rt := range types {
		if routeType == rt {
			return true
		}
	}
	return false
}

func (pricing *ESPricing) ToService(generator serviceIDGenerator, weight int, carrier model.ShippingProvider) (*model.AvailableShippingService, error) {
	pRuleDetail := GetPriceRuleDetail(weight, pricing.Details)
	if pRuleDetail == nil {
		return nil, cm.Error(cm.Internal, "Không có bảng giá phù hợp", nil)
	}
	price, err := GetPriceByPricingDetail(weight, pRuleDetail)
	if err != nil {
		return nil, err
	}

	serviceID, err := generator.GenerateEtopServiceCode(pricing.Type)
	if err != nil {
		return nil, err
	}
	return &model.AvailableShippingService{
		Name:              pricing.Type,
		ProviderServiceID: serviceID,
		ServiceFee:        price,
		ShippingFeeMain:   price,
		Provider:          carrier,
		Source:            model.TypeShippingSourceEtop,
	}, nil
}

func GetPriceRuleDetail(weight int, priceRuleDetails map[int]*ESPricingDetail) *ESPricingDetail {
	var weightIndex []int
	for wIndex, _ := range priceRuleDetails {
		weightIndex = append(weightIndex, wIndex)
	}
	sort.Ints(weightIndex) // increase
	var index int
	for _, wIndex := range weightIndex {
		index = wIndex
		if weight <= wIndex {
			break
		}
	}

	return priceRuleDetails[index]
}

func GetPriceByPricingDetail(weight int, pRuleDetail *ESPricingDetail) (int, error) {
	if pRuleDetail.PriceStep == 0 && weight > pRuleDetail.Weight {
		// can not apply this rule
		return 0, cm.Error(cm.InvalidArgument, "Can not apply to this rule", nil)
	}
	addWeight := weight - pRuleDetail.Weight
	if addWeight <= 0 {
		return pRuleDetail.Price, nil
	}
	s := float64(addWeight) / float64(pRuleDetail.WeightStep)
	step := int(math.Ceil(s))
	price := pRuleDetail.Price + step*pRuleDetail.PriceStep
	return price, nil
}

type serviceIDGenerator struct {
	rd *rand.Rand
}

func newServiceIDGenerator(seed int64, carrier model.ShippingProvider) serviceIDGenerator {
	// make sure generate difference code with difference carrier
	switch carrier {
	case model.TypeGHTK:
		seed += 1
	case model.TypeVTPost:
		seed += 2
	default:

	}

	src := rand.NewSource(seed)
	rd := rand.New(src)
	return serviceIDGenerator{rd}
}

// GenerateEtopServiceCode: generate service ID for ETOP shipping
// format ex: 7ETOP20N
func (c serviceIDGenerator) GenerateEtopServiceCode(shippingType string) (string, error) {
	n := c.rd.Uint64()
	v := gencode.Alphabet32.EncodeReverse(n, 8)
	code := string(v[:8])
	// code := gencode.GenerateCode(gencode.Alphabet32, 8)
	switch shippingType {
	case model.ShippingServiceNameStandard:
		code = code[:7] + "C"
	case model.ShippingServiceNameFaster:
		code = code[:7] + "N"
	default:
		return "", cm.Error(cm.InvalidArgument, "Shipping service type is invalid", nil).WithMeta("type", shippingType)
	}
	return code[:1] + "ETOP" + code[5:], nil
}

func ParseEtopServiceCode(serviceCode string) (shippingType string, ok bool) {
	if strings.Index(serviceCode, "ETOP") != 1 {
		return "", false
	}
	switch serviceCode[7] {
	case 'C':
		return model.ShippingServiceNameStandard, true
	case 'N':
		return model.ShippingServiceNameFaster, true
	default:
		return "", false
	}
}

func DecodeShippingServiceName(code string) (name string, ok bool) {
	return ParseEtopServiceCode(code)
}

func FillInfoEtopServices(providerServices []*model.AvailableShippingService, etopServices []*model.AvailableShippingService) ([]*model.AvailableShippingService, error) {
	if len(etopServices) == 0 {
		return nil, nil
	}
	if len(providerServices) == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Không đủ thông tin", nil)
	}

	serviceTypeIndex := make(map[string]*model.AvailableShippingService)
	for _, service := range providerServices {
		key := fmt.Sprintf("%v_%v", string(service.Provider), service.Name)
		if serviceTypeIndex[key] == nil {
			serviceTypeIndex[key] = service
		}
	}

	for _, service := range etopServices {
		if service.Source != model.TypeShippingSourceEtop {
			continue
		}
		key := fmt.Sprintf("%v_%v", string(service.Provider), service.Name)
		if serviceTypeIndex[key] != nil {
			s := serviceTypeIndex[key]
			service.ExpectedPickAt = s.ExpectedPickAt
			service.ExpectedDeliveryAt = s.ExpectedDeliveryAt
			// update additionFee for etopService (insurance_fee, ...)
			additionFee := s.ServiceFee - s.ShippingFeeMain
			service.ServiceFee = service.ShippingFeeMain + additionFee
		}
	}
	return etopServices, nil
}

func GetShippingRouteTypes(fromProvince *location.Province, toProvince *location.Province) []model.ShippingRouteType {
	var res []model.ShippingRouteType
	if fromProvince.Code == toProvince.Code {
		res = append(res, model.RouteSameProvince)
	}
	if fromProvince.Region == toProvince.Region {
		res = append(res, model.RouteSameRegion)
	}
	res = append(res, model.RouteNationWide)
	return res
}

func GetShippingDistrictType(district *location.District) model.ShippingDistrictType {
	switch district.UrbanType {
	case location.Urban:
		return model.ShippingDistrictTypeUrban
	case location.Suburban1:
		return model.ShippingDistrictTypeSubUrban1
	case location.Suburban2:
		return model.ShippingDistrictTypeSubUrban2
	default:
		return model.ShippingDistrictTypeSubUrban2
	}
}
