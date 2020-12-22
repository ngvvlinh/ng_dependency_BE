package model

import (
	"time"

	"o.o/api/etelecom/call_log_direction"
	"o.o/api/etelecom/call_state"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status5"
	"o.o/capi/dot"
)

// +sqlgen
type Hotline struct {
	ID dot.ID
	// OwnerID - userID chủ tài khoản
	OwnerID          dot.ID
	Name             string
	Hotline          string
	Network          string
	ConnectionID     dot.ID
	ConnectionMethod connection_type.ConnectionMethod
	CreatedAt        time.Time `sq:"create"`
	UpdatedAt        time.Time `sq:"update"`
	DeletedAt        time.Time
	Status           status3.Status
	Description      string
	IsFreeCharge     dot.NullBool
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
	TenantDomain      string
	ExternalData      *ExtensionExternalData `json:"external_data"`
	CreatedAt         time.Time              `sq:"create"`
	UpdatedAt         time.Time              `sq:"update"`
	DeletedAt         time.Time
}

type ExtensionExternalData struct {
	ID string `json:"id"`
}

// +sqlgen
type CallLog struct {
	ID         dot.ID    `paging:"id"`
	ExternalID string    // sessionID
	AccountID  dot.ID    // shopID from extension
	StartedAt  time.Time `paging:"started_at"`
	EndedAt    time.Time
	Duration   int      // second
	Caller     string   // sdt người gọi
	Callee     string   // sdt người nghe
	AudioURLs  []string // file ghi âm
	// VHT
	// direction: ext/in/out/in2out
	// • ext: means extension to extension // gọi nội bộ: máy nhánh với máy nhánh
	// • in: means call from trunk // Gọi từ sim vào tổng đài
	// • out: means call from extension to trunk // từ tổng đài ra sim khách hàng
	// • in2out: means a call is coming from trunk, but was forwarded to external number
	ExternalDirection  string
	Direction          call_log_direction.CallLogDirection
	ExtensionID        dot.ID
	HotlineID          dot.ID
	ExternalCallStatus string
	ContactID          dot.ID
	CreatedAt          time.Time `sq:"create"`
	UpdatedAt          time.Time `sq:"update"`
	CallState          call_state.CallState
	CallStatus         status5.Status
	DurationPostage    int
	Postage            int

	// depends on connectionID to get extension
	// accountID get from extension (above)
	// get contactID from caller/callee + accountID
}
