package query

import (
	"context"

	"etop.vn/api/main/moneytx"
	"etop.vn/backend/com/main/moneytx/sqlstore"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi/dot"
)

var _ moneytx.QueryService = &MoneyTxQuery{}

type MoneyTxQuery struct {
	moneyTxShippingStore         sqlstore.MoneyTxShippingStoreFactory
	moneyTxShippingExternalStore sqlstore.MoneyTxShippingExternalStoreFactory
}

func NewMoneyTxQuery(db *cmsql.Database) *MoneyTxQuery {
	return &MoneyTxQuery{
		moneyTxShippingStore:         sqlstore.NewMoneyTxShippingStore(db),
		moneyTxShippingExternalStore: sqlstore.NewMoneyTxShippingExternalStore(db),
	}
}

func (q *MoneyTxQuery) MessageBus() moneytx.QueryBus {
	b := bus.New()
	return moneytx.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *MoneyTxQuery) GetMoneyTransactionShippingByID(ctx context.Context, args *moneytx.GetMoneyTxByIDQueryArgs) (*moneytx.MoneyTransactionShipping, error) {
	return q.moneyTxShippingStore(ctx).ID(args.ID).OptionalShopID(args.ShopID).GetMoneyTxShipping()
}

func (q *MoneyTxQuery) GetMoneyTxShippingByID(context.Context, *moneytx.GetMoneyTxByIDQueryArgs) (*moneytx.MoneyTransactionShippingExtended, error) {
	panic("implement me")
}

func (q *MoneyTxQuery) ListMoneyTxShippings(context.Context, *moneytx.ListMoneyTxArgs) (*moneytx.ListMoneyTxShippingsResponse, error) {
	panic("implement me")
}

func (q *MoneyTxQuery) ListMoneyTxShippingsByMoneyTxShippingExternalID(ctx context.Context, MoneyTxShippingExternalID dot.ID) ([]*moneytx.MoneyTransactionShippingExtended, error) {
	panic("implement me")
}

func (q *MoneyTxQuery) GetMoneyTxShippingExternal(ctx context.Context, id dot.ID) (*moneytx.MoneyTransactionShippingExternalExtended, error) {
	return q.moneyTxShippingExternalStore(ctx).ID(id).GetMoneyTxShippingExternalExtended()
}

func (q *MoneyTxQuery) ListMoneyTxShippingExternals(context.Context, *moneytx.ListMoneyTxShippingExternalsArgs) (*moneytx.ListMoneyTxShippingExternalsResponse, error) {
	panic("implement me")
}

func (q *MoneyTxQuery) GetMoneyTxShippingEtop(ctx context.Context, ID dot.ID) (*moneytx.MoneyTransactionShippingEtopExtended, error) {
	panic("implement me")
}

func (q *MoneyTxQuery) ListMoneyTxShippingEtops(context.Context, *moneytx.ListMoneyTxShippingEtopsArgs) (*moneytx.ListMoneyTxShippingEtopsResponse, error) {
	panic("implement me")
}
