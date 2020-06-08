package receipting

import (
	"github.com/google/wire"

	"o.o/backend/com/main/receipting/aggregate"
	"o.o/backend/com/main/receipting/pm"
	"o.o/backend/com/main/receipting/query"
)

var WireSet = wire.NewSet(
	pm.New,
	aggregate.NewReceiptAggregate, aggregate.ReceiptAggregateMessageBus,
	query.NewReceiptQuery, query.ReceiptQueryMessageBus,
)
