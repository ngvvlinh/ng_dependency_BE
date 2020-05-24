// +build wireinject

package subscriptionbill

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewSubrBillAggregate, SubrBillAggregateMessageBus,
	NewSubrBillQuery, SubrBillQueryMessageBus,
)
