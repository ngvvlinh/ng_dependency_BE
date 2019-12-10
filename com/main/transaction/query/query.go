package query

import (
	"context"

	"etop.vn/api/main/transaction"
	"etop.vn/backend/com/main/transaction/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/capi/dot"
)

var _ transaction.QueryService = &QueryService{}

type QueryService struct {
	store sqlstore.TransactionStoreFactory
}

func NewQueryService(db *cmsql.Database) *QueryService {
	return &QueryService{
		store: sqlstore.NewTransactionStore(db),
	}
}

func (q *QueryService) MessageBus() transaction.QueryBus {
	b := bus.New()
	return transaction.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q QueryService) GetTransactionByID(ctx context.Context, tranID dot.ID, userID dot.ID) (*transaction.Transaction, error) {
	return q.store(ctx).ID(tranID).AccountID(userID).GetTransaction()
}

func (q QueryService) ListTransactions(ctx context.Context, args *transaction.GetTransactionsArgs) (*transaction.TransactionResponse, error) {
	query := q.store(ctx).AccountID(args.AccountID)
	transactions, err := query.Paging(args.Paging).ListTransactions()
	if err != nil {
		return nil, err
	}
	count, err := query.Count()
	if err != nil {
		return nil, err
	}
	return &transaction.TransactionResponse{
		Count:        count,
		Paging:       query.GetPaging(),
		Transactions: transactions,
	}, nil
}

func (q QueryService) GetBalance(ctx context.Context, args *transaction.GetBalanceArgs) (int, error) {
	if args.AccountID == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing AccountID")
	}
	return q.store(ctx).AccountID(args.AccountID).OptionalTransactionType(args.TransactionType).ByConfirmedTransaction().GetBalance()
}
