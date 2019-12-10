package payment_source

// +enum
// +enum:zero=null
type PaymentSource int

const (
	// +enum=unknown
	Unknown PaymentSource = 0

	// +enum=order
	PaymentSourceOrder PaymentSource = 1
)
