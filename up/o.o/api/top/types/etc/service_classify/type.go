package service_classify

// +enum
type ServiceClassify int

type NullServiceClassify struct {
	Enum  ServiceClassify
	Valid bool
}

const (
	// +enum=shipping
	Shipping ServiceClassify = 0

	// +enum=telecom
	Telecom ServiceClassify = 1

	// +enum=all
	All ServiceClassify = 9
)
