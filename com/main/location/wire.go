// +build wireinject

package location

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	NewAggregate, AggregateMessageBus,
	New, QueryMessageBus,
)
