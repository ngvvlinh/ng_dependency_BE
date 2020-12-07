package build

import (
	"context"
	"fmt"
	"net/http"

	"o.o/api/etelecom"
	"o.o/api/main/connectioning"
	"o.o/backend/cmd/telecom-sync-service/config"
	"o.o/backend/com/etelecom/provider"
	"o.o/backend/com/etelecom/provider/types"
	com "o.o/backend/com/main"
	connectioningpm "o.o/backend/com/main/connectioning/pm"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/lifecycle"
	"o.o/backend/pkg/common/metrics"
	vhtsync "o.o/backend/pkg/integration/telecom/vht/sync"
	"o.o/common/l"
)

type Output struct {
	Servers      []lifecycle.HTTPServer
	Health       *health.Service
	TelecomSyncs []types.TelecomSync

	// pm
	_connectionPM *connectioningpm.ProcessManager
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

func BuildSyncs(
	ctx context.Context, dbMain com.MainDB,
	telecomManager *provider.TelecomManager,
	telecomQS etelecom.QueryBus,
	telecomA etelecom.CommandBus,
	connectionA connectioning.CommandBus,
) []types.TelecomSync {
	_vhtSync := vhtsync.New(dbMain, telecomManager, telecomQS, telecomA, connectionA)
	if err := _vhtSync.Init(ctx); err != nil {
		panic(fmt.Sprintf("Can't init VHT sync: %v", err.Error()))
	}
	go func() {
		_vhtSync.Start(ctx)
	}()

	return []types.TelecomSync{_vhtSync}
}
