package main

import (
	"context"
	"time"

	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/lifecycle"
	"o.o/backend/zexp/sample/calc3/build"
	"o.o/backend/zexp/sample/calc3/config"
	"o.o/common/l"
)

var ll = l.New()

func main() {
	cc.InitFlags()
	cc.ParseFlags()

	// load config
	cfg, err := config.Load()
	ll.Must(err, "can not load config")
	cmenv.SetEnvironment(cfg.Env)
	if cmenv.IsDev() {
		ll.Info("config", l.Object("cfg", cfg))
	}

	sd, ctxCancel := lifecycle.WithCancel(context.Background())
	defer sd.Wait()
	lifecycle.ListenForSignal(ctxCancel, 30*time.Second)

	server, err := build.Build(cfg)
	cancelHTTP := lifecycle.StartHTTP(ctxCancel, server)
	sd.Register(cancelHTTP)
}