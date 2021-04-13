package transaction

import (
	"context"

	transactioning "o.o/api/main/transaction"
	api "o.o/api/top/int/shop"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type TransactionService struct {
	session.Session
	TransactionAggr  transactioning.CommandBus
	TransactionQuery transactioning.QueryBus
}

func (s *TransactionService) GetTransaction(ctx context.Context, request *api.GetTransactionRequest) (*api.Transaction, error) {
	query := &transactioning.GetTransactionByIDQuery{
		TrxnID:    request.ID,
		AccountID: s.SS.Shop().ID,
	}
	if err := s.TransactionQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := convertpb.Convert_core_Transaction_To_api_Transaction(query.Result)
	return result, nil
}

func (s *TransactionService) Clone() api.TransactionService { res := *s; return &res }

func (s *TransactionService) GetTransactions(ctx context.Context, r *api.GetTransactionsRequest) (*api.GetTransactionsResponse, error) {
	paging, err := cmapi.CMCursorPaging(r.Paging)
	if err != nil {
		return nil, err
	}

	query := &transactioning.ListTransactionsQuery{
		AccountID: s.SS.Shop().ID,
		Paging:    *paging,
	}
	if r.Filter != nil {
		query.RefID = r.Filter.RefID
		query.RefType = r.Filter.RefType
		query.DateTo = r.Filter.DateTo
		query.DateFrom = r.Filter.DateFrom
	}

	if err := s.TransactionQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	res := query.Result

	apiTransactions := convertpb.Convert_core_Transactions_To_api_Transactions(res.Transactions)
	result := &api.GetTransactionsResponse{
		Transactions: apiTransactions,
		Paging:       cmapi.PbCursorPageInfo(paging, &query.Result.Paging),
	}
	return result, nil
}
