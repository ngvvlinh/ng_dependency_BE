package ledgering

import (
	"time"

	"etop.vn/api/main/identity"
	"etop.vn/api/top/types/etc/ledger_type"
	dot "etop.vn/capi/dot"
)

type ShopLedger struct {
	ID          dot.ID
	ShopID      dot.ID
	Name        string
	BankAccount *identity.BankAccount
	Note        string
	Type        ledger_type.LedgerType
	Status      int
	CreatedBy   dot.ID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
