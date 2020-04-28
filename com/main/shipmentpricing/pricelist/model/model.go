package model

import (
	"time"

	"o.o/capi/dot"
)

// +sqlgen
type ShipmentPriceList struct {
	ID          dot.ID
	Name        string
	Description string
	IsActive    bool
	CreatedAt   time.Time `sq:"create"`
	UpdatedAt   time.Time `sq:"update"`
	DeletedAt   time.Time
	WLPartnerID dot.ID
}
