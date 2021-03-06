package build

import (
	"net/http"
	"strings"

	"o.o/api/main/catalog"
	"o.o/api/main/location"
	subscriptioncore "o.o/api/subscripting/subscription"
	webserverinternal "o.o/api/webserver"
	"o.o/backend/cmd/etop-server/config"
	config_server "o.o/backend/cogs/config/_server"
	server_admin "o.o/backend/cogs/server/admin"
	server_kpay "o.o/backend/cogs/server/kpay"
	_main "o.o/backend/cogs/server/main"
	server_shop "o.o/backend/cogs/server/shop"
	server_vtpay "o.o/backend/cogs/server/vtpay"
	_ghtk "o.o/backend/cogs/shipment/ghtk"
	_ntx "o.o/backend/cogs/shipment/ntx"
	_vtpost "o.o/backend/cogs/shipment/vtpost"
	shipmentwebhookall "o.o/backend/cogs/shipment/webhook/_all"
	etelecompm "o.o/backend/com/etelecom/pm"
	portsippm "o.o/backend/com/etelecom/pm/portsip_pm"
	etelecomusersettingpm "o.o/backend/com/etelecom/usersetting/pm"
	catalogpm "o.o/backend/com/main/catalog/pm"
	connectioningpm "o.o/backend/com/main/connectioning/pm"
	creditpm "o.o/backend/com/main/credit/pm"
	identitypm "o.o/backend/com/main/identity/pm"
	inventorypm "o.o/backend/com/main/inventory/pm"
	invitationpm "o.o/backend/com/main/invitation/pm"
	invoicingpm "o.o/backend/com/main/invoicing/pm"
	ledgerpm "o.o/backend/com/main/ledgering/pm"
	moneytxpm "o.o/backend/com/main/moneytx/pm"
	orderingpm "o.o/backend/com/main/ordering/pm"
	purchaseorderpm "o.o/backend/com/main/purchaseorder/pm"
	purchaserefundpm "o.o/backend/com/main/purchaserefund/pm"
	receiptpm "o.o/backend/com/main/receipting/pm"
	refundpm "o.o/backend/com/main/refund/pm"
	pricelistpm "o.o/backend/com/main/shipmentpricing/pricelist/pm"
	shipnowpm "o.o/backend/com/main/shipnow/pm"
	shippingpm "o.o/backend/com/main/shipping/pm"
	transactionpm "o.o/backend/com/main/transaction/pm"
	"o.o/backend/com/report/reportserver"
	affiliatepm "o.o/backend/com/services/affiliate/pm"
	customerpm "o.o/backend/com/shopping/customering/pm"
	traderpm "o.o/backend/com/shopping/tradering/pm"
	subscriptionpm "o.o/backend/com/subscripting/subscription/pm"
	ticketpm "o.o/backend/com/supporting/ticket/pm"
	"o.o/backend/com/web/ecom/webserver"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/apifw/servedoc"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/common/lifecycle"
	"o.o/backend/pkg/common/metrics"
	"o.o/backend/pkg/common/projectpath"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/sqltrace"
	"o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/authorize/authetop"
	"o.o/backend/pkg/etop/authorize/middleware"
	middlewarewebphone "o.o/backend/pkg/etop/authorize/middleware/webphone"
	ahamoveserver "o.o/backend/pkg/integration/shipnow/ahamove/server"
	"o.o/capi/httprpc"
	"o.o/common/jsonx"
	"o.o/common/l"
)

var ll = l.New()

type Output struct {
	Servers []lifecycle.HTTPServer
	Health  *health.Service

	// pm
	_identityPM            *identitypm.ProcessManager
	_inventoryPM           *inventorypm.ProcessManager
	_invitationPM          *invitationpm.ProcessManager
	_catalogPM             *catalogpm.ProcessManager
	_ledgerPM              *ledgerpm.ProcessManager
	_moneytxPM             *moneytxpm.ProcessManager
	_orderPM               *orderingpm.ProcessManager
	_purchaseOrderPM       *purchaseorderpm.ProcessManager
	_purchaseRefundPM      *purchaserefundpm.ProcessManager
	_receiptPM             *receiptpm.ProcessManager
	_refundPM              *refundpm.ProcessManager
	_shipnowPM             *shipnowpm.ProcessManager
	_shippingPM            *shippingpm.ProcessManager
	_affiliatePM           *affiliatepm.ProcessManager
	_traderPM              *traderpm.ProcessManager
	_connectionPM          *connectioningpm.ProcessManager
	_pricelistPM           *pricelistpm.ProcessManager
	_customerPM            *customerpm.ProcessManager
	_etelecomPM            *etelecompm.ProcessManager
	_subscriptionPM        *subscriptionpm.ProcessManager
	_invoicingPM           *invoicingpm.ProcessManager
	_transactionPM         *transactionpm.ProcessManager
	_etelecomUserSettingPM *etelecomusersettingpm.ProcessManager
	_ticketPM              *ticketpm.ProcessManager
	_portsipPM             *portsippm.ProcessManager
	_creditPM              *creditpm.ProcessManager
}

func BuildServers(
	etopServer MainServer,
	webServer WebServer,
	shipmentWebhookServer shipmentwebhookall.ShipmentWebhookServer,
	ghtkServer _ghtk.GHTKWebhookServer,
	vtpostServer _vtpost.VTPostWebhookServer,
	ntxServer _ntx.NTXWebhookServer,
	ahamoveServer ahamoveserver.AhamoveWebhookServer,
) []lifecycle.HTTPServer {
	svrs := []lifecycle.HTTPServer{
		{"Main   ", etopServer},
		{"Web    ", webServer},
		{"Shipment webhook    ", shipmentWebhookServer},
		{"GHTK   ", ghtkServer},
		{"VTPOST ", vtpostServer},
		{"NTX ", ntxServer},
		{"AHAMOVE", ahamoveServer},
	}
	return svrs
}

type MainServer *http.Server

func BuildMainServer(
	healthService *health.Service,
	intHandlers _main.IntHandlers,
	extHandlers _main.ExtHandlers,
	authxHandler _main.AuthxHandler,
	oidcHandler _main.OIDCHandler,
	portsipHandler _main.PortSipHandler,
	sharedCfg config_server.SharedConfig,
	cfg config.Config,
	adminImport server_admin.ImportServer,
	shopImport server_shop.ImportHandler,
	eventStream server_shop.EventStreamHandler,
	downloadHandler server_shop.DownloadHandler,
	vtpayServer server_vtpay.VTPayHandler,
	kpayServer server_kpay.KPayHandler,
	reportServer reportserver.ReportServer,
) MainServer {
	mux := http.NewServeMux()
	l.RegisterHTTPHandler(mux)
	metrics.RegisterHTTPHandler(mux)
	healthService.RegisterHTTPHandler(mux)
	jsonx.RegisterHTTPHandler(mux)
	sqltrace.RegisterHTTPHandler(mux)

	mwares := httpx.Compose(
		headers.ForwardHeadersX(headers.Config{AllowQueryAuthorization: true}),
		bus.Middleware,
	)
	intHandlers = httprpc.WithPrefix("/api/", intHandlers)
	extHandlers = httprpc.WithPrefix("/v1/", extHandlers)
	mux.Handle("/api/", http.NotFoundHandler())
	mux.Handle("/v1/", http.NotFoundHandler())

	var handlers []httprpc.Server
	handlers = append(handlers, intHandlers...)
	handlers = append(handlers, extHandlers...)
	for _, h := range handlers {
		mux.Handle(h.PathPrefix(), mwares(h))
	}

	mux.Handle(adminImport.PathPrefix(), mwares(adminImport))
	mux.Handle(shopImport.PathPrefix(), mwares(shopImport))
	mux.Handle(eventStream.PathPrefix(), eventStream)
	mux.Handle(downloadHandler.PathPrefix(), downloadHandler)
	mux.Handle(vtpayServer.PathPrefix(), vtpayServer)
	mux.Handle(kpayServer.PathPrefix(), kpayServer)
	mux.Handle(reportServer.PathPrefix(), mwares(reportServer))
	mux.Handle(authxHandler.PathPrefix(), mwares(authxHandler))
	mux.Handle(oidcHandler.PathPrefix(), mwares(oidcHandler))
	mux.Handle(portsipHandler.PathPrefix(), mwares(portsipHandler))

	if sharedCfg.ServeDoc {
		mux.Handle("/", http.RedirectHandler("/doc/etop", http.StatusTemporaryRedirect))
		mux.Handle("/doc", http.RedirectHandler("/doc/etop", http.StatusTemporaryRedirect))
		for _, s := range strings.Split("etop/sadmin,etop/admin,etop/shop,etop/etelecom,etop/integration,etop/affiliate,services/crm,services/affiliate,fabo", ",") {
			swaggerPath := "/doc/" + s + "/swagger.json"
			mux.Handle("/doc/"+s, servedoc.RedocHandler())
			mux.Handle(swaggerPath, servedoc.SwaggerHandler(s+"/swagger.json"))
		}
		mux.Handle("/doc/etop", servedoc.RedocHandler())
		mux.Handle("/doc/etop/swagger.json", servedoc.SwaggerHandler("etop/swagger.json"))

	} else {
		ll.Warn("DOCUMENTATION IS DISABLED (config.serve_doc = false)")
		mux.Handle("/doc", http.NotFoundHandler())
	}

	// always serve partner documentation
	mux.Handle("/doc/ext", http.RedirectHandler("/doc", http.StatusTemporaryRedirect))
	for _, s := range strings.Split("shop,partner,partnercarrier,mc/vnp,mc/vht", ",") {
		mux.Handle("/doc/ext/"+s, servedoc.RedocHandler())
		mux.Handle("/doc/ext/"+s+"/swagger.json", servedoc.SwaggerHandler("external/"+s+"/swagger.json"))
	}

	// middleware webphone -> global
	h := middlewarewebphone.CORS(cfg.WebphonePublicKey)(middleware.CORS(mux))
	svr := &http.Server{
		Addr:    sharedCfg.HTTP.Address(),
		Handler: h,
	}
	return svr
}

type WebServer *http.Server

func BuildWebServer(cfg config.Config, webServerQuery webserverinternal.QueryBus, catalogQuery catalog.QueryBus, subscriptionQuery subscriptioncore.QueryBus, rd redis.Store, locationQueryBus location.QueryBus) WebServer {
	ecom := cfg.Ecom
	c := webserver.Config{
		MainSite: ecom.MainSite,
		CoreSite: cfg.URL.MainSite,
		RootPath: projectpath.GetPath(),
	}
	handler, err := webserver.New(c, webServerQuery, catalogQuery, rd, locationQueryBus, subscriptionQuery, nil) // TODO: order logic
	if err != nil {
		ll.S.Panicf("error starting web server: %v", err)
	}

	svr := &http.Server{
		Addr:    ecom.HTTP.Address(),
		Handler: handler,
	}
	return svr
}

func ProvidePolicy() auth.Policy { return authetop.Policy }
