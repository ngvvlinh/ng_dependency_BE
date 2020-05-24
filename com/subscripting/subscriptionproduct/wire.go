// +build wireinject

package subscriptionproduct

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewSubrProductAggregate, SubrProductAggregateMessageBus,
	NewSubrProductQuery, SubrProductQueryMessageBus,
)
