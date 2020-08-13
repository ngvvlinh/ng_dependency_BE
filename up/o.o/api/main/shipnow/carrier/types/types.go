package types

// +enum
type ShipnowCarrier int

type NullShipnowCarrier struct {
	Enum  ShipnowCarrier
	Valid bool
}

const (
	// +enum=default
	Default ShipnowCarrier = 0

	// +enum=ahamove
	Ahamove ShipnowCarrier = 1
)
