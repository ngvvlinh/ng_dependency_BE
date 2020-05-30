// +build wireinject

package inventory

import (
	"github.com/google/wire"

	"o.o/backend/com/main/inventory/aggregate"
	"o.o/backend/com/main/inventory/query"
)

var WireSet = wire.NewSet(
	aggregate.NewAggregateInventory, aggregate.InventoryAggregateMessageBus,
	query.NewQueryInventory, query.InventoryQueryServiceMessageBus,
)
