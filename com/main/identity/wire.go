// +build wireinject

package identity

import (
	"github.com/google/wire"

	"o.o/backend/com/main/identity/pm"
)

var WireSet = wire.NewSet(
	pm.New,
	NewAggregate, AggregateMessageBus,
	NewQueryService, QueryServiceMessageBus,
)
