package model

import (
	"time"

	"etop.vn/capi/dot"
)

// +sqlgen
type ShopTrader struct {
	ID        dot.ID
	ShopID    dot.ID
	Type      string
	DeletedAt time.Time

	Rid dot.ID
}
