// +build wireinject

package customering

import (
	"github.com/google/wire"
	"o.o/backend/com/shopping/customering/aggregate"
	"o.o/backend/com/shopping/customering/query"
)

var WireSet = wire.NewSet(
	aggregate.NewCustomerAggregate, aggregate.CustomerAggregateMessageBus,
	query.NewCustomerQuery, query.CustomerQueryMessageBus,
)
