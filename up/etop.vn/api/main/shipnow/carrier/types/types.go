package types

// +enum
type Carrier int

type NullCarrier struct {
	Enum  Carrier
	Valid bool
}

const (
	// +enum=default
	Default Carrier = 0

	// +enum=ahamove
	Ahamove Carrier = 1
)
