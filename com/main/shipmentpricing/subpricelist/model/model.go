package model

import (
	"time"

	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +sqlgen
type ShipmentSubPriceList struct {
	ID           dot.ID
	Name         string
	Description  string
	Status       status3.Status
	ConnectionID dot.ID
	WLPartnerID  dot.ID
	CreatedAt    time.Time `sq:"create"`
	UpdatedAt    time.Time `sq:"update"`
	DeletedAt    time.Time
}
