package etelecom

import (
	"github.com/google/wire"
	"o.o/backend/com/etelecom/aggregate"
	"o.o/backend/com/etelecom/query"
)

var WireSet = wire.NewSet(
	aggregate.NewEtelecomAggregate, aggregate.AggregateMessageBus,
	query.NewQueryService, query.QueryServiceMessageBus,
)
