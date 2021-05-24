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

	NameOrDomainIncorrect PortsipErrorCode = "1040001"
	DomainError           PortsipErrorCode = "1050005"
)

type PortsipAccountCfg struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Token        string `json:"token"`
	TenantHost   string `json:"tenant_host"`
	TenantToken  string `json:"tenant_token"`
	TenantDomain string `json:"tenant_domain"`
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
	Profile           *ExtensionProfile            `json:"profile,omitempty"`
	Options           *OptionsCreateExtension      `json:"options,omitempty"`
	ForwardRules      *ForwardRulesCreateExtension `json:"forward_rules,omitempty"`
}

type ExtensionProfile struct {
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
	RingTime     Int               `json:"ring_time"`
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

type CreateTenantRequest struct {
	// required
	Name string `json:"name"`
	// required
	Domain string `json:"domain"`
	// required
	Password   string            `json:"password"`
	Enabled    bool              `json:"enabled"`
	Profile    *TenantProfile    `json:"profile"`
	Capability *TenantCapability `json:"capability"`
	Quota      *TenantQuata      `json:"quota"`
}

type TenantProfile struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	CompanyName string `json:"company_name"`
	// Vietnam
	Region string `json:"region"`
	// Asia/Ho_Chi_Minh
	Timezone string `json:"timezone"`
	// VND
	Currency                      string `json:"currency"`
	EnableExtensionChangePassword bool   `json:"enable_extension_change_password"`
	EnableExtensionVideoRecording bool   `json:"enable_extension_video_recording"`
	EnableExtensionAudioRecording bool   `json:"enable_extension_audio_recording"`
}

type TenantCapability struct {
	MaxExtensions           Int `json:"max_extensions"`            // 100,
	MaxConcurrentCalls      Int `json:"max_concurrent_calls"`      // 100,
	MaxRingGroups           Int `json:"max_ring_groups"`           // 10,
	MaxVirtualReceptionists Int `json:"max_virtual_receptionists"` // 10,
	MaxCallQueues           Int `json:"max_call_queues"`           // 10,
	MaxConferenceRooms      Int `json:"max_conference_rooms"`      // 10
}

type TenantQuata struct {
	MaxRecordingsQuota       Int `json:"max_recordings_quota"`         // 0
	MaxVoicemailQuota        Int `json:"max_voicemail_quota"`          // 0
	MaxCallReportQuota       Int `json:"max_call_report_quota"`        // 0
	AutoCleanRecordingsDays  Int `json:"auto_clean_recordings_days"`   // 30
	AutoCleanVoicemailDays   Int `json:"auto_clean_voicemail_days"`    // 30
	AutoCleanCallReportsDays Int `json:"auto_clean_call_reports_days"` // 30
}

type CreateTenantResponse struct {
	ID String `json:"id"` // portsip tenant ID
}

type CreateOutboundRuleRequest struct {
	Name                string               `json:"name"`                     // "test",
	NumberPrefix        string               `json:"number_prefix,omitempty"`  // "9",
	NumberLength        int                  `json:"number_length,omitempty"`  // 15,
	FromExtension       string               `json:"from_extension,omitempty"` // "101",
	FromExtensionGroups *ExtensionGroup      `json:"from_extension_groups"`
	Routes              []*OutboundRuleRoute `json:"routes"`
}

type ExtensionGroup struct {
	ID           String `json:"id"`
	GroupName    String `json:"group_name"`
	MembersCount String `json:"members_count,omitempty"`
}

type OutboundRuleRoute struct {
	ID          string `json:"id"`
	Provider    string `json:"provider,omitempty"`     // "Sample provider 1",
	StripDigits int    `json:"strip_digits,omitempty"` // 1,
	Prepend     int    `json:"prepend,omitempty"`      // "Sample prepend 1",
	Blocked     bool   `json:"blocked,omitempty"`      // false
}

type CreateOutboundRuleResponse struct {
	ID String `json:"id"` // portsip Outbound Rule ID
}

type ListOutboundRulesResponse struct {
	CommonResponse
	Rules []*OutboundRule `json:"rules"`
}

type OutboundRule struct {
	ID                  String            `json:"id"`
	Name                String            `json:"name"`
	Enable              Bool              `json:"enable"`
	FromExtension       String            `json:"from_extension"`
	FromExtensionGroups []*ExtensionGroup `json:"from_extension_groups"`
}

type UpdateTrunkProviderRequest struct {
	// id trunk provider: aarenet
	ID string `json:"id"`
	// DidPool: sẽ update lại (tức mất data cũ nếu không truyền lên) khi gọi api này
	DidPool []*TrunkProviderDidPool `json:"did_pool"`
}

type TrunkProviderDidPool struct {
	// portsip tenant id
	TenantID string `json:"tenant_id"`
	Name     string `json:"name"`
	// Số hotline
	NumberMask string `json:"number_mask"`
}

type GetTrunkProviderRequest struct {
	ID string `json:"id"`
}

type TrunkProvider struct {
	Name               String                  `json:"name"`
	ID                 String                  `json:"id"`
	Hostname           String                  `json:"hostname"`
	Username           String                  `json:"username"`
	AuthID             String                  `json:"auth_id"`
	Port               Int                     `json:"port"`
	OutboundServer     String                  `json:"outbound_server"`
	Protocol           String                  `json:"protocol"`
	OutboundServerPort Int                     `json:"outbound_server_port"`
	ReregisterInterval Int                     `json:"reregister_interval"`
	Password           String                  `json:"password"`
	AuthMode           String                  `json:"auth_mode"`
	ProviderInLan      Bool                    `json:"provider_in_lan"`
	SingleViaHeader    Bool                    `json:"single_via_header"`
	Ips                []String                `json:"ips"`
	DidPool            []*TrunkProviderDidPool `json:"did_pool"`
}

type CommonResponse struct {
	Pagination Int    `json:"pagination"`
	Pagesize   Int    `json:"pagesize"`
	SortBy     String `json:"sort_by"`
	Count      String `json:"count"`
}

type GetExtensionGroupsResponse struct {
	CommonResponse
	Groups []*ExtensionGroup `json:"groups"`
}

type CommonListRequest struct {
	Pagination int `json:"pagination"`
	Pagesize   int `json:"pagesize,omitempty"`
}
