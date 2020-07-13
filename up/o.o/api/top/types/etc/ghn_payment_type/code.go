package ghn_payment_type

// +enum
// +enum:zero=null
type GHNPaymentType int

type NullGHNPaymentType struct {
	Enum  GHNPaymentType
	Valid bool
}

const (
	// +enum=seller
	SELLER GHNPaymentType = 1

	// +enum=buyer
	BUYER GHNPaymentType = 2
)
