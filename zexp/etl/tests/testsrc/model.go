package testsrc

import "etop.vn/capi/dot"

// +sqlgen
type Account struct {
	ID        dot.ID
	FirstName string
	LastName  string

	Rid dot.ID
}
