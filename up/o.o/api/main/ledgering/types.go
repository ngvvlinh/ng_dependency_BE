package ledgering

import (
	"time"

	identitytypes "o.o/api/main/identity/types"
	"o.o/api/top/types/etc/ledger_type"
	dot "o.o/capi/dot"
)

type ShopLedger struct {
	ID          dot.ID
	ShopID      dot.ID
	Name        string
	BankAccount *identitytypes.BankAccount
	Note        string
	Type        ledger_type.LedgerType
	Status      int
	CreatedBy   dot.ID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
