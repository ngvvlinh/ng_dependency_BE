package model

import (
	"time"

	"etop.vn/api/top/types/etc/ledger_type"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenShopLedger(&ShopLedger{})

type ShopLedger struct {
	ID          dot.ID
	ShopID      dot.ID
	Name        string
	BankAccount *model.BankAccount
	Note        string
	Type        ledger_type.LedgerType
	Status      int
	CreatedBy   dot.ID
	CreatedAt   time.Time `sq:"create"`
	UpdatedAt   time.Time `sq:"update"`
	DeletedAt   time.Time
}
