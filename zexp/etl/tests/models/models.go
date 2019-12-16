package models

import "etop.vn/capi/dot"

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenAccount(&Account{})

type Account struct {
	ID      dot.ID
	OwnerID dot.ID
	Name    string

	NewName string
}
