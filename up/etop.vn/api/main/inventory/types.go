package inventory

// +gen:event:topic=event/inventory

type InventoryVoucherCreatingEvent struct {
	ShopID             int64
	InventoryVoucherID int64
	Line               []*InventoryVoucherItem
}

type InventoryVoucherUpdatingEvent struct {
	ShopID             int64
	InventoryVoucherID int64
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

type InventoryRefType string

const (
	RefTypeReturns       InventoryRefType = "return"
	RefTypeStockTake     InventoryRefType = "stocktake"
	RefTypePurchaseOrder InventoryRefType = "purchase_order"
	RefTypeOrder         InventoryRefType = "order"
)

type InventoryVoucherType string

const (
	InventoryVoucherTypeIn  InventoryVoucherType = "in"
	InventoryVoucherTypeOut InventoryVoucherType = "out"
)

type AutoInventoryVoucher string

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
