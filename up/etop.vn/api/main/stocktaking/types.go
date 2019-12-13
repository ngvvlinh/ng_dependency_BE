package stocktaking

import (
	"time"

	"etop.vn/api/main/inventory"
	"etop.vn/api/top/types/etc/status3"
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
	Attributes []*Attribute
}

type Attribute struct {
	Name  string
	Value string
}
