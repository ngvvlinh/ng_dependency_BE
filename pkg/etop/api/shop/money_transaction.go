package shop

import (
	"context"

	"o.o/api/main/moneytx"
	"o.o/api/top/int/types"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
)

type MoneyTransactionService struct {
	MoneyTxQuery moneytx.QueryBus
}

func (s *MoneyTransactionService) Clone() *MoneyTransactionService { res := *s; return &res }

func (s *MoneyTransactionService) GetMoneyTransaction(ctx context.Context, q *GetMoneyTransactionEndpoint) error {
	query := &moneytx.GetMoneyTxShippingByIDQuery{
		MoneyTxShippingID: q.Id,
		ShopID:            q.Context.Shop.ID,
	}
	if err := s.MoneyTxQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTxShipping(query.Result)
	return nil
}

func (s *MoneyTransactionService) GetMoneyTransactions(ctx context.Context, q *GetMoneyTransactionsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &moneytx.ListMoneyTxShippingsQuery{
		ShopID:  q.Context.Shop.ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := s.MoneyTxQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &types.MoneyTransactionsResponse{
		Paging:            cmapi.PbMetaPageInfo(query.Result.Paging),
		MoneyTransactions: convertpb.PbMoneyTxShippings(query.Result.MoneyTxShippings),
	}
	return nil
}
