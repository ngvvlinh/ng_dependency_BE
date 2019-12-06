package status5

// +enum
type Status int

const (
	// +enum=Z
	Z Status = 0

	// +enum=P
	P Status = 1

	// +enum=S
	S Status = 2

	// +enum=N
	N Status = 127

	// +enum=NS
	NS Status = 126
)
