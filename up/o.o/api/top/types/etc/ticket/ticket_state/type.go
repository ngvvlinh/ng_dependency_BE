package ticket_state

// +enum
// +enum:sql=int
type TicketState int

type NullTicketState struct {
	Enum  TicketState
	Valid bool
}

const (
	// +enum=unknown
	Unknown TicketState = 0

	// +enum=new
	New TicketState = 45

	// +enum=received
	Received TicketState = 21

	// +enum=processing
	Processing TicketState = 54

	// +enum=success
	Success TicketState = 71

	// +enum=fail
	Fail TicketState = 84

	// +enum=ignore
	Ignore TicketState = 27

	// +enum=cancel
	Cancel TicketState = 68
)
