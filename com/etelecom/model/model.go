package model

import (
	"time"

	"o.o/api/top/types/etc/connection_type"
	"o.o/capi/dot"
)

// +sqlgen
type Hotline struct {
	ID               dot.ID
	UserID           dot.ID
	Hotline          string
	Network          string
	Provider         string
	ConnectionID     dot.ID
	ConnectionMethod connection_type.ConnectionMethod
	CreatedAt        time.Time `sq:"create"`
	UpdatedAt        time.Time `sq:"update"`
	DeletedAt        time.Time
}

// +sqlgen
type Extension struct {
	ID                dot.ID
	UserID            dot.ID
	AccountID         dot.ID
	ExtensionNumber   string
	ExtensionPassword string
	ExternalData      *ExtensionExternalData
	ConnectionID      dot.ID
	ConnectionMethod  connection_type.ConnectionMethod
	CreatedAt         time.Time `sq:"create"`
	UpdatedAt         time.Time `sq:"update"`
	DeletedAt         time.Time
}

type ExtensionExternalData struct {
	ID string
}

// +sqlgen
type Summary struct {
	ID             dot.ID
	ExtensionID    dot.ID
	Date           string
	TotalPhoneCall int
	TotalCallTime  int
	CreatedAt      time.Time `sq:"create"`
	UpdatedAt      time.Time `sq:"update"`
	DeletedAt      time.Time
}
