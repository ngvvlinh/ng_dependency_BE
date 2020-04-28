package model

import (
	"time"

	"o.o/capi/dot"
)

// +sqlgen
type InventoryVariant struct {
	ShopID         dot.ID
	VariantID      dot.ID
	QuantityOnHand int
	QuantityPicked int
	CostPrice      int

	CreatedAt time.Time
	UpdatedAt time.Time

	Rid dot.ID
}
