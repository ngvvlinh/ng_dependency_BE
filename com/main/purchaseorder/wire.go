package purchaseorder

import (
	"github.com/google/wire"

	"o.o/backend/com/main/purchaseorder/aggregate"
	"o.o/backend/com/main/purchaseorder/pm"
	"o.o/backend/com/main/purchaseorder/query"
)

var WireSet = wire.NewSet(
	pm.New,
	aggregate.NewPurchaseOrderAggregate, aggregate.PurchaseOrderAggregateMessageBus,
	query.NewPurchaseOrderQuery, query.PurchaseOrderQueryMessageBus,
)
