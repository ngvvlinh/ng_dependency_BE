package build

import (
	"net/http"

	"o.o/backend/cmd/fabo-server/config"
	config_server "o.o/backend/cogs/config/_server"
	server_fabo "o.o/backend/cogs/server/fabo"
	server_shop "o.o/backend/cogs/server/shop"
	ghnv2 "o.o/backend/cogs/shipment/ghn/v2"
	fabopublisher "o.o/backend/com/eventhandler/fabo/publisher"
	"o.o/backend/com/eventhandler/handler"
	fbmessagetemplatepm "o.o/backend/com/fabo/main/fbmessagetemplate/pm"
	"o.o/backend/com/fabo/main/fbmessaging"
	fbuserpm "o.o/backend/com/fabo/main/fbuser/pm"
	"o.o/backend/com/fabo/pkg/fbclient"
	fbwebhook "o.o/backend/com/fabo/pkg/webhook"
	catalogpm "o.o/backend/com/main/catalog/pm"
	identitypm "o.o/backend/com/main/identity/pm"
	orderingpm "o.o/backend/com/main/ordering/pm"
	shippingpm "o.o/backend/com/main/shipping/pm"
	"o.o/backend/pkg/common/apifw/captcha"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/common/lifecycle"
	"o.o/backend/pkg/common/metrics"
	"o.o/backend/pkg/common/sql/sqltrace"
	apiroot "o.o/backend/pkg/etop/api/root"
	"o.o/backend/pkg/etop/api/sadmin"
	"o.o/backend/pkg/etop/api/shop"
	"o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/authorize/authfabo"
	"o.o/backend/pkg/etop/authorize/middleware"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/eventstream"
	"o.o/backend/pkg/etop/middlewares"
	"o.o/backend/pkg/fabo"
	"o.o/backend/tools/pkg/acl"
	"o.o/capi/httprpc"
	"o.o/common/jsonx"
	"o.o/common/l"
)

type Output struct {
	Servers     []lifecycle.HTTPServer
	EventStream *eventstream.EventStream
	Health      *health.Service

	// handlers
	Handler   *handler.Handler
	Publisher *fabopublisher.Publisher

	// fbClient
	FbClient *fbclient.FbClient

	// pm
	_catalogPM           *catalogpm.ProcessManager
	_identityPM          *identitypm.ProcessManager
	_orderPM             *orderingpm.ProcessManager
	_shippingPM          *shippingpm.ProcessManager
	_fbuserPM            *fbuserpm.ProcessManager
	_fbMessagingPM       *fbmessaging.ProcessManager
	_fbMessageTemplatePM *fbmessagetemplatepm.ProcessManager
}

func BuildServers(
	mainServer MainServer,
	ghnServer ghnv2.GHNWebhookServer,
	fbWebhook FBWebhookServer,
) []lifecycle.HTTPServer {
	svrs := []lifecycle.HTTPServer{
		{"Main   ", mainServer},
		{"GHN    ", ghnServer},
		{"Webhook", fbWebhook},
	}
	return svrs
}

type IntHandlers []httprpc.Server

func BuildIntHandlers(
	rootServers apiroot.Servers,
	shopServers shop.Servers,
	faboServers fabo.Servers,
	sadminServers sadmin.Servers,
	c *captcha.Captcha,
) (hs IntHandlers, _ error) {
	logging := middlewares.NewLogging()
	ssHooks, err := session.NewHook(acl.GetACL(), session.OptCaptcha(c))
	if err != nil {
		return nil, err
	}

	hs = append(hs, rootServers...)
	hs = append(hs, shopServers...)
	hs = append(hs, faboServers...)
	hs = append(hs, sadminServers...)
	hs = httprpc.WithHooks(hs, ssHooks, logging)
	return
}

type MainServer *http.Server

func BuildMainServer(
	healthService *health.Service,
	intHandlers IntHandlers,
	cfg config_server.SharedConfig,
	shopImport server_shop.ImportHandler,
	eventStream server_shop.EventStreamHandler,
	downloadHandler server_shop.DownloadHandler,
	faboImageHandler server_fabo.FaboImageHandler,
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

	mux.Handle(shopImport.PathPrefix(), mwares(shopImport))
	mux.Handle(eventStream.PathPrefix(), eventStream)
	mux.Handle(downloadHandler.PathPrefix(), downloadHandler)
	mux.Handle(faboImageHandler.PathPrefix(), faboImageHandler)

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

func ProvidePolicy() auth.Policy { return authfabo.Policy }
