package main

import (
	"context"
	"time"

	"o.o/backend/cmd/etop-event-handler/build"
	"o.o/backend/cmd/etop-event-handler/config"
	notihandler "o.o/backend/com/eventhandler/notifier/handler"
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
	cmenv.SetEnvironment("event-handler", cfg.Env)
	cm.SetMainSiteBaseURL(cfg.URL.MainSite)

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

	ll.Must(output.WhSender.Load(), "can not load webhook")
	output.WhSender.Start(sdCtx)
	notihandler.Init(output.Notifier)

	// start servers
	cancelHTTP := lifecycle.StartHTTP(ctxCancel, output.Servers...)
	sdCtx.Register(cancelHTTP)
	sdCtx.Register(cancelServer)
	for _, w := range output.Waiters {
		sdCtx.Register(w.Wait)
	}

	// intctrl consuming
	output.IntctlHandler.ConsumeAndHandle(sdCtx)

	defer output.Health.Shutdown()
	defer sdCtx.Wait()
	output.Health.MarkReady()
}
