package types

import (
	"time"

	"o.o/api/etelecom/call_direction"
	"o.o/api/etelecom/call_state"
	"o.o/api/etelecom/mobile_network"
	"o.o/api/top/types/common"
	"o.o/api/top/types/etc/charge_type"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/payment_method"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status5"
	"o.o/capi/dot"
	"o.o/common/jsonx"
	"o.o/common/xerrors"
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
	Hotlines []*Hotline             `json:"hotlines"`
	Paging   *common.CursorPageInfo `json:"paging"`
}

func (m *GetHotLinesResponse) String() string { return jsonx.MustMarshalToString(m) }

type Extension struct {
	ID                dot.ID    `json:"id"`
	UserID            dot.ID    `json:"user_id"`
	AccountID         dot.ID    `json:"account_id"`
	ExtensionNumber   string    `json:"extension_number"`
	ExtensionPassword string    `json:"extension_password"`
	TenantID          dot.ID    `json:"tenant_id"`
	TenantDomain      string    `json:"tenant_domain"`
	HotlineID         dot.ID    `json:"hotline_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	ExpiresAt         time.Time `json:"expires_at"`
	SubscriptionID    dot.ID    `json:"subscription_id"`
}

func (m *Extension) String() string { return jsonx.MustMarshalToString(m) }

type ExtensionExternalData struct {
	ID dot.ID `json:"id"`
}

type GetExtensionsRequest struct {
	Filter *GetExtensionsFilter `json:"filter"`
}

func (m *GetExtensionsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetExtensionsFilter struct {
	HotlineID        dot.ID   `json:"hotline_id"`
	ExtensionNumbers []string `json:"extension_numbers"`
}

type GetExtensionsResponse struct {
	Extensions []*Extension `json:"extensions"`
}

func (m *GetExtensionsResponse) String() string { return jsonx.MustMarshalToString(m) }

type CreateExtensionRequest struct {
	// user_id: nh??n vi??n c???a shop, ng?????i ???????c g??n v??o extension
	UserID          dot.ID `json:"user_id"`
	HotlineID       dot.ID `json:"hotline_id"`
	ExtensionNumber int    `json:"extension_number"`
}

func (m *CreateExtensionRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateExtensionBySubscriptionRequest struct {
	// Nh??n vi??n c???a shop, ng?????i ???????c g??n v??o extension
	UserID    dot.ID `json:"user_id"`
	HotlineID dot.ID `json:"hotline_id"`

	SubscriptionID     dot.ID                       `json:"subscription_id"`
	SubscriptionPlanID dot.ID                       `json:"subscription_plan_id"`
	PaymentMethod      payment_method.PaymentMethod `json:"payment_method"`
	ExtensionNumber    int                          `json:"extension_number"`
}

func (m *CreateExtensionBySubscriptionRequest) String() string { return jsonx.MustMarshalToString(m) }

type ExtendExtensionRequest struct {
	ExtensionID dot.ID `json:"extension_id"`
	// Nh??n vi??n c???a shop, ng?????i ???????c g??n v??o extension
	UserID dot.ID `json:"user_id"`

	// B??? tr???ng n???u mu???n gia h???n g??i c??
	SubscriptionID dot.ID `json:"subscription_id"`

	// B??? tr???ng n???u mu???n gia h???n g??i c??
	SubscriptionPlanID dot.ID                       `json:"subscription_plan_id"`
	PaymentMethod      payment_method.PaymentMethod `json:"payment_method"`
}

func (m *ExtendExtensionRequest) String() string { return jsonx.MustMarshalToString(m) }

type CallLog struct {
	ID                 dot.ID                       `json:"id"`
	ExternalID         string                       `json:"external_id"`
	AccountID          dot.ID                       `json:"account_id"`
	HotlineID          dot.ID                       `json:"hotline_id"`
	OwnerID            dot.ID                       `json:"owner_id"`
	UserID             dot.ID                       `json:"user_id"`
	StartedAt          time.Time                    `json:"started_at"`
	EndedAt            time.Time                    `json:"ended_at"`
	Duration           int                          `json:"duration"`
	Caller             string                       `json:"caller"`
	Callee             string                       `json:"callee"`
	AudioURLs          []string                     `json:"audio_urls"`
	ExternalDirection  string                       `json:"external_direction"`
	Direction          call_direction.CallDirection `json:"direction"`
	ExtensionID        dot.ID                       `json:"extension_id"`
	ExternalCallStatus string                       `json:"external_call_status"`
	ContactID          dot.ID                       `json:"contact_id"`
	CreatedAt          time.Time                    `json:"created_at"`
	UpdatedAt          time.Time                    `json:"updated_at"`
	CallState          call_state.CallState         `json:"call_state"`
	CallStatus         status5.Status               `json:"call_status"`
	// ????n v???: ph??t
	DurationPostage int           `json:"duration_postage"`
	Postage         int           `json:"postage"`
	Note            string        `json:"note"`
	CallTargets     []*CallTarget `json:"call_targets"`
}

func (m *CallLog) String() string { return jsonx.MustMarshalToString(m) }

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

func (m *CallTarget) String() string { return jsonx.MustMarshalToString(m) }

type GetCallLogsRequest struct {
	Paging *common.CursorPaging `json:"paging"`
	Filter *CallLogsFilter      `json:"filter"`
}

func (m *GetCallLogsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CallLogsFilter struct {
	HotlineIDs   []dot.ID `json:"hotline_ids"`
	ExtensionIDs []dot.ID `json:"extension_ids"`
	UserID       dot.ID   `json:"user_id"`
	// Caller or callee
	CallNumber string                       `json:"call_number"`
	CallState  call_state.CallState         `json:"call_state"`
	DateFrom   time.Time                    `json:"date_from"`
	DateTo     time.Time                    `json:"date_to"`
	Direction  call_direction.CallDirection `json:"direction"`
}

type GetCallLogsResponse struct {
	CallLogs []*CallLog             `json:"call_logs"`
	Paging   *common.CursorPageInfo `json:"paging"`
}

func (m *GetCallLogsResponse) String() string { return jsonx.MustMarshalToString(m) }

type CreateHotlineRequest struct {
	OwnerID      dot.ID                       `json:"owner_id"`
	Name         string                       `json:"name"`
	Hotline      string                       `json:"hotline"`
	Network      mobile_network.MobileNetwork `json:"network"`
	Description  string                       `json:"description"`
	IsFreeCharge dot.NullBool                 `json:"is_free_charge"`
}

func (m *CreateHotlineRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateHotlineRequest struct {
	ID           dot.ID                       `json:"id"`
	Name         string                       `json:"name"`
	Description  string                       `json:"description"`
	Status       status3.NullStatus           `json:"status"`
	IsFreeCharge dot.NullBool                 `json:"is_free_charge"`
	Network      mobile_network.MobileNetwork `json:"network"`
}

func (m *UpdateHotlineRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetHotLinesRequest struct {
	Paging *common.CursorPaging `json:"paging"`
	Filter *HotlinesFilter      `json:"filter"`
}

func (m *GetHotLinesRequest) String() string { return jsonx.MustMarshalToString(m) }

type HotlinesFilter struct {
	OwnerID  dot.ID `json:"owner_id"`
	TenantID dot.ID `json:"tenant_id"`
}

type EtelecomUserSetting struct {
	// User ID
	ID                  dot.ID                 `json:"id"`
	ExtensionChargeType charge_type.ChargeType `json:"extension_charge_type"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
}

func (r *EtelecomUserSetting) String() string {
	return jsonx.MustMarshalToString(r)
}

type UserSettingsResponse struct {
	UserSettings []*EtelecomUserSetting `json:"user_settings"`
	Paging       *common.CursorPageInfo `json:"paging"`
}

func (r *UserSettingsResponse) String() string {
	return jsonx.MustMarshalToString(r)
}

type GetUserSettingsRequest struct {
	UserIDs []dot.ID             `json:"user_ids"`
	Paging  *common.CursorPaging `json:"paging"`
}

func (m *GetUserSettingsRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateUserSettingRequest struct {
	UserID              dot.ID                 `json:"user_id"`
	ExtensionChargeType charge_type.ChargeType `json:"extension_charge_type"`
}

func (m *UpdateUserSettingRequest) String() string { return jsonx.MustMarshalToString(m) }

type Tenant struct {
	ID               dot.ID                           `json:"id"`
	OwnerID          dot.ID                           `json:"owner_id"`
	Name             string                           `json:"name"`
	Domain           string                           `json:"domain"`
	CreatedAt        time.Time                        `json:"created_at"`
	UpdatedAt        time.Time                        `json:"updated_at"`
	Status           status3.NullStatus               `json:"status"`
	ConnectionID     dot.ID                           `json:"connection_id"`
	ConnectionMethod connection_type.ConnectionMethod `json:"connection_method"`
}

func (m *Tenant) String() string {
	return jsonx.MustMarshalToString(m)
}

type ActivateTenantRequest struct {
	OwnerID  dot.ID `json:"owner_id"`
	TenantID dot.ID `json:"tenant_id"`
	// @required
	HotlineID dot.ID `json:"hotline_id"`
	// Portsip direct connection (defautl value)
	ConnectionID dot.ID `json:"connection_id"`
}

func (m *ActivateTenantRequest) String() string { return jsonx.MustMarshalToString(m) }

type RemoveHotlineOutOfTenantRequest struct {
	OwnerID   dot.ID `json:"owner_id"`
	HotlineID dot.ID `json:"hotline_id"`
}

func (m *RemoveHotlineOutOfTenantRequest) String() string { return jsonx.MustMarshalToString(m) }

type AddHotlineToTenantRequest struct {
	HotlineID dot.ID `json:"hotline_id"`
	TenantID  dot.ID `json:"tenant_id"`
}

func (m *AddHotlineToTenantRequest) String() string { return jsonx.MustMarshalToString(m) }

func (r *AddHotlineToTenantRequest) Validate() error {
	if r.HotlineID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing hotline_id")
	}
	if r.TenantID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing tenant_id")
	}
	return nil
}

type AdminCreateTenantRequest struct {
	OwnerID   dot.ID `json:"owner_id"`
	AccountID dot.ID `json:"account_id"`
	// Support connection portsip direct
	ConnectionID dot.ID `json:"connection_id"`
}

func (m *AdminCreateTenantRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetTenantsRequest struct {
	Paging *common.CursorPaging `json:"paging"`
	Filter *TenantsFilter       `json:"filter"`
}

func (m *GetTenantsRequest) String() string { return jsonx.MustMarshalToString(m) }

type TenantsFilter struct {
	OwnerID dot.ID `json:"owner_id"`
}

type GetTenantsResponse struct {
	Tenants []*Tenant              `json:"tenants"`
	Paging  *common.CursorPageInfo `json:"paging"`
}

func (m *GetTenantsResponse) String() string { return jsonx.MustMarshalToString(m) }

type ActionCallRequest struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Input       string `json:"input"`
	FromName    string `json:"from_name"`
	AccountName string `json:"account_name"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

func (m *ActionCallRequest) String() string { return jsonx.MustMarshalToString(m) }

type ActionCallResponse struct {
	StatusCode  int    `json:"status_code"`
	Action      string `json:"action"`
	Destination string `json:"destination"`
}

func (m *ActionCallResponse) String() string { return jsonx.MustMarshalToString(m) }
