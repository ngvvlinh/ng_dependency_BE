package payment_provider

// +enum
// +enum:zero=null
type PaymentProvider int

type NullPaymentProvider struct {
	Enum  PaymentProvider
	Valid bool
}

const (
	// +enum=unknown
	Unknown PaymentProvider = 0

	// +enum=vtpay
	VTPay PaymentProvider = 1
)
