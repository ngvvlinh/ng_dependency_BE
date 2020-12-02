package client

import "o.o/backend/pkg/common/apifw/httpreq"

type (
	Bool   = httpreq.Bool
	Float  = httpreq.Float
	Int    = httpreq.Int
	String = httpreq.String
	Time   = httpreq.Time
)

type VHTAccountCfg struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Token       string `json:"token"`
	TenantHost  string `json:"tenant_host"`
	TenantToken string `json:"tenant_token"`
}

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID          String `json:"id"`
	Expires     Int    `json:"expires"`
	ApiVersion  String `json:"api_version"`
	Role        String `json:"role"`
	AccessToken String `json:"access_token"`
	Name        String `json:"name"`
}

type CreateExtensionsRequest struct {
	ExtensionNumber   string                       `json:"extension_number"`
	Password          string                       `json:"password"`
	WebAccessPassword string                       `json:"web_access_password"`
	Profile           *ProfileCreateExtension      `json:"profile,omitempty"`
	Options           *OptionsCreateExtension      `json:"options,omitempty"`
	ForwardRules      *ForwardRulesCreateExtension `json:"forward_rules,omitempty"`
}

type ProfileCreateExtension struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Email       string `json:"email,omitempty"`
	MobilePhone string `json:"mobile_phone,omitempty"`
	WorkPhone   string `json:"work_phone,omitempty"`
	HomePhone   string `json:"home_phone,omitempty"`
	Twitter     string `json:"twitter,omitempty"`
	Facebook    string `json:"facebook,omitempty"`
	Linkedin    string `json:"linkedin,omitempty"`
	Instagram   string `json:"instagram,omitempty"`
	Description string `json:"description,omitempty"`
}

type OptionsCreateExtension struct {
	EnableAudioRecordCalls bool   `json:"enable_audio_record_calls"`
	EnableVideoRecordCalls bool   `json:"enable_video_record_calls"`
	EnableExtension        bool   `json:"enable_extension"`
	OutboundCallerID       string `json:"outbound_caller_id"`
}

type VoiceMailCreateExtension struct {
	EnableVoicemail    bool   `json:"enable_voicemail"`
	PromptLanguage     string `json:"prompt_language"`
	EnableVmPinAuth    bool   `json:"enable_vm_pin_auth"`
	VoicemailPin       string `json:"voicemail_pin"`
	EnablePlayCallerID bool   `json:"enable_play_caller_id"`
	MsgReadOutDatetime string `json:"msg_read_out_datetime"`
}

type AvailableForwardRules struct {
	NoAnswerTimeval     int    `json:"no_answer_timeval"`
	NoAnswerAction      string `json:"no_answer_action"`
	NoAnswerActionValue string `json:"no_answer_action_value"`
	BusyAction          string `json:"busy_action"`
	BusyActionValue     string `json:"busy_action_value"`
}

type OfflineForwardRules struct {
	OfficeHoursAction             string `json:"office_hours_action"`
	OfficeHoursActionValue        string `json:"office_hours_action_value"`
	OutsideOfficeHoursAction      string `json:"outside_office_hours_action"`
	OutsideOfficeHoursActionValue string `json:"outside_office_hours_action_value"`
}

type DndForwardRules struct {
	OfficeHoursAction             string `json:"office_hours_action"`
	OfficeHoursActionValue        string `json:"office_hours_action_value"`
	OutsideOfficeHoursAction      string `json:"outside_office_hours_action"`
	OutsideOfficeHoursActionValue string `json:"outside_office_hours_action_value"`
}

type AwayForwardRules struct {
	OfficeHoursAction             string `json:"office_hours_action"`
	OfficeHoursActionValue        string `json:"office_hours_action_value"`
	OutsideOfficeHoursAction      string `json:"outside_office_hours_action"`
	OutsideOfficeHoursActionValue string `json:"outside_office_hours_action_value"`
}

type ForwardRulesCreateExtension struct {
	Available *AvailableForwardRules `json:"available,omitempty"`
	Offline   *OfflineForwardRules   `json:"offline,omitempty"`
	Dnd       *DndForwardRules       `json:"dnd,omitempty"`
	Away      *AwayForwardRules      `json:"away,omitempty"`
}

type CreateExtensionResponse struct {
	ID String `json:"id"`
}

type GetCallLogsResponse struct {
	ScrollID String            `json:"scroll_id"`
	Sessions []*SessionCallLog `json:"sessions"`
	Total    Int               `json:"total"`
}

type SessionCallLog struct {
	AppID              Int              `json:"app_id"` // 3000001
	AudioURLs          []String         `json:"audio_urls"`
	CallID             String           `json:"call_id"`             // "9a3thmvujh498hkv35m9-gw"
	CallStatus         String           `json:"call_status"`         // FAIL
	Callee             String           `json:"callee"`              // "0943630091"
	CalleeDomain       String           `json:"callee_domain"`       // "etop-dev.vht.com.vn"
	Caller             String           `json:"caller"`              // "2611"
	CallerDisplayName  String           `json:"caller_display_name"` // "2611"
	CallerDomain       String           `json:"caller_domain"`       // "etop-dev.vht.com.vn"
	Customer           *CustomerSession `json:"customer"`
	Direction          String           `json:"direction"`  // "ext"
	EndReason          String           `json:"end_reason"` // "Unknown"
	EndTime            String           `json:"end_time"`   // "2020-12-01T18:08:28+07:00"
	Order              *OrderSession    `json:"order"`
	RequestDescription String           `json:"request_description"`
	SessionID          String           `json:"session_id"`    // "386111305045643264"
	StartTime          String           `json:"start_time"`    // "2020-12-01T18:08:28+07:00"
	TaskDuration       Int              `json:"task_duration"` // 0
	TenantID           String           `json:"tenant_id"`     // "373302079663509504"
	TenantName         String           `json:"tenant_name"`   // "Etop-dev"
	Type               Int              `json:"type"`          // 1
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

type ErrorResponse struct {
	ErrCode String `json:"err_code"`
	Msg     String `json:"msg"`
	Message String `json:"message"`
}

func (e *ErrorResponse) Error() string {
	if e.Msg.String() != "" {
		return e.Msg.String()
	}
	return e.Message.String()
}

func URL(baseUrl, path string) string {
	return baseUrl + path
}
