// +build wireinject

package build

import (
	"context"

	"github.com/google/wire"

	"o.o/backend/cmd/fabo-sync-service/config"
	_base "o.o/backend/cogs/base"
	fabo_min"o.o/backend/com/fabo/cogs/_min"
	"o.o/backend/com/fabo/main/fbmessaging"
	"o.o/backend/com/fabo/main/fbpage"
	"o.o/backend/com/fabo/main/fbuser"
	com "o.o/backend/com/main"
	"o.o/backend/com/shopping/customering"
	"o.o/backend/pkg/common/bus"
	"o.o/capi"
)

func Build(ctx context.Context, cfg config.Config) (Output, func(), error) {
	panic(wire.Build(
		wire.FieldsOf(&cfg,
			"Redis",
			"Databases",
			"SyncConfig",
			"FacebookApp",
		),
		wire.Struct(new(Output), "*"),
		_base.WireSet,

		fbuser.WireSet,
		fbmessaging.WireSet,
		fbpage.WireSet,
		fabo_min.WireSet,
		customering.WireSet,

		com.BuildDatabaseMain,
		wire.Bind(new(capi.EventBus), new(bus.Bus)),
		wire.Bind(new(bus.EventRegistry), new(bus.Bus)),

		BuildServers,
	))
}
