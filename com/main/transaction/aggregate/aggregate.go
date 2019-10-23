package aggregate

import (
	"context"

	etoptypes "etop.vn/api/main/etop"
	"etop.vn/api/main/transaction"
	"etop.vn/backend/com/main/transaction/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
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

func (a *Aggregate) MessageBus() transaction.CommandBus {
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

func (a *Aggregate) ConfirmTransaction(ctx context.Context, trxnID int64, accountID int64) (*transaction.Transaction, error) {
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
		Status:    etoptypes.S3Positive,
	}
	return a.store(ctx).UpdateTransactionStatus(update)
}

func canTransactionChange(trxn *transaction.Transaction) bool {
	return trxn.Status == etoptypes.S3Zero
}

func (a *Aggregate) CancelTransaction(ctx context.Context, trxnID int64, accountID int64) (*transaction.Transaction, error) {
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
		Status:    etoptypes.S3Negative,
	}
	return a.store(ctx).UpdateTransactionStatus(update)
}
