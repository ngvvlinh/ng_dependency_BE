package model

import (
	"time"

	"o.o/capi/dot"
)

// +sqlgen
type CustomRegion struct {
	ID            dot.ID
	Name          string
	Description   string
	ProvinceCodes []string
	CreatedAt     time.Time `sq:"create"`
	UpdatedAt     time.Time `sq:"update"`
	DeletedAt     time.Time
	WLPartnerID   dot.ID
}
