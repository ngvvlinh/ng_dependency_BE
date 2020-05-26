package invitation

import (
	"github.com/google/wire"
	"o.o/backend/com/main/invitation/aggregate"
	"o.o/backend/com/main/invitation/pm"
	"o.o/backend/com/main/invitation/query"
)

var WireSet = wire.NewSet(
	aggregate.NewInvitationAggregate, aggregate.InvitationAggregateMessageBus,
	query.NewInvitationQuery, query.InvitationQueryMessageBus,
	pm.New,
)
