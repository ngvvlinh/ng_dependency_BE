package model

import (
	"etop.vn/api/top/types/etc/account_type"
	"etop.vn/capi/dot"
)

// +sqlgen
type Account struct {
	ID       dot.ID
	OwnerID  dot.ID
	Name     string
	Type     account_type.AccountType
	ImageURL string
	URLSlug  string
	Rid      dot.ID
}
