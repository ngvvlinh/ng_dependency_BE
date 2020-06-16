package pricelist

import (
	"github.com/google/wire"

	"o.o/backend/com/main/shipmentpricing/pricelist/pm"
)

var WireSet = wire.NewSet(
	pm.New,
	NewAggregate, AggregateMessageBus,
	NewQueryService, QueryServiceMessageBus,
)
