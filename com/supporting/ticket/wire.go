package ticket

import (
	"github.com/google/wire"
	"o.o/backend/com/supporting/ticket/aggregate"
	"o.o/backend/com/supporting/ticket/query"
	"o.o/backend/com/supporting/ticket/webhook"
)

var WireSet = wire.NewSet(
	aggregate.NewTicketAggregate, aggregate.TicketAggregateMessageBus,
	query.NewTicketQuery, query.TicketQueryMessageBus, webhook.New,
)
