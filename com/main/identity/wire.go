// +build wireinject

package identity

import (
	"github.com/google/wire"
	"o.o/backend/com/main/identity/pm"
	"o.o/backend/com/main/identity/pm_etelecom"
)

var WireSet = wire.NewSet(
	pm.New,
	pm_etelecom.New,
	NewAggregate, AggregateMessageBus,
	NewQueryService, QueryServiceMessageBus,
)
