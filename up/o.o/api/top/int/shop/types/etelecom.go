package types

import (
	"time"

	"o.o/api/etelecom/call_log_direction"
	"o.o/api/etelecom/call_state"
	"o.o/api/etelecom/mobile_network"
	"o.o/api/top/types/common"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status5"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

type Hotline struct {
	ID               dot.ID                           `json:"id"`
	OwnerID          dot.ID                           `json:"owner_id"`
	Name             string                           `json:"name"`
	Hotline          string                           `json:"hotline"`
	Network          mobile_network.MobileNetwork     `json:"network"`
	ConnectionID     dot.ID                           `json:"connection_id"`
	ConnectionMethod connection_type.ConnectionMethod `json:"connection_method"`
	CreatedAt        time.Time                        `json:"created_at"`
	UpdatedAt        time.Time                        `json:"updated_at"`
	Status           status3.Status                   `json:"status"`
	Description      string                           `json:"description"`
	IsFreeCharge     dot.NullBool                     `json:"is_free_charge"`
}

func (m *Hotline) String() string { return jsonx.MustMarshalToString(m) }

type GetHotLinesResponse struct {
	Hotlines []*Hotline `json:"hotlines"`
}

func (m *GetHotLinesResponse) String() string { return jsonx.MustMarshalToString(m) }

type Extension struct {
	ID                dot.ID    `json:"id"`
	UserID            dot.ID    `json:"user_id"`
	AccountID         dot.ID    `json:"account_id"`
	ExtensionNumber   string    `json:"extension_number"`
	ExtensionPassword string    `json:"extension_password"`
	HotlineID         dot.ID    `json:"hotline_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (m *Extension) String() string { return jsonx.MustMarshalToString(m) }

type ExtensionExternalData struct {
	ID dot.ID `json:"id"`
}

type GetExtensionsRequest struct {
	HotlineID dot.ID `json:"hotline_id"`
}

func (m *GetExtensionsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetExtensionsResponse struct {
	Extensions []*Extension `json:"extensions"`
}

func (m *GetExtensionsResponse) String() string { return jsonx.MustMarshalToString(m) }

type CreateExtensionRequest struct {
	// user_id: nhân viên của shop, người được gán vào extension
	UserID    dot.ID `json:"user_id"`
	HotlineID dot.ID `json:"hotline_id"`
}

func (m *CreateExtensionRequest) String() string { return jsonx.MustMarshalToString(m) }

type CallLog struct {
	ID                 dot.ID                              `json:"id"`
	ExternalID         string                              `json:"external_id"`
	AccountID          dot.ID                              `json:"account_id"`
	HotlineID          dot.ID                              `json:"hotline_id"`
	StartedAt          time.Time                           `json:"started_at"`
	EndedAt            time.Time                           `json:"ended_at"`
	Duration           int                                 `json:"duration"`
	Caller             string                              `json:"caller"`
	Callee             string                              `json:"callee"`
	AudioURLs          []string                            `json:"audio_urls"`
	ExternalDirection  string                              `json:"external_direction"`
	Direction          call_log_direction.CallLogDirection `json:"direction"`
	ExtensionID        dot.ID                              `json:"extension_id"`
	ExternalCallStatus string                              `json:"external_call_status"`
	ContactID          dot.ID                              `json:"contact_id"`
	CreatedAt          time.Time                           `json:"created_at"`
	UpdatedAt          time.Time                           `json:"updated_at"`
	CallState          call_state.CallState                `json:"call_state"`
	CallStatus         status5.Status                      `json:"call_status"`
	// Đơn vị: phút
	DurationPostage int `json:"duration_postage"`
	Postage         int `json:"postage"`
}

func (m *CallLog) String() string { return jsonx.MustMarshalToString(m) }

type GetCallLogsRequest struct {
	Paging *common.CursorPaging `json:"paging"`
	Filter *CallLogsFilter      `json:"filter"`
}

func (m *GetCallLogsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CallLogsFilter struct {
	HotlineIDs   []dot.ID `json:"hotline_ids"`
	ExtensionIDs []dot.ID `json:"extension_ids"`
}

type GetCallLogsResponse struct {
	CallLogs []*CallLog             `json:"call_logs"`
	Paging   *common.CursorPageInfo `json:"paging"`
}

func (m *GetCallLogsResponse) String() string { return jsonx.MustMarshalToString(m) }
