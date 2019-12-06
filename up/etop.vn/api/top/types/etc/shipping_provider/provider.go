package shipping_provider

// +enum
type ShippingProvider int

const (
	// +enum=unknown
	Unknown ShippingProvider = 0

	// +enum=all
	All ShippingProvider = 22

	// +enum=manual
	Manual ShippingProvider = 20

	// +enum=ghn
	Ghn ShippingProvider = 19

	// +enum=ghtk
	Ghtk ShippingProvider = 21

	// +enum=vtpost
	Vtpost ShippingProvider = 23
)
