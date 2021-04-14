package invoice_type

// +enum
type InvoiceType int

type NullInvoiceType struct {
	Enum  InvoiceType
	Valid bool
}

const (
	// +enum=default
	Default InvoiceType = 0

	// +enum=in
	In InvoiceType = 1

	// +enum=out
	Out InvoiceType = 9
)
