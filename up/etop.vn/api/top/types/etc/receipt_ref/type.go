package receipt_ref

// +enum
// +enum:zero=null
type ReceiptRef int

type NullReceiptRef struct {
	Enum  ReceiptRef
	Valid bool
}

const (
	// +enum=unknown
	Unknown ReceiptRef = 0

	// +enum=order
	ReceiptRefTypeOrder ReceiptRef = 1

	// +enum=fulfillment
	ReceiptRefTypeFulfillment ReceiptRef = 2

	// +enum=purchase_order
	ReceiptRefTypePurchaseOrder ReceiptRef = 3

	// +enum=refund
	ReceiptRefTypeRefund ReceiptRef = 4
)
