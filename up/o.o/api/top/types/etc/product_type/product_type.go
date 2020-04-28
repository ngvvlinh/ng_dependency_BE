package product_type

// +enum
// +enum:zero=null
type ProductType int

type NullProductType struct {
	Enum  ProductType
	Valid bool
}

const (
	// +enum=unknown
	Unknown ProductType = 0

	// +enum=services
	Services ProductType = 1

	// +enum=goods
	Goods ProductType = 2
)
