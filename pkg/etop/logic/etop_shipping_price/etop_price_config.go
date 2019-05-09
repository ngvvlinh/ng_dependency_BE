package etop_shipping_price

import "etop.vn/backend/pkg/etop/model"

var ESPriceRules = []*ESPriceRule{
	{
		Carrier: model.TypeGHTK,
		Pricings: []*ESPricing{
			{
				Type: model.ShippingServiceNameFaster,
				FromProvince: &FromProvinceDetail{
					IncludeCode: []string{"01", "79"},
				},
				RouteType: RouteTypeDetail{
					Include: model.RouteSameProvince,
				},
				DistrictTypes: []model.ShippingDistrictType{model.ShippingDistrictTypeUrban, model.ShippingDistrictTypeSubUrban1},
				Details: map[int]*ESPricingDetail{
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
						ID:         3000,
						Weight:     3000,
						Price:      20000,
						WeightStep: 500,
						PriceStep:  4500,
					},
				},
			},
			{
				Type: model.ShippingServiceNameFaster,
				FromProvince: &FromProvinceDetail{
					IncludeCode: []string{"01", "79"},
				},
				RouteType: RouteTypeDetail{
					Include: model.RouteSameProvince,
				},
				DistrictTypes: []model.ShippingDistrictType{model.ShippingDistrictTypeSubUrban2},
				Details: map[int]*ESPricingDetail{
					100: {
						ID:     100,
						Weight: 100,
						Price:  32000,
					},
					300: {
						ID:     300,
						Weight: 300,
						Price:  32000,
					},
					500: {
						ID:     500,
						Weight: 500,
						Price:  32000,
					},
					3000: {
						ID:         3000,
						Weight:     3000,
						Price:      32000,
						WeightStep: 500,
						PriceStep:  4500,
					},
				},
			},
			{
				Type: model.ShippingServiceNameStandard,
				FromProvince: &FromProvinceDetail{
					IncludeCode: []string{"01", "79"},
				},
				RouteType: RouteTypeDetail{
					Include: model.RouteNationWide,
					Exclude: model.RouteSameProvince,
				},
				DistrictTypes: nil,
				Details: map[int]*ESPricingDetail{
					100: {
						ID:     100,
						Weight: 100,
						Price:  24000,
					},
					300: {
						ID:     300,
						Weight: 300,
						Price:  26000,
					},
					500: {
						ID:         500,
						Weight:     500,
						Price:      29000,
						WeightStep: 500,
						PriceStep:  4500,
					},
				},
			},
			{
				Type: model.ShippingServiceNameFaster,
				FromProvince: &FromProvinceDetail{
					IncludeCode: []string{"01", "79"},
				},
				RouteType: RouteTypeDetail{
					Include: model.RouteNationWide,
					Exclude: model.RouteSameRegion,
				},
				DistrictTypes: nil,
				Details: map[int]*ESPricingDetail{
					100: {
						ID:     100,
						Weight: 100,
						Price:  27500,
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
				},
			},
		},
	},
}
