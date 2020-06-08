// +build wireinject

package build

import (
	"context"

	"github.com/google/wire"

	"o.o/backend/cmd/etop-server/config"
	_base "o.o/backend/cogs/base"
	_producer "o.o/backend/cogs/base/producer"
	config_server "o.o/backend/cogs/config/_server"
	_core "o.o/backend/cogs/core"
	database_all "o.o/backend/cogs/database/_all"
	payment_all "o.o/backend/cogs/payment/_all"
	server_admin "o.o/backend/cogs/server/admin"
	server_max "o.o/backend/cogs/server/main/_max"
	server_shop "o.o/backend/cogs/server/shop"
	server_vtpay "o.o/backend/cogs/server/vtpay"
	shipment_all "o.o/backend/cogs/shipment/_all"
	shipnow_all "o.o/backend/cogs/shipnow/_all"
	sms_all "o.o/backend/cogs/sms/_all"
	_uploader "o.o/backend/cogs/uploader"
	paymentmanager "o.o/backend/com/external/payment/manager"
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
	"o.o/backend/com/services/affiliate"
	"o.o/backend/com/shopping/carrying"
	"o.o/backend/com/shopping/customering"
	"o.o/backend/com/shopping/suppliering"
	"o.o/backend/com/shopping/tradering"
	"o.o/backend/com/subscripting"
	"o.o/backend/com/summary"
	"o.o/backend/com/web/webserver"
	"o.o/backend/pkg/common/apifw/captcha"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/admin"
	admin_all "o.o/backend/pkg/etop/api/admin/_all"
	affapi "o.o/backend/pkg/etop/api/affiliate"
	"o.o/backend/pkg/etop/api/export"
	"o.o/backend/pkg/etop/api/integration"
	"o.o/backend/pkg/etop/api/sadmin"
	"o.o/backend/pkg/etop/api/shop"
	shop_all "o.o/backend/pkg/etop/api/shop/_all"
	"o.o/backend/pkg/etop/apix/partner"
	"o.o/backend/pkg/etop/apix/partnercarrier"
	"o.o/backend/pkg/etop/apix/partnerimport"
	xshipping "o.o/backend/pkg/etop/apix/shipping"
	xshop "o.o/backend/pkg/etop/apix/shop"
	"o.o/backend/pkg/etop/apix/shopping"
	"o.o/backend/pkg/etop/apix/webhook"
	"o.o/backend/pkg/etop/authorize/middleware"
	"o.o/backend/pkg/etop/eventstream"
	logicorder "o.o/backend/pkg/etop/logic/orders"
	orderimcsv "o.o/backend/pkg/etop/logic/orders/imcsv"
	productimcsv "o.o/backend/pkg/etop/logic/products/imcsv"
	logicsummary "o.o/backend/pkg/etop/logic/summary"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/backend/pkg/integration/email"
	saffapi "o.o/backend/pkg/services/affiliate/api"
	"o.o/capi"
)

func Build(
	ctx context.Context,
	cfg config.Config,
	eventBus bus.Bus,
	healthServer *health.Service,
	partnerAuthURL partner.AuthURL,
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
			"ahamove",
			"vtpay",
			"redis",
			"export",
			"captcha",
			"upload",
			"WhiteLabel",
			"SharedConfig",
			"AhamoveWebhook",
			"FlagApplyShipmentPrice",
			"FlagEnableNewLinkInvitation",
		),
		wire.Struct(new(Output), "*"),
		_base.WireSet,
		shipnow_all.WireSet,
		payment_all.WireSet,
		shipment_all.WireSet,
		database_all.WireSet,
		sms_all.WireSet,
		config_server.WireSet,
		_producer.WireSet,
		_uploader.WireSet,
		_core.WireSet,
		server_max.WireSet,
		server_vtpay.WireSet,
		server_shop.WireSet,
		server_admin.WireSet,
		admin_all.WireSet,
		shop_all.WireSet,
		shop.WireSet,

		email.WireSet,
		affiliate.WireSet,
		logicorder.WireSet,
		moneytx.WireSet,
		orderimcsv.WireSet,
		productimcsv.WireSet,
		eventstream.WireSet,
		shopping.WireSet,
		api.WireSet,
		location.WireSet,
		catalog.WireSet,
		customering.WireSet,
		ordering.WireSet,
		stocktaking.WireSet,
		refund.WireSet,
		shipnow.WireSet,
		identity.WireSet,
		address.WireSet,
		paymentmanager.WireSet,
		suppliering.WireSet,
		carrying.WireSet,
		tradering.WireSet,
		receipting.WireSet,
		inventory.WireSet,
		ledgering.WireSet,
		purchaseorder.WireSet,
		summary.WireSet,
		purchaserefund.WireSet,
		connectioning.WireSet,
		webserver.WireSet,
		subscripting.WireSet,
		saffapi.WireSet,
		affapi.WireSet,
		partner.WireSet,
		partnercarrier.WireSet,
		xshop.WireSet,
		xshipping.WireSet,
		partnerimport.WireSet,
		admin.WireSet,
		sadmin.WireSet,
		integration.WireSet,
		export.WireSet,
		webhook.WireSet,
		middleware.WireSet,
		logicsummary.WireSet,
		wire.Bind(new(bus.EventRegistry), new(bus.Bus)),
		wire.Bind(new(capi.EventBus), new(bus.Bus)),
		wire.Bind(new(eventstream.Publisher), new(*eventstream.EventStream)),
		sqlstore.WireSet,
		captcha.WireSet,

		BuildServers,
		BuildMainServer,
		BuildWebServer,
	))
}
