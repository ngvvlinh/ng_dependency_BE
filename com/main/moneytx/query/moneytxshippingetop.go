package query

import (
	"context"

	"o.o/api/main/moneytx"
	"o.o/capi/dot"
)

func (q *MoneyTxQuery) GetMoneyTxShippingEtop(ctx context.Context, id dot.ID) (*moneytx.MoneyTransactionShippingEtop, error) {
	return q.moneyTxShippingEtopStore(ctx).ID(id).GetMoneyTxShippingEtop()
}

func (q *MoneyTxQuery) ListMoneyTxShippingEtops(ctx context.Context, args *moneytx.ListMoneyTxShippingEtopsArgs) (*moneytx.ListMoneyTxShippingEtopsResponse, error) {
	query := q.moneyTxShippingEtopStore(ctx).Filters(args.Filter).WithPaging(args.Paging)
	if len(args.MoneyTxShippingEtopIDs) > 0 {
		query = query.IDs(args.MoneyTxShippingEtopIDs...)
	}
	if args.Status.Valid {
		query = query.Status(args.Status.Enum)
	}
	moneyTxs, err := query.ListMoneyTxShippingEtops()
	if err != nil {
		return nil, err
	}
	return &moneytx.ListMoneyTxShippingEtopsResponse{
		MoneyTxShippingEtops: moneyTxs,
		Paging:               query.GetPaging(),
	}, nil
}
