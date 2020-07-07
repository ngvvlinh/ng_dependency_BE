package build

import (
	"net/http"

	shipnowcarrier "o.o/api/main/shipnow/carrier"
	"o.o/backend/cmd/fabo-server/config"
	config_server "o.o/backend/cogs/config/_server"
	server_admin "o.o/backend/cogs/server/admin"
	server_shop "o.o/backend/cogs/server/shop"
	_ghn "o.o/backend/cogs/shipment/ghn"
	_ghtk "o.o/backend/cogs/shipment/ghtk"
	_vtpost "o.o/backend/cogs/shipment/vtpost"
	fabopublisher "o.o/backend/com/eventhandler/fabo/publisher"
	"o.o/backend/com/eventhandler/handler"
	"o.o/backend/com/fabo/main/fbmessaging"
	fbuserpm "o.o/backend/com/fabo/main/fbuser/pm"
	fbwebhook "o.o/backend/com/fabo/pkg/webhook"
	catalogpm "o.o/backend/com/main/catalog/pm"
	identitypm "o.o/backend/com/main/identity/pm"
	inventorypm "o.o/backend/com/main/inventory/pm"
	invitationpm "o.o/backend/com/main/invitation/pm"
	ledgerpm "o.o/backend/com/main/ledgering/pm"
	moneytxpm "o.o/backend/com/main/moneytx/pm"
	orderingpm "o.o/backend/com/main/ordering/pm"
	receiptpm "o.o/backend/com/main/receipting/pm"
	refundpm "o.o/backend/com/main/refund/pm"
	shippingpm "o.o/backend/com/main/shipping/pm"
	"o.o/backend/pkg/common/apifw/captcha"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/common/lifecycle"
	"o.o/backend/pkg/common/metrics"
	"o.o/backend/pkg/common/sql/sqltrace"
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/admin"
	"o.o/backend/pkg/etop/api/sadmin"
	"o.o/backend/pkg/etop/api/shop"
	"o.o/backend/pkg/etop/authorize/middleware"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/eventstream"
	"o.o/backend/pkg/etop/middlewares"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/backend/pkg/fabo"
	"o.o/backend/tools/pkg/acl"
	"o.o/capi/httprpc"
	"o.o/common/jsonx"
	"o.o/common/l"
)

type Output struct {
	Servers     []lifecycle.HTTPServer
	EventStream *eventstream.EventStream

	// handlers
	Handler   *handler.Handler
	Publisher *fabopublisher.Publisher

	// pm
	_catalogPM     *catalogpm.ProcessManager
	_identityPM    *identitypm.ProcessManager
	_inventoryPM   *inventorypm.ProcessManager
	_invitationPM  *invitationpm.ProcessManager
	_ledgerPM      *ledgerpm.ProcessManager
	_moneytxPM     *moneytxpm.ProcessManager
	_orderPM       *orderingpm.ProcessManager
	_receiptPM     *receiptpm.ProcessManager
	_refundPM      *refundpm.ProcessManager
	_shippingPM    *shippingpm.ProcessManager
	_fbuserPM      *fbuserpm.ProcessManager
	_fbMessagingPM *fbmessaging.ProcessManager

	// _affiliatePM      *affiliatepm.ProcessManager
	// _purchaseOrderPM  *purchaseorderpm.ProcessManager
	// _purchaseRefundPM *purchaserefundpm.ProcessManager
	// _shipnowPM  *shipnowpm.ProcessManager
	// _traderPM         *traderpm.ProcessManager

	// inject
	_s *sqlstore.Store
	_m middleware.Middleware
	_c *captcha.Captcha
}

func BuildServers(
	mainServer MainServer,
	ghnServer _ghn.GHNWebhookServer,
	ghtkServer _ghtk.GHTKWebhookServer,
	vtpostServer _vtpost.VTPostWebhookServer,
	fbWebhook FBWebhookServer,
) []lifecycle.HTTPServer {
	svrs := []lifecycle.HTTPServer{
		{"Main   ", mainServer},
		{"GHN    ", ghnServer},
		{"GHTK   ", ghtkServer},
		{"VTPOST ", vtpostServer},
		{"Webhook", fbWebhook},
	}
	return svrs
}

type IntHandlers []httprpc.Server

func BuildIntHandlers(
	rootServers api.Servers,
	shopServers shop.Servers,
	adminServers admin.Servers,
	sadminServers sadmin.Servers,
	faboServers fabo.Servers,
) (hs IntHandlers) {
	logging := middlewares.NewLogging()
	ssHooks := session.NewHook(acl.GetACL())

	hs = append(hs, rootServers...)
	hs = append(hs, shopServers...)
	hs = append(hs, adminServers...)
	hs = append(hs, faboServers...)
	hs = append(hs, sadminServers...)
	hs = httprpc.WithHooks(hs, ssHooks, logging)
	return hs
}

type MainServer *http.Server

func BuildMainServer(
	healthService *health.Service,
	intHandlers IntHandlers,
	cfg config_server.SharedConfig,
	adminImport server_admin.ImportServer,
	shopImport server_shop.ImportHandler,
	eventStream server_shop.EventStreamHandler,
	downloadHandler server_shop.DownloadHandler,
) MainServer {
	mux := http.NewServeMux()
	l.RegisterHTTPHandler(mux)
	metrics.RegisterHTTPHandler(mux)
	healthService.RegisterHTTPHandler(mux)
	jsonx.RegisterHTTPHandler(mux)
	sqltrace.RegisterHTTPHandler(mux)

	mwares := httpx.Compose(
		headers.ForwardHeadersX(),
		bus.Middleware,
	)
	intHandlers = httprpc.WithPrefix("/api/", intHandlers)
	mux.Handle("/api/", http.NotFoundHandler())

	var handlers []httprpc.Server
	handlers = append(handlers, intHandlers...)
	for _, h := range handlers {
		mux.Handle(h.PathPrefix(), mwares(h))
	}

	mux.Handle(adminImport.PathPrefix(), mwares(adminImport))
	mux.Handle(shopImport.PathPrefix(), mwares(shopImport))
	mux.Handle(eventStream.PathPrefix(), eventStream)
	mux.Handle(downloadHandler.PathPrefix(), downloadHandler)

	h := middleware.CORS(mux)
	svr := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: h,
	}
	return svr
}

type FBWebhookServer *http.Server

func BuildWebhookServer(
	cfg config.WebhookConfig,
	webhook *fbwebhook.Webhook,
) FBWebhookServer {
	rt := httpx.New()
	rt.Use(httpx.RecoverAndLog(true))
	webhook.Register(rt)

	mux := http.NewServeMux()
	mux.Handle("/", rt)
	svr := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: mux,
	}
	return svr
}

func SupportedShipnowManager() shipnowcarrier.Manager {
	return nil
}
