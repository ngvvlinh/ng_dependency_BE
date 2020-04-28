package affiliate

import (
	etop "o.o/api/top/int/etop"
	"o.o/common/jsonx"
)

type RegisterAffiliateRequest struct {
	// @required
	Name        string            `json:"name"`
	Phone       string            `json:"phone"`
	Email       string            `json:"email"`
	BankAccount *etop.BankAccount `json:"bank_account"`
}

func (m *RegisterAffiliateRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateAffiliateRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

func (m *UpdateAffiliateRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateAffiliateBankAccountRequest struct {
	BankAccount *etop.BankAccount `json:"bank_account"`
}

func (m *UpdateAffiliateBankAccountRequest) String() string { return jsonx.MustMarshalToString(m) }
