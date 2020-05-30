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

	cmenv.SetEnvironment(cfg.Env)
	cm.SetMainSiteBaseURL(cfg.URL.MainSite) // TODO(vu): refactor
	sqltrace.Init()
	wl.Init(cmenv.Env())
	eventBus := bus.New()
	healthService := health.New()

	// TODO(vu): refactor
	model.GetShippingServiceRegistry().Initialize()

	// lifecyle
	sd, ctxCancel := lifecycle.WithCancel(context.Background())
	defer sd.Wait()
	lifecycle.ListenForSignal(ctxCancel, 30*time.Second)

	// build servers
	servers, cancelServer, err := build.Servers(sd, cfg, eventBus, healthService, cfg.URL.Auth)
	ll.Must(err, "can not build server")

	// start servers
	cancelHTTP := lifecycle.StartHTTP(ctxCancel, servers...)
	sd.Register(cancelHTTP)
	sd.Register(cancelServer)
	healthService.MarkReady()
}
