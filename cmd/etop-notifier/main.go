package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shopify/sarama"

	"o.o/backend/cmd/etop-notifier/config"
	"o.o/backend/com/eventhandler/etop/handler"
	"o.o/backend/com/eventhandler/notifier"
	notihandler "o.o/backend/com/eventhandler/notifier/handler"
	servicelocation "o.o/backend/com/main/location"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/common/metrics"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/common/l"
)

var (
	ll  = l.New()
	cfg config.Config
	ctx context.Context

	ctxCancel     context.CancelFunc
	healthservice = health.New()
)

func main() {
	cc.InitFlags()
	cc.ParseFlags()

	var err error
	cfg, err = config.Load()
	if err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}

	cmenv.SetEnvironment(cfg.Env)
	cm.SetMainSiteBaseURL(cfg.URL.MainSite)
	ll.Info("Service started with config", l.String("commit", cm.CommitMessage()))
	if cmenv.IsDev() {
		ll.Info("config", l.Object("cfg", cfg))
	}
	wl.Init(cmenv.Env())

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

	if cmenv.IsDev() {
		ll.Warn("DEVELOPMENT MODE ENABLED")
	}

	cfg.TelegramBot.MustRegister()

	db, err := cmsql.Connect(cfg.Postgres)
	if err != nil {
		ll.Fatal("Unable to connect to Postgres", l.Error(err))
	}

	locationBus := servicelocation.QueryMessageBus(servicelocation.New(nil))

	dbNotifier, err := cmsql.Connect(cfg.PostgresNotifier)
	if err != nil {
		ll.Fatal("Unable to connect to Postgres Notifier", l.Error(err))
	}
	sqlstore.New(db, dbNotifier, locationBus, nil)
	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	{
		consumer, err := mq.NewKafkaConsumer(cfg.Kafka.Brokers, notihandler.ConsumerGroup)
		if err != nil {
			ll.Fatal("Unable to connect to Kafka", l.Error(err))
		}
		hMain, hNotifier := notihandler.New(db, dbNotifier, consumer, cfg.Kafka)
		hMain.StartConsuming(ctx, handler.GetTopics(notihandler.TopicsAndHandlersEtop()), notihandler.TopicsAndHandlersEtop())
		hNotifier.StartConsuming(ctx, handler.GetTopics(notihandler.TopicsAndHandlerNotifier()), notihandler.TopicsAndHandlerNotifier())
	}
	{
		if cfg.Onesignal.ApiKey != "" {
			if err := notifier.Init(cfg.Onesignal); err != nil {
				ll.Fatal("Unable to connect to Onesignal", l.Error(err))
			}
		} else {
			if cmenv.IsDev() {
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
	mux.Handle("/api/", headers.ForwardHeaders(apiMux))
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

	ll.SendMessage(fmt.Sprintf("–––\n✨ etop-notification started on %v✨\n%v", cmenv.Env(), cm.CommitMessage()))
	defer ll.SendMessage("👹 etop-notification stopped 👹\n–––")

	// Wait for OS signal or any error from services
	<-ctx.Done()

	_ = svr.Shutdown(context.Background())
	ll.Info("Gracefully stopped!")
}
