package shipping

import (
	"github.com/google/wire"

	"o.o/backend/com/main/shipping/aggregate"
	"o.o/backend/com/main/shipping/pm"
	"o.o/backend/com/main/shipping/query"
)

var WireSet = wire.NewSet(
	pm.New,
	aggregate.NewAggregate, aggregate.AggregateMessageBus,
	query.NewQueryService, query.QueryServiceMessageBus,
)
