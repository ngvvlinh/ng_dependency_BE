package address_type

// +enum
// +enum:zero=null
type AddressType int

const (
	// +enum=unknown
	Unknown AddressType = 0

	// +enum=general
	General AddressType = 1

	// +enum=warehouse
	Warehouse AddressType = 2

	// +enum=shipto
	Shipto AddressType = 3

	// +enum=shipfrom
	Shipfrom AddressType = 4
)
