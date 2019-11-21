package purchaseorder

import (
	"time"

	"etop.vn/api/main/catalog"
	"etop.vn/api/main/etop"
	"etop.vn/api/main/inventory"
	"etop.vn/api/meta"
)

// +gen:event:topic=event/purchaseorder

type PurchaseOrder struct {
	ID               int64
	ShopID           int64
	SupplierID       int64
	Supplier         *PurchaseOrderSupplier
	InventoryVoucher *inventory.InventoryVoucher
	BasketValue      int64
	TotalDiscount    int64
	TotalAmount      int64
	Code             string
	CodeNorm         int32
	Note             string
	Status           etop.Status3
	VariantIDs       []int64
	Lines            []*PurchaseOrderLine
	PaidAmount       int64
	CreatedBy        int64
	CancelledReason  string
	ConfirmedAt      time.Time
	CancelledAt      time.Time
	CreatedAt        time.Time `sq:"create"`
	UpdatedAt        time.Time `sq:"update"`
}

type PurchaseOrderLine struct {
	VariantID   int64
	Quantity    int64
	Price       int64
	ProductID   int64
	ProductName string
	Code        string
	ImageUrl    string
	Attributes  []*catalog.Attribute
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

func (r *PurchaseOrder) GetVariantIDs() []int64 {
	ids := make([]int64, 0, len(r.Lines))
	for _, line := range r.Lines {
		if line.VariantID != 0 {
			ids = append(ids, line.VariantID)
		}
	}
	return ids
}

type PurchaseOrderConfirmedEvent struct {
	meta.EventMeta
	ShopID               int64
	PurchaseOrderCode    string
	UserID               int64
	PurchaseOrderID      int64
	TraderID             int64
	TotalAmount          int64
	AutoInventoryVoucher inventory.AutoInventoryVoucher

	Lines []*inventory.InventoryVoucherItem
}
