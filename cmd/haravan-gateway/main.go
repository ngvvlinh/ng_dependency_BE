package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"etop.vn/backend/pkg/integration/ghtk"
	"etop.vn/backend/pkg/integration/vtpost"

	"etop.vn/api/main/location"

	"etop.vn/backend/pkg/integration/ghn"

	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/etop/sqlstore"

	"etop.vn/backend/cmd/haravan-gateway/config"
	cm "etop.vn/backend/pkg/common"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/health"
	"etop.vn/backend/pkg/common/telebot"
	servicelocation "etop.vn/backend/pkg/services/location"
	"etop.vn/common/l"
)

var (
	ll            = l.New()
	cfg           config.Config
	ctx           context.Context
	bot           *telebot.Channel
	db            cmsql.Database
	ghnCarrier    *ghn.Carrier
	ghtkCarrier   *ghtk.Carrier
	vtpostCarrier *vtpost.Carrier
	locationBus   location.QueryBus

	ctxCancel     context.CancelFunc
	healthservice = health.New()
)

func main() {
	cc.InitFlags()
	flag.Parse()

	var err error
	cfg, err = config.Load()
	if err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}

	cm.SetEnvironment(cfg.Env)
	cm.SetMainSiteBaseURL(cfg.URL.MainSite)
	ll.Info("Service start with config", l.String("commit", cm.Commit()))
	if cm.IsDev() {
		ll.Info("config", l.Object("cfg", cfg))
	}
	ctx, ctxCancel = context.WithCancel(context.Background())
	go func() {
		osSignal := make(chan os.Signal, 1)
		signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
		ll.Info("Receive OS signal", l.Stringer("signal", <-osSignal))
		ctxCancel()

		// Wait for maximun 15s
		timer := time.NewTicker(15 * time.Second)
		<-timer.C
		ll.Fatal("Force shutdown due to timeout!")
	}()

	if cm.IsDev() {
		ll.Warn("DEVELOPMENT MODE ENABLED")
	}

	bot, err := cfg.TelegramBot.ConnectDefault()
	if err != nil {
		ll.Fatal("Unable to connect to Telegram", l.Error(err))
	}
	db, err = cmsql.Connect(cfg.Postgres)
	if err != nil {
		ll.Fatal("Unable to connect to Postgres", l.Error(err))
	}
	sqlstore.Init(db)

	locationBus = servicelocation.New().MessageBus()
	if cfg.GHN.AccountDefault.Token != "" {
		ghnCarrier = ghn.New(cfg.GHN, locationBus)
		if err := ghnCarrier.InitAllClients(ctx); err != nil {
			ll.Fatal("Unable to connect to GHN", l.Error(err))
		}
	} else {
		if cm.IsDev() {
			ll.Warn("DEVELOPMENT. Skip connecting to GHN")
		} else {
			ll.Fatal("GHN: No token")
		}
	}
	if cfg.GHTK.AccountDefault.Token != "" {
		ghtkCarrier = ghtk.New(cfg.GHTK, locationBus)
		if err := ghtkCarrier.InitAllClients(ctx); err != nil {
			ll.Fatal("Unable to connect to GHTK", l.Error(err))
		}
	} else {
		if cm.IsDev() {
			ll.Warn("DEVELOPMENT. Skip connecting to GHTK.")
		} else {
			ll.Fatal("GHTK: No token")
		}
	}

	if cfg.VTPost.AccountDefault.Username != "" {
		vtpostCarrier = vtpost.New(cfg.VTPost, locationBus)
		if err := vtpostCarrier.InitAllClients(ctx); err != nil {
			ll.Fatal("Unable to connect to VTPost", l.Error(err))
		}
	} else {
		if cm.IsDev() {
			ll.Warn("DEVELOPMENT. Skip connecting to VTPost.")
		} else {
			ll.Fatal("VTPost: No token")
		}
	}

	svr := startServers()
	if bot != nil {
		bot.SendMessage("–––\n✨ haravan-gateway started ✨\n" + cm.Commit())
		defer bot.SendMessage("👹 haravan-gateway stopped 👹\n–––")
	}

	healthservice.MarkReady()
	// Wait for OS signal or any error from services
	<-ctx.Done()

	_ = svr.Shutdown(context.Background())
	ll.Info("Gracefully stopped!")
}
