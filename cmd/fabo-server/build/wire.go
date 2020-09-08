// +build wireinject

package build

import (
	"context"

	"github.com/google/wire"

	"o.o/api/services/affiliate"
	"o.o/api/shopping/tradering"
	"o.o/backend/cmd/fabo-server/config"
	_base "o.o/backend/cogs/base"
	config_server "o.o/backend/cogs/config/_server"
	_core "o.o/backend/cogs/core"
	database_min "o.o/backend/cogs/database/_min"
	server_fabo "o.o/backend/cogs/server/fabo"
	server_shop "o.o/backend/cogs/server/shop"
	shipment_fabo "o.o/backend/cogs/shipment/_fabo"
	sms_min "o.o/backend/cogs/sms/_min"
	storage_all "o.o/backend/cogs/storage/_all"
	_uploader "o.o/backend/cogs/uploader"
	fabopublisher "o.o/backend/com/eventhandler/fabo/publisher"
	"o.o/backend/com/eventhandler/handler"
	"o.o/backend/com/eventhandler/notifier"
	comfabo "o.o/backend/com/fabo"
	"o.o/backend/com/main/address"
	"o.o/backend/com/main/catalog"
	"o.o/backend/com/main/connectioning"
	"o.o/backend/com/main/identity"
	"o.o/backend/com/main/inventory/aggregatex"
	"o.o/backend/com/main/location"
	"o.o/backend/com/main/ordering"
	"o.o/backend/com/main/receipting"
	"o.o/backend/com/main/stocktaking"
	"o.o/backend/com/shopping/carrying"
	"o.o/backend/com/shopping/customering"
	"o.o/backend/com/summary"
	"o.o/backend/com/supporting/ticket"
	"o.o/backend/pkg/common/apifw/captcha"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/export"
	sadmin_fabo "o.o/backend/pkg/etop/api/sadmin/_fabo"
	shop_min "o.o/backend/pkg/etop/api/shop/_min"
	shop_wire "o.o/backend/pkg/etop/api/shop/_wire"
	"o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/authorize/middleware"
	"o.o/backend/pkg/etop/eventstream"
	fulfillmentcsv "o.o/backend/pkg/etop/logic/fulfillments/imcsv"
	logicorder "o.o/backend/pkg/etop/logic/orders"
	orderimcsv "o.o/backend/pkg/etop/logic/orders/imcsv"
	productimcsv "o.o/backend/pkg/etop/logic/products/imcsv"
	logicsummary "o.o/backend/pkg/etop/logic/summary"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/backend/pkg/fabo"
	"o.o/backend/pkg/integration/email"
	"o.o/capi"
)

func Build(
	ctx context.Context,
	cfg config.Config,
	consumer mq.KafkaConsumer,
) (Output, func(), error) {
	panic(wire.Build(
		wire.FieldsOf(&cfg,
			"email",
			"smtp",
			"sms",
			"databases",
			"invitation",
			"kafka",
			"shipment",
			"redis",
			"captcha",
			"ExportDirs",
			"UploadDirs",
			"StorageDriver",
			"FacebookApp",
			"SharedConfig",
			"Webhook",
			"FlagEnableNewLinkInvitation",
			"FlagFaboOrderAutoConfirmPaymentStatus",
		),
		wire.Struct(new(Output), "*"),
		_base.WireSet,
		database_min.WireSet,
		sms_min.WireSet,
		config_server.WireSet,
		_uploader.WireSet,
		_core.WireSet,
		server_shop.WireSet,
		server_fabo.WireSet,
		shop_min.WireSet,
		shop_wire.WireSet,
		storage_all.WireSet,

		email.WireSet,
		logicorder.WireSet,
		orderimcsv.WireSet,
		productimcsv.WireSet,
		fulfillmentcsv.WireSet,
		eventstream.WireSet,
		api.WireSet,
		location.WireSet,
		catalog.WireSet,
		customering.WireSet,
		ordering.WireSet,
		stocktaking.WireSet,
		identity.WireSet,
		notifier.WireSet,

		address.WireSet,
		carrying.WireSet,

		receipting.WireSet,
		aggregatex.WireSet,
		summary.WireSet,
		connectioning.WireSet,
		shipment_fabo.WireSet,
		export.WireSet,
		middleware.WireSet,
		logicsummary.WireSet,
		sadmin_fabo.WireSet,
		wire.Bind(new(bus.EventRegistry), new(bus.Bus)),
		wire.Bind(new(capi.EventBus), new(bus.Bus)),
		wire.Bind(new(eventstream.Publisher), new(*eventstream.EventStream)),
		sqlstore.WireSet,
		captcha.WireSet,
		ticket.WireSet, // TODO(vu): remove

		ProvidePolicy,
		auth.WireSet,

		// fabo
		handler.WireSet,
		fabopublisher.WireSet,
		fabo.WireSet,
		comfabo.WireSet,

		// TODO(vu): remove
		wire.Value(tradering.QueryBus{}),
		wire.Value(affiliate.CommandBus{}),

		BuildIntHandlers,
		BuildMainServer,
		BuildWebhookServer,
		BuildServers,
		SupportedShipnowManager,
	))
}
