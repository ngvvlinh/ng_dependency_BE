package client

import (
	"o.o/api/etelecom/call_state"
	"o.o/backend/pkg/common/apifw/httpreq"
)

type (
	Bool              = httpreq.Bool
	Float             = httpreq.Float
	Int               = httpreq.Int
	String            = httpreq.String
	Time              = httpreq.Time
	PortsipCallStatus string
	PortsipErrorCode  string
)

const (
	PortsipCallStatusAnswered   PortsipCallStatus = "ANSWERED"
	PortsipCallStatusNoAnswered PortsipCallStatus = "NOANSWER"
	PortsipCallStatusFail       PortsipCallStatus = "FAIL"
	PortsipCallStatusNone       PortsipCallStatus = "NONE"
)

func (s PortsipCallStatus) ToCallState() call_state.CallState {
	switch s {
	case PortsipCallStatusAnswered:
		return call_state.Answered
	case PortsipCallStatusNoAnswered:
		return call_state.NotAnswered
	case PortsipCallStatusNone, PortsipCallStatusFail:
		return call_state.NotAnswered
	default:
		return call_state.Unknown
	}
}

type GetCallLogsRequest struct {
	ScrollID  string `url:"scroll_id,omitempty"`
	StartTime int64  `url:"start_time,omitempty"`
	EndTime   int64  `url:"end_time,omitempty"`
}

type GetCallLogsResponse struct {
	ScrollID String            `json:"scroll_id"`
	Sessions []*SessionCallLog `json:"sessions"`
	Total    Int               `json:"total"`
}

type SessionCallLog struct {
	AppID              Int               `json:"app_id"` // 3000001
	RecordingFileURL   String            `json:"recording_file_url"`
	CallID             String            `json:"call_id"`             // "9a3thmvujh498hkv35m9-gw"
	CallStatus         PortsipCallStatus `json:"call_status"`         // FAIL
	CallTargets        []*CallTarget     `json:"call_targets"`        // người nhận cuộc gọi, sẽ có nhiều CallTarget nếu cuộc gọi chuyển tiếp
	Callee             String            `json:"callee"`              // "0943630091"
	CalleeDomain       String            `json:"callee_domain"`       // "etop-dev.vht.com.vn"
	Caller             String            `json:"caller"`              // "2611"
	CallerDisplayName  String            `json:"caller_display_name"` // "2611"
	CallerDomain       String            `json:"caller_domain"`       // "etop-dev.vht.com.vn"
	Customer           *CustomerSession  `json:"customer"`
	Direction          String            `json:"direction"`    // "ext"
	EndedReason        String            `json:"ended_reason"` // "Unknown"
	EndedTime          Time              `json:"ended_time"`   // "2020-12-01T18:08:28+07:00"
	Order              *OrderSession     `json:"order"`
	RequestDescription String            `json:"request_description"`
	SessionID          String            `json:"session_id"`    // "386111305045643264"
	StartTime          Time              `json:"start_time"`    // "2020-12-01T18:08:28+07:00"
	TalkDuration       Int               `json:"talk_duration"` // 0
	TenantID           String            `json:"tenant_id"`     // "373302079663509504"
	TenantName         String            `json:"tenant_name"`   // "Etop-dev"
	Type               Int               `json:"type"`          // 1
	// DidCid: hotline number when direction = in
	DidCid String `json:"did_cid"`
	// OutboundCallerID: hotline number when direction = out
	// Case direction = ext: can not recognize hotline number
	OutboundCallerID String `json:"outbound_caller_id"`
}

type CallTarget struct {
	AddTime      Time              `json:"add_time"`
	AnsweredTime Time              `json:"answered_time"`
	EndReason    String            `json:"end_reason"`
	EndedTime    Time              `json:"ended_time"`
	FailCode     Int               `json:"fail_code"`
	RingDuration Int               `json:"ring_duration"`
	RingTime     Time              `json:"ring_time"`
	Status       PortsipCallStatus `json:"status"`
	TalkDuration Int               `json:"talk_duration"`
	TargetNumber String            `json:"target_number"` // số callee
	TrunkName    String            `json:"trunk_name"`
}

type CustomerSession struct {
	CustomerMessage  String `json:"customer_message"`
	CustomerResponse Int    `json:"customer_response"`
}

type OrderSession struct {
	LadingCode     String `json:"lading_code"`
	PoCode         String `json:"po_code"`
	PoDistrictCode String `json:"po_district_code"`
	PoProvinceCode String `json:"po_province_code"`
	PostmanCode    String `json:"postman_code"`
	RouteCode      String `json:"route_code"`
}

type ConfigTenantCDRRequest struct {
	// username tenant portsip
	Name string `json:"name"`
	// password tenant portsip
	Password string `json:"password"`
}

type ConfigTenantCDRResponse struct {
	ID String `json:"id"`
}
