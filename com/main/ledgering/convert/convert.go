package convert

import (
	"time"

	"etop.vn/api/main/identity"
	"etop.vn/api/main/ledgering"
	"etop.vn/backend/com/main/ledgering/model"
	cm "etop.vn/backend/pkg/common"
	etopmodel "etop.vn/backend/pkg/etop/model"
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

func shopLedger(args *model.ShopLedger, out *ledgering.ShopLedger) {
	convert_ledgeringmodel_ShopLedger_ledgering_ShopLedger(args, out)
	if args.BankAccount != nil {
		out.BankAccount = &identity.BankAccount{
			Name:          args.BankAccount.Name,
			Province:      args.BankAccount.Province,
			Branch:        args.BankAccount.Branch,
			AccountNumber: args.BankAccount.AccountNumber,
			AccountName:   args.BankAccount.AccountName,
		}
	}
}

func shopLedgerDB(args *ledgering.ShopLedger, out *model.ShopLedger) {
	convert_ledgering_ShopLedger_ledgeringmodel_ShopLedger(args, out)
	if args.BankAccount != nil {
		out.BankAccount = &etopmodel.BankAccount{
			Name:          args.BankAccount.Name,
			Province:      args.BankAccount.Province,
			Branch:        args.BankAccount.Branch,
			AccountNumber: args.BankAccount.AccountNumber,
			AccountName:   args.BankAccount.AccountName,
		}
	}
}
