package model

import "etop.vn/capi/dot"

type UpdateAccountURLSlugCommand struct {
	AccountID dot.ID
	URLSlug   string
}

type GetAccountAuthQuery struct {
	AuthKey     string
	AccountType AccountType
	AccountID   dot.ID

	Result struct {
		AccountAuth *AccountAuth
		Account     AccountInterface
	}
}
