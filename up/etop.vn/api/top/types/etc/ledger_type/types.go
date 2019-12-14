package ledger_type

// +enum
// +enum:zero=null
type LedgerType int

type NullLedgerType struct {
	Enum  LedgerType
	Valid bool
}

const (
	// +enum=unknown
	Unknown LedgerType = 0

	// +enum=cash
	LedgerTypeCash LedgerType = 1

	// +enum=bank
	LedgerTypeBank LedgerType = 2
)
