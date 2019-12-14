package shipnow_state

import "etop.vn/api/top/types/etc/status5"

// +enum
type State int

type NullState struct {
	Enum  State
	Valid bool
}

const (
	// +enum=default
	StateDefault State = 0

	// +enum=created
	StateCreated State = 1

	// +enum=assigning
	StateAssigning State = 2

	// +enum=picking
	StatePicking State = 3

	// +enum=delivering
	StateDelivering State = 4

	// +enum=delivered
	StateDelivered State = 5

	// +enum=returning
	StateReturning State = 6

	// +enum=returned
	StateReturned State = 7

	// +enum=unknown
	StateUnknown State = 101

	// +enum=undeliverable
	StateUndeliverable State = 126

	// +enum=cancelled
	StateCancelled State = 127
)

func (s State) ToStatus5() status5.Status {
	switch s {
	case StateDefault, StateCreated:
		return status5.Z
	case StateCancelled:
		return status5.N
	case StateReturned, StateReturning:
		return status5.NS
	case StateDelivered:
		return status5.P
	default:
		return status5.S
	}
}
