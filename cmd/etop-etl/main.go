package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"o.o/backend/cmd/etop-etl/config"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/redis"
	etlutil "o.o/backend/zexp/etl/util"
	"o.o/common/l"
)

var ll = l.New()

func main() {
	cc.InitFlags()
	resetDB := flag.Bool("reset-db", false, "drop all tables (only dev)")
	cc.ParseFlags()

	cfg, err := config.Load()
	if err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}

	cmenv.SetEnvironment("etop-etl", cfg.Env)
	ll.Info("Service started with config", l.String("commit", cm.CommitMessage()))
	if cmenv.IsDev() {
		ll.Info("config", l.Object("cfg", cfg))
	}

	ctx, ctxCancel := context.WithCancel(context.Background())
	go func() {
		osSignal := make(chan os.Signal, 1)
		signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
		ll.Info("Received OS signal", l.Stringer("signal", <-osSignal))
		ctxCancel()

		// Wait for maximum 15s
		timer := time.NewTicker(15 * time.Second)
		<-timer.C
		ll.Fatal("Force shutdown due to timeout!")
	}()

	if cmenv.IsDev() {
		ll.Warn("DEVELOPMENT MODE ENABLED")
	}

	cfg.TelegramBot.MustRegister(ctx)
	util := etlutil.New(cfg.MapDB, *resetDB)
	util.HandleETL(ctx)

	redisStore := redis.Connect(cfg.Redis)
	healthService := health.New(redisStore)
	healthService.MarkReady()
	defer healthService.Shutdown()

	mux := http.NewServeMux()
	healthService.RegisterHTTPHandler(mux)
	svr := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: mux,
	}
	go func() {
		defer ctxCancel()
		ll.S.Infof("HTTP server listening at %v", cfg.HTTP.Address())
		err := svr.ListenAndServe()
		if err != http.ErrServerClosed {
			ll.Error("HTTP server", l.Error(err))
		}
		ll.Sync()
	}()

	// Wait for OS signal or any error from services
	<-ctx.Done()
	ll.Info("Gracefully stopped!")

	// Graceful stop
	svr.Shutdown(context.Background())
	ll.Info("Gracefully stopped!")
}
