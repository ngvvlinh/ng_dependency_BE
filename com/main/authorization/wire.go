package authorization

import (
	"github.com/google/wire"
	"o.o/backend/com/main/authorization/aggregate"
	"o.o/backend/com/main/authorization/query"
)

var WireSet = wire.NewSet(
	aggregate.NewAuthorizationAggregate, aggregate.AuthorizationAggregateMessageBus,
	query.NewAuthorizationQuery, query.AuthorizationQueryMessageBus,
)
