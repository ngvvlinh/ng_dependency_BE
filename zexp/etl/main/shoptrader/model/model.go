package model

import (
	"etop.vn/capi/dot"
)

// +sqlgen
type ShopTrader struct {
	ID     dot.ID
	ShopID dot.ID
	Type   string

	Rid dot.ID
}
