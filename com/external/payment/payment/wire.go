// +build wireinject

package payment

import (
	"github.com/google/wire"
	"o.o/backend/com/external/payment/payment/aggregate"
	"o.o/backend/com/external/payment/payment/query"
)

var WireSet = wire.NewSet(
	aggregate.NewAggregate, aggregate.AggregateMessageBus,
	query.NewQueryService, query.QueryServiceMessageBus,
)
