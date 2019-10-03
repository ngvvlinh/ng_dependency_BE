package trading

type TradingOrderCreatingEvent struct {
	ReferralCode string
	UserID       int64
}

type TradingOrderCreatedEvent struct {
	OrderID      int64
	ReferralCode string
}
