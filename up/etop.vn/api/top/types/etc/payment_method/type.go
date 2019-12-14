package payment_method

// +enum
// +enum:zero=null
type PaymentMethod int

type NullPaymentMethod struct {
	Enum  PaymentMethod
	Valid bool
}

const (
	// +enum=unknown
	Unknown PaymentMethod = 0

	// +enum=cod
	COD PaymentMethod = 1

	// +enum=bank
	Bank PaymentMethod = 2

	// +enum=other
	Other PaymentMethod = 3

	// +enum=vtpay
	VTPay PaymentMethod = 4
)
