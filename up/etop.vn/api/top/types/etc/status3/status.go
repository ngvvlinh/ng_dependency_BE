package status3

// +enum
// +enum:sql=int
type Status int

const (
	// +enum=Z
	Z Status = 0

	// +enum=P
	P Status = 1

	// +enum=N
	N Status = -1
)
