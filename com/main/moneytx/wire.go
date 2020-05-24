// +build wireinject

package moneytx

import (
	"github.com/google/wire"
	"o.o/backend/com/main/moneytx/aggregate"
	"o.o/backend/com/main/moneytx/query"
)

var WireSet = wire.NewSet(
	aggregate.NewMoneyTxAggregate, aggregate.MoneyTxAggregateMessageBus,
	query.NewMoneyTxQuery, query.MoneyTxQueryMessageBus,
)
