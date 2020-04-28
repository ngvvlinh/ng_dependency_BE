package stocktaking

import (
	"time"

	catalogtype "o.o/api/main/catalog/types"
	"o.o/api/top/types/etc/inventory_auto"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/stocktake_type"
	"o.o/capi/dot"
)

// +gen:event:topic=event/stocktake

type StocktakeConfirmedEvent struct {
	StocktakeID          dot.ID
	ShopID               dot.ID
	Overstock            bool
	ConfirmedBy          dot.ID
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher
}

type ShopStocktake struct {
	ID            dot.ID
	ShopID        dot.ID
	TotalQuantity int

	CreatedBy    dot.ID
	UpdatedBy    dot.ID
	CancelReason string

	Note     string
	Code     string
	CodeNorm int

	CreatedAt   time.Time
	UpdatedAt   time.Time
	ConfirmedAt time.Time
	CancelledAt time.Time

	Lines  []*StocktakeLine
	Status status3.Status
	Type   stocktake_type.StocktakeType
}

type StocktakeLine struct {
	ProductID   dot.ID
	ProductName string

	VariantID   dot.ID
	OldQuantity int
	NewQuantity int
	VariantName string
	Code        string
	ImageURL    string

	CostPrice  int
	Attributes []*catalogtype.Attribute
}
