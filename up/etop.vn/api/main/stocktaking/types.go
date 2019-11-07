package stocktaking

import (
	"time"

	"etop.vn/api/main/inventory"

	"etop.vn/api/main/etop"
)

// +gen:event:topic=event/stocktake

type StocktakeConfirmedEvent struct {
	Stocktake            *ShopStocktake
	Overstock            bool
	Note                 string
	ConfirmedBy          int64
	AutoInventoryVoucher inventory.AutoInventoryVoucher
}

type ShopStocktake struct {
	ID            int64
	ShopID        int64
	TotalQuantity int32

	CreatedBy int64
	UpdatedBy int64

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
	Price       int32
	ImageURL    string
	Attributes  []*Attribute
}

type Attribute struct {
	Name  string
	Value string
}
