package receipt_ref

// +enum
// +enum:zero=null
type ReceiptRef int

type NullReceiptRef struct {
	Enum  ReceiptRef
	Valid bool
}

const (
	// +enum=none
	None ReceiptRef = 0

	// +enum=order
	Order ReceiptRef = 1

	// +enum=fulfillment
	Fulfillment ReceiptRef = 2

	// +enum=purchase_order
	PurchaseOrder ReceiptRef = 3

	// +enum=refund
	Refund ReceiptRef = 4

	// +enum=purchase_refund
	PurchaseRefund ReceiptRef = 5
)
