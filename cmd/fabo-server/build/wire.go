// +build wireinject

package build

import (
	"context"

	"github.com/google/wire"

	"o.o/backend/cmd/fabo-server/config"
	_base "o.o/backend/cogs/base"
	config_server "o.o/backend/cogs/config/_server"
	_core "o.o/backend/cogs/core"
	database_min "o.o/backend/cogs/database/_min"
	server_admin "o.o/backend/cogs/server/admin"
	server_shop "o.o/backend/cogs/server/shop"
	shipment_all "o.o/backend/cogs/shipment/_all"
	sms_all "o.o/backend/cogs/sms/_all"
	_uploader "o.o/backend/cogs/uploader"
	fabopublisher "o.o/backend/com/eventhandler/fabo/publisher"
	"o.o/backend/com/eventhandler/handler"
	comfabo "o.o/backend/com/fabo"
	"o.o/backend/com/main/address"
	"o.o/backend/com/main/catalog"
	"o.o/backend/com/main/connectioning"
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
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/admin"
	admin_min "o.o/backend/pkg/etop/api/admin/_min"
	"o.o/backend/pkg/etop/api/export"
	"o.o/backend/pkg/etop/api/sadmin"
	"o.o/backend/pkg/etop/api/shop"
	shop_min "o.o/backend/pkg/etop/api/shop/_min"
	"o.o/backend/pkg/etop/apix/partner"
	"o.o/backend/pkg/etop/authorize/middleware"
	"o.o/backend/pkg/etop/eventstream"
	logicorder "o.o/backend/pkg/etop/logic/orders"
	orderimcsv "o.o/backend/pkg/etop/logic/orders/imcsv"
	productimcsv "o.o/backend/pkg/etop/logic/products/imcsv"
	logicsummary "o.o/backend/pkg/etop/logic/summary"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/backend/pkg/fabo"
	"o.o/capi"
)

func Servers(
	ctx context.Context,
	cfg config.Config,
	eventBus capi.EventBus,
	healthServer *health.Service,
	partnerAuthURL partner.AuthURL,
	consumer mq.KafkaConsumer,
) (Output, func(), error) {
	panic(wire.Build(
		wire.FieldsOf(&cfg,
			"email",
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
			"WhiteLabel",
			"SharedConfig",
			"Webhook",
			"FlagApplyShipmentPrice",
			"FlagEnableNewLinkInvitation",
		),
		wire.Struct(new(Output), "*"),
		_base.WireSet,
		shipment_all.WireSet,
		database_min.WireSet,
		sms_all.WireSet,
		config_server.WireSet,
		_uploader.WireSet,
		_core.WireSet,
		server_shop.WireSet,
		server_admin.WireSet,
		shop_min.WireSet,
		shop.WireSet,
		admin_min.WireSet,

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
		wire.InterfaceValue(new(eventstream.Publisher), new(eventstream.EventStream)),
		sqlstore.WireSet,
		captcha.WireSet,

		// fabo
		handler.WireSet,
		fabopublisher.WireSet,
		fabo.WireSet,
		comfabo.WireSet,

		BuildIntHandlers,
		BuildMainServer,
		BuildWebhookServer,
		BuildServers,
		SupportedShipnowManager,
	))
}
