// +build wireinject

package aggregatex

import (
	"github.com/google/wire"

	"o.o/backend/com/main/inventory/pm"
	"o.o/backend/com/main/inventory/query"
)

var WireSet = wire.NewSet(
	pm.New,
	NewAggregateInventory, InventoryAggregateMessageBus,
	query.NewQueryInventory, query.InventoryQueryServiceMessageBus,
)
