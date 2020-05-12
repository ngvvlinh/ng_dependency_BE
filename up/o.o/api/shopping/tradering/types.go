package tradering

import (
	"o.o/api/meta"
	"o.o/capi/dot"
)

// +gen:event:topic=event/tradering

var (
	CustomerType = "customer"
	SupplierType = "supplier"
	CarrierType  = "carrier"
)

type ShopTrader struct {
	ID       dot.ID
	ShopID   dot.ID
	Type     string
	FullName string
	Phone    string
}

type TraderDeletedEvent struct {
	meta.EventMeta
	ShopID      dot.ID
	TraderID    dot.ID
	TradingType string
}
