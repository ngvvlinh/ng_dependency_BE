package testdst

import "etop.vn/capi/dot"

// +sqlgen
type Account struct {
	ID        dot.ID
	FirstName string
	LastName  string
	FullName  string
	Rid       dot.ID
}
