package convert

import (
	"time"

	"o.o/api/main/ledgering"
	cm "o.o/backend/pkg/common"
)

// +gen:convert: o.o/backend/com/main/ledgering/model  -> o.o/api/main/ledgering
// +gen:convert: o.o/api/main/ledgering

func createShopLedger(args *ledgering.CreateLedgerArgs, out *ledgering.ShopLedger) {
	apply_ledgering_CreateLedgerArgs_ledgering_ShopLedger(args, out)
	out.ID = cm.NewID()
}

func updateShopLedger(args *ledgering.UpdateLedgerArgs, out *ledgering.ShopLedger) {
	apply_ledgering_UpdateLedgerArgs_ledgering_ShopLedger(args, out)
	out.UpdatedAt = time.Now()
}
