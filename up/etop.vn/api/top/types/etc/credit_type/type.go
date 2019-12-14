package credit_type

// +enum
// +enum:zero=null
type CreditType int

type NullCreditType struct {
	Enum  CreditType
	Valid bool
}

const (
	// +enum=shop
	Shop CreditType = 1
)
