package refund

import (
	"github.com/google/wire"

	"o.o/backend/com/main/refund/aggregate"
	"o.o/backend/com/main/refund/pm"
	"o.o/backend/com/main/refund/query"
)

var WireSet = wire.NewSet(
	pm.New,
	aggregate.NewRefundAggregate, aggregate.RefundAggregateMessageBus,
	query.NewQueryRefund, query.RefundQueryServiceMessageBus,
)
