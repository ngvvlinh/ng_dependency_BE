package tradering

import "etop.vn/api/meta"

// +gen:event:topic=event/tradering

var (
	CustomerType = "customer"
	SupplierType = "supplier"
	CarrierType  = "carrier"
)

type ShopTrader struct {
	ID       int64
	ShopID   int64
	Type     string
	FullName string
	Phone    string
}

type TraderDeletedEvent struct {
	meta.EventMeta
	ShopID   int64
	TraderID int64
}
