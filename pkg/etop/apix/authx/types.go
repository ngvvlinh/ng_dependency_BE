package authx

import (
	"o.o/api/top/types/etc/account_type"
	"o.o/capi/dot"
)

type AuthxResponse struct {
	User    *AuthxUser    `json:"user"`
	Account *AuthxAccount `json:"account"`
}

type AuthxUser struct {
	ID dot.ID `json:"id"`
}

type AuthxAccount struct {
	ID      dot.ID                   `json:"id"`
	Name    string                   `json:"name"`
	OwnerID dot.ID                   `json:"owner_id"`
	Type    account_type.AccountType `json:"type"`
}
