package status3

// +enum
// +enum:sql=int
type Status int

type NullStatus struct {
	Enum  Status
	Valid bool
}

const (
	// +enum=Z
	Z Status = 0

	// +enum=P
	P Status = 1

	// +enum=N
	N Status = -1
)
