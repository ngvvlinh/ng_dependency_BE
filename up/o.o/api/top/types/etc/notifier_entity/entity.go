package notifier_entity

// +enum
// +enum:zero=null
type NotifierEntity int

type NullNotifierEntity struct {
	Enum  NotifierEntity
	Valid bool
}

const (
	// +enum=unknown
	Unknown NotifierEntity = 0

	// +enum=fulfillment
	Fulfillment NotifierEntity = 1

	// +enum=money_transaction_shipping
	MoneyTransactionShipping NotifierEntity = 2

	// +enum=fb_external_comment
	FaboComment NotifierEntity = 30

	// +enum=fb_external_message
	FaboMessage NotifierEntity = 36
)
