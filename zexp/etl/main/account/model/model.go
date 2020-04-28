package model

import (
	"o.o/api/top/types/etc/account_type"
	"o.o/capi/dot"
)

// +sqlgen
type Account struct {
	ID       dot.ID
	OwnerID  dot.ID
	Name     string
	Type     account_type.AccountType `sql_type:"enum(account_type)"`
	ImageURL string
	URLSlug  string
	Rid      dot.ID
}
