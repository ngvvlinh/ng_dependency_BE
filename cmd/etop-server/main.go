package main

import (
	"context"
	"time"

	"o.o/backend/cmd/etop-server/build"
	"o.o/backend/cmd/etop-server/config"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/lifecycle"
	"o.o/backend/pkg/common/sql/sqltrace"
	"o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/authorize/authetop"
	"o.o/backend/pkg/etop/model"
	"o.o/common/l"
)

var ll = l.New()

func main() {
	cc.InitFlags()
	cc.ParseFlags()
	auth.Init(authetop.Policy)

	// load config
	cfg, err := config.Load(false)
	ll.Must(err, "can not load config")
	cmenv.SetEnvironment(cfg.SharedConfig.Env)
	if cmenv.IsDev() {
		cfg.SMS.Telegram = true
		cfg.SMS.Enabled = true
		ll.Info("config", l.Object("cfg", cfg))
	}

	cm.SetMainSiteBaseURL(cfg.URL.MainSite) // TODO(vu): refactor
	sqltrace.Init()
	wl.Init(cmenv.Env(), wl.EtopServer)
	eventBus := bus.New()
	healthService := health.New()

	// TODO(vu): refactor
	model.GetShippingServiceRegistry().Initialize()

	// lifecycle
	sdCtx, ctxCancel := lifecycle.WithCancel(context.Background())
	defer ll.SendMessagef("🎃 etop-server on %v stopped 🎃", cmenv.Env())
	defer sdCtx.Wait()
	lifecycle.ListenForSignal(ctxCancel, 30*time.Second)
	cfg.TelegramBot.MustRegister(sdCtx)

	// build servers
	output, cancelServer, err := build.Build(sdCtx, cfg, eventBus, healthService, cfg.URL.Auth)
	ll.Must(err, "can not build server")

	// start servers
	cancelHTTP := lifecycle.StartHTTP(ctxCancel, output.Servers...)
	sdCtx.Register(cancelHTTP)
	sdCtx.Register(cancelServer)
	healthService.MarkReady()

	ll.SendMessagef("✨ etop-server on %v started ✨\n%v", cmenv.Env(), cm.CommitMessage())
}
