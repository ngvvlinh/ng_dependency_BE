// +build wireinject

package tradering

import (
	"github.com/google/wire"

	"o.o/backend/com/shopping/tradering/aggregate"
	"o.o/backend/com/shopping/tradering/query"
)

var WireSet = wire.NewSet(
	aggregate.NewTraderAgg, aggregate.TraderAggMessageBus,
	query.NewTraderQuery, query.TraderQueryMessageBus,
)
