package aggregate

import (
	"context"

	"o.o/api/main/transaction"
	"o.o/api/top/types/etc/status3"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/transaction/convert"
	"o.o/backend/com/main/transaction/model"
	"o.o/backend/com/main/transaction/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/capi/dot"
)

var _ transaction.Aggregate = &Aggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type Aggregate struct {
	store sqlstore.TransactionStoreFactory
}

func NewAggregate(db com.MainDB) *Aggregate {
	return &Aggregate{
		store: sqlstore.NewTransactionStore(db),
	}
}

func AggregateMessageBus(a *Aggregate) transaction.CommandBus {
	b := bus.New()
	return transaction.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreateTransaction(ctx context.Context, args *transaction.CreateTransactionArgs) (*transaction.Transaction, error) {
	if err := args.Validate(); err != nil {
		return nil, err
	}
	var trxn transaction.Transaction
	if err := scheme.Convert(args, &trxn); err != nil {
		return nil, err
	}
	return a.store(ctx).CreateTransaction(&trxn)
}

func (a *Aggregate) ConfirmTransaction(ctx context.Context, trxnID dot.ID, accountID dot.ID) (*transaction.Transaction, error) {
	if trxnID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing TransactionID")
	}
	if accountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing AccountID")
	}
	trxn, err := a.store(ctx).ID(trxnID).AccountID(accountID).GetTransaction()
	if err != nil {
		return nil, err
	}
	if !canTransactionChange(trxn) {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Can not confirm this transaction")
	}
	update := &sqlstore.UpdateTransactionStatusArgs{
		ID:        trxnID,
		AccountID: accountID,
		Status:    status3.P,
	}
	return a.store(ctx).UpdateTransactionStatus(update)
}

func canTransactionChange(trxn *transaction.Transaction) bool {
	return trxn.Status == status3.Z
}

func (a *Aggregate) CancelTransaction(ctx context.Context, trxnID dot.ID, accountID dot.ID) (*transaction.Transaction, error) {
	if trxnID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing TransactionID")
	}
	if accountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing AccountID")
	}
	trxn, err := a.store(ctx).ID(trxnID).AccountID(accountID).GetTransaction()
	if err != nil {
		return nil, err
	}

	if !canTransactionChange(trxn) {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Can not cancel this transaction")
	}
	update := &sqlstore.UpdateTransactionStatusArgs{
		ID:        trxnID,
		AccountID: accountID,
		Status:    status3.N,
	}
	return a.store(ctx).UpdateTransactionStatus(update)
}

func (a *Aggregate) DeleteTransaction(ctx context.Context, trxnID, accountID dot.ID) error {
	if trxnID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing TransactionID")
	}
	if accountID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing AccountID")
	}
	return a.store(ctx).ID(trxnID).AccountID(accountID).DeleteTransaction()
}

func (a *Aggregate) ForceCreateTransaction(ctx context.Context, trxn *model.Transaction) error {
	if trxn.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	if trxn.AccountID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing accountID")
	}
	return a.store(ctx).CreateTransactionDB(trxn)
}
