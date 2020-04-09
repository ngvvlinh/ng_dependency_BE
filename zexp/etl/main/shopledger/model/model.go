package model

import (
	"time"

	"etop.vn/api/top/types/etc/ledger_type"
	identitysharemodel "etop.vn/backend/com/main/identity/sharemodel"
	"etop.vn/capi/dot"
)

// +sqlgen
type ShopLedger struct {
	ID          dot.ID
	ShopID      dot.ID
	Name        string
	BankAccount *identitysharemodel.BankAccount
	Note        string
	Type        ledger_type.LedgerType `sql_type:"enum(shop_ledger_type)"`
	Status      int
	CreatedBy   dot.ID
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Rid dot.ID
}
