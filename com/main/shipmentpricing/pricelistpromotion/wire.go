package pricelistpromotion

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	NewAggregate, AggregateMessageBus,
	NewQueryService, QueryServiceMessageBus,
)
