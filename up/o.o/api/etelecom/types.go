package etelecom

import (
	"time"

	"o.o/api/top/types/etc/connection_type"
	"o.o/capi/dot"
)

// +gen:event:topic=event/etelecom

type Hotline struct {
	ID               dot.ID
	OwnerID          dot.ID
	Hotline          string
	Network          string
	ConnectionID     dot.ID
	ConnectionMethod connection_type.ConnectionMethod
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        time.Time
}

type Extension struct {
	ID                dot.ID
	UserID            dot.ID
	AccountID         dot.ID
	HotlineID         dot.ID
	ExtensionNumber   string
	ExtensionPassword string
	ExternalData      *ExtensionExternalData
	ConnectionID      dot.ID
	ConnectionMethod  connection_type.ConnectionMethod
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         time.Time
}

type ExtensionExternalData struct {
	ID string
}

type Summary struct {
	ID             dot.ID
	ExtensionID    dot.ID
	Date           string
	TotalPhoneCall int
	TotalCallTime  int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
}

type ExtensionCreatingEvent struct {
	// ID nhân viên shop
	UserID dot.ID

	// Shop ID
	AccountID dot.ID
}
