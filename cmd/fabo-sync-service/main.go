package main

import (
	"context"
	"time"

	"o.o/backend/cmd/fabo-sync-service/build"
	"o.o/backend/cmd/fabo-sync-service/config"
	cm "o.o/backend/pkg/common"
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
	cmenv.SetEnvironment("fabo-sync-service", cfg.Env)
	ll.Info("service starting", l.String("commit", cm.CommitMessage()))
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

	// start sync
	ll.Must(output.Sync.Init(), "can not init synchronizer")
	go func() { defer cm.RecoverAndLog(); output.Sync.Start() }()

	// start servers
	cancelHTTP := lifecycle.StartHTTP(ctxCancel, output.Servers...)
	sdCtx.Register(cancelHTTP)
	sdCtx.Register(cancelServer)

	defer output.Health.Shutdown()
	defer sdCtx.Wait()
	output.Health.MarkReady()
}
