package moneytx

import (
	"github.com/google/wire"

	"o.o/backend/com/main/moneytx/aggregate"
	"o.o/backend/com/main/moneytx/pm"
	"o.o/backend/com/main/moneytx/query"
)

var WireSet = wire.NewSet(
	pm.New,
	aggregate.NewMoneyTxAggregate, aggregate.MoneyTxAggregateMessageBus,
	query.NewMoneyTxQuery, query.MoneyTxQueryMessageBus,
)
