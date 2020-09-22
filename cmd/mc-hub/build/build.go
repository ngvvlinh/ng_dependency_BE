package build

import (
	"net/http"

	"o.o/api/top/external/mc/vnp"
	"o.o/backend/cmd/mc-hub/config"
	"o.o/backend/cmd/mc-hub/service"
	"o.o/backend/cmd/mc-hub/service/client"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/common/lifecycle"
	"o.o/backend/pkg/common/metrics"
	"o.o/backend/pkg/common/sql/sqltrace"
	"o.o/capi/httprpc"
	"o.o/common/jsonx"
	"o.o/common/l"
)

type Output struct {
	Servers []lifecycle.HTTPServer
	Health  *health.Service
}

func BuildServers(mainServer MainServer) []lifecycle.HTTPServer {
	servers := []lifecycle.HTTPServer{
		{"Main", mainServer},
	}
	return servers
}

type MainServer *http.Server

func BuildMainServer(
	cfg config.Config,
	healthService *health.Service,
	mcService service.MCShipnowService,
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
	handlers := httprpc.MustNewServers(
		mcService.Clone,
	)
	handlers = httprpc.WithPrefix("/v1/", handlers)
	for _, h := range handlers {
		mux.Handle(h.PathPrefix(), mwares(h))
	}
	svr := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: mux,
	}
	return svr
}

func BuildMCShipnow(cfg config.Config) (vnp.ShipnowService, error) {
	return client.New(cfg.Endpoints.Main)
}
