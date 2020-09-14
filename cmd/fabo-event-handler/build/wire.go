// +build wireinject

package build

import (
	"context"

	"github.com/google/wire"

	"o.o/backend/cmd/fabo-event-handler/config"
	_base "o.o/backend/cogs/base"
	handlerapi "o.o/backend/com/eventhandler/handler/api"
	"o.o/backend/com/eventhandler/webhook/sender"
	"o.o/backend/com/eventhandler/webhook/storage"
	"o.o/backend/com/fabo/main/fbmessaging"
	"o.o/backend/com/fabo/main/fbpage"
	"o.o/backend/com/fabo/main/fbuser"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/identity"
	"o.o/backend/com/shopping/customering"
	"o.o/backend/pkg/etop/sqlstore"
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
		customering.WireSet,
		sqlstore.WireSet,
		fbuser.WireSet,
		fbmessaging.WireSet,
		fbpage.WireSet,
		identity.WireSet,

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
