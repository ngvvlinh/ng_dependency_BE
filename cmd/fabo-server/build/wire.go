// +build wireinject

package build

import (
	"context"

	"github.com/google/wire"

	"o.o/api/services/affiliate"
	"o.o/backend/cmd/fabo-server/config"
	_base "o.o/backend/cogs/base"
	config_server "o.o/backend/cogs/config/_server"
	_core "o.o/backend/cogs/core"
	database_min "o.o/backend/cogs/database/_min"
	server_admin "o.o/backend/cogs/server/admin"
	server_shop "o.o/backend/cogs/server/shop"
	shipment_fabo "o.o/backend/cogs/shipment/_fabo"
	sms_min "o.o/backend/cogs/sms/_min"
	_uploader "o.o/backend/cogs/uploader"
	fabopublisher "o.o/backend/com/eventhandler/fabo/publisher"
	"o.o/backend/com/eventhandler/handler"
	comfabo "o.o/backend/com/fabo"
	"o.o/backend/com/main/address"
	"o.o/backend/com/main/catalog"
	"o.o/backend/com/main/connectioning"
	"o.o/backend/com/main/credit"
	"o.o/backend/com/main/identity"
	"o.o/backend/com/main/inventory"
	"o.o/backend/com/main/ledgering"
	"o.o/backend/com/main/location"
	"o.o/backend/com/main/moneytx"
	"o.o/backend/com/main/ordering"
	"o.o/backend/com/main/purchaseorder"
	"o.o/backend/com/main/purchaserefund"
	"o.o/backend/com/main/receipting"
	"o.o/backend/com/main/refund"
	"o.o/backend/com/main/shipnow"
	"o.o/backend/com/main/stocktaking"
	"o.o/backend/com/shopping/carrying"
	"o.o/backend/com/shopping/customering"
	"o.o/backend/com/shopping/suppliering"
	"o.o/backend/com/shopping/tradering"
	"o.o/backend/com/summary"
	"o.o/backend/pkg/common/apifw/captcha"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/admin"
	admin_min "o.o/backend/pkg/etop/api/admin/_min"
	"o.o/backend/pkg/etop/api/export"
	"o.o/backend/pkg/etop/api/sadmin"
	"o.o/backend/pkg/etop/api/shop"
	shop_min "o.o/backend/pkg/etop/api/shop/_min"
	"o.o/backend/pkg/etop/authorize/middleware"
	"o.o/backend/pkg/etop/eventstream"
	hotfixmoneytx "o.o/backend/pkg/etop/logic/hotfix"
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
			"secret",
			"kafka",
			"shipment",
			"redis",
			"export",
			"captcha",
			"upload",
			"FacebookApp",
			"SharedConfig",
			"Webhook",
			"FlagApplyShipmentPrice",
			"FlagEnableNewLinkInvitation",
			"FlagFaboOrderAutoConfirmPaymentStatus",
		),
		wire.Struct(new(Output), "*"),
		_base.WireSet,
		shipment_fabo.WireSet,
		database_min.WireSet,
		hotfixmoneytx.WireSet,
		sms_min.WireSet,
		config_server.WireSet,
		_uploader.WireSet,
		_core.WireSet,
		server_shop.WireSet,
		server_admin.WireSet,
		shop_min.WireSet,
		shop.WireSet,
		admin_min.WireSet,

		email.WireSet,
		logicorder.WireSet,
		moneytx.WireSet,
		orderimcsv.WireSet,
		productimcsv.WireSet,
		eventstream.WireSet,
		api.WireSet,
		location.WireSet,
		catalog.WireSet,
		customering.WireSet,
		ordering.WireSet,
		stocktaking.WireSet,
		refund.WireSet,
		shipnow.WireSet, // TODO(vu): remove
		identity.WireSet,
		address.WireSet,
		suppliering.WireSet, // TODO(vu): remove
		carrying.WireSet,
		tradering.WireSet, // TODO(vu): remove
		receipting.WireSet,
		inventory.WireSet,
		ledgering.WireSet,
		purchaseorder.WireSet,
		summary.WireSet,
		purchaserefund.WireSet,
		connectioning.WireSet,
		admin.WireSet,
		sadmin.WireSet,
		export.WireSet,
		middleware.WireSet,
		logicsummary.WireSet,
		wire.Bind(new(bus.EventRegistry), new(bus.Bus)),
		wire.Bind(new(capi.EventBus), new(bus.Bus)),
		wire.Bind(new(eventstream.Publisher), new(*eventstream.EventStream)),
		sqlstore.WireSet,
		captcha.WireSet,
		credit.WireSet,

		// fabo
		handler.WireSet,
		fabopublisher.WireSet,
		fabo.WireSet,
		comfabo.WireSet,

		// TODO(vu): remove
		wire.Value(affiliate.CommandBus{}),

		BuildIntHandlers,
		BuildMainServer,
		BuildWebhookServer,
		BuildServers,
		SupportedShipnowManager,
	))
}
