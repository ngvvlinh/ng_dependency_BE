package stocktaking

import (
	"time"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/inventory"
	"etop.vn/capi/dot"
)

// +gen:event:topic=event/stocktake

type StocktakeConfirmedEvent struct {
	StocktakeID          dot.ID
	ShopID               dot.ID
	Overstock            bool
	ConfirmedBy          dot.ID
	AutoInventoryVoucher inventory.AutoInventoryVoucher
}

type ShopStocktake struct {
	ID            dot.ID
	ShopID        dot.ID
	TotalQuantity int32

	CreatedBy    dot.ID
	UpdatedBy    dot.ID
	CancelReason string

	Note     string
	Code     string
	CodeNorm int32

	CreatedAt   time.Time
	UpdatedAt   time.Time
	ConfirmedAt time.Time
	CancelledAt time.Time

	Lines  []*StocktakeLine
	Status etop.Status3
}

type StocktakeLine struct {
	ProductID   dot.ID
	ProductName string

	VariantID   dot.ID
	OldQuantity int32
	NewQuantity int32
	VariantName string
	Code        string
	ImageURL    string

	CostPrice  int32
	Attributes []*Attribute
}

type Attribute struct {
	Name  string
	Value string
}
