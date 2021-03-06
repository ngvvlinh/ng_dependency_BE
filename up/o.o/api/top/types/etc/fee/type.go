package fee

// +enum
// +enum:zero=null
type FeeType int

type NullFeeType struct {
	Enum  FeeType
	Valid bool
}

const (
	// +enum=other
	Other FeeType = 0

	// +enum=shipping
	Shipping FeeType = 1

	// +enum=tax
	Tax FeeType = 2
)
