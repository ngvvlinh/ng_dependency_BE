package fbuser

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewFbUserAggregate, FbUserAggregateMessageBus,
	NewFbUserQuery, FbUserQueryMessageBus,
)
