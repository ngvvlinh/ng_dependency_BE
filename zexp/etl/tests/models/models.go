package models

import "o.o/capi/dot"

// +sqlgen
type Account struct {
	ID      dot.ID
	OwnerID dot.ID
	Name    string

	NewName string
}
