package shipnow

import (
	"github.com/google/wire"

	"o.o/backend/com/main/shipnow/pm"
)

var WireSet = wire.NewSet(
	pm.New,
	NewAggregate, AggregateMessageBus,
	NewQueryService, QueryServiceMessageBus,
)
