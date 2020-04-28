package model

import (
	"time"

	"o.o/api/top/types/etc/ledger_type"
	identitysharemodel "o.o/backend/com/main/identity/sharemodel"
	"o.o/capi/dot"
)

// +sqlgen
type ShopLedger struct {
	ID          dot.ID
	ShopID      dot.ID
	Name        string
	BankAccount *identitysharemodel.BankAccount
	Note        string
	Type        ledger_type.LedgerType
	Status      int
	CreatedBy   dot.ID
	CreatedAt   time.Time `sq:"create"`
	UpdatedAt   time.Time `sq:"update"`
	DeletedAt   time.Time

	Rid dot.ID
}
