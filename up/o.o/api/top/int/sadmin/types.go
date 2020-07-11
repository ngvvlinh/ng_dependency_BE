package sadmin

import (
	etop "o.o/api/top/int/etop"
	"o.o/api/top/types/etc/webhook_type"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

type SAdminResetPasswordRequest struct {
	UserId   dot.ID `json:"user_id"`
	Password string `json:"password"`
	Confirm  string `json:"confirm"`
}

func (m *SAdminResetPasswordRequest) String() string { return jsonx.MustMarshalToString(m) }

type SAdminCreateUserRequest struct {
	Info        *etop.CreateUserRequest `json:"info"`
	IsEtopAdmin bool                    `json:"is_etop_admin"`
	Permission  *etop.Permission        `json:"permission"`
}

func (m *SAdminCreateUserRequest) String() string { return jsonx.MustMarshalToString(m) }

type LoginAsAccountRequest struct {
	UserId    dot.ID `json:"user_id"`
	AccountId dot.ID `json:"account_id"`
}

func (m *LoginAsAccountRequest) String() string { return jsonx.MustMarshalToString(m) }

func (m *SAdminCreateUserRequest) Censor() {
	if m.Info != nil && m.Info.Password != "" {
		m.Info.Password = "..."
	}
}

func (m *SAdminResetPasswordRequest) Censor() {
	if m.Password != "" {
		m.Password = "..."
	}
	if m.Confirm != "" {
		m.Confirm = "..."
	}
}

type SAdminRegisterWebhookRequest struct {
	CallbackURL string                   `json:"callback_url"`
	Type        webhook_type.WebhookType `json:"type"`
}

func (m *SAdminRegisterWebhookRequest) String() string { return jsonx.MustMarshalToString(m) }

type SAdminUnregisterWebhookRequest struct {
	CallbackURL string                   `json:"callback_url"`
	RemoveAll   bool                     `json:"remove_all"`
	Type        webhook_type.WebhookType `json:"type"`
}

func (m *SAdminUnregisterWebhookRequest) String() string { return jsonx.MustMarshalToString(m) }
