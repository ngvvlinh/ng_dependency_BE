package ticket_type

// +enum
// +enum:sql=int
type TicketType int

type NullTicketType struct {
	Enum  TicketType
	Valid bool
}

const (
	// +enum=unknown
	Unknown TicketType = 0

	// +enum=internal
	Internal TicketType = 12

	// +enum=system
	System TicketType = 74
)
