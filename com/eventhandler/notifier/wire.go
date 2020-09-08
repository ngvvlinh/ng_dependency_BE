// +build wireinject

package notifier

import (
	"github.com/google/wire"

	"o.o/backend/com/eventhandler/notifier/sqlstore"
)

var WireSet = wire.NewSet(
	NewOneSignalNotifier,
	NewQueryService,
	NewNotifyAggregate,
	QueryServiceNotifyBus,
	NewNotifyAggregateMessageBus,

	sqlstore.NewDeviceStore,
	sqlstore.NewNotificationStore,
)
