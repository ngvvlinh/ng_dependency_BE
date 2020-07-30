package ticket_source

// +enum
// +enum:sql=int
type TicketSource int

type NullTicketSource struct {
	Enum  TicketSource
	Valid bool
}

const (
	// +enum=pos_web
	POSWeb TicketSource = 65

	// +enum=pos_app
	POSApp TicketSource = 12

	// +enum=shipment_app
	ShipmentApp TicketSource = 74

	// +enum=admin
	Admin TicketSource = 38

	// +enum=system
	System TicketSource = 42
)
