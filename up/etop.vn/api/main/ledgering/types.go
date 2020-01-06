package ledgering

import (
	"time"

	identitytypes "etop.vn/api/main/identity/types"
	"etop.vn/api/top/types/etc/ledger_type"
	dot "etop.vn/capi/dot"
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
