package model

import (
	"time"

	"etop.vn/api/top/types/etc/filter_type"
	"etop.vn/api/top/types/etc/location_type"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
)

// +sqlgen
type ShipmentService struct {
	ID                 dot.ID
	ConnectionID       dot.ID
	Name               string
	EdCode             string
	ServiceIDs         []string
	Description        string
	CreatedAt          time.Time `sq:"create"`
	UpdatedAt          time.Time `sq:"update"`
	DeletedAt          time.Time
	WLPartnerID        dot.ID
	ImageURL           string
	Status             status3.Status
	AvailableLocations []*AvailableLocation
	BlacklistLocations []*BlacklistLocation
	OtherCondition     *OtherCondition
}

type AvailableLocation struct {
	FilterType           filter_type.FilterType             `json:"filter_type"`
	ShippingLocationType location_type.ShippingLocationType `json:"shipping_location_type"`
	RegionTypes          []location_type.RegionType         `json:"region_types"`
	CustomRegionIDs      []dot.ID                           `json:"custom_region_ids"`
	ProvinceCodes        []string                           `json:"province_codes"`
}

type BlacklistLocation struct {
	ShippingLocationType location_type.ShippingLocationType `json:"shipping_location_type"`
	ProvinceCodes        []string                           `json:"province_codes"`
	DistrictCodes        []string                           `json:"district_codes"`
	WardCodes            []string                           `json:"ward_codes"`
	Reason               string                             `json:"reason"`
}

type OtherCondition struct {
	MinWeight int `json:"min_weight"`
	MaxWeight int `json:"max_weight"`
}
