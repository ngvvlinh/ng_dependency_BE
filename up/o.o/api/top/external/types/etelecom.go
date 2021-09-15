package types

import (
	"time"

	"o.o/api/etelecom/call_direction"
	"o.o/api/etelecom/call_state"
	"o.o/api/top/types/common"
	"o.o/api/top/types/etc/status5"
	"o.o/capi/dot"
	"o.o/common/jsonx"
	"o.o/common/xerrors"
)

type GetExtensionInfoRequest struct {
	ExtensionNumber string `json:"extension_number"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
}

func (m *GetExtensionInfoRequest) String() string { return jsonx.MustMarshalToString(m) }

func (r *GetExtensionInfoRequest) Validate() error {
	if r.ExtensionNumber == "" && r.Email == "" && r.Phone == "" {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing required params")
	}
	return nil
}

type ExtensionInfo struct {
	ExtensionNumber   string    `json:"extension_number"`
	ExtensionPassword string    `json:"extension_password"`
	TenantDomain      string    `json:"tenant_domain"`
	ExpiresAt         time.Time `json:"expires_at"`
}

func (m *ExtensionInfo) String() string { return jsonx.MustMarshalToString(m) }

type ListCallLogsRequest struct {
	Filter *CallLogFilter       `json:"filter"`
	Paging *common.CursorPaging `json:"paging"`
}

func (m *ListCallLogsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CallLogsResponse struct {
	CallLogs []*ShopCallLog         `json:"call_logs"`
	Paging   *common.CursorPageInfo `json:"paging"`
}

func (m *CallLogsResponse) String() string { return jsonx.MustMarshalToString(m) }

type ShopCallLog struct {
	ID                dot.ID                       `json:"id"`
	ExternalSessionID string                       `json:"external_session_id"`
	UserID            dot.ID                       `json:"user_id"`
	StartedAt         time.Time                    `json:"started_at"`
	EndedAt           time.Time                    `json:"ended_at"`
	Duration          int                          `json:"duration"`
	Caller            string                       `json:"caller"`
	Callee            string                       `json:"callee"`
	RecordingURLs     []string                     `json:"recording_urls"`
	Direction         call_direction.CallDirection `json:"direction"`
	ExtensionID       dot.ID                       `json:"extension_id"`
	ContactID         dot.ID                       `json:"contact_id"`
	CreatedAt         time.Time                    `json:"created_at"`
	UpdatedAt         time.Time                    `json:"updated_at"`
	CallState         call_state.CallState         `json:"call_state"`
	CallStatus        status5.Status               `json:"call_status"`
	Note              string                       `json:"note"`
	CallTargets       []*CallTarget                `json:"call_targets"`
}

func (m *ShopCallLog) String() string { return jsonx.MustMarshalToString(m) }

type CallLogFilter struct {
	HotlineIDs   []dot.ID `json:"hotline_ids"`
	ExtensionIDs []dot.ID `json:"extension_ids"`
	UserID       dot.ID   `json:"user_id"`
	// Caller or callee
	CallNumber string `json:"call_number"`
}

type CallTarget struct {
	AddTime      time.Time `json:"add_time"`
	AnsweredTime time.Time `json:"answered_time"`
	EndReason    string    `json:"end_reason"`
	EndedTime    time.Time `json:"ended_time"`
	FailCode     int       `json:"fail_code"`
	RingDuration int       `json:"ring_duration"`
	RingTime     time.Time `json:"ring_time"`
	Status       string    `json:"status"`
	TargetNumber string    `json:"target_number"`
	TrunkName    string    `json:"trunk_name"`
}
