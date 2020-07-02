package shop

import (
	"context"

	"o.o/api/main/moneytx"
	api "o.o/api/top/int/shop"
	inttypes "o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type MoneyTransactionService struct {
	session.Session

	MoneyTxQuery moneytx.QueryBus
}

func (s *MoneyTransactionService) Clone() api.MoneyTransactionService { res := *s; return &res }

func (s *MoneyTransactionService) GetMoneyTransaction(ctx context.Context, q *pbcm.IDRequest) (*inttypes.MoneyTransaction, error) {
	query := &moneytx.GetMoneyTxShippingByIDQuery{
		MoneyTxShippingID: q.Id,
		ShopID:            s.SS.Shop().ID,
	}
	if err := s.MoneyTxQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := convertpb.PbMoneyTxShipping(query.Result)
	return result, nil
}

func (s *MoneyTransactionService) GetMoneyTransactions(ctx context.Context, q *api.GetMoneyTransactionsRequest) (*inttypes.MoneyTransactionsResponse, error) {
	paging := cmapi.CMPaging(q.Paging)
	query := &moneytx.ListMoneyTxShippingsQuery{
		ShopID:  s.SS.Shop().ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := s.MoneyTxQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &inttypes.MoneyTransactionsResponse{
		Paging:            cmapi.PbMetaPageInfo(query.Result.Paging),
		MoneyTransactions: convertpb.PbMoneyTxShippings(query.Result.MoneyTxShippings),
	}
	return result, nil
}
