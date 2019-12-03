package convert

import (
	"time"

	"etop.vn/api/main/ledgering"
	cm "etop.vn/backend/pkg/common"
)

// +gen:convert: etop.vn/backend/com/main/ledgering/model -> etop.vn/api/main/ledgering
// +gen:convert: etop.vn/api/main/ledgering

func createShopLedger(args *ledgering.CreateLedgerArgs, out *ledgering.ShopLedger) {
	apply_ledgering_CreateLedgerArgs_ledgering_ShopLedger(args, out)
	out.ID = cm.NewID()
	out.BankAccount = args.BankAccount
}

func updateShopLedger(args *ledgering.UpdateLedgerArgs, out *ledgering.ShopLedger) {
	apply_ledgering_UpdateLedgerArgs_ledgering_ShopLedger(args, out)
	out.BankAccount = args.BankAccount
	out.UpdatedAt = time.Now()
}
