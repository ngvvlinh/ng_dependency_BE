package ledgering

import (
	"time"

	"etop.vn/api/main/identity"
)

type LedgerType string

const (
	LedgerTypeCash LedgerType = "cash"
	LedgerTypeBank LedgerType = "bank"
)

type ShopLedger struct {
	ID          int64
	ShopID      int64
	Name        string
	BankAccount *identity.BankAccount
	Note        string
	Type        string
	Status      int32
	CreatedBy   int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
