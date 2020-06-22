package shipping_fee_type

// +enum
// +enum:zero=null
type ShippingFeeType int

type NullShippingFeeType struct {
	Enum  ShippingFeeType
	Valid bool
}

const (
	// +enum=unknown
	Unknown ShippingFeeType = 0

	// +enum=main
	Main ShippingFeeType = 1

	// +enum=return
	Return ShippingFeeType = 2

	// +enum=adjustment
	Adjustment ShippingFeeType = 3

	// +enum=insurance
	Insurance ShippingFeeType = 4

	// +enum=tax
	Tax ShippingFeeType = 5

	// +enum=other
	Other ShippingFeeType = 6

	// +enum=cods
	Cods ShippingFeeType = 7

	// +enum=address_change
	AddressChange ShippingFeeType = 8

	// +enum=discount
	Discount ShippingFeeType = 9

	// +enum=redelivery
	Redelivery ShippingFeeType = 10
)

func Contain(list []ShippingFeeType, _type ShippingFeeType) bool {
	for _, item := range list {
		if item == _type {
			return true
		}
	}
	return false
}
