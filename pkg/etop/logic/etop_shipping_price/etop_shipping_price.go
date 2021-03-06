package etop_shipping_price

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"

	"o.o/api/main/location"
	"o.o/api/top/types/etc/shipping_provider"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/code/gencode"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
	"o.o/common/l"
)

type ESPriceRule struct {
	Carrier  shipping_provider.ShippingProvider
	Pricings []*ESPricing
}

type ESPricing struct {
	// Type: Nhanh | Chuẩn
	Type          string
	FromProvince  *FromProvinceDetail
	RouteType     RouteTypeDetail
	DistrictTypes []location.UrbanType
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
	ID         int
	Weight     int
	Price      int
	Overweight []*ESPricingDetailOverweightPrice
}

type ESPricingDetailOverweightPrice struct {
	MinWeight int
	MaxWeight int
	// WeightStep: trọng lượng tăng thêm (g)
	WeightStep int
	// PriceStep: Mức giá của phần trọng lượng tăng thêm
	// VD: Step: 500, PriceStep: 4500
	// Cứ 500g tăng thêm 4500đ
	PriceStep int
}

var (
	ll             = l.New()
	priceRuleIndex = make(map[shipping_provider.ShippingProvider][]*ESPricing)
)

func init() {
	model.GetShippingServiceRegistry().RegisterNameFunc(shipping_provider.Etop, DecodeShippingServiceName)
	for _, rule := range ESPriceRules {
		key := rule.Carrier
		if priceRuleIndex[key] != nil {
			ll.Warn("Etop Shipping rule: Duplicate rule", l.Any("key", key))
			continue
		}
		priceRuleIndex[key] = rule.Pricings
	}
}

type GetEtopShippingServicesArgs struct {
	ArbitraryID  dot.ID
	Carrier      shipping_provider.ShippingProvider
	FromProvince *location.Province
	ToProvince   *location.Province
	ToDistrict   *location.District
	Weight       int
}

func GetEtopShippingServices(args *GetEtopShippingServicesArgs) []*shippingsharemodel.AvailableShippingService {
	var res []*shippingsharemodel.AvailableShippingService
	pricings := priceRuleIndex[args.Carrier]
	pricingsMatch := GetESPricingsMatch(pricings, args.FromProvince, args.ToProvince, args.ToDistrict)

	generator := newServiceIDGenerator(int64(args.ArbitraryID), args.Carrier)
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

func (pricing *ESPricing) ToService(generator serviceIDGenerator, weight int, carrier shipping_provider.ShippingProvider) (*shippingsharemodel.AvailableShippingService, error) {
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
	return &shippingsharemodel.AvailableShippingService{
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
	for wIndex := range priceRuleDetails {
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
	if (pRuleDetail.Overweight == nil || len(pRuleDetail.Overweight) == 0) && weight > pRuleDetail.Weight {
		// can not apply this rule
		return 0, cm.Error(cm.InvalidArgument, "Can not apply to this rule", nil)
	}

	addWeight := weight - pRuleDetail.Weight
	if addWeight <= 0 {
		return pRuleDetail.Price, nil
	}

	price := pRuleDetail.Price
	for _, ov := range pRuleDetail.Overweight {
		ovPrice := GetOverweightPrice(weight, ov)
		price += ovPrice
	}
	return price, nil
}

func GetOverweightPrice(weight int, ov *ESPricingDetailOverweightPrice) int {
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

type serviceIDGenerator struct {
	rd *rand.Rand
}

func newServiceIDGenerator(seed int64, carrier shipping_provider.ShippingProvider) serviceIDGenerator {
	// make sure generate difference code with difference carrier
	switch carrier {
	case shipping_provider.GHTK:
		seed += 1
	case shipping_provider.VTPost:
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

func FillInfoEtopServices(providerServices []*shippingsharemodel.AvailableShippingService, etopServices []*shippingsharemodel.AvailableShippingService) ([]*shippingsharemodel.AvailableShippingService, error) {
	if len(etopServices) == 0 {
		return nil, nil
	}
	if len(providerServices) == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Không đủ thông tin", nil)
	}

	serviceTypeIndex := make(map[string]*shippingsharemodel.AvailableShippingService)
	for _, service := range providerServices {
		key := fmt.Sprintf("%v_%v", service.Provider.String(), service.Name)
		if serviceTypeIndex[key] == nil {
			serviceTypeIndex[key] = service
		}
	}

	for _, service := range etopServices {
		if service.Source != model.TypeShippingSourceEtop {
			continue
		}
		key := fmt.Sprintf("%v_%v", service.Provider.String(), service.Name)
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

func GetShippingDistrictType(district *location.District) location.UrbanType {
	switch district.UrbanType {
	case location.Urban:
		return location.Urban
	case location.Suburban1:
		return location.Suburban1
	case location.Suburban2:
		return location.Suburban2
	default:
		return location.Suburban2
	}
}
