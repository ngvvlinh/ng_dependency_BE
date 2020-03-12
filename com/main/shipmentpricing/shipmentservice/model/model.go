package model

import (
	"time"

	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
)

// +sqlgen
type ShipmentService struct {
	ID           dot.ID
	ConnectionID dot.ID
	Name         string
	EdCode       string
	ServiceIDs   []string
	Description  string
	CreatedAt    time.Time `sq:"create"`
	UpdatedAt    time.Time `sq:"update"`
	DeletedAt    time.Time
	WLPartnerID  dot.ID
	ImageURL     string
	Status       status3.Status
}
