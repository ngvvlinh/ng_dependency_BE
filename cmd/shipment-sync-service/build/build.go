package build

import (
	"context"
	"fmt"
	"net/http"

	"o.o/api/main/shipping"
	"o.o/backend/cmd/shipment-sync-service/config"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/shipping/carrier"
	shippingpm "o.o/backend/com/main/shipping/pm"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/lifecycle"
	"o.o/backend/pkg/common/metrics"
	dhlSync "o.o/backend/pkg/integration/shipping/dhl/sync"
	"o.o/common/l"
)

type Output struct {
	Servers []lifecycle.HTTPServer
	Health  *health.Service
	DHLSync *dhlSync.DHLSync

	// pm
	_shippingPM *shippingpm.ProcessManager
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
	ctx context.Context, db com.MainDB, shipmentManager *carrier.ShipmentManager,
	shippingQS shipping.QueryBus, shippingAggr shipping.CommandBus,
) *dhlSync.DHLSync {
	_dhlSync := dhlSync.New(db, shipmentManager, shippingQS, shippingAggr)
	if err := _dhlSync.Init(ctx); err != nil {
		panic(fmt.Sprintf("Can't init DHL sync: %v", err.Error()))
	}
	go func() {
		_dhlSync.Start(ctx)
	}()

	return _dhlSync
}
