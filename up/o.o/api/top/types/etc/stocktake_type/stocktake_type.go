package stocktake_type

// +enum
// +enum:zero=null
type StocktakeType int

type NullStocktakeType struct {
	Enum  StocktakeType
	Valid bool
}

const (
	// +enum=balance
	Balance StocktakeType = 0

	// +enum=discard
	Discard StocktakeType = 1
)
