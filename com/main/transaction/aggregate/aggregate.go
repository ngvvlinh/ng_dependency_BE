package aggregate

import (
	"context"

	"o.o/api/main/transaction"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/main/transaction/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

var _ transaction.Aggregate = &Aggregate{}

type Aggregate struct {
	store sqlstore.TransactionStoreFactory
}

func NewAggregate(db *cmsql.Database) *Aggregate {
	return &Aggregate{
		store: sqlstore.NewTransactionStore(db),
	}
}

func AggregateMessageBus(a *Aggregate) transaction.CommandBus {
	b := bus.New()
	return transaction.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreateTransaction(ctx context.Context, args *transaction.CreateTransactionArgs) (*transaction.Transaction, error) {
	trxn := &transaction.Transaction{
		ID:        args.ID,
		Amount:    args.Amount,
		AccountID: args.AccountID,
		Status:    args.Status,
		Type:      args.Type,
		Metadata:  args.Metadata,
	}
	return a.store(ctx).CreateTransaction(trxn)
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
