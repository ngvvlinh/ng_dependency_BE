package authorization

import (
	"github.com/google/wire"

	"o.o/backend/com/main/authorization/aggregate"
	"o.o/backend/com/main/authorization/query"
)

var WireSet = wire.NewSet(
	wire.Struct(new(aggregate.AuthorizationAggregate), "*"),
	wire.Struct(new(query.AuthorizationQuery), "*"),
	aggregate.AuthorizationAggregateMessageBus,
	query.AuthorizationQueryMessageBus,
)
