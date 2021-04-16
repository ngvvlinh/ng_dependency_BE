package admin

import (
	"context"

	"o.o/api/main/transaction"
	api "o.o/api/top/int/admin"
	"o.o/api/top/int/types"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type TransactionService struct {
	session.Session
	TransactionQuery transaction.QueryBus
	TransactionAggr  transaction.CommandBus
}

func (s *TransactionService) Clone() api.TransactionService { res := *s; return &res }

func (s *TransactionService) GetTransactions(ctx context.Context, r *types.GetAdminTransactionsRequest) (*types.GetTransactionsResponse, error) {
	paging, err := cmapi.CMCursorPaging(r.Paging)
	if err != nil {
		return nil, err
	}

	query := &transaction.ListTransactionsQuery{
		Paging: *paging,
	}

	if r.Filter != nil {
		query.RefID = r.Filter.RefID
		query.RefType = r.Filter.RefType
		query.DateTo = r.Filter.DateTo
		query.DateFrom = r.Filter.DateFrom
		query.AccountID = r.Filter.AccountID
	}

	if err := s.TransactionQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	res := query.Result

	apiTransactions := convertpb.Convert_core_Transactions_To_api_Transactions(res.Transactions)
	result := &types.GetTransactionsResponse{
		Transactions: apiTransactions,
		Paging:       cmapi.PbCursorPageInfo(paging, &query.Result.Paging),
	}
	return result, nil
}
