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

	"o.o/backend/cmd/etop-event-handler/config"
	pgeventapi "o.o/backend/cmd/pgevent-forwarder/api"
	"o.o/backend/com/eventhandler"
	etophandler "o.o/backend/com/eventhandler/etop/handler"
	fabohandler "o.o/backend/com/eventhandler/fabo/handler"
	handler "o.o/backend/com/eventhandler/handler"
	handlerapi "o.o/backend/com/eventhandler/handler/api"
	"o.o/backend/com/eventhandler/handler/intctl"
	"o.o/backend/com/eventhandler/pgevent"
	webhooksender "o.o/backend/com/eventhandler/webhook/sender"
	"o.o/backend/com/eventhandler/webhook/storage"
	servicefbmessaging "o.o/backend/com/fabo/main/fbmessaging"
	servicefbpage "o.o/backend/com/fabo/main/fbpage"
	servicefbuser "o.o/backend/com/fabo/main/fbuser"
	catalogquery "o.o/backend/com/main/catalog/query"
	serviceidentity "o.o/backend/com/main/identity"
	inventoryquery "o.o/backend/com/main/inventory/query"
	servicelocation "o.o/backend/com/main/location"
	stocktakequery "o.o/backend/com/main/stocktaking/query"
	customerquery "o.o/backend/com/shopping/customering/query"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/common/metrics"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/model"
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

	cfg.TelegramBot.MustRegister()

	redisStore := redis.ConnectWithStr(cfg.Redis.ConnectionString())
	db, err := cmsql.Connect(cfg.Postgres)
	if err != nil {
		ll.Fatal("Unable to connect to Postgres", l.Error(err))
	}
	sqlstore.New(db, nil, servicelocation.QueryMessageBus(servicelocation.New(nil)), nil)

	dbWebhook, err := cmsql.Connect(cfg.PostgresWebhook)
	if err != nil {
		ll.Fatal("Unable to connect to Postgres (webhook)", l.Error(err))
	}
	changesStore := storage.NewChangesStore(dbWebhook)

	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	catalogQuery := catalogquery.QueryServiceMessageBus(catalogquery.New(db))
	customerQuery := customerquery.CustomerQueryMessageBus(customerquery.NewCustomerQuery(db))
	stocktakeQuery := stocktakequery.StocktakeQueryMessageBus(stocktakequery.NewQueryStocktake(db))
	inventoryQuery := inventoryquery.InventoryQueryServiceMessageBus(inventoryquery.NewQueryInventory(stocktakeQuery, bus.New(), db))
	locationBus := servicelocation.QueryMessageBus(servicelocation.New(db))
	addressQuery := customerquery.AddressQueryMessageBus(customerquery.NewAddressQuery(db))

	var intctlHandler *intctl.Handler
	var webhookSender *webhooksender.WebhookSender
	var waiters []interface{ Wait() }
	{
		// intctl handlerpkg
		consumer, err := mq.NewKafkaConsumer(cfg.Kafka.Brokers, intctl.ConsumerGroup)
		if err != nil {
			ll.Fatal("Unable to connect to Kafka", l.Error(err))
		}

		intctlHandler = intctl.New(consumer, cfg.Kafka.TopicPrefix)
		waiters = append(waiters, intctlHandler)
	}
	{
		// webhook handlers
		consumer, err := mq.NewKafkaConsumer(cfg.Kafka.Brokers, etophandler.ConsumerGroup, kafkaCfg)
		if err != nil {
			ll.Fatal("Unable to connect to Kafka", l.Error(err))
		}

		webhookSender = webhooksender.New(db, redisStore, changesStore)
		waiters = append(waiters, webhookSender)
		if err := webhookSender.Load(); err != nil {
			ll.Fatal("Error loading webhooks", l.Error(err))
		}
		producer, err := mq.NewKafkaProducer(ctx, cfg.Kafka.Brokers)
		if err != nil {
			ll.Fatal("Error while connecting to Kafka", l.Error(err))
		}

		// TODO(vu): remove dependence on pgevent
		topics := []eventhandler.TopicDef{}
		topics = append(topics, etophandler.Topics()...)
		topics = append(topics, fabohandler.Topics()...)
		sMain, err := pgevent.NewService(ctx, model.DBMain, cfg.Postgres, producer, cfg.Kafka.TopicPrefix, topics)
		if err != nil {
			ll.Fatal("Error while listening to Postgres")
		}

		fbMessagingQuery := servicefbmessaging.FbMessagingQueryMessageBus(servicefbmessaging.NewFbMessagingQuery(db))
		fbPageQuery := servicefbpage.FbPageQueryMessageBus(servicefbpage.NewFbPageQuery(db))
		fbUserQuery := servicefbuser.FbUserQueryMessageBus(servicefbuser.NewFbUserQuery(db, customerQuery))
		identityQuery := serviceidentity.QueryServiceMessageBus(serviceidentity.NewQueryService(db))

		pgeventapi.Init(&sMain)
		faboHandler := fabohandler.New(db, consumer, producer, cfg.Kafka.TopicPrefix, fbUserQuery, fbMessagingQuery, fbPageQuery, identityQuery)
		etopHandler := etophandler.New(db, webhookSender, catalogQuery, customerQuery, inventoryQuery, addressQuery, locationBus)
		etopHandler.RegisterTo(intctlHandler)

		h := handler.New(consumer, cfg.Kafka)
		h.StartConsuming(ctx, etophandler.Topics(), etopHandler.TopicsAndHandlers())
		h.StartConsuming(ctx, fabohandler.Topics(), faboHandler.TopicsAndHandlers())
		waiters = append(waiters, h)
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

	ll.SendMessage(fmt.Sprintf("–––\n✨ etop-handler started on %v✨\n%v", cmenv.Env(), cm.CommitMessage()))
	defer ll.SendMessage("👹 etop-handler stopped 👹\n–––")

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
