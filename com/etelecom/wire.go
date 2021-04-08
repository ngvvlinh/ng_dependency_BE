package etelecom

import (
	"github.com/google/wire"
	"o.o/backend/com/etelecom/aggregate"
	etelecompm "o.o/backend/com/etelecom/pm"
	portsippm "o.o/backend/com/etelecom/pm/portsip_pm"
	"o.o/backend/com/etelecom/query"
)

var WireSet = wire.NewSet(
	etelecompm.New,
	aggregate.NewEtelecomAggregate, aggregate.AggregateMessageBus,
	query.NewQueryService, query.QueryServiceMessageBus,
	portsippm.New,
)
