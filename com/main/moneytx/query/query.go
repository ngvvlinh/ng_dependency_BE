package query

import (
	"o.o/api/main/moneytx"
	"o.o/api/main/shipping"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/moneytx/sqlstore"
	"o.o/backend/pkg/common/bus"
)

var _ moneytx.QueryService = &MoneyTxQuery{}

type MoneyTxQuery struct {
	moneyTxShippingStore         sqlstore.MoneyTxShippingStoreFactory
	moneyTxShippingExternalStore sqlstore.MoneyTxShippingExternalStoreFactory
	moneyTxShippingEtopStore     sqlstore.MoneyTxShippingEtopStoreFactory
	shippingQuery                shipping.QueryBus
}

func NewMoneyTxQuery(db com.MainDB, shippingQuery shipping.QueryBus) *MoneyTxQuery {
	return &MoneyTxQuery{
		moneyTxShippingStore:         sqlstore.NewMoneyTxShippingStore(db),
		moneyTxShippingExternalStore: sqlstore.NewMoneyTxShippingExternalStore(db),
		moneyTxShippingEtopStore:     sqlstore.NewMoneyTxShippingEtopStore(db),
		shippingQuery:                shippingQuery,
	}
}

func MoneyTxQueryMessageBus(q *MoneyTxQuery) moneytx.QueryBus {
	b := bus.New()
	return moneytx.NewQueryServiceHandler(q).RegisterHandlers(b)
}
