package etop_shipping_price

import "etop.vn/backend/pkg/etop/model"

var ESPriceRules = []*ESPriceRule{
	{
		Carrier: model.TypeGHTK,
		Pricings: []*ESPricing{
			{
				// gói chuẩn nội miền và toàn quốc: 200g
				// áp dụng cho nội thành và ngoại thành 1
				Type: model.ShippingServiceNameStandard,
				RouteType: RouteTypeDetail{
					Include: model.RouteNationWide,
					Exclude: []model.ShippingRouteType{model.RouteSameProvince},
				},
				DistrictTypes: []model.ShippingDistrictType{model.ShippingDistrictTypeUrban, model.ShippingDistrictTypeSubUrban1},
				Details: map[int]*ESPricingDetail{
					200: {
						ID:     200,
						Weight: 200,
						Price:  24000,
					},
				},
			},
		},
	},

	{
		Carrier: model.TypeGHN,
		Pricings: []*ESPricing{
			{
				// gói chuẩn toàn quốc <= 1kg
				// trừ khu vực nội tỉnh
				Type: model.ShippingServiceNameStandard,
				RouteType: RouteTypeDetail{
					Include: model.RouteNationWide,
					Exclude: []model.ShippingRouteType{model.RouteSameProvince},
				},
				DistrictTypes: nil,
				Details: map[int]*ESPricingDetail{
					1000: {
						ID:     1000,
						Weight: 1000,
						Price:  32000,
					},
				},
			},
		},
	},
}
