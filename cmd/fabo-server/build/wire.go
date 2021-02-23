// +build wireinject

package build

import (
	"context"

	"github.com/google/wire"

	"o.o/api/main/accountshipnow"
	"o.o/api/main/transaction"
	"o.o/api/services/affiliate"
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
	comfabo "o.o/backend/com/fabo/cogs/_publisher"
	"o.o/backend/com/fabo/main/fbcustomerconversationsearch"
	"o.o/backend/com/fabo/main/fbmessagetemplate"
	"o.o/backend/com/fabo/main/fbmessaging"
	"o.o/backend/com/fabo/main/fbpage"
	"o.o/backend/com/fabo/main/fbuser"
	"o.o/backend/com/main/address"
	"o.o/backend/com/main/catalog"
	"o.o/backend/com/main/connectioning"
	"o.o/backend/com/main/identity"
	"o.o/backend/com/main/inventory/aggregatex"
	"o.o/backend/com/main/location"
	"o.o/backend/com/main/ordering"
	"o.o/backend/com/main/receipting"
	"o.o/backend/com/main/shipnow"
	"o.o/backend/com/main/stocktaking"
	"o.o/backend/com/shopping/customering"
	"o.o/backend/com/shopping/setting"
	fabosummary "o.o/backend/com/summary/fabo"
	"o.o/backend/pkg/common/apifw/captcha"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/etop/api/export"
	apiroot_fabo "o.o/backend/pkg/etop/api/root/fabo"
	sadmin_fabo "o.o/backend/pkg/etop/api/sadmin/_fabo"
	shop_min "o.o/backend/pkg/etop/api/shop/_min/fabo"
	shop_wire "o.o/backend/pkg/etop/api/shop/_wire/fabo"
	"o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/authorize/middleware"
	"o.o/backend/pkg/etop/eventstream"
	fulfillmentcsv "o.o/backend/pkg/etop/logic/fulfillments/imcsv"
	logicorder "o.o/backend/pkg/etop/logic/orders"
	orderimcsv "o.o/backend/pkg/etop/logic/orders/imcsv"
	productimcsv "o.o/backend/pkg/etop/logic/products/imcsv"
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
			"WebphonePublicKey",
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
		apiroot_fabo.WireSet,
		location.WireSet,
		catalog.WireSet,
		customering.WireSet,
		setting.WireSet,
		ordering.WireSet,
		stocktaking.WireSet,
		identity.WireSet,
		notifier.WireSet,

		address.WireSet,
		receipting.WireSet,
		aggregatex.WireSet,
		fabosummary.WireSet,
		connectioning.WireSet,
		shipment_fabo.WireSet,
		export.WireSet,
		middleware.WireSet,
		sadmin_fabo.WireSet,
		wire.Bind(new(bus.EventRegistry), new(bus.Bus)),
		wire.Bind(new(capi.EventBus), new(bus.Bus)),
		wire.Bind(new(eventstream.Publisher), new(*eventstream.EventStream)),
		sqlstore.WireSet,
		captcha.WireSet,

		ProvidePolicy,
		auth.WireSet,

		// fabo
		handler.WireSet,
		fabopublisher.WireSet,
		fbmessaging.WireSet,
		fbpage.WireSet,
		fbuser.WireSet,
		fabo.WireSet,
		comfabo.WireSet,
		fbcustomerconversationsearch.WireSet,
		fbmessagetemplate.WireSet,
		shipnow.WireSet,

		// TODO(vu): remove
		wire.Value(affiliate.CommandBus{}),
		wire.Value(accountshipnow.CommandBus{}),
		wire.Value(accountshipnow.QueryBus{}),
		wire.Value(transaction.QueryBus{}),

		BuildPgProducer,
		BuildIntHandlers,
		BuildMainServer,
		BuildWebhookServer,
		BuildServers,
	))
}
