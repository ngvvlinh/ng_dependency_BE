package testsrc

import "o.o/capi/dot"

// +sqlgen
type Account struct {
	ID        dot.ID
	FirstName string
	LastName  string

	Rid dot.ID
}
