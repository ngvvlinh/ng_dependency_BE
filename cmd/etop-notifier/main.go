package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shopify/sarama"

	"etop.vn/backend/cmd/etop-notifier/config"
	"etop.vn/backend/com/handler/notifier"
	notihandler "etop.vn/backend/com/handler/notifier/handler"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/health"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/extservice/telebot"
	"etop.vn/backend/pkg/common/metrics"
	"etop.vn/backend/pkg/common/mq"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/etop/authorize/middleware"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/common/l"
)

var (
	ll  = l.New()
	cfg config.Config
	ctx context.Context
	bot *telebot.Channel

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
	ll.Info("Service started with config", l.String("commit", cm.CommitMessage()))
	if cm.IsDev() {
		ll.Info("config", l.Object("cfg", cfg))
	}

	ctx, ctxCancel = context.WithCancel(context.Background())
	go func() {
		osSignal := make(chan os.Signal, 1)
		signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
		ll.Info("Received OS signal", l.Stringer("signal", <-osSignal))
		ctxCancel()

		// Wait for maximum 15s
		timer := time.NewTimer(15 * time.Second)
		<-timer.C
		ll.Fatal("Force shutdown due to timeout!")
	}()

	if cm.IsDev() {
		ll.Warn("DEVELOPMENT MODE ENABLED")
	}

	bot, err = cfg.TelegramBot.ConnectDefault()
	if err != nil {
		ll.Fatal("Unable to connect to Telegram", l.Error(err))
	}
	db, err := cmsql.Connect(cfg.Postgres)
	if err != nil {
		ll.Fatal("Unable to connect to Postgres", l.Error(err))
	}
	sqlstore.Init(db)

	dbNotifier, err := cmsql.Connect(cfg.PostgresNotifier)
	if err != nil {
		ll.Fatal("Unable to connect to Postgres Notifier", l.Error(err))
	}
	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	{
		consumer, err := mq.NewKafkaConsumer(cfg.Kafka.Brokers, notihandler.ConsumerGroup)
		if err != nil {
			ll.Fatal("Unable to connect to Kafka", l.Error(err))
		}
		hMain, hNotifier := notihandler.New(db, dbNotifier, bot, consumer, cfg.Kafka.TopicPrefix)
		hMain.ConsumeAndHandleAllTopics(ctx)
		hNotifier.ConsumeAndHandleAllTopics(ctx)
	}
	{
		if cfg.Onesignal.ApiKey != "" {
			if err := notifier.Init(db, cfg.Onesignal); err != nil {
				ll.Fatal("Unable to connect to Onesignal", l.Error(err))
			}
		} else {
			if cm.IsDev() {
				ll.Warn("DEVELOPMENT. Skip connect to Onesignal")
			} else {
				ll.Fatal("Onesignal: No apikey")
			}
		}
	}

	apiMux := http.NewServeMux()
	apiMux.Handle("/api/", http.NotFoundHandler())
	// wraphandler.NewHandlerServer(apiMux, nil, cfg.Secret)

	mux := http.NewServeMux()
	mux.Handle("/api/", middleware.ForwardHeaders(apiMux))
	svr := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: mux,
	}

	l.RegisterHTTPHandler(mux)
	metrics.RegisterHTTPHandler(mux)
	healthservice.RegisterHTTPHandler(mux)
	healthservice.MarkReady()
	go func() {
		defer ctxCancel()
		err := svr.ListenAndServe()
		if err != http.ErrServerClosed {
			ll.Error("HTTP server", l.Error(err))
		}
		ll.Sync()
	}()

	if bot != nil {
		bot.SendMessage("–––\n✨ etop-notification started ✨\n" + cm.CommitMessage())
		defer bot.SendMessage("👹 etop-notification stopped 👹\n–––")
	}

	// Wait for OS signal or any error from services
	<-ctx.Done()

	_ = svr.Shutdown(context.Background())
	ll.Info("Gracefully stopped!")
}
