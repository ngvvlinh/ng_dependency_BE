package customer_type

// +enum
// +enum:zero=null
type CustomerType int

type NullCustomerType struct {
	Enum  CustomerType
	Valid bool
}

const (
	// +enum=unknown
	Unknown CustomerType = 0

	// +enum=individual
	Individual CustomerType = 1

	// +enum=organization
	Organization CustomerType = 2

	// +enum=anonymous,independent
	Independent CustomerType = 3
)
