package model

import (
	"time"

	"etop.vn/capi/dot"
)

// +sqlgen
type ShopCustomerGroup struct {
	ID        dot.ID
	PartnerID dot.ID
	Name      string
	ShopID    dot.ID

	CreatedAt time.Time
	UpdatedAt time.Time

	Rid dot.ID
}
