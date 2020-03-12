package model

import (
	"time"

	"etop.vn/api/top/types/etc/route_type"
	"etop.vn/capi/dot"
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
	PriorityPoint       int
	CreatedAt           time.Time `sq:"create"`
	UpdatedAt           time.Time `sq:"update"`
	DeletedAt           time.Time
	WLPartnerID         dot.ID
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
