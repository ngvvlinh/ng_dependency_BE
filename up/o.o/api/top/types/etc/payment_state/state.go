package payment_state

// +enum
// +enum:zero=null
type PaymentState int

type NullPaymentState struct {
	Enum  PaymentState
	Valid bool
}

const (
	// +enum=unknown
	Unknown PaymentState = 0

	// +enum=default
	Default PaymentState = 1

	// +enum=created
	Created PaymentState = 2

	// +enum=pending
	Pending PaymentState = 3

	// +enum=success
	Success PaymentState = 4

	// +enum=failed
	Failed PaymentState = 5

	// +enum=cancelled
	Cancelled PaymentState = 6
)
