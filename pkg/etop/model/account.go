package model

import (
	"etop.vn/api/top/types/etc/account_type"
	"etop.vn/capi/dot"
)

type UpdateAccountURLSlugCommand struct {
	AccountID dot.ID
	URLSlug   string
}

type GetAccountAuthQuery struct {
	AuthKey     string
	AccountType account_type.AccountType
	AccountID   dot.ID

	Result struct {
		AccountAuth *AccountAuth
		Account     AccountInterface
	}
}
