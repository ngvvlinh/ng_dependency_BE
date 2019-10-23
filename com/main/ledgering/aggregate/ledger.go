package aggregate

import (
	"context"

	"etop.vn/api/main/identity"
	"etop.vn/api/main/ledgering"
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
	store sqlstore.LedgerStoreFactory
}

func NewLedgerAggregate(db *cmsql.Database) *LedgerAggregate {
	return &LedgerAggregate{
		store: sqlstore.NewLedgerStore(db),
	}
}

func (a *LedgerAggregate) MessageBus() ledgering.CommandBus {
	b := bus.New()
	return ledgering.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *LedgerAggregate) CreateLedger(
	ctx context.Context, args *ledgering.CreateLedgerArgs,
) (*ledgering.ShopLedger, error) {
	if args.Type != string(ledgering.LedgerTypeCash) && args.Type != string(ledgering.LedgerTypeBank) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "type %v không hợp lệ", args.Type)
	}
	if args.Name == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "name không được bỏ trống")
	}

	if args.Type == string(ledgering.LedgerTypeCash) {
		args.BankAccount = nil
	} else {
		if err := verifyBankAccount(args.BankAccount); err != nil {
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
			Wrap(cm.NotFound, "không tìm thấy sổ quỹ").
			Throw()
	}

	// ignore bankAccount where type equal "cash"
	if shopLedger.Type == string(ledgering.LedgerTypeCash) {
		args.BankAccount = nil
	} else {
		if err := verifyBankAccount(args.BankAccount); err != nil {
			return nil, err
		}
	}

	if err := scheme.Convert(args, shopLedger); err != nil {
		return nil, err
	}

	shopLedgerDB := new(model.ShopLedger)
	if err := scheme.Convert(shopLedger, shopLedgerDB); err != nil {
		return nil, err
	}
	err = a.store(ctx).UpdateLedgerDB(shopLedgerDB)
	return shopLedger, err
}

func verifyBankAccount(bankAccount *identity.BankAccount) error {
	if bankAccount.Name == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "tên ngân hàng không được để trống")
	}
	if bankAccount.AccountName == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "tên chủ tài khoản không được để trống")
	}
	if bankAccount.AccountNumber == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "số tài khoản không được để trống")
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

	deleted, err = a.store(ctx).ID(ID).ShopID(shopID).SoftDelete()
	return
}
