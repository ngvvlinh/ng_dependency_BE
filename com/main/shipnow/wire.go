// +build wireinject

package shipnow

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewAggregate, AggregateMessageBus,
	NewQueryService, QueryServiceMessageBus,
)
