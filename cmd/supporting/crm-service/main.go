package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"etop.vn/backend/cmd/supporting/crm-service/config"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/health"
	"etop.vn/backend/pkg/common/metrics"
	"etop.vn/backend/pkg/common/redis"
	"etop.vn/backend/pkg/etop/authorize/middleware"
	"etop.vn/backend/pkg/etop/authorize/tokens"
	servicecrm "etop.vn/backend/pkg/services/crm-service/service"
	wrapcrm "etop.vn/backend/wrapper/services/crmservice"
	"etop.vn/common/l"
)

var (
	ll  = l.New()
	cfg config.Config
	ctx context.Context

	ctxCancel     context.CancelFunc
	healthservice = health.New()
)

func main() {
	cc.InitFlags()
	flag.Parse()

	var err error
	cfg, err = config.Load()
	if err != nil {
		ll.Fatal("error while loading config", l.Error(err))
	}

	cm.SetEnvironment(cfg.Env)
	if cm.IsDev() {
		ll.Info("config", l.Object("cfg", cfg))
	}

	ctx, ctxCancel = context.WithCancel(context.Background())
	go func() {
		osSignal := make(chan os.Signal, 1)
		signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
		ll.Info("Received OS signal", l.Stringer("signal", <-osSignal))
		ctxCancel()

		// Wait for maximum 15s
		timer := time.NewTimer(15 * time.Second)
		<-timer.C
		ll.Fatal("Force shutdown due to timeout!")
	}()

	redisStore := redis.Connect(cfg.Redis.ConnectionString())
	tokens.Init(redisStore)

	db, err := cmsql.Connect(cfg.Postgres)
	if err != nil {
		ll.Fatal("error while connecting to postgres", l.Error(err))
	}

	configMap, err := config.ReadMappingFile(cfg.MappingFile)
	if err != nil {
		ll.Fatal("error while reading field map file", l.String("file", cfg.MappingFile), l.Error(err))
	}

	s := servicecrm.NewService(db, cfg.Vtiger, configMap)
	s.Register()

	apiMux := http.NewServeMux()
	wrapcrm.NewCrmserviceServer(apiMux, nil)

	mux := http.NewServeMux()
	mux.Handle("/api/", middleware.ForwardHeaders(apiMux))
	svr := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: mux,
	}

	metrics.RegisterHTTPHandler(mux)
	healthservice.RegisterHTTPHandler(mux)
	healthservice.MarkReady()
	go func() {
		defer ctxCancel()
		err := svr.ListenAndServe()
		if err != http.ErrServerClosed {
			ll.Error("HTTP server", l.Error(err))
		}
		ll.Sync()
	}()

	<-ctx.Done()
	_ = svr.Shutdown(context.Background())
	ll.Info("Waiting for all requests to finish")
	ll.Info("Gracefully stopped!")
}
