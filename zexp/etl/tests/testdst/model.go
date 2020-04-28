package testdst

import "o.o/capi/dot"

// +sqlgen
type Account struct {
	ID        dot.ID
	FirstName string
	LastName  string
	FullName  string
	Rid       dot.ID
}
