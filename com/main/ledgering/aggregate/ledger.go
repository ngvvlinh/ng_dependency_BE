package aggregate

import (
	"context"

	"etop.vn/api/main/receipting"

	"etop.vn/api/main/identity"
	"etop.vn/api/main/ledgering"
	identityconvert "etop.vn/backend/com/main/identity/convert"
	"etop.vn/backend/com/main/ledgering/convert"
	"etop.vn/backend/com/main/ledgering/model"
	"etop.vn/backend/com/main/ledgering/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/conversion"
)

var _ ledgering.Aggregate = &LedgerAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type LedgerAggregate struct {
	store        sqlstore.LedgerStoreFactory
	receiptQuery *receipting.QueryBus
}

func NewLedgerAggregate(db *cmsql.Database, receiptQuery *receipting.QueryBus) *LedgerAggregate {
	return &LedgerAggregate{
		store:        sqlstore.NewLedgerStore(db),
		receiptQuery: receiptQuery,
	}
}

func (a *LedgerAggregate) MessageBus() ledgering.CommandBus {
	b := bus.New()
	return ledgering.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *LedgerAggregate) CreateLedger(
	ctx context.Context, args *ledgering.CreateLedgerArgs,
) (*ledgering.ShopLedger, error) {
	if args.Type != ledgering.LedgerTypeCash && args.Type != ledgering.LedgerTypeBank {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "type %v không hợp lệ", args.Type)
	}
	if args.Name == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "name không được bỏ trống")
	}

	if args.Type == ledgering.LedgerTypeCash {
		args.BankAccount = nil
	} else {
		if args.BankAccount == nil {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Tài khoản thanh toán chuyển khoản không được bỏ trống thông tin tài khoản")
		}
		if err := a.verifyBankAccount(ctx, 0, args.ShopID, args.BankAccount); err != nil {
			return nil, err
		}
	}

	shopLedger := new(ledgering.ShopLedger)
	if err := scheme.Convert(args, shopLedger); err != nil {
		return nil, err
	}
	err := a.store(ctx).CreateLedger(shopLedger)
	return shopLedger, err
}

func (a *LedgerAggregate) UpdateLedger(
	ctx context.Context, args *ledgering.UpdateLedgerArgs,
) (*ledgering.ShopLedger, error) {
	shopLedger, err := a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetLedger()
	if err != nil {
		return nil, cm.MapError(err).
			Wrap(cm.NotFound, "không tìm thấy tài khoản thanh toán").
			Throw()
	}

	// ignore bankAccount where type equal "cash"
	if shopLedger.Type == string(ledgering.LedgerTypeCash) {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "bạn không được sửa tài khoản thanh toán mặc định")
	} else {
		if err := a.verifyBankAccount(ctx, args.ID, args.ShopID, args.BankAccount); err != nil {
			return nil, err
		}
		if args.BankAccount == nil {
			args.BankAccount = shopLedger.BankAccount
		}
	}

	if err := scheme.Convert(args, shopLedger); err != nil {
		return nil, err
	}

	shopLedgerDB := new(model.ShopLedger)
	if err := scheme.Convert(shopLedger, shopLedgerDB); err != nil {
		return nil, err
	}
	shopLedgerDB.BankAccount = identityconvert.BankAccountDB(shopLedger.BankAccount)

	err = a.store(ctx).UpdateLedgerDB(shopLedgerDB)
	return shopLedger, err
}

func (a *LedgerAggregate) verifyBankAccount(ctx context.Context, ledgerID, shopID int64, bankAccount *identity.BankAccount) error {
	if bankAccount == nil {
		return nil
	}
	if bankAccount.Name == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "tên ngân hàng không được để trống")
	}
	if bankAccount.AccountName == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "tên chủ tài khoản không được để trống")
	}
	if bankAccount.AccountNumber == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "số tài khoản không được để trống")
	}
	ledger, err := a.store(ctx).ShopID(shopID).AccountNumber(bankAccount.AccountNumber).GetLedger()
	switch cm.ErrorCode(err) {
	case cm.NoError:
		if ledger.ID != ledgerID {
			return cm.Errorf(cm.FailedPrecondition, nil, "Số tài khoản này đã tồn tại trong shop")
		}
	case cm.NotFound:
	// no-op
	default:
		return err
	}

	return nil
}

func (a *LedgerAggregate) DeleteLedger(
	ctx context.Context, ID, shopID int64,
) (deleted int, err error) {
	shopLedger, err := a.store(ctx).ID(ID).ShopID(shopID).GetLedger()
	if err != nil {
		return 0, cm.MapError(err).
			Wrap(cm.NotFound, "không tìm thấy số quỹ").
			Throw()
	}
	if shopLedger.Type == string(ledgering.LedgerTypeCash) {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "không thể xóa số quỹ mặc định")
	}

	// Check ledger_id exists into any receipts
	query := &receipting.ListReceiptsByLedgerIDsQuery{
		ShopID:    shopID,
		LedgerIDs: []int64{ID},
	}
	if err := a.receiptQuery.Dispatch(ctx, query); err != nil {
		return 0, err
	}
	if len(query.Result.Receipts) != 0 {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Tài khoản thanh toán đã được được sử dụng, không thể xóa")
	}

	deleted, err = a.store(ctx).ID(ID).ShopID(shopID).SoftDelete()
	return
}
