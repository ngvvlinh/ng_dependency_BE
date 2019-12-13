package refund

import (
	"time"

	catalog "etop.vn/api/main/catalog/types"
	"etop.vn/api/main/inventory"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
)

// +gen:event:topic=event/refund

type Refund struct {
	ID       dot.ID
	ShopID   dot.ID
	OrderID  dot.ID
	Note     string
	Code     string
	CodeNorm int
	Lines    []*RefundLine

	CreatedAt   time.Time
	UpdatedAt   time.Time
	CancelledAt time.Time
	ConfirmedAt time.Time

	CreatedBy    dot.ID
	UpdatedBy    dot.ID
	CancelReason string
	Discount     int
	Status       status3.Status
	CustomerID   dot.ID

	TotalAmount int
	BasketValue int
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
}

type ConfirmedRefundEvent struct {
	ShopID               dot.ID
	RefundID             dot.ID
	UpdatedBy            dot.ID
	AutoInventoryVoucher inventory.AutoInventoryVoucher
}

type CheckReceiptLinesResponse struct {
	CustomerID  dot.ID
	RefundLine  []*RefundLine
	BasketValue int
}
