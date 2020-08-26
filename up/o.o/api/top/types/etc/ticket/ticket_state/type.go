package ticket_state

import "o.o/api/top/types/etc/status5"

// +enum
// +enum:sql=int
type TicketState int

type NullTicketState struct {
	Enum  TicketState
	Valid bool
}

const (
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

func (s TicketState) ToStatus5() status5.Status {
	switch s {
	case New:
		return status5.Z
	case Received, Processing:
		return status5.S
	case Success:
		return status5.P
	case Fail:
		return status5.NS
	case Ignore, Cancel:
		return status5.N
	default:
		return status5.Z
	}
}
