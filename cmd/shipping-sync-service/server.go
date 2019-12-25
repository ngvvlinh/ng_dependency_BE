package main

import (
	"net/http"

	servicelocation "etop.vn/backend/com/main/location"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/metrics"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/etop/logic/shipping_provider"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/backend/pkg/integration/shipping/ghn"
	"etop.vn/backend/pkg/integration/shipping/ghtk"
	"etop.vn/common/l"
)

func startServers() []*http.Server {
	return []*http.Server{
		startServiceServer(),
	}
}

func startServiceServer() *http.Server {
	mux := http.NewServeMux()
	l.RegisterHTTPHandler(mux)
	metrics.RegisterHTTPHandler(mux)
	healthservice.RegisterHTTPHandler(mux)

	db, err := cmsql.Connect(cfg.Postgres)
	if err != nil {
		ll.Fatal("Unable to connect to Postgres", l.Error(err))
	}
	sqlstore.Init(db)

	locationBus := servicelocation.New().MessageBus()
	var ghnCarrier, ghtkCarrier shipping_provider.ShippingProvider

	if cfg.GHN.AccountDefault.Token != "" {
		ghnCarrier = ghn.New(cfg.GHN, locationBus)
		if err := ghnCarrier.(*ghn.Carrier).InitAllClients(ctx); err != nil {
			ll.Fatal("Unable to connect to GHN", l.Error(err))
		}
		ll.S.Info("GHN: connect success")
	} else {
		if cm.IsDev() {
			ll.Warn("DEVELOPMENT. Skip connecting to GHN")
		} else {
			ll.Fatal("GHN: No token")
		}
	}
	if cfg.GHTK.AccountDefault.Token != "" {
		ghtkCarrier = ghtk.New(cfg.GHTK, locationBus)
		if err := ghtkCarrier.(*ghtk.Carrier).InitAllClients(ctx); err != nil {
			ll.Fatal("Unable to connect to GHTK", l.Error(err))
		}
		ll.S.Info("GHTK: connect success")
	} else {
		ll.Fatal("GHTK: No token")
	}

	svr := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: mux,
	}
	ll.S.Infof("HTTP server listening at %v", cfg.HTTP.Address())

	go func() {
		defer ctxCancel()
		err := svr.ListenAndServe()
		if err != http.ErrServerClosed {
			ll.Error("HTTP server", l.Error(err))
		}
		ll.Sync()
	}()
	return svr
}
