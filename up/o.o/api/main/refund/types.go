package refund

import (
	"time"

	catalog "o.o/api/main/catalog/types"
	"o.o/api/top/int/types"
	"o.o/api/top/types/etc/inventory_auto"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +gen:event:topic=event/refund

type Refund struct {
	ID              dot.ID
	ShopID          dot.ID
	OrderID         dot.ID
	Note            string
	Code            string
	CodeNorm        int
	Lines           []*RefundLine
	CreatedAt       time.Time
	UpdatedAt       time.Time
	CancelledAt     time.Time
	ConfirmedAt     time.Time
	CreatedBy       dot.ID
	UpdatedBy       dot.ID
	CancelReason    string
	AdjustmentLines []*types.AdjustmentLine
	TotalAdjustment int
	Status          status3.Status
	CustomerID      dot.ID
	TotalAmount     int
	BasketValue     int
}

type RefundLine struct {
	VariantID   dot.ID
	Quantity    int
	Code        string
	ImageURL    string
	RetailPrice int
	ProductID   dot.ID
	ProductName string
	Attributes  []*catalog.Attribute
	Adjustment  int
}

type RefundConfirmedEvent struct {
	ShopID               dot.ID
	RefundID             dot.ID
	UpdatedBy            dot.ID
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher
}

type RefundCancelledEvent struct {
	ShopID               dot.ID
	RefundID             dot.ID
	UpdatedBy            dot.ID
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher
}

type CheckReceiptLinesResponse struct {
	CustomerID  dot.ID
	RefundLine  []*RefundLine
	BasketValue int
}
