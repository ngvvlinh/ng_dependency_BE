package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"o.o/backend/cmd/shipping-sync-service/config"
	servicelocation "o.o/backend/com/main/location"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/integration/shipping/ghn"
	ghnsync "o.o/backend/pkg/integration/shipping/ghn/sync"
	"o.o/common/l"
)

var (
	ll  = l.New()
	ctx context.Context
	cfg config.Config

	ctxCancel     context.CancelFunc
	healthservice = health.New()
)

func main() {
	cc.InitFlags()
	cc.ParseFlags()

	var err error
	cfg, err = config.Load()
	if err != nil {
		ll.Fatal("Unable to load config", l.Error(err))
	}

	cmenv.SetEnvironment(cfg.Env)
	ll.Info("Service started with config", l.String("commit", cm.CommitMessage()))
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

		ll.SendMessage("👻 shipping-sync-service stopped (forced) 👻\n–––")

		ll.Fatal("Force shutdown due to timeout!")
	}()

	cfg.TelegramBot.MustRegister()
	svrs := startServers()
	ll.SendMessage(fmt.Sprintf("–––\n✨ shipping-sync-service started on %v✨\n%v", cmenv.Env(), cm.CommitMessage()))
	defer ll.SendMessage("👻 shipping-sync-service stopped 👻\n–––")

	locationBus := servicelocation.QueryMessageBus(servicelocation.New(nil))
	ghnCarrier := ghn.New(cfg.GHN, locationBus)
	ghnSynchronizer := ghnsync.New(ghnCarrier)

	ll.Info("Start ghn sync order")
	go func() {
		ghnSynchronizer.Start()
	}()
	go func() {
		SyncUnCompleteFfms()
	}()

	healthservice.MarkReady()

	// Wait for OS signal or any error from services
	<-ctx.Done()
	ll.Info("Waiting for all requests to finish")

	// Graceful stop
	var wg sync.WaitGroup
	wg.Add(len(svrs))
	for _, svr := range svrs {
		go func(svr *http.Server) {
			defer wg.Done()
			ll.Info("Stop ghn sync order")
			svr.Shutdown(context.Background())
		}(svr)
	}
	wg.Wait()
	ghnSynchronizer.Stop()
	StopSync()
	ll.Info("Gracefully stopped!")
}
