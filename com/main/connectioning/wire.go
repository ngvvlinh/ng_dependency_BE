// +build wireinject

package connectioning

import (
	"github.com/google/wire"
	"o.o/backend/com/main/connectioning/aggregate"
	"o.o/backend/com/main/connectioning/query"
)

var WireSet = wire.NewSet(
	aggregate.NewConnectionAggregate, aggregate.ConnectionAggregateMessageBus,
	query.NewConnectionQuery, query.ConnectionQueryMessageBus,
)
