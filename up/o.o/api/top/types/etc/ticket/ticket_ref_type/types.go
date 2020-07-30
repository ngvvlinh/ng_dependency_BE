package ticket_ref_type

// +enum
// +enum:sql=int
type TicketRefType int

type NullTicketRefType struct {
	Enum  TicketRefType
	Valid bool
}

const (
	// +enum=order_trading
	OrderTrading TicketRefType = 34

	// +enum=ffm
	FFM TicketRefType = 42

	// +enum=money_transaction
	MoneyTransaction TicketRefType = 95

	// +enum=other
	Other TicketRefType = 31
)
