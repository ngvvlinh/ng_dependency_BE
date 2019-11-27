package purchaseorder

import (
	"time"

	"etop.vn/api/main/catalog"
	"etop.vn/api/main/etop"
	"etop.vn/api/main/inventory"
	"etop.vn/api/meta"
	dot "etop.vn/capi/dot"
)

// +gen:event:topic=event/purchaseorder

type PurchaseOrder struct {
	ID               dot.ID
	ShopID           dot.ID
	SupplierID       dot.ID
	Supplier         *PurchaseOrderSupplier
	InventoryVoucher *inventory.InventoryVoucher
	BasketValue      int
	TotalDiscount    int
	TotalAmount      int
	Code             string
	CodeNorm         int
	Note             string
	Status           etop.Status3
	VariantIDs       []dot.ID
	Lines            []*PurchaseOrderLine
	PaidAmount       int
	CreatedBy        dot.ID
	CancelledReason  string
	ConfirmedAt      time.Time
	CancelledAt      time.Time
	CreatedAt        time.Time `sq:"create"`
	UpdatedAt        time.Time `sq:"update"`
}

type PurchaseOrderLine struct {
	VariantID    dot.ID
	Quantity     int
	PaymentPrice int
	ProductID    dot.ID
	ProductName  string
	Code         string
	ImageUrl     string
	Attributes   []*catalog.Attribute
}

type PurchaseOrderSupplier struct {
	FullName           string
	Phone              string
	Email              string
	CompanyName        string
	TaxNumber          string
	HeadquarterAddress string
	Deleted            bool
}

func (r *PurchaseOrder) GetVariantIDs() []dot.ID {
	ids := make([]dot.ID, 0, len(r.Lines))
	for _, line := range r.Lines {
		if line.VariantID != 0 {
			ids = append(ids, line.VariantID)
		}
	}
	return ids
}

type PurchaseOrderConfirmedEvent struct {
	meta.EventMeta
	ShopID               dot.ID
	PurchaseOrderCode    string
	UserID               dot.ID
	PurchaseOrderID      dot.ID
	TraderID             dot.ID
	TotalAmount          int
	AutoInventoryVoucher inventory.AutoInventoryVoucher

	Lines []*inventory.InventoryVoucherItem
}
