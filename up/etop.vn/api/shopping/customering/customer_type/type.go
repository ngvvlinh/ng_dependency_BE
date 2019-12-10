package customer_type

// +enum
// +enum:zero=null
type CustomerType int

const (
	// +enum=unknown
	Unknown CustomerType = 0

	// +enum=individual
	Individual CustomerType = 1

	// +enum=organization
	Organization CustomerType = 2

	// +enum=independent
	Independent CustomerType = 3
)
