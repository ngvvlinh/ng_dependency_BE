// +build wireinject

package ordering

import (
	"github.com/google/wire"

	"o.o/backend/com/main/ordering/pm"
)

var WireSet = wire.NewSet(
	pm.New,
	NewAggregate, AggregateMessageBus,
	NewQueryService, QueryServiceMessageBus,
)
