package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"etop.vn/backend/cmd/shipping-sync-service/config"
	cm "etop.vn/backend/pkg/common"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/health"
	"etop.vn/backend/pkg/common/telebot"
	"etop.vn/backend/pkg/integration/ghn"
	ghnsync "etop.vn/backend/pkg/integration/ghn/sync"
	servicelocation "etop.vn/backend/pkg/services/location"
	"etop.vn/common/l"
)

var (
	ll  = l.New()
	ctx context.Context
	cfg config.Config
	bot *telebot.Channel

	ctxCancel     context.CancelFunc
	healthservice = health.New()
)

func main() {
	cc.InitFlags()
	flag.Parse()

	var err error
	cfg, err = config.Load()
	if err != nil {
		ll.Fatal("Unable to load config", l.Error(err))
	}

	cm.SetEnvironment(cfg.Env)
	ll.Info("Service started with config", l.String("commit", cm.Commit()))
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
		if bot != nil {
			bot.SendMessage("👻 shipping-sync-service stopped (forced) 👻\n–––")
		}
		ll.Fatal("Force shutdown due to timeout!")
	}()

	bot, err = cfg.TelegramBot.ConnectDefault()
	if err != nil {
		ll.Fatal("Unable to connect to Telegram", l.Error(err))
	}

	svrs := startServers()
	if bot != nil {
		bot.SendMessage("–––\n✨ shipping-sync-service started ✨\n" + cm.Commit())
		defer bot.SendMessage("👻 shipping-sync-service stopped 👻\n–––")
	}

	locationBus := servicelocation.New().MessageBus()
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
