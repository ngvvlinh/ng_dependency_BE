package accountshipnow

import (
	"github.com/google/wire"
	"o.o/backend/com/main/accountshipnow/aggregate"
	"o.o/backend/com/main/accountshipnow/query"
)

var WireSet = wire.NewSet(
	aggregate.NewAggregate, aggregate.AggregateMessageBus,
	query.NewQueryService, query.QueryServiceMessageBus,
)
