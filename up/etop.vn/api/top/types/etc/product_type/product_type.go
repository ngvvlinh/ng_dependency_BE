package product_type

// +enum
type ProductType int

const (
	// +enum=unknown
	Unknown ProductType = 0

	// +enum=services
	Services ProductType = 1

	// +enum=goods
	Goods ProductType = 2
)
