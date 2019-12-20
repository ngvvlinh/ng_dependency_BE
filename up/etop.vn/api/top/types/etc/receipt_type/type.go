package receipt_type

// +enum
// +enum:zero=null
type ReceiptType int

type NullReceiptType struct {
	Enum  ReceiptType
	Valid bool
}

const (
	// +enum=unknown
	Unknown ReceiptType = 0

	// +enum=receipt
	Receipt ReceiptType = 1

	// +enum=payment
	Payment ReceiptType = 2
)
