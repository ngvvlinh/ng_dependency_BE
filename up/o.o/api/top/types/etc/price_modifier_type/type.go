package price_modifier_type

// +enum
// +enum:zero=null
type PriceModifierType int

type NullPriceModifierType struct {
	Enum  PriceModifierType
	Valid bool
}

const (
	// +enum=percentage
	Percentage PriceModifierType = 1

	// +enum=fixed_amount
	FixedAmount PriceModifierType = 2
)
