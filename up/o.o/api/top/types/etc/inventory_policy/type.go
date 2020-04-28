package inventory_policy

// +enum
// +enum:zero=null
type InventoryPolicy int

type NullInventoryPolicy struct {
	Enum  InventoryPolicy
	Valid bool
}

const (
	// +enum=unknown
	Unknown InventoryPolicy = 0

	// +enum=obey
	Obey InventoryPolicy = 1

	// +enum=ignore
	Ignore InventoryPolicy = 2
)
