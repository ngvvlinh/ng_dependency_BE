// +build wireinject

package shipmentprice

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewAggregate, AggregateMessageBus,
	NewQueryService, QueryServiceMessageBus,
)
