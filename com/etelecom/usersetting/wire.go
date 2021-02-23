package usersetting

import (
	"github.com/google/wire"
	"o.o/backend/com/etelecom/usersetting/pm"
)

var WireSet = wire.NewSet(
	pm.New,
	NewUserSettingAggregate, AggerateMessageBus,
	NewQueryService, QueryServiceMessageBus,
)
