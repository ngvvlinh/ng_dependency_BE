package receipt_mode

// +enum
// +enum:zero=null
type ReceiptMode int

type NullReceiptMode struct {
	Enum  ReceiptMode
	Valid bool
}

const (
	// +enum=unknown
	Unknown ReceiptMode = 0

	// +enum=manual
	Manual ReceiptMode = 1

	// +enum=auto
	Auto ReceiptMode = 2
)
