package model

import (
	"time"

	"o.o/api/top/types/etc/additional_fee_base_value"
	"o.o/api/top/types/etc/calculation_method"
	"o.o/api/top/types/etc/price_modifier_type"
	"o.o/api/top/types/etc/route_type"
	"o.o/api/top/types/etc/shipping_fee_type"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +sqlgen
type ShipmentPrice struct {
	ID                  dot.ID
	ShipmentPriceListID dot.ID
	ShipmentServiceID   dot.ID
	Name                string
	CustomRegionTypes   []route_type.CustomRegionRouteType
	CustomRegionIDs     []dot.ID
	RegionTypes         []route_type.RegionRouteType
	ProvinceTypes       []route_type.ProvinceRouteType
	UrbanTypes          []route_type.UrbanType
	Details             []*PricingDetail
	AdditionalFees      []*AdditionalFee
	PriorityPoint       int
	CreatedAt           time.Time `sq:"create"`
	UpdatedAt           time.Time `sq:"update"`
	DeletedAt           time.Time
	WLPartnerID         dot.ID
	Status              status3.Status
}

type PricingDetail struct {
	Weight     int                        `json:"weight"`
	Price      int                        `json:"price"`
	Overweight []*PricingDetailOverweight `json:"overweight"`
}

type PricingDetailOverweight struct {
	MinWeight  int `json:"min_weight"`
	MaxWeight  int `json:"max_weight"`
	WeightStep int `json:"weight_step"`
	PriceStep  int `json:"price_step"`
}

type AdditionalFee struct {
	FeeType           shipping_fee_type.ShippingFeeType        `json:"fee_type"`
	CalculationMethod calculation_method.CalculationMethodType `json:"calculation_method"`
	BaseValueType     additional_fee_base_value.BaseValueType  `json:"base_value_type"`
	Rules             []*AdditionalFeeRule                     `json:"rules"`
}

type AdditionalFeeRule struct {
	MinValue          int                                   `json:"min_value"`
	MaxValue          int                                   `json:"max_value"`
	PriceModifierType price_modifier_type.PriceModifierType `json:"price_modifier_type"`
	Amount            float64                               `json:"amount"`
	MinPrice          int                                   `json:"min_price"`
	StartValue        int                                   `json:"start_value"`
}
