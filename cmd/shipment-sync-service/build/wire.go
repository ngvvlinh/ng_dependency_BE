// +build wireinject

package build

import (
	"context"

	"github.com/google/wire"

	"o.o/backend/cmd/shipment-sync-service/config"
	_base "o.o/backend/cogs/base"
	shipment_all "o.o/backend/cogs/shipment/_all"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/address"
	"o.o/backend/com/main/connectioning"
	"o.o/backend/com/main/identity"
	"o.o/backend/com/main/location"
	"o.o/backend/com/main/ordering"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi"
)

func Build(
	ctx context.Context,
	cfg config.Config,
) (Output, func(), error) {
	panic(wire.Build(
		wire.FieldsOf(&cfg,
			"Redis",
			"Databases",
			"shipment",
		),
		wire.Struct(new(Output), "*"),
		_base.WireSet,
		address.WireSet,
		ordering.WireSet,
		shipment_all.WireSet,
		location.WireSet,
		identity.WireSet,
		connectioning.WireSet,

		com.BuildDatabaseMain,
		wire.Bind(new(capi.EventBus), new(bus.Bus)),
		wire.Bind(new(bus.EventRegistry), new(bus.Bus)),
		sqlstore.WireSet,

		BuildSyncs,
		BuildServers,
	))
}
