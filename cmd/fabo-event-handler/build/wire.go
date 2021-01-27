// +build wireinject

package build

import (
	"context"

	"github.com/google/wire"

	"o.o/backend/cmd/fabo-event-handler/config"
	_base "o.o/backend/cogs/base"
	database_min "o.o/backend/cogs/database/_min"
	shipment_fabo "o.o/backend/cogs/shipment/_fabo"
	handlerapi "o.o/backend/com/eventhandler/handler/api"
	"o.o/backend/com/eventhandler/webhook/sender"
	"o.o/backend/com/eventhandler/webhook/storage"
	comfabo "o.o/backend/com/fabo"
	"o.o/backend/com/fabo/main/fbmessaging"
	"o.o/backend/com/fabo/main/fbpage"
	"o.o/backend/com/fabo/main/fbuser"
	"o.o/backend/com/main/connectioning"
	"o.o/backend/com/main/identity"
	"o.o/backend/com/main/location"
	"o.o/backend/com/main/ordering"
	"o.o/backend/com/shopping/customering"
	"o.o/backend/com/shopping/setting"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi"
)

func Build(ctx context.Context, cfg config.Config) (Output, func(), error) {
	panic(wire.Build(
		wire.FieldsOf(&cfg,
			"Redis",
			"Databases",
			"OneSignal",
			"FacebookApp",
			"Webhook",
			"kafka",
		),
		wire.Struct(new(Output), "*"),
		_base.WireSet,
		database_min.WireSet,
		handlerapi.WireSet,
		sender.WireSet,
		storage.WireSet,
		customering.WireSet,
		sqlstore.WireSet,
		fbuser.WireSet,
		fbmessaging.WireSet,
		fbpage.WireSet,
		identity.WireSet,
		ordering.WireSet,
		location.WireSet,
		connectioning.WireSet,
		shipment_fabo.WireSet,
		comfabo.WireSet,
		setting.WireSet,

		wire.Bind(new(bus.EventRegistry), new(bus.Bus)),
		wire.Bind(new(capi.EventBus), new(bus.Bus)),

		BuildProducer,
		BuildPgEventService,
		BuildIntHandler,
		BuildWebhookHandler,
		BuildWaiters,
		BuildOneSignal,
		BuildHandlers,
		BuildServers,
		BuildMainServer,
	))
}
