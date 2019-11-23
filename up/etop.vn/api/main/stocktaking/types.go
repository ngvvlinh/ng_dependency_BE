package stocktaking

import (
	"time"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/inventory"
)

// +gen:event:topic=event/stocktake

type StocktakeConfirmedEvent struct {
	StocktakeID          int64
	ShopID               int64
	Overstock            bool
	ConfirmedBy          int64
	AutoInventoryVoucher inventory.AutoInventoryVoucher
}

type ShopStocktake struct {
	ID            int64
	ShopID        int64
	TotalQuantity int32

	CreatedBy    int64
	UpdatedBy    int64
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
	ProductID   int64
	ProductName string

	VariantID   int64
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
