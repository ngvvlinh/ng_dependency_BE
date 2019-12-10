package payment_provider

// +enum
// +enum:zero=null
type PaymentProvider int

const (
	// +enum=unknown
	Unknown PaymentProvider = 0

	// +enum=vtpay
	Vtpay PaymentProvider = 1
)
