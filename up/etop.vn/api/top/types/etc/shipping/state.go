package shipping

// +enum
type State int

const (
	// +enum=default
	Default State = 0

	// +enum=created
	Created State = 1

	// +enum=confirmed
	Confirmed State = 2

	// +enum=processing
	Processing State = 3

	// +enum=picking
	Picking State = 4

	// +enum=holding
	Holding State = 5

	// +enum=returning
	Returning State = 6

	// +enum=returned
	Returned State = 7

	// +enum=delivering
	Delivering State = 8

	// +enum=delivered
	Delivered State = 9

	// +enum=unknown
	Unknown State = 101

	// +enum=undeliverable
	Undeliverable State = 126

	// +enum=cancelled
	Cancelled State = 127
)
