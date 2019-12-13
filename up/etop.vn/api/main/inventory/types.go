package inventory

import "etop.vn/capi/dot"

// +gen:event:topic=event/inventory

type InventoryVoucherCreatingEvent struct {
	ShopID             dot.ID
	InventoryVoucherID dot.ID
	Line               []*InventoryVoucherItem
}

type InventoryVoucherUpdatingEvent struct {
	ShopID             dot.ID
	InventoryVoucherID dot.ID
	Line               []*InventoryVoucherItem
}

type InventoryVoucherRefName string

const (
	RefNameReturns       InventoryVoucherRefName = "Trả hàng"
	RefNameStockTake     InventoryVoucherRefName = "Kiểm kho"
	RefNamePurchaseOrder InventoryVoucherRefName = "Nhập hàng"
	RefNameOrder         InventoryVoucherRefName = "Bán hàng"
	RefNameCancelOrder   InventoryVoucherRefName = "Hủy đơn hàng"
)

func (s InventoryVoucherRefName) String() string { return string(s) }

type InventoryRefType string

const (
	RefTypeRefund        InventoryRefType = "refund"
	RefTypeStockTake     InventoryRefType = "stocktake"
	RefTypePurchaseOrder InventoryRefType = "purchase_order"
	RefTypeOrder         InventoryRefType = "order"
)

func (s InventoryRefType) String() string { return string(s) }

type InventoryVoucherType string

func (s InventoryVoucherType) String() string { return string(s) }

const (
	InventoryVoucherTypeIn  InventoryVoucherType = "in"
	InventoryVoucherTypeOut InventoryVoucherType = "out"
)

type AutoInventoryVoucher string

func (s AutoInventoryVoucher) String() string { return string(s) }

const (
	AutoCreateInventory           AutoInventoryVoucher = "create"
	AutoCreateAndConfirmInventory AutoInventoryVoucher = "confirm"
)

func (s AutoInventoryVoucher) ValidateAutoInventoryVoucher() bool {
	if s == AutoCreateInventory || s == AutoCreateAndConfirmInventory {
		return true
	}
	return false
}
