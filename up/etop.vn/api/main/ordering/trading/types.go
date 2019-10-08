package trading

import "etop.vn/api/meta"

// +gen:event:topic=event/trading_order

type TradingOrderCreatingEvent struct {
	meta.EventMeta

	OrderID      int64
	ReferralCode string
	UserID       int64
}

type TradingOrderCreatedEvent struct {
	meta.EventMeta

	OrderID      int64
	ReferralCode string
}
