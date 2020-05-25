package shop

import (
	"context"

	apitypes "o.o/api/top/int/types"
	moneymodelx "o.o/backend/com/main/moneytx/modelx"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api/convertpb"
)

type MoneyTransactionService struct{}

func (s *MoneyTransactionService) Clone() *MoneyTransactionService { res := *s; return &res }

func (s *MoneyTransactionService) GetMoneyTransaction(ctx context.Context, q *GetMoneyTransactionEndpoint) error {
	query := &moneymodelx.GetMoneyTransaction{
		ShopID: q.Context.Shop.ID,
		ID:     q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTransactionExtended(query.Result)
	return nil
}

func (s *MoneyTransactionService) GetMoneyTransactions(ctx context.Context, q *GetMoneyTransactionsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &moneymodelx.GetMoneyTransactions{
		ShopID:  q.Context.Shop.ID,
		Paging:  paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &apitypes.MoneyTransactionsResponse{
		MoneyTransactions: convertpb.PbMoneyTransactionExtendeds(query.Result.MoneyTransactions),
		Paging:            cmapi.PbPageInfo(paging),
	}
	return nil
}
