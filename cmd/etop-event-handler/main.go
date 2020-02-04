package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shopify/sarama"

	"etop.vn/backend/cmd/etop-event-handler/config"
	haravanhandler "etop.vn/backend/com/external/haravan/handler"
	handler "etop.vn/backend/com/handler/etop-handler"
	handlerapi "etop.vn/backend/com/handler/etop-handler/api"
	"etop.vn/backend/com/handler/etop-handler/intctl"
	webhooksender "etop.vn/backend/com/handler/etop-handler/webhook/sender"
	"etop.vn/backend/com/handler/etop-handler/webhook/storage"
	catalogquery "etop.vn/backend/com/main/catalog/query"
	inventoryquery "etop.vn/backend/com/main/inventory/query"
	servicelocation "etop.vn/backend/com/main/location"
	customerquery "etop.vn/backend/com/shopping/customering/query"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/health"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmenv"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/extservice/telebot"
	"etop.vn/backend/pkg/common/headers"
	"etop.vn/backend/pkg/common/metrics"
	"etop.vn/backend/pkg/common/mq"
	"etop.vn/backend/pkg/common/redis"
	"etop.vn/backend/pkg/common/sql/cmsql"
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
	cc.ParseFlags()

	var err error
	cfg, err = config.Load()
	if err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}

	cmenv.SetEnvironment(cfg.Env)
	ll.Info("Service started with config", l.String("commit", cm.CommitMessage()))
	if cmenv.IsDev() {
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

	if cmenv.IsDev() {
		ll.Warn("DEVELOPMENT MODE ENABLED")
	}

	bot, err = cfg.TelegramBot.ConnectDefault()
	if err != nil {
		ll.Fatal("Unable to connect to Telegram", l.Error(err))
	}

	redisStore := redis.Connect(cfg.Redis.ConnectionString())
	db, err := cmsql.Connect(cfg.Postgres)
	if err != nil {
		ll.Fatal("Unable to connect to Postgres", l.Error(err))
	}
	sqlstore.Init(db)

	dbWebhook, err := cmsql.Connect(cfg.PostgresWebhook)
	if err != nil {
		ll.Fatal("Unable to connect to Postgres (webhook)", l.Error(err))
	}
	changesStore := storage.NewChangesStore(dbWebhook)

	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	catalogQuery := catalogquery.New(db).MessageBus()
	customerQuery := customerquery.NewCustomerQuery(db).MessageBus()
	inventoryquery := inventoryquery.NewQueryInventory(bus.New(), db).MessageBus()
	locationBus := servicelocation.New().MessageBus()
	addressQuery := customerquery.NewAddressQuery(db).MessageBus()

	var intctlHandler *intctl.Handler
	var webhookSender *webhooksender.WebhookSender
	var waiters []interface{ Wait() }
	{
		// intctl handler
		consumer, err := mq.NewKafkaConsumer(cfg.Kafka.Brokers, intctl.ConsumerGroup)
		if err != nil {
			ll.Fatal("Unable to connect to Kafka", l.Error(err))
		}

		intctlHandler = intctl.New(bot, consumer, cfg.Kafka.TopicPrefix)
		waiters = append(waiters, intctlHandler)
	}
	{
		// webhook handlers
		consumer, err := mq.NewKafkaConsumer(cfg.Kafka.Brokers, handler.ConsumerGroup, kafkaCfg)
		if err != nil {
			ll.Fatal("Unable to connect to Kafka", l.Error(err))
		}

		webhookSender = webhooksender.New(db, redisStore, changesStore)
		waiters = append(waiters, webhookSender)
		if err := webhookSender.Load(); err != nil {
			ll.Fatal("Error loading webhooks", l.Error(err))
		}

		h := handler.New(db, webhookSender, bot, consumer, cfg.Kafka.TopicPrefix, catalogQuery, customerQuery, inventoryquery, addressQuery, locationBus)
		h.RegisterTo(intctlHandler)
		h.ConsumeAndHandleAllTopics(ctx)
		waiters = append(waiters, h)
	}
	{
		// Haravan carrier service synced status
		consumer, err := mq.NewKafkaConsumer(cfg.Kafka.Brokers, haravanhandler.ConsumerGroup)
		if err != nil {
			ll.Fatal("Unable to connect to Kafka", l.Error(err))
		}
		h := haravanhandler.New(db, bot, consumer, cfg.Kafka.TopicPrefix)
		h.ConsumeAndHandleAllTopics(ctx)
	}

	intctlHandler.ConsumeAndHandle(ctx)
	webhookSender.Start(ctx)

	handlerapi.Init(webhookSender)

	apiMux := http.NewServeMux()
	apiMux.Handle("/api/", http.NotFoundHandler())
	handlerapi.NewHandlerServer(apiMux, cfg.Secret)

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

	if bot != nil {
		bot.SendMessage("–––\n✨ etop-handler started ✨\n" + cm.CommitMessage())
		defer bot.SendMessage("👹 etop-handler stopped 👹\n–––")
	}

	// Wait for OS signal or any error from services
	<-ctx.Done()
	_ = svr.Shutdown(context.Background())
	ll.Info("Waiting for all requests to finish")

	// Graceful stop
	for _, h := range waiters {
		h.Wait()
	}
	ll.Info("Gracefully stopped!")
}
