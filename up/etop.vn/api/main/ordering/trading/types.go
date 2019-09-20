package trading

type CheckTradingOrderValidEvent struct {
	OrderID      int64
	ReferralCode string
}

type TradingOrderCreatedEvent struct {
	OrderID      int64
	ReferralCode string
}
