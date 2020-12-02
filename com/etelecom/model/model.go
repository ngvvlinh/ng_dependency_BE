package model

import (
	"time"

	"o.o/api/top/types/etc/connection_type"
	"o.o/capi/dot"
)

// +sqlgen
type Hotline struct {
	ID dot.ID
	// OwnerID - userID chủ tài khoản
	OwnerID          dot.ID
	Hotline          string
	Network          string
	ConnectionID     dot.ID
	ConnectionMethod connection_type.ConnectionMethod
	CreatedAt        time.Time `sq:"create"`
	UpdatedAt        time.Time `sq:"update"`
	DeletedAt        time.Time
}

// +sqlgen
type Extension struct {
	ID dot.ID
	// UserID - userID nhân viên được gán với extension
	UserID            dot.ID
	AccountID         dot.ID
	HotlineID         dot.ID
	ExtensionNumber   string
	ExtensionPassword string
	ExternalData      *ExtensionExternalData `json:"external_data"`
	ConnectionID      dot.ID
	ConnectionMethod  connection_type.ConnectionMethod
	CreatedAt         time.Time `sq:"create"`
	UpdatedAt         time.Time `sq:"update"`
	DeletedAt         time.Time
}

type ExtensionExternalData struct {
	ID string `json:"id"`
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
