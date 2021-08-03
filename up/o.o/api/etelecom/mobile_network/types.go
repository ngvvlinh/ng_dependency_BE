package mobile_network

// +enum
// +enum:zero=null
type MobileNetwork int

type NullMobileNetwork struct {
	Enum  MobileNetwork
	Valid bool
}

const (
	// +enum=unknown
	Unknown MobileNetwork = 0

	// +enum=mobifone
	MobiFone MobileNetwork = 1

	// +enum=vinaphone
	Vinaphone MobileNetwork = 2

	// +enum=viettel
	Viettel MobileNetwork = 3

	// +enum=other
	Other MobileNetwork = 4
)
