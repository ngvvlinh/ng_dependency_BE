package purchaserefund

import (
	"time"

	catalog "etop.vn/api/main/catalog/types"
	"etop.vn/api/top/types/etc/inventory_auto"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
)

// +gen:event:topic=event/purchaserefund

type PurchaseRefund struct {
	ID              dot.ID
	ShopID          dot.ID
	PurchaseOrderID dot.ID
	Note            string
	Code            string
	CodeNorm        int
	Lines           []*PurchaseRefundLine

	CreatedAt   time.Time
	UpdatedAt   time.Time
	CancelledAt time.Time
	ConfirmedAt time.Time

	CreatedBy    dot.ID
	UpdatedBy    dot.ID
	CancelReason string
	Discount     int
	Status       status3.Status
	SupplierID   dot.ID

	TotalAmount int
	BasketValue int
}

type PurchaseRefundLine struct {
	VariantID    dot.ID
	Quantity     int
	Code         string
	ImageURL     string
	PaymentPrice int
	ProductID    dot.ID
	ProductName  string
	Attributes   []*catalog.Attribute
}

type ConfirmedPurchaseRefundEvent struct {
	ShopID               dot.ID
	PurchaseRefundID     dot.ID
	UpdatedBy            dot.ID
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher
	InventoryOverStock   bool
}

type CheckReceiptLinesResponse struct {
	SupplierID         dot.ID
	PurchaseRefundLine []*PurchaseRefundLine
	BasketValue        int
}

type PurchaseRefundCancelledEvent struct {
	PurchaseRefundID     dot.ID
	ShopID               dot.ID
	UpdatedBy            dot.ID
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher
	InventoryOverStock   bool
}