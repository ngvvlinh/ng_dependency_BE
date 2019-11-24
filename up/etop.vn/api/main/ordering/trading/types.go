package trading

import (
	"etop.vn/api/meta"
	"etop.vn/capi/dot"
)

// +gen:event:topic=event/trading_order

type TradingOrderCreatingEvent struct {
	meta.EventMeta

	OrderID      dot.ID
	ReferralCode string
	UserID       dot.ID
}

type TradingOrderCreatedEvent struct {
	meta.EventMeta

	OrderID      dot.ID
	ReferralCode string
}
