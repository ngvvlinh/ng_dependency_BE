// +build wireinject

package carrying

import (
	"github.com/google/wire"
	"o.o/backend/com/shopping/carrying/aggregate"
	"o.o/backend/com/shopping/carrying/query"
)

var WireSet = wire.NewSet(
	aggregate.NewCarrierAggregate, aggregate.CarrierAggregateMessageBus,
	query.NewCarrierQuery, query.CarrierQueryMessageBus,
)
