package query

import (
	"context"

	"etop.vn/api/main/moneytx"
	"etop.vn/capi/dot"
)

func (q *MoneyTxQuery) GetMoneyTxShippingExternal(ctx context.Context, id dot.ID) (*moneytx.MoneyTransactionShippingExternalFtLine, error) {
	return q.moneyTxShippingExternalStore(ctx).ID(id).GetMoneyTxShippingExternalFtLine()
}

func (q *MoneyTxQuery) ListMoneyTxShippingExternals(ctx context.Context, args *moneytx.ListMoneyTxShippingExternalsArgs) (*moneytx.ListMoneyTxShippingExternalsResponse, error) {
	query := q.moneyTxShippingExternalStore(ctx).Filters(args.Filters)
	if len(args.MoneyTxShippingExternalIDs) > 0 {
		query = query.IDs(args.MoneyTxShippingExternalIDs...)
	}

	moneyTxs, err := query.WithPaging(args.Paging).ListMoneyTxShippingExternalsFtLine()
	if err != nil {
		return nil, err
	}
	return &moneytx.ListMoneyTxShippingExternalsResponse{
		MoneyTxShippingExternals: moneyTxs,
		Paging:                   query.GetPaging(),
	}, nil
}
