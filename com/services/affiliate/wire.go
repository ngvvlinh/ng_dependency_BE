// +build wireinject

package affiliate

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	NewAggregate, AggregateMessageBus,
	NewQuery, QueryServiceMessageBus,
)
