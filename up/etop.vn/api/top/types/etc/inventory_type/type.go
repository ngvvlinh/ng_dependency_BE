package inventory_type

// +enum
// +enum:zero=null
type InventoryVoucherType int

type NullInventoryVoucherType struct {
	Enum  InventoryVoucherType
	Valid bool
}

const (
	// +enum=unknown
	Unknown InventoryVoucherType = 0

	// +enum=in
	In InventoryVoucherType = 1

	// +enum=out
	Out InventoryVoucherType = 2
)
