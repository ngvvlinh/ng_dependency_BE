package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"o.o/backend/cmd/supporting/etl/config"
	identitymodel "o.o/backend/com/main/identity/model"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/metrics"
	"o.o/backend/pkg/common/sql/cmsql"
	etl2 "o.o/backend/zexp/etl"
	convert "o.o/backend/zexp/etl/tests/convert"
	"o.o/backend/zexp/etl/tests/models"
	"o.o/common/l"
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

	cmenv.SetEnvironment(cfg.Env)
	if cmenv.IsDev() {
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

	db, err := cmsql.Connect(cfg.Postgres)
	if err != nil {
		ll.Fatal("error while connecting to postgres", l.Error(err))
	}

	dbTest, err := cmsql.Connect(cfg.PostgresTest)
	if err != nil {
		ll.Fatal("error while connecting to postgres", l.Error(err))
	}

	etl := etl2.NewETLEngine(nil)
	etl.Register(db, (*identitymodel.Accounts)(nil), dbTest, (*models.Accounts)(nil))
	etl.RegisterConversion(convert.RegisterConversions)

	etl.Run()

	mux := http.NewServeMux()
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
