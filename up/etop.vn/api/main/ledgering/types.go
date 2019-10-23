package ledgering

import (
	"time"

	"etop.vn/api/main/identity"
	"etop.vn/api/meta"
)

// +gen:event:topic=event/ledgering

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

type AccountCreatedEvent struct {
	meta.EventMeta
	ShopID int64
	UserID int64
}
