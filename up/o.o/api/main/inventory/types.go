package inventory

import "o.o/capi/dot"

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
