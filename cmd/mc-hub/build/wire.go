// +build wireinject

package build

import (
	"context"

	"github.com/google/wire"

	"o.o/backend/cmd/mc-hub/config"
	"o.o/backend/cmd/mc-hub/service"
	_base "o.o/backend/cogs/base"
)

func Build(
	ctx context.Context,
	cfg config.Config,
) (Output, func(), error) {
	panic(wire.Build(
		wire.FieldsOf(&cfg, "redis"),
		wire.Struct(new(Output), "*"),

		_base.WireSet,
		service.WireSet,

		BuildServers,
		BuildMainServer,
		BuildMCShipnow,
	))
}
