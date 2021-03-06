package etop_shipping_price

import (
	"o.o/api/main/location"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/backend/pkg/etop/model"
)

var ESPriceRules = []*ESPriceRule{
	{
		Carrier: shipping_provider.GHTK,
		Pricings: []*ESPricing{
			{
				// gói chuẩn nội miền và toàn quốc: 200g
				// áp dụng cho nội thành và ngoại thành 1
				Type: model.ShippingServiceNameStandard,
				RouteType: RouteTypeDetail{
					Include: model.RouteNationWide,
					Exclude: []model.ShippingRouteType{model.RouteSameProvince},
				},
				DistrictTypes: []location.UrbanType{location.Urban, location.Suburban1},
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
		Carrier: shipping_provider.GHN,
		Pricings: []*ESPricing{
			{
				Type: model.ShippingServiceNameStandard,
				RouteType: RouteTypeDetail{
					Include: model.RouteSameProvince,
					Exclude: nil,
				},
				DistrictTypes: []location.UrbanType{
					location.Urban,
				},
				Details: map[int]*ESPricingDetail{
					3000: {
						ID:     3000,
						Weight: 3000,
						Price:  15500,
						Overweight: []*ESPricingDetailOverweightPrice{
							{
								MinWeight:  3000,
								MaxWeight:  4000,
								WeightStep: 500,
								PriceStep:  2500,
							},
							{
								MinWeight:  4000,
								MaxWeight:  -1,
								WeightStep: 500,
								PriceStep:  4000,
							},
						},
					},
				},
			},
			{
				Type: model.ShippingServiceNameStandard,
				RouteType: RouteTypeDetail{
					Include: model.RouteSameProvince,
					Exclude: nil,
				},
				DistrictTypes: []location.UrbanType{
					location.Suburban1, location.Suburban2,
				},
				Details: map[int]*ESPricingDetail{
					3000: {
						ID:     3000,
						Weight: 3000,
						Price:  28000,
						Overweight: []*ESPricingDetailOverweightPrice{
							{
								MinWeight:  3000,
								MaxWeight:  4000,
								WeightStep: 500,
								PriceStep:  2500,
							},
							{
								MinWeight:  4000,
								MaxWeight:  -1,
								WeightStep: 500,
								PriceStep:  4000,
							},
						},
					},
				},
			},
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
						Overweight: []*ESPricingDetailOverweightPrice{
							{
								MinWeight:  1000,
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
					},
				},
			},
		},
	},
}
