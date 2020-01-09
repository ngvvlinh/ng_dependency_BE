package query

import (
	"context"

	"etop.vn/api/main/moneytx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/capi/dot"
)

func (q *MoneyTxQuery) GetMoneyTxShippingByID(ctx context.Context, args *moneytx.GetMoneyTxByIDQueryArgs) (*moneytx.MoneyTransactionShipping, error) {
	return q.moneyTxShippingStore(ctx).ID(args.MoneyTxShippingID).OptionalShopID(args.ShopID).GetMoneyTxShipping()
}

func (q *MoneyTxQuery) ListMoneyTxShippings(ctx context.Context, args *moneytx.ListMoneyTxShippingArgs) (*moneytx.ListMoneyTxShippingsResponse, error) {
	query := q.moneyTxShippingStore(ctx).Filters(args.Filters).WithPaging(args.Paging)
	query = query.OptionalShopID(args.ShopID).OptionalMoneyTxShippingEtopID(args.MoneyTxShippingEtopID)
	if len(args.MoneyTxShippingIDs) > 0 {
		query = query.IDs(args.MoneyTxShippingIDs...)
	}

	moneyTxs, err := query.ListMoneyTxShippings()
	if err != nil {
		return nil, err
	}
	return &moneytx.ListMoneyTxShippingsResponse{
		MoneyTxShippings: moneyTxs,
		Paging:           query.GetPaging(),
	}, nil
}

func (q *MoneyTxQuery) ListMoneyTxShippingsByMoneyTxShippingExternalID(ctx context.Context, moneyTxShippingExternalID dot.ID) ([]*moneytx.MoneyTransactionShipping, error) {
	if moneyTxShippingExternalID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing money_tx_shipping_external_id")
	}
	return q.moneyTxShippingStore(ctx).MoneyTxShippingExternalID(moneyTxShippingExternalID).ListMoneyTxShippings()
}
