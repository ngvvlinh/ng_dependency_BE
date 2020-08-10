package shipping_payment_type

// +enum
// +enum:sql=int
// +enum:zero=null
type ShippingPaymentType int

type NullShippingPaymentType struct {
	Enum  ShippingPaymentType
	Valid bool
}

const (
	// +enum=none
	None ShippingPaymentType = 0

	// người bán trả tiền
	// +enum=seller
	Seller ShippingPaymentType = 1

	// người mua trả tiền
	// +enum=buyer
	Buyer ShippingPaymentType = 2
)
