package transaction

import (
	"github.com/google/wire"

	"o.o/backend/com/main/transaction/aggregate"
	transactionpm "o.o/backend/com/main/transaction/pm"
	"o.o/backend/com/main/transaction/query"
)

var WireSet = wire.NewSet(
	transactionpm.New,
	aggregate.NewAggregate, aggregate.AggregateMessageBus,
	query.NewQueryService, query.QueryServiceMessageBus,
)
