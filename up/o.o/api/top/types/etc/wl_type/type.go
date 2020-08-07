package wl_type

// +enum
// +enum:zero=null
type WhiteLabelType int

type NullWhiteLabelType struct {
	Enum  WhiteLabelType
	Valid bool
}

const (
	// +enum=unknown
	Unknown WhiteLabelType = 0

	// +enum=pos
	POS WhiteLabelType = 267

	// +enum=shipship
	ShipShip WhiteLabelType = 453
)
