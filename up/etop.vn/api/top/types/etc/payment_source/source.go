package payment_source

// +enum
// +enum:zero=null
type PaymentSource int

type NullPaymentSource struct {
	Enum  PaymentSource
	Valid bool
}

const (
	// +enum=unknown
	Unknown PaymentSource = 0

	// +enum=order
	PaymentSourceOrder PaymentSource = 1
)
