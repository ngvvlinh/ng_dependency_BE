package main

import (
	"context"
	"time"

	"o.o/backend/cmd/etop-server/build"
	"o.o/backend/cmd/etop-server/config"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/lifecycle"
	"o.o/backend/pkg/common/sql/sqltrace"
	"o.o/backend/pkg/etop/model"
	"o.o/common/l"
)

var ll = l.New()

func main() {
	cc.InitFlags()
	cc.ParseFlags()

	// load config
	cfg, err := config.Load(false)
	ll.Must(err, "can not load config")
	cmenv.SetEnvironment("etop-server", cfg.SharedConfig.Env)
	if cmenv.IsDev() {
		cfg.SMS.Telegram = true
		cfg.SMS.Enabled = true
		ll.Info("config", l.Object("cfg", cfg))
	}

	cm.SetMainSiteBaseURL(cfg.URL.MainSite) // TODO(vu): refactor
	sqltrace.Init()
	wl.Init(cmenv.Env(), wl.EtopServer) // TODO(vu): refactor

	// TODO(vu): refactor
	model.GetShippingServiceRegistry().Initialize()

	// lifecycle
	sdCtx, ctxCancel := lifecycle.WithCancel(context.Background())
	lifecycle.ListenForSignal(ctxCancel, 30*time.Second)
	cfg.TelegramBot.MustRegister(sdCtx)

	// build servers
	output, cancelServer, err := build.Build(sdCtx, cfg, cfg.URL.Auth)
	ll.Must(err, "can not build server")

	// start servers
	cancelHTTP := lifecycle.StartHTTP(ctxCancel, output.Servers...)
	sdCtx.Register(cancelHTTP)
	sdCtx.Register(cancelServer)

	defer output.Health.Shutdown()
	defer sdCtx.Wait()
	output.Health.MarkReady()
}
