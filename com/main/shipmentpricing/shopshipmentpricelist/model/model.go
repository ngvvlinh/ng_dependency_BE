package model

import (
	"time"

	"o.o/capi/dot"
)

// +sqlgen
type ShopShipmentPriceList struct {
	ShopID              dot.ID
	ShipmentPriceListID dot.ID
	Note                string
	CreatedAt           time.Time `sq:"create"`
	UpdatedAt           time.Time `sq:"update"`
	DeletedAt           time.Time
	UpdatedBy           dot.ID
	WLPartnerID         dot.ID
}
