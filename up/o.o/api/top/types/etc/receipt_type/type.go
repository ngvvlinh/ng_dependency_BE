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
	// +enum:RefName:thu
	Receipt ReceiptType = 1

	// +enum=payment
	// +enum:RefName:chi
	Payment ReceiptType = 2
)
