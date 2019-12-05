package affiliate

import (
	etop "etop.vn/api/top/int/etop"
	"etop.vn/common/jsonx"
)

type RegisterAffiliateRequest struct {
	// @required
	Name        string            `json:"name"`
	Phone       string            `json:"phone"`
	Email       string            `json:"email"`
	BankAccount *etop.BankAccount `json:"bank_account"`
}

func (m *RegisterAffiliateRequest) Reset()         { *m = RegisterAffiliateRequest{} }
func (m *RegisterAffiliateRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateAffiliateRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

func (m *UpdateAffiliateRequest) Reset()         { *m = UpdateAffiliateRequest{} }
func (m *UpdateAffiliateRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateAffiliateBankAccountRequest struct {
	BankAccount *etop.BankAccount `json:"bank_account"`
}

func (m *UpdateAffiliateBankAccountRequest) Reset()         { *m = UpdateAffiliateBankAccountRequest{} }
func (m *UpdateAffiliateBankAccountRequest) String() string { return jsonx.MustMarshalToString(m) }
