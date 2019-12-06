package notifier_entity

// +enum
type NotifierEntity int

const (
	// +enum=unknown
	Unknown NotifierEntity = 0

	// +enum=fulfillment
	Fulfillment NotifierEntity = 1

	// +enum=money_transaction_shipping
	MoneyTransactionShipping NotifierEntity = 2
)
