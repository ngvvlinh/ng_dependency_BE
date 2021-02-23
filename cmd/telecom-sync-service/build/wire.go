// +build wireinject

package build

import (
	"context"

	"github.com/google/wire"

	"o.o/backend/cmd/telecom-sync-service/config"
	_base "o.o/backend/cogs/base"
	payment_all "o.o/backend/cogs/payment/_all"
	shipment_all "o.o/backend/cogs/shipment/_all"
	telecom_all "o.o/backend/cogs/telecom/_all"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/connectioning"
	"o.o/backend/com/main/contact"
	"o.o/backend/com/main/identity"
	"o.o/backend/com/subscripting"
	"o.o/backend/pkg/common/bus"
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
		),
		wire.Struct(new(Output), "*"),
		_base.WireSet,
		identity.WireSet,
		connectioning.WireSet,
		contact.WireSet,
		payment_all.WireSet,
		subscripting.WireSet,
		telecom_all.WireSet,

		com.BuildDatabaseEtelecomDB,
		com.BuildDatabaseMain,
		wire.Bind(new(capi.EventBus), new(bus.Bus)),
		wire.Bind(new(bus.EventRegistry), new(bus.Bus)),

		shipment_all.SupportedShipmentServices,
		BuildSyncs,
		BuildServers,
	))
}
