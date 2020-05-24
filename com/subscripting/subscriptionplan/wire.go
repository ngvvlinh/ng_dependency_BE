// +build wireinject

package subscriptionplan

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewSubrPlanAggregate, SubrPlanAggregateMessageBus,
	NewSubrPlanQuery, SubrPlanQueryMessageBus,
)
