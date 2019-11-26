package sadmin

import (
	_ "etop.vn/api/pb/common"
	etop "etop.vn/api/pb/etop"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
)

type SAdminResetPasswordRequest struct {
	UserId   dot.ID `json:"user_id"`
	Password string `json:"password"`
	Confirm  string `json:"confirm"`
}

func (m *SAdminResetPasswordRequest) Reset()         { *m = SAdminResetPasswordRequest{} }
func (m *SAdminResetPasswordRequest) String() string { return jsonx.MustMarshalToString(m) }

type SAdminCreateUserRequest struct {
	Info        *etop.CreateUserRequest `json:"info"`
	IsEtopAdmin bool                    `json:"is_etop_admin"`
	Permission  *etop.Permission        `json:"permission"`
}

func (m *SAdminCreateUserRequest) Reset()         { *m = SAdminCreateUserRequest{} }
func (m *SAdminCreateUserRequest) String() string { return jsonx.MustMarshalToString(m) }

type LoginAsAccountRequest struct {
	UserId    dot.ID `json:"user_id"`
	AccountId dot.ID `json:"account_id"`
}

func (m *LoginAsAccountRequest) Reset()         { *m = LoginAsAccountRequest{} }
func (m *LoginAsAccountRequest) String() string { return jsonx.MustMarshalToString(m) }
