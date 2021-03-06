// +build wireinject

package connectioning

import (
	"github.com/google/wire"
	"o.o/backend/com/main/connectioning/aggregate"
	connectionmanager "o.o/backend/com/main/connectioning/manager"
	connectioningpm "o.o/backend/com/main/connectioning/pm"
	"o.o/backend/com/main/connectioning/query"
)

var WireSet = wire.NewSet(
	connectioningpm.New,
	aggregate.NewConnectionAggregate, aggregate.ConnectionAggregateMessageBus,
	query.NewConnectionQuery, query.ConnectionQueryMessageBus,
	connectionmanager.NewConnectionManager,
)
