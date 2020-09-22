package main

import (
	"context"
	"time"

	"o.o/backend/cmd/mc-hub/build"
	"o.o/backend/cmd/mc-hub/config"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/lifecycle"
	"o.o/common/l"
)

var ll = l.New()

func main() {
	cc.InitFlags()
	cc.ParseFlags()
	cfg, err := config.Load()
	ll.Must(err, "can not load config")
	cmenv.SetEnvironment("mc-hub", cfg.Env)
	if cmenv.IsDev() {
		ll.Info("config", l.Object("cfg", cfg))
	}

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
