// +build wireinject

package build

import (
	"context"

	"github.com/google/wire"

	"o.o/backend/cmd/etop-event-handler/config"
	_base "o.o/backend/cogs/base"
	handlerapi "o.o/backend/com/eventhandler/handler/api"
	"o.o/backend/com/eventhandler/webhook/sender"
	"o.o/backend/com/eventhandler/webhook/storage"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/catalog"
	"o.o/backend/com/main/inventory"
	"o.o/backend/com/main/location"
	"o.o/backend/com/main/shipnow"
	"o.o/backend/com/main/stocktaking"
	"o.o/backend/com/shopping/customering"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi"
)

func Build(ctx context.Context, cfg config.Config) (Output, func(), error) {
	panic(wire.Build(
		wire.FieldsOf(&cfg,
			"Redis",
			"Databases",
		),
		wire.Struct(new(Output), "*"),
		_base.WireSet,
		handlerapi.WireSet,
		sender.WireSet,
		storage.WireSet,
		catalog.WireSet,
		customering.WireSet,
		inventory.WireSet,
		stocktaking.WireSet,
		location.WireSet,
		shipnow.WireSet,
		sqlstore.WireSet,

		wire.Bind(new(capi.EventBus), new(bus.Bus)),

		com.BuildDatabaseWebhook,
		com.BuildDatabaseMain,

		BuildPgEventService,
		BuildIntHandler,
		BuildWebhookHandler,
		BuildWaiters,
		BuildServers,
		BuildMainServer,
	))
}
