package transaction_type

// +enum
// +enum:zero=null
type TransactionType int

type NullTransactionType struct {
	Enum  TransactionType
	Valid bool
}

const (
	// +enum=credit
	Credit TransactionType = 9

	// +enum=invoice
	Invoice TransactionType = 13
)
