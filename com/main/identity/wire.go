// +build wireinject

package identity

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewAggregate, AggregateMessageBus,
	NewQueryService, QueryServiceMessageBus,
)
