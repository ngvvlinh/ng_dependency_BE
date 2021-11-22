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
	server_kpay "o.o/backend/cogs/server/kpay"
	server_max "o.o/backend/cogs/server/main/_max"
	server_shop "o.o/backend/cogs/server/shop"
	server_vtpay "o.o/backend/cogs/server/vtpay"
	shipment_all "o.o/backend/cogs/shipment/_all"
	shipnow_all "o.o/backend/cogs/shipnow/_all"
	sms_all "o.o/backend/cogs/sms/_all"
	storage_all "o.o/backend/cogs/storage/_all"
	telecom_all "o.o/backend/cogs/telecom/_all"
	ticket_all "o.o/backend/cogs/ticket/_all"
	_uploader "o.o/backend/cogs/uploader"
	"o.o/backend/com/eventhandler/notifier"
	paymentmanager "o.o/backend/com/external/payment/manager"
	minfabo "o.o/backend/com/fabo/cogs/_min"
	"o.o/backend/com/fabo/main/fbmessaging"
	"o.o/backend/com/main/address"
	"o.o/backend/com/main/catalog"
	"o.o/backend/com/main/connectioning"
	"o.o/backend/com/main/contact"
	credit "o.o/backend/com/main/credit"
	"o.o/backend/com/main/identity"
	"o.o/backend/com/main/department"
	"o.o/backend/com/main/inventory"
	"o.o/backend/com/main/invoicing"
	"o.o/backend/com/main/ledgering"
	"o.o/backend/com/main/location"
	"o.o/backend/com/main/moneytx"
	"o.o/backend/com/main/ordering"
	"o.o/backend/com/main/purchaseorder"
	"o.o/backend/com/main/purchaserefund"
	"o.o/backend/com/main/receipting"
	"o.o/backend/com/main/refund"
	"o.o/backend/com/main/reporting"
	"o.o/backend/com/main/shipnow"
	"o.o/backend/com/main/stocktaking"
	"o.o/backend/com/main/transaction"
	"o.o/backend/com/report"
	"o.o/backend/com/services/affiliate"
	"o.o/backend/com/shopping/carrying"
	"o.o/backend/com/shopping/customering"
	"o.o/backend/com/shopping/setting"
	"o.o/backend/com/shopping/suppliering"
	"o.o/backend/com/shopping/tradering"
	"o.o/backend/com/subscripting"
	etelecomsummary "o.o/backend/com/summary/etelecom"
	etopsummary "o.o/backend/com/summary/etop"
	"o.o/backend/com/supporting/ticket"
	"o.o/backend/com/web/webserver"
	"o.o/backend/pkg/common/apifw/captcha"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api/admin"
	admin_all "o.o/backend/pkg/etop/api/admin/_all"
	affapi "o.o/backend/pkg/etop/api/affiliate"
	"o.o/backend/pkg/etop/api/export"
	"o.o/backend/pkg/etop/api/integration"
	apiroot_all "o.o/backend/pkg/etop/api/root/_all"
	"o.o/backend/pkg/etop/api/sadmin"
	shop_all "o.o/backend/pkg/etop/api/shop/_all"
	shop_wire "o.o/backend/pkg/etop/api/shop/_wire/etop"
	"o.o/backend/pkg/etop/apix/mc/vht"
	"o.o/backend/pkg/etop/apix/mc/vnp"
	"o.o/backend/pkg/etop/apix/partner"
	"o.o/backend/pkg/etop/apix/partnercarrier"
	"o.o/backend/pkg/etop/apix/partnerimport"
	xshipping "o.o/backend/pkg/etop/apix/shipping"
	xshop "o.o/backend/pkg/etop/apix/shop"
	"o.o/backend/pkg/etop/apix/shopping"
	"o.o/backend/pkg/etop/apix/webhook"
	"o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/authorize/middleware"
	"o.o/backend/pkg/etop/eventstream"
	fulfillmentcsv "o.o/backend/pkg/etop/logic/fulfillments/imcsv"
	hotfix "o.o/backend/pkg/etop/logic/hotfix"
	logicorder "o.o/backend/pkg/etop/logic/orders"
	orderimcsv "o.o/backend/pkg/etop/logic/orders/imcsv"
	productimcsv "o.o/backend/pkg/etop/logic/products/imcsv"
	logicsummary "o.o/backend/pkg/etop/logic/summary"
	"o.o/backend/pkg/etop/sqlstore"
	telecomstore "o.o/backend/pkg/etop/sqlstore/telecom"
	"o.o/backend/pkg/integration/email"
	"o.o/backend/pkg/integration/jira"
	saffapi "o.o/backend/pkg/services/affiliate/api"
	"o.o/capi"
)

func Build(
	ctx context.Context,
	cfg config.Config,
	partnerAuthURL partner.AuthURL,
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
			"ahamove",
			"vtpay",
			"kpay",
			"redis",
			"captcha",
			"ExportDirs",
			"UploadDirs",
			"StorageDriver",
			"WhiteLabel",
			"SharedConfig",
			"AhamoveWebhook",
			"FlagFaboOrderAutoConfirmPaymentStatus",
			"WebphonePublicKey",
			"AdminPortsip",
			"Jira",
		),
		wire.Struct(new(Output), "*"),
		_base.WireSet,
		payment_all.WireSet,
		database_all.WireSet,
		storage_all.WireSet,
		hotfix.WireSet,
		sms_all.WireSet,
		config_server.WireSet,
		_producer.WireSet,
		_uploader.WireSet,
		_core.WireSet,
		server_max.WireSet,
		server_vtpay.WireSet,
		server_kpay.WireSet,
		server_shop.WireSet,
		server_admin.WireSet,
		shop_wire.WireSet,
		vnp.WireSet,
		vht.WireSet,

		email.WireSet,
		affiliate.WireSet,
		logicorder.WireSet,
		moneytx.WireSet,
		orderimcsv.WireSet,
		productimcsv.WireSet,
		fulfillmentcsv.WireSet,
		eventstream.WireSet,
		shopping.WireSet,
		apiroot_all.WireSet,
		location.WireSet,
		catalog.WireSet,
		customering.WireSet,
		setting.WireSet,
		ordering.WireSet,
		stocktaking.WireSet,
		refund.WireSet,
		identity.WireSet,
		notifier.WireSet,
		address.WireSet,
		paymentmanager.WireSet,
		suppliering.WireSet,
		carrying.WireSet,
		tradering.WireSet,
		receipting.WireSet,
		inventory.WireSet,
		ledgering.WireSet,
		purchaseorder.WireSet,
		etopsummary.WireSet,
		etelecomsummary.WireSet,
		purchaserefund.WireSet,
		connectioning.WireSet,
		shipment_all.WireSet,
		shipnow_all.WireSet,
		shipnow.WireSet,
		webserver.WireSet,
		subscripting.WireSet,
		invoicing.WireSet,
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
		reporting.WireSet,
		report.WireSet,
		minfabo.WireSet,
		fbmessaging.WireSet,
		wire.Bind(new(bus.EventRegistry), new(bus.Bus)),
		wire.Bind(new(capi.EventBus), new(bus.Bus)),
		wire.Bind(new(eventstream.Publisher), new(*eventstream.EventStream)),
		sqlstore.WireSet,
		telecomstore.WireSet,
		captcha.WireSet,
		credit.WireSet,
		transaction.WireSet,
		ticket.WireSet,
		contact.WireSet,
		department.WireSet,

		admin_all.WireSet,
		shop_all.WireSet,
		ticket_all.WireSet,
		telecom_all.WireSet,

		jira.WireSet,

		ProvidePolicy,
		auth.WireSet,

		BuildServers,
		BuildMainServer,
		BuildWebServer,
	))
}
