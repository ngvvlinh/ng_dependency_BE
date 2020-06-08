package purchaserefund

import (
	"github.com/google/wire"

	"o.o/backend/com/main/purchaserefund/aggregate"
	"o.o/backend/com/main/purchaserefund/pm"
	"o.o/backend/com/main/purchaserefund/query"
)

var WireSet = wire.NewSet(
	pm.New,
	aggregate.NewPurchaseRefundAggregate, aggregate.PurchaseRefundAggregateMessageBus,
	query.NewQueryPurchasePurchaseRefund, query.PurchaseRefundQueryServiceMessageBus,
)
