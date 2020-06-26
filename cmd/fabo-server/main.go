package main

import (
	"context"
	"time"

	"github.com/Shopify/sarama"

	"o.o/backend/cmd/fabo-server/build"
	"o.o/backend/cmd/fabo-server/config"
	fabopublisher "o.o/backend/com/eventhandler/fabo/publisher"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/lifecycle"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/common/sql/sqltrace"
	"o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/authorize/authfabo"
	"o.o/backend/pkg/etop/model"
	"o.o/common/l"
)

var ll = l.New()

func main() {
	cc.InitFlags()
	cc.ParseFlags()

	// load config
	cfg, err := config.Load()
	ll.Must(err, "can not load config")
	cmenv.SetEnvironment(cfg.SharedConfig.Env)
	if cmenv.IsDev() {
		cfg.SMS.Telegram = true
		cfg.SMS.Enabled = true
		ll.Info("config", l.Object("cfg", cfg))
	}

	cm.SetMainSiteBaseURL(cfg.URL.MainSite) // TODO(vu): refactor
	sqltrace.Init()
	wl.Init(cmenv.Env(), wl.FaboServer)
	auth.Init(authfabo.Policy)
	cfg.TelegramBot.MustRegister()
	eventBus := bus.New()
	healthService := health.New()

	// TODO(vu): refactor
	model.GetShippingServiceRegistry().Initialize()

	// lifecycle
	sdCtx, ctxCancel := lifecycle.WithCancel(context.Background())
	defer sdCtx.Wait()
	lifecycle.ListenForSignal(ctxCancel, 30*time.Second)

	// kafka
	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := mq.NewKafkaConsumer(cfg.Kafka.Brokers, "handler/fabo-server", kafkaCfg)
	if err != nil {
		ll.Fatal("can not connect to Kafka", l.Error(err))
	}

	// build servers
	output, cancel, err := build.Build(sdCtx, cfg, eventBus, healthService, consumer)
	ll.Must(err, "can not build server")

	// start forwarder
	go output.EventStream.RunForwarder()
	h, fp := output.Handler, output.Publisher
	h.StartConsuming(sdCtx, fabopublisher.GetTopics(fp.TopicsAndHandlers()), fp.TopicsAndHandlers())

	// start servers
	cancelHTTP := lifecycle.StartHTTP(ctxCancel, output.Servers...)
	sdCtx.Register(cancelHTTP)
	sdCtx.Register(cancel)
	sdCtx.Register(func() { ll.SendMessagef("🎃 fabo-server on %v stopped 🎃", cmenv.Env()) })
	healthService.MarkReady()

	ll.SendMessagef("✨ fabo-server on %v started ✨\n%v", cmenv.Env(), cm.CommitMessage())
}
