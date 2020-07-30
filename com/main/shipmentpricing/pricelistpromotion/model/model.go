package model

import (
	"time"

	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
	"o.o/capi/filter"
)

// +sqlgen
type ShipmentPriceListPromotion struct {
	ID            dot.ID
	PriceListID   dot.ID
	Name          string
	Description   string
	Status        status3.Status
	DateFrom      time.Time
	DateTo        time.Time
	AppliedRules  *AppliedRules
	CreatedAt     time.Time `sq:"create"`
	UpdatedAt     time.Time `sq:"update"`
	DeletedAt     time.Time
	WLPartnerID   dot.ID
	ConnectionID  dot.ID
	PriorityPoint int
}

type AppliedRules struct {
	FromCustomRegionIDs []dot.ID    `json:"custom_region_ids"`
	ShopCreatedDate     filter.Date `json:"shop_created_date"`
	UserCreatedDate     filter.Date `json:"user_created_date"`
	UsingPriceListIDs   []dot.ID    `json:"using_price_list_ids"`
}
