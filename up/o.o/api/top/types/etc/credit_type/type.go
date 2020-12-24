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

// +enum
type CreditClassify int

type NullCreditClassify struct {
	Enum  CreditClassify
	Valid bool
}

const (
	// +enum=shipping
	CreditClassifyShipping CreditClassify = 0

	// +enum=telecom
	CreditClassifyTelecom CreditClassify = 1
)
