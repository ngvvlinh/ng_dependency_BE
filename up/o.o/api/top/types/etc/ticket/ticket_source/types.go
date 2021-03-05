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

	// +enum=telecom_app_call
	TelecomAppCall TicketSource = 101

	// +enum=telecom_web
	TelecomWeb TicketSource = 77

	// +enum=telecom_ext_call
	TelecomExtCall TicketSource = 79

	// +enum=telecom_web_call
	TelecomWebCall TicketSource = 42

	// +enum=shipment_web
	ShipmentWeb TicketSource = 85

	// +enum=admin
	Admin TicketSource = 38

	// +enum=system
	System TicketSource = 42

	// +enum=webphone
	WebPhone TicketSource = 98
)
