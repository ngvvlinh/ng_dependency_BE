// +build wireinject

package receipting

import (
	"github.com/google/wire"

	"o.o/backend/com/main/receipting/aggregate"
	"o.o/backend/com/main/receipting/query"
)

var WireSet = wire.NewSet(
	aggregate.NewReceiptAggregate, aggregate.ReceiptAggregateMessageBus,
	query.NewReceiptQuery, query.ReceiptQueryMessageBus,
)
