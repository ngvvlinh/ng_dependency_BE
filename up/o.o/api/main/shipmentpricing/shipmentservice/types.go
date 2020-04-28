package shipmentservice

import (
	"time"

	"o.o/api/top/types/etc/filter_type"
	"o.o/api/top/types/etc/location_type"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

type ShipmentService struct {
	ID                 dot.ID               `json:"id"`
	ConnectionID       dot.ID               `json:"connection_id"`
	Name               string               `json:"name"`
	EdCode             string               `json:"ed_code"`
	ServiceIDs         []string             `json:"service_ids"`
	Description        string               `json:"description"`
	CreatedAt          time.Time            `json:"-"`
	UpdatedAt          time.Time            `json:"-"`
	DeletedAt          time.Time            `json:"-"`
	WLPartnerID        dot.ID               `json:"-"`
	ImageURL           string               `json:"image_url"`
	Status             status3.Status       `json:"status"`
	AvailableLocations []*AvailableLocation `json:"available_locations"`
	BlacklistLocations []*BlacklistLocation `json:"blacklist_locations"`
	OtherCondition     *OtherCondition      `json:"other_condition"`
}

type AvailableLocation struct {
	FilterType           filter_type.FilterType             `json:"filter_type"`
	ShippingLocationType location_type.ShippingLocationType `json:"shipping_location_type"`
	RegionTypes          []location_type.RegionType         `json:"regions"`
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
