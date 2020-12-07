package main

import (
	"context"
	"time"

	"o.o/backend/cmd/telecom-sync-service/build"
	"o.o/backend/cmd/telecom-sync-service/config"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/lifecycle"
	"o.o/backend/pkg/common/sql/sqltrace"
	"o.o/common/l"
)

var ll = l.New()

func main() {
	cc.InitFlags()
	cc.ParseFlags()

	// load config
	cfg, err := config.Load()
	ll.Must(err, "can not load config")

	cmenv.SetEnvironment("telecom-sync-service", cfg.Env)
	if cmenv.IsDev() {
		ll.Info("config", l.Object("cfg", cfg))
	}

	sqltrace.Init()
	wl.Init(cmenv.Env(), wl.EtopServer) // TODO(vu): refactor

	// lifecycle
	sdCtx, ctxCancel := lifecycle.WithCancel(context.Background())
	lifecycle.ListenForSignal(ctxCancel, 30*time.Second)
	cfg.TelegramBot.MustRegister(sdCtx)

	// build servers
	output, cancelServer, err := build.Build(sdCtx, cfg)
	ll.Must(err, "can not build server")

	// start servers
	cancelHTTP := lifecycle.StartHTTP(ctxCancel, output.Servers...)
	sdCtx.Register(cancelHTTP)
	sdCtx.Register(cancelServer)

	defer output.Health.Shutdown()
	defer sdCtx.Wait()
	output.Health.MarkReady()
}
