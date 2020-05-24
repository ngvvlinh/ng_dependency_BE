// +build wireinject

package transaction

import (
	"github.com/google/wire"
	"o.o/backend/com/main/transaction/aggregate"
	"o.o/backend/com/main/transaction/query"
)

var WireSet = wire.NewSet(
	aggregate.NewAggregate, aggregate.AggregateMessageBus,
	query.NewQueryService, query.QueryServiceMessageBus,
)
