package etop_shipping_price

import (
	"context"
	"testing"

	"etop.vn/api/main/location"
	servicelocation "etop.vn/backend/com/main/location"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/common/l"

	"github.com/stretchr/testify/assert"
)

var locationBus = servicelocation.New().MessageBus()

func TestGetPriceRuleDetail(t *testing.T) {
	priceRuleDetail := map[int]*ESPricingDetail{
		100: {
			ID:     100,
			Weight: 100,
			Price:  20000,
		},
		300: {
			ID:     300,
			Weight: 300,
			Price:  20000,
		},
		500: {
			ID:     500,
			Weight: 500,
			Price:  20000,
		},
		3000: {
			ID:     3000,
			Weight: 3000,
			Price:  20000,
			Overweight: []*ESPricingDetailOverweightPrice{
				{
					MinWeight:  3000,
					MaxWeight:  -1,
					WeightStep: 500,
					PriceStep:  4500,
				},
			},
		},
	}
	priceRuleDetail2 := map[int]*ESPricingDetail{
		100: {
			Weight: 100,
			Price:  27500,
			ID:     100,
		},
		300: {
			ID:     300,
			Weight: 300,
			Price:  30000,
		},
		500: {
			ID:     500,
			Weight: 500,
			Price:  34000,
		},
	}

	for _, tt := range []struct {
		weight    int
		priceRule map[int]*ESPricingDetail
		expect    int
	}{
		{0, priceRuleDetail, 100},
		{100, priceRuleDetail, 100},
		{200, priceRuleDetail, 300},
		{300, priceRuleDetail, 300},
		{400, priceRuleDetail, 500},
		{500, priceRuleDetail, 500},
		{1200, priceRuleDetail, 3000},
		{1800, priceRuleDetail, 3000},
		{3000, priceRuleDetail, 3000},
		{3500, priceRuleDetail, 3000},
		{5400, priceRuleDetail, 3000},
		{0, priceRuleDetail2, 100},
		{100, priceRuleDetail2, 100},
		{200, priceRuleDetail2, 300},
		{300, priceRuleDetail2, 300},
		{400, priceRuleDetail2, 500},
		{500, priceRuleDetail2, 500},
		{1200, priceRuleDetail2, 500},
		{1800, priceRuleDetail2, 500},
		{3000, priceRuleDetail2, 500},
		{4300, priceRuleDetail2, 500},
	} {

		t.Run("TestGetPriceRuleDetail", func(t *testing.T) {
			rule := GetPriceRuleDetail(tt.weight, tt.priceRule)
			assert.EqualValues(t, tt.expect, rule.ID)
		})
	}
}

func TestGetPriceByPricingDetail(t *testing.T) {
	ruleDetail1 := &ESPricingDetail{
		ID:     300,
		Weight: 300,
		Price:  26000,
	}
	ruleDetail2 := &ESPricingDetail{
		ID:     500,
		Weight: 500,
		Price:  20000,
		Overweight: []*ESPricingDetailOverweightPrice{
			{
				MinWeight:  500,
				MaxWeight:  4000,
				WeightStep: 500,
				PriceStep:  5000,
			},
			{
				MinWeight:  4000,
				MaxWeight:  -1,
				WeightStep: 500,
				PriceStep:  7000,
			},
		},
	}
	ruleDetail3 := &ESPricingDetail{
		ID:     1000,
		Weight: 1000,
		Price:  35000,
		Overweight: []*ESPricingDetailOverweightPrice{
			{
				MinWeight:  4000,
				MaxWeight:  -1,
				WeightStep: 500,
				PriceStep:  7000,
			},
			{
				MinWeight:  1000,
				MaxWeight:  4000,
				WeightStep: 500,
				PriceStep:  5000,
			},
		},
	}
	for _, tt := range []struct {
		weight  int
		rDetail *ESPricingDetail
		expect  int
	}{
		{0, ruleDetail1, 26000},
		{100, ruleDetail1, 26000},
		{200, ruleDetail1, 26000},
		{300, ruleDetail1, 26000},
		{400, ruleDetail1, 0},
		{500, ruleDetail1, 0},
		{1200, ruleDetail1, 0},
		{0, ruleDetail2, 20000},
		{100, ruleDetail2, 20000},
		{300, ruleDetail2, 20000},
		{500, ruleDetail2, 20000},
		{1200, ruleDetail2, 30000},
		{1800, ruleDetail2, 35000},
		{3000, ruleDetail2, 45000},
		{3600, ruleDetail2, 55000},
		{4000, ruleDetail2, 55000},
		{5400, ruleDetail2, 76000},
		{4000, ruleDetail3, 65000},
		{600, ruleDetail3, 35000},
		{5000, ruleDetail3, 79000},
		{4100, ruleDetail3, 72000},
	} {
		t.Run("TestGetPriceByPricingDetail2", func(t *testing.T) {
			price, _ := GetPriceByPricingDetail(tt.weight, tt.rDetail)
			assert.EqualValues(t, tt.expect, price)
		})
	}
}

func FindLocation(province, district, ward string) *location.LocationQueryResult {
	query := &location.FindLocationQuery{
		Province: province,
		District: district,
		Ward:     ward,
	}
	if err := locationBus.Dispatch(context.Background(), query); err != nil {
		ll.Panic("unexpected", l.Error(err))
	}
	return query.Result
}

func TestGetEtopShippingServices(t *testing.T) {
	hcm1 := FindLocation("Ho Chi Minh", "quan 10", "")
	hcm2 := FindLocation("Ho Chi Minh", "Binh Thanh", "")
	hn := FindLocation("Ha Noi", "Huyện Ba Vì", "")
	dn := FindLocation("Da Nang", "Quận Cẩm Lệ", "")
	ag := FindLocation("An Giang", "Thành phố Long Xuyên", "")

	_hcm := hcm1.Province
	_quan10 := hcm1.District
	_binhthanh := hcm2.District
	_hn := hn.Province
	_bavi := hn.District
	_dn := dn.Province
	_camle := dn.District
	_ag := ag.Province
	_longxuyen := ag.District

	for _, tt := range []struct {
		carrier            model.ShippingProvider
		fromProvince       *location.Province
		toProvince         *location.Province
		toDistrict         *location.District
		weight             int
		expectServiceTypes []string
		expectPrices       []int
	}{
		{model.TypeGHTK, _hcm, _hn, _bavi, 100, nil, nil},
		{model.TypeGHTK, _hcm, _ag, _longxuyen, 200, []string{model.ShippingServiceNameStandard}, []int{24000}},
		{model.TypeGHTK, _hcm, _dn, _camle, 400, nil, nil},
		{model.TypeGHTK, _hcm, _dn, _camle, 1400, nil, nil},
		{model.TypeGHTK, _hcm, _hcm, _binhthanh, 3400, nil, nil},
		{model.TypeGHTK, _dn, _hcm, _quan10, 700, nil, nil},
		{model.TypeGHTK, _dn, _hn, _bavi, 200, nil, nil},
		{model.TypeGHTK, _ag, _dn, _longxuyen, 400, nil, nil},
		{model.TypeGHTK, _ag, _hn, _bavi, 400, nil, nil},
		{model.TypeGHTK, _hn, _dn, _camle, 3400, nil, nil},
		{model.TypeGHTK, _hn, _ag, _longxuyen, 100, []string{model.ShippingServiceNameStandard}, []int{24000}},

		{model.TypeGHN, _hn, _dn, _camle, 200, []string{model.ShippingServiceNameStandard}, []int{32000}},
		{model.TypeGHN, _dn, _hcm, _quan10, 700, []string{model.ShippingServiceNameStandard}, []int{32000}},
		{model.TypeGHN, _hcm, _ag, _longxuyen, 200, []string{model.ShippingServiceNameStandard}, []int{32000}},
		{model.TypeGHN, _hn, _dn, _camle, 2200, []string{model.ShippingServiceNameStandard}, []int{47000}},
		{model.TypeGHN, _hcm, _hcm, _quan10, 2200, []string{model.ShippingServiceNameStandard}, []int{15500}},
	} {
		t.Run("GetEtopShippingServices", func(t *testing.T) {
			args := &GetEtopShippingServicesArgs{
				ArbitraryID:  0,
				Carrier:      tt.carrier,
				FromProvince: tt.fromProvince,
				ToProvince:   tt.toProvince,
				ToDistrict:   tt.toDistrict,
				Weight:       tt.weight,
			}
			eServices := GetEtopShippingServices(args)

			if len(tt.expectServiceTypes) > 0 {
				assert.Contains(t, tt.expectServiceTypes, eServices[0].Name)
				assert.Contains(t, tt.expectPrices, eServices[0].ServiceFee)
			} else {
				assert.EqualValues(t, 0, len(eServices))
			}
		})
	}
}
