package etelecom

import (
	"o.o/api/etelecom/call_direction"
	"o.o/api/etelecom/call_state"
	"o.o/api/top/int/shop"
	"o.o/api/top/types/etc/charge_type"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

type SummaryEtelecomRequest struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

func (r *SummaryEtelecomRequest) String() string { return jsonx.MustMarshalToString(r) }

type SummaryEtelecomResponse struct {
	Tables []*shop.SummaryTable `json:"tables"`
}

func (r *SummaryEtelecomResponse) String() string { return jsonx.MustMarshalToString(r) }

type CreateUserAndAssignExtensionRequest struct {
	FullName  string `json:"full_name"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	HotlineID dot.ID `json:"hotline_id"`
}

func (r *CreateUserAndAssignExtensionRequest) String() string { return jsonx.MustMarshalToString(r) }

type CreateCallLogRequest struct {
	ExternalSessionID string                       `json:"external_session_id"`
	Direction         call_direction.CallDirection `json:"direction"`
	Caller            string                       `json:"caller"`
	Callee            string                       `json:"callee"`
	ExtensionID       dot.ID                       `json:"extension_id"`
	ContactID         dot.ID                       `json:"contact_id"`
	CallState         call_state.CallState         `json:"call_state"`
}

func (r *CreateCallLogRequest) String() string {
	return jsonx.MustMarshalToString(r)
}

type UpdateUserSettingRequest struct {
	ExtensionChargeType charge_type.ChargeType `json:"extension_charge_type"`
}

func (r *UpdateUserSettingRequest) String() string {
	return jsonx.MustMarshalToString(r)
}
