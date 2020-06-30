package shipmentprice

import (
	"time"

	"o.o/api/top/types/etc/price_modifier_type"
	"o.o/api/top/types/etc/route_type"
	"o.o/api/top/types/etc/shipping_fee_type"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

const MaximumValue = -1

type ShipmentPrice struct {
	ID                  dot.ID                             `json:"id"`
	ShipmentPriceListID dot.ID                             `json:"shipment_price_list_id"`
	ShipmentServiceID   dot.ID                             `json:"shipment_service_id"`
	Name                string                             `json:"name"`
	CustomRegionTypes   []route_type.CustomRegionRouteType `json:"custom_region_types"`
	CustomRegionIDs     []dot.ID                           `json:"custom_region_ids"`
	RegionTypes         []route_type.RegionRouteType       `json:"region_types"`
	ProvinceTypes       []route_type.ProvinceRouteType     `json:"province_types"`
	UrbanTypes          []route_type.UrbanType             `json:"urban_types"`
	Details             []*PricingDetail                   `json:"details"`
	AdditionalFees      []*AdditionalFee                   `json:"additional_fees"`
	PriorityPoint       int                                `json:"priority_point"`
	CreatedAt           time.Time                          `json:"-"`
	UpdatedAt           time.Time                          `json:"-"`
	DeletedAt           time.Time                          `json:"-"`
	WLPartnerID         dot.ID                             `json:"-"`
	Status              status3.Status                     `json:"status"`
}

// PricingDetail
//
// Use for setting main fees
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
	FeeType shipping_fee_type.ShippingFeeType `json:"fee_type"`
	Rules   []*AdditionalFeeRule              `json:"rules"`
}

type AdditionalFeeRule struct {
	MinValue          int                                   `json:"min_value"`
	MaxValue          int                                   `json:"max_value"`
	PriceModifierType price_modifier_type.PriceModifierType `json:"price_modifier_type"`
	Amount            float64                               `json:"amount"`
	MinPrice          int                                   `json:"min_price"`
}
