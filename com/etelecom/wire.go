package etelecom

import (
	"github.com/google/wire"
	"o.o/backend/com/etelecom/aggregate"
	etelecompm "o.o/backend/com/etelecom/pm"
	"o.o/backend/com/etelecom/query"
)

var WireSet = wire.NewSet(
	etelecompm.New,
	aggregate.NewEtelecomAggregate, aggregate.AggregateMessageBus,
	query.NewQueryService, query.QueryServiceMessageBus,
)
