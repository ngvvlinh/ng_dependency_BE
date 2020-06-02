package build

import (
	"net/http"

	shipnowcarrier "o.o/api/main/shipnow/carrier"
	config_server "o.o/backend/cogs/config/_server"
	server_admin "o.o/backend/cogs/server/admin"
	_main "o.o/backend/cogs/server/main"
	server_shop "o.o/backend/cogs/server/shop"
	_ghn "o.o/backend/cogs/shipment/ghn"
	_ghtk "o.o/backend/cogs/shipment/ghtk"
	_vtpost "o.o/backend/cogs/shipment/vtpost"
	"o.o/backend/pkg/common/apifw/captcha"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/common/lifecycle"
	"o.o/backend/pkg/common/metrics"
	"o.o/backend/pkg/common/sql/sqltrace"
	"o.o/backend/pkg/etop/authorize/middleware"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi/httprpc"
	"o.o/common/jsonx"
	"o.o/common/l"
)

func BuildServers(
	_ *sqlstore.Store, // inject
	_ middleware.Middleware, // inject
	_ *captcha.Captcha, // inject

	mainServer MainServer,
	ghnServer _ghn.GHNWebhookServer,
	ghtkServer _ghtk.GHTKWebhookServer,
	vtpostServer _vtpost.VTPostWebhookServer,
) []lifecycle.HTTPServer {
	svrs := []lifecycle.HTTPServer{
		{"Main   ", mainServer},
		{"GHN    ", ghnServer},
		{"GHTK   ", ghtkServer},
		{"VTPOST ", vtpostServer},
	}
	return svrs
}

type MainServer *http.Server

func BuildMainServer(
	healthService *health.Service,
	intHandlers _main.IntHandlers,
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

	middlewares := httpx.Compose(
		headers.ForwardHeadersX(),
		bus.Middleware,
	)
	intHandlers = httprpc.WithPrefix("/api/", intHandlers)

	var handlers []httprpc.Server
	handlers = append(handlers, intHandlers...)
	for _, h := range handlers {
		mux.Handle(h.PathPrefix(), middlewares(h))
	}

	mux.Handle(adminImport.PathPrefix(), adminImport)
	mux.Handle(shopImport.PathPrefix(), shopImport)
	mux.Handle(eventStream.PathPrefix(), eventStream)
	mux.Handle(downloadHandler.PathPrefix(), downloadHandler)

	handler := middleware.CORS(mux)
	svr := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: handler,
	}
	return svr
}

func SupportedShipnowManager() shipnowcarrier.Manager {
	return nil
}
