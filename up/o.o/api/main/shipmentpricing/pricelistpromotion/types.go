package pricelistpromotion

import (
	"time"

	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
	"o.o/capi/filter"
)

type ShipmentPriceListPromotion struct {
	ID            dot.ID
	PriceListID   dot.ID
	Name          string
	Description   string
	Status        status3.Status
	DateFrom      time.Time
	DateTo        time.Time
	AppliedRules  *AppliedRules
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
	WLPartnerID   dot.ID
	ConnectionID  dot.ID
	PriorityPoint int
}

type AppliedRules struct {
	FromCustomRegionIDs []dot.ID
	ShopCreatedDate     filter.Date
	UserCreatedDate     filter.Date
	UsingPriceListIDs   []dot.ID
}
