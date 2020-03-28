package models

import "etop.vn/capi/dot"

// +sqlgen
type Account struct {
	ID      dot.ID
	OwnerID dot.ID
	Name    string

	NewName string
}
