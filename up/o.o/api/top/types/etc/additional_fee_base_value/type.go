package additional_fee_base_value

// +enum
// +enum:zero=null
type BaseValueType int

type NullBaseValueType struct {
	Enum  BaseValueType
	Valid bool
}

const (
	// +enum=unknown
	Unknown BaseValueType = 0

	// +enum=cod_amount
	CODAmount BaseValueType = 1

	// +enum=main_fee
	MainFee BaseValueType = 2

	// +enum=basket_value
	BasketValue BaseValueType = 3
)
