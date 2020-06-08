package affiliate

import (
	"github.com/google/wire"

	"o.o/backend/com/services/affiliate/pm"
)

var WireSet = wire.NewSet(
	pm.New,
	NewAggregate, AggregateMessageBus,
	NewQuery, QueryServiceMessageBus,
)
