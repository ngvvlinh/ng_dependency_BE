package shipping_fee_type

// +enum
type ShippingFeeType int

const (
	// +enum=main
	Main ShippingFeeType = 0

	// +enum=return
	Return ShippingFeeType = 1

	// +enum=adjustment
	Adjustment ShippingFeeType = 2

	// +enum=insurance
	Insurance ShippingFeeType = 3

	// +enum=tax
	Tax ShippingFeeType = 4

	// +enum=other
	Other ShippingFeeType = 5

	// +enum=cods
	Cods ShippingFeeType = 6

	// +enum=address_change
	AddressChange ShippingFeeType = 7

	// +enum=discount
	Discount ShippingFeeType = 8

	// +enum=unknown
	Unknown ShippingFeeType = 127
)
