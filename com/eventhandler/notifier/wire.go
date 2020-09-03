// +build wireinject

package notifier

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	NewQueryService,
	NewNotifyAggregate,
	QueryServiceNotifyBus,
	NewNotifyAggregateMessageBus,
)
