package payment_provider

// +enum
type PaymentProvider int

const (
	// +enum=unknown
	Unknown PaymentProvider = 0

	// +enum=vtpay
	Vtpay PaymentProvider = 1
)
