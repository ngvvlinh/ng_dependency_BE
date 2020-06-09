package fbuser

import (
	"github.com/google/wire"
	"o.o/backend/com/fabo/main/fbuser/pm"
)

var WireSet = wire.NewSet(
	pm.New,
	NewFbUserAggregate, FbUserAggregateMessageBus,
	NewFbUserQuery, FbUserQueryMessageBus,
)
