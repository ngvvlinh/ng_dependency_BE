package build

import (
	"net/http"

	"o.o/backend/cmd/fabo-sync-service/config"
	servicefbmessaging "o.o/backend/com/fabo/main/fbmessaging"
	"o.o/backend/com/fabo/pkg/sync"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/lifecycle"
	"o.o/backend/pkg/common/metrics"
	"o.o/common/l"
)

type Output struct {
	Servers []lifecycle.HTTPServer
	Sync    *sync.Synchronizer
	Health  *health.Service

	// pm
	_fbmessaging *servicefbmessaging.ProcessManager
}

func BuildServers(cfg config.Config, healthService *health.Service) []lifecycle.HTTPServer {
	mux := http.NewServeMux()
	l.RegisterHTTPHandler(mux)
	metrics.RegisterHTTPHandler(mux)
	healthService.RegisterHTTPHandler(mux)

	svr := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: mux,
	}
	return []lifecycle.HTTPServer{
		{"Main", svr},
	}
}
