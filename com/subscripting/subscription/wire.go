// +build wireinject

package subscription

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewSubscriptionAggregate, SubscriptionAggregateMessageBus,
	NewSubscriptionQuery, SubscriptionQueryMessageBus,
)
