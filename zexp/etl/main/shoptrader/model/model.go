package model

import (
	"etop.vn/capi/dot"
)

// +sqlgen
type ShopTrader struct {
	ID     dot.ID
	ShopID dot.ID
	Type   string `sql_type:"enum(trader_type)"`

	Rid dot.ID
}
