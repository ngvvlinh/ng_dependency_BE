// +build wireinject

package purchaserefund

import (
	"github.com/google/wire"
	"o.o/backend/com/main/purchaseorder/query"
	"o.o/backend/com/main/purchaserefund/aggregate"
)

var WireSet = wire.NewSet(
	aggregate.NewPurchaseRefundAggregate, aggregate.PurchaseRefundAggregateMessageBus,
	query.NewPurchaseOrderQuery, query.PurchaseOrderQueryMessageBus,
)
