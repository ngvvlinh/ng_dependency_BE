package testdst

import "etop.vn/capi/dot"

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenAccount(&Account{})

type Account struct {
	ID        dot.ID
	FirstName string
	LastName  string
	FullName  string
	Rid       dot.ID
}
