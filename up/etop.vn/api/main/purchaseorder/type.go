package purchaseorder

import (
	"time"

	catalogtypes "etop.vn/api/main/catalog/types"
	"etop.vn/api/main/inventory"
	"etop.vn/api/meta"
	"etop.vn/api/top/int/types"
	"etop.vn/api/top/types/etc/inventory_auto"
	"etop.vn/api/top/types/etc/status3"
	dot "etop.vn/capi/dot"
)

// +gen:event:topic=event/purchaseorder

type PurchaseOrder struct {
	ID               dot.ID
	ShopID           dot.ID
	SupplierID       dot.ID
	Supplier         *PurchaseOrderSupplier
	InventoryVoucher *inventory.InventoryVoucher
	DiscountLines    []*types.DiscountLine
	TotalDiscount    int
	FeeLines         []*types.FeeLine
	TotalFee         int
	TotalAmount      int
	BasketValue      int
	Code             string
	CodeNorm         int
	Note             string
	Status           status3.Status
	VariantIDs       []dot.ID
	Lines            []*PurchaseOrderLine
	PaidAmount       int
	CreatedBy        dot.ID
	CancelReason     string
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
	Attributes   []*catalogtypes.Attribute
	Discount     int
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
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher

	Lines []*inventory.InventoryVoucherItem
}

type PurchaseOrderCancelledEvent struct {
	PurchaseOrderID      dot.ID
	ShopID               dot.ID
	UpdatedBy            dot.ID
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher
	InventoryOverStock   bool
}
