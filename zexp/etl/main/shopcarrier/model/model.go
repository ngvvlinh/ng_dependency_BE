package model

import (
	"time"

	"etop.vn/capi/dot"
)

// +sqlgen
type ShopCarrier struct {
	ID        dot.ID
	ShopID    dot.ID
	FullName  string
	Note      string
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time

	Rid dot.ID
}
