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

	cm "etop.vn/backend/pkg/common"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/health"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/common/telebot"
	"etop.vn/backend/pkg/integration/kiotviet"
	"etop.vn/backend/pkg/integration/kiotviet/cmd/kiotviet-sync-service/config"
)

var (
	flWebhookURL = flag.String("webhook", "", "Expose webhook url")

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

	if cm.IsDev() {
		ll.Warn("DEVELOPMENT MODE ENABLED")
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
			bot.SendMessage("👻 kiotviet-sync-service stopped (forced) 👻\n–––")
		}
		ll.Fatal("Force shutdown due to timeout!")
	}()

	bot, err = cfg.TelegramBot.ConnectDefault()
	if err != nil {
		ll.Fatal("Unable to connect to Telegram", l.Error(err))
	}

	kiotviet.Init(cfg.Sync)
	svrs := startServers()
	if bot != nil {
		bot.SendMessage("–––\n✨ kiotviet-sync-service started ✨\n" + cm.Commit())
		defer bot.SendMessage("👻 kiotviet-sync-service stopped 👻\n–––")
	}

	ll.Info("Start Kiotviet sync job")
	kiotviet.StartSync()

	healthservice.MarkReady()

	// Wait for OS signal or any error from services
	<-ctx.Done()
	ll.Info("Waiting for all requests to finish")

	// Graceful stop
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		kiotviet.StopSync()
		ll.Info("Stop sync")
	}()

	wg.Add(len(svrs))
	for _, svr := range svrs {
		go func(svr *http.Server) {
			defer wg.Done()
			svr.Shutdown(context.Background())
		}(svr)
	}
	wg.Wait()
	ll.Info("Gracefully stopped!")
}
