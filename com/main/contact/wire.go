package contact

import (
	"github.com/google/wire"
	"o.o/backend/com/main/contact/aggregate"
	"o.o/backend/com/main/contact/query"
)

var WireSet = wire.NewSet(
	aggregate.NewContactAggregate, aggregate.ContactAggregateMessageBus,
	query.NewContactQuery, query.ContactQueryMessageBus,
)
