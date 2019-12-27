package inventory_auto

// +enum
// +enum:zero=null
type AutoInventoryVoucher int

type NullAutoInventoryVoucher struct {
	Enum  AutoInventoryVoucher
	Valid bool
}

const (
	// +enum=unknown
	Unknown AutoInventoryVoucher = 0

	// +enum=create
	Create AutoInventoryVoucher = 1

	// +enum=confirm
	Confirm AutoInventoryVoucher = 2
)
