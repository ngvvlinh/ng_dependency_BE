package inventory

import (
	"github.com/google/wire"

	"o.o/backend/com/main/inventory/aggregate"
	"o.o/backend/com/main/inventory/pm"
	"o.o/backend/com/main/inventory/query"
)

var WireSet = wire.NewSet(
	pm.New,
	aggregate.NewAggregateInventory, aggregate.InventoryAggregateMessageBus,
	query.NewQueryInventory, query.InventoryQueryServiceMessageBus,
)
