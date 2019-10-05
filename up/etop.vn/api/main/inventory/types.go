package inventory

// +gen:event:topic=event/inventory

type InventoryVoucherCreatedEvent struct {
	ShopID             int64
	InventoryVoucherID int64
	Line               []*InventoryVoucherItem
}
