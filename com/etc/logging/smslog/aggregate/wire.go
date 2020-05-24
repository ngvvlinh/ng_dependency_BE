// +build wireinject

package aggregate

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewSmsLogAggregate, SmsLogAggregateMessageBus,
)
