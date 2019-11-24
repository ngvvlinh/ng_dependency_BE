package ledgering

import (
	"time"

	"etop.vn/api/main/identity"
	dot "etop.vn/capi/dot"
)

type LedgerType string

const (
	LedgerTypeCash LedgerType = "cash"
	LedgerTypeBank LedgerType = "bank"
)

type ShopLedger struct {
	ID          dot.ID
	ShopID      dot.ID
	Name        string
	BankAccount *identity.BankAccount
	Note        string
	Type        string
	Status      int32
	CreatedBy   dot.ID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
