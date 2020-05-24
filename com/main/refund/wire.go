// +build wireinject

package refund

import (
	"github.com/google/wire"
	"o.o/backend/com/main/refund/aggregate"
	"o.o/backend/com/main/refund/query"
)

var WireSet = wire.NewSet(
	aggregate.NewRefundAggregate, aggregate.RefundAggregateMessageBus,
	query.NewQueryRefund, query.RefundQueryServiceMessageBus,
)
