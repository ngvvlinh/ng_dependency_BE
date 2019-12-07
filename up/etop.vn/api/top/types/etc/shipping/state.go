package shipping

import (
	"etop.vn/api/top/types/etc/status4"
	"etop.vn/api/top/types/etc/status5"
)

// +enum
type State int

const (
	// +enum=unknown
	Unknown State = 0

	// +enum=default
	Default State = 1

	// +enum=created
	Created State = 2

	// +enum=confirmed
	Confirmed State = 3

	// +enum=processing
	Processing State = 4

	// +enum=picking
	Picking State = 5

	// +enum=holding
	Holding State = 6

	// +enum=returning
	Returning State = 7

	// +enum=returned
	Returned State = 8

	// +enum=delivering
	Delivering State = 9

	// +enum=delivered
	Delivered State = 10

	// +enum=cancelled
	Cancelled State = -1

	// +enum=undeliverable
	Undeliverable State = -2
)

func (s State) Text() string {
	name := s.Name()
	if name == "" {
		return "Không xác định"
	}
	return name
}

func (s State) ToStatus4() status4.Status {
	switch s {
	case Default:
		return status4.Z
	case Cancelled, Returned:
		return status4.N
	case Delivered:
		return status4.P
	}
	return status4.S
}

func (s State) ToStatus5() status5.Status {
	switch s {
	case Default:
		return status5.Z
	case Cancelled:
		return status5.N
	case Returning, Returned:
		return status5.NS
	case Delivered:
		return status5.P
	default:
		return status5.S
	}
}
