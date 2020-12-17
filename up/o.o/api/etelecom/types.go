package etelecom

import (
	"time"

	"o.o/api/etelecom/call_log_direction"
	"o.o/api/etelecom/call_state"
	"o.o/api/etelecom/mobile_network"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status5"
	"o.o/capi/dot"
)

// +gen:event:topic=event/etelecom

type Hotline struct {
	ID               dot.ID
	OwnerID          dot.ID
	Name             string
	Hotline          string
	Network          mobile_network.MobileNetwork
	ConnectionID     dot.ID
	ConnectionMethod connection_type.ConnectionMethod
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        time.Time
	Status           status3.Status
	Description      string
	IsFreeCharge     dot.NullBool
}

type Extension struct {
	ID                dot.ID
	UserID            dot.ID
	AccountID         dot.ID
	HotlineID         dot.ID
	ExtensionNumber   string
	ExtensionPassword string
	ExternalData      *ExtensionExternalData
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

type CallLog struct {
	ID                 dot.ID
	ExternalID         string
	AccountID          dot.ID
	StartedAt          time.Time
	EndedAt            time.Time
	Duration           int
	Caller             string
	Callee             string
	AudioURLs          []string
	ExternalDirection  string
	ExternalCallStatus string
	CallState          call_state.CallState
	CallStatus         status5.Status
	Direction          call_log_direction.CallLogDirection
	ExtensionID        dot.ID
	HotlineID          dot.ID
	ContactID          dot.ID
	CreatedAt          time.Time
	UpdatedAt          time.Time
	// DurationForPostage: minute
	DurationPostage int
	Postage         int
}

type CallLogCreatedEvent struct {
	ID         dot.ID
	Direction  call_log_direction.CallLogDirection
	Callee     string
	Duration   int
	CallStatus status5.Status
	HotlineID  dot.ID
}
