package tradering

import (
	"github.com/google/wire"

	"o.o/backend/com/shopping/tradering/aggregate"
	"o.o/backend/com/shopping/tradering/pm"
	"o.o/backend/com/shopping/tradering/query"
)

var WireSet = wire.NewSet(
	pm.New,
	aggregate.NewTraderAgg, aggregate.TraderAggMessageBus,
	query.NewTraderQuery, query.TraderQueryMessageBus,
)
