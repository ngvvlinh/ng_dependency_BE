package ticket

import (
	"github.com/google/wire"
	"o.o/backend/com/supporting/ticket/aggregate"
	"o.o/backend/com/supporting/ticket/pm"
	"o.o/backend/com/supporting/ticket/query"
)

var WireSet = wire.NewSet(
	aggregate.NewTicketAggregate, aggregate.TicketAggregateMessageBus,
	query.NewTicketQuery, query.TicketQueryMessageBus,
	pm.NewProcessManager,
)
