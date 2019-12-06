package status3

// +enum
type Status int

const (
	// +enum=Z
	Z Status = 0

	// +enum=P
	P Status = 1

	// +enum=N
	N Status = 127
)
