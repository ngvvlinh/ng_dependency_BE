package model

import (
	"time"

	"etop.vn/capi/dot"
)

// +sqlgen
type ShopVariantSupplier struct {
	ShopID     dot.ID
	SupplierID dot.ID
	VariantID  dot.ID
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Rid dot.ID
}
