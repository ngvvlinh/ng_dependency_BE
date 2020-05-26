// +build wireinject

package purchaserefund

import (
	"github.com/google/wire"
	"o.o/backend/com/main/purchaserefund/aggregate"
	"o.o/backend/com/main/purchaserefund/query"
)

var WireSet = wire.NewSet(
	aggregate.NewPurchaseRefundAggregate, aggregate.PurchaseRefundAggregateMessageBus,
	query.NewQueryPurchasePurchaseRefund, query.PurchaseRefundQueryServiceMessageBus,
)
