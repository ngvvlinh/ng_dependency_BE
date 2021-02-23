package subscription

import (
	"github.com/google/wire"
	subscriptionpm "o.o/backend/com/subscripting/subscription/pm"
)

var WireSet = wire.NewSet(
	NewSubscriptionAggregate, SubscriptionAggregateMessageBus,
	NewSubscriptionQuery, SubscriptionQueryMessageBus,
	subscriptionpm.New,
)
