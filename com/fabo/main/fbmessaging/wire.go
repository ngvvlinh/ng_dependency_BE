// +build wireinject

package fbmessaging

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewFbExternalMessagingAggregate, FbExternalMessagingAggregateMessageBus,
	NewFbMessagingQuery, FbMessagingQueryMessageBus,
)
