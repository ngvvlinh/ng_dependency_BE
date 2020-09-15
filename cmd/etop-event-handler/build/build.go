package build

import (
	"context"
	"net/http"

	"github.com/Shopify/sarama"

	"o.o/api/main/catalog"
	"o.o/api/main/inventory"
	"o.o/api/main/location"
	"o.o/api/main/shipnow"
	"o.o/api/shopping/addressing"
	"o.o/api/shopping/customering"
	"o.o/backend/cmd/etop-event-handler/config"
	etophandler "o.o/backend/com/eventhandler/etop/handler"
	"o.o/backend/com/eventhandler/handler"
	handlerapi "o.o/backend/com/eventhandler/handler/api"
	"o.o/backend/com/eventhandler/handler/intctl"
	"o.o/backend/com/eventhandler/notifier"
	notihandler "o.o/backend/com/eventhandler/notifier/handler"
	"o.o/backend/com/eventhandler/pgevent"
	"o.o/backend/com/eventhandler/webhook/sender"
	com "o.o/backend/com/main"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/common/lifecycle"
	"o.o/backend/pkg/common/metrics"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/common/sql/sqltrace"
	"o.o/backend/pkg/etc/dbdecl"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/middlewares"
	"o.o/backend/tools/pkg/acl"
	"o.o/capi/httprpc"
	"o.o/common/jsonx"
	"o.o/common/l"
)

var ll = l.New()

type Output struct {
	Servers   []lifecycle.HTTPServer
	Waiters   []lifecycle.Waiter
	PgService *pgevent.Service
	WhSender  *sender.WebhookSender
	Notifier  *notifier.Notifier
	Handlers  []*handler.Handler
	Health    *health.Service
}

func BuildServers(
	main MainServer,
) []lifecycle.HTTPServer {
	svrs := []lifecycle.HTTPServer{
		{"Main", main},
	}
	return svrs
}

type MainServer *http.Server

func BuildMainServer(
	cfg config.Config,
	healthService *health.Service,
	handlerServers handlerapi.Servers,
) (MainServer, error) {
	mux := http.NewServeMux()
	l.RegisterHTTPHandler(mux)
	metrics.RegisterHTTPHandler(mux)
	healthService.RegisterHTTPHandler(mux)
	jsonx.RegisterHTTPHandler(mux)
	sqltrace.RegisterHTTPHandler(mux)

	mwares := httpx.Compose(
		headers.ForwardHeadersX(),
		bus.Middleware,
	)
	mux.Handle("/api/", http.NotFoundHandler())

	logging := middlewares.NewLogging()
	ssHooks, err := session.NewHook(acl.GetACL(), session.OptSecret(cfg.Secret))
	if err != nil {
		return nil, err
	}
	handlers := httprpc.WithHooks(handlerServers, ssHooks, logging)
	for _, h := range handlers {
		mux.Handle(h.PathPrefix(), mwares(h))
	}

	svr := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: mux,
	}
	return svr, nil
}

func BuildPgEventService(
	ctx context.Context,
	cfg config.Config,
) (*pgevent.Service, error) {
	producer, err := mq.NewKafkaProducer(ctx, cfg.Kafka.Brokers)
	if err != nil {
		return nil, err
	}
	dbcfg := *cfg.Databases["postgres"]
	s, err := pgevent.NewService(ctx, dbdecl.DBMain, dbcfg, producer, cfg.Kafka.TopicPrefix, etophandler.Topics())
	return s, err
}

func BuildIntHandler(
	ctx context.Context,
	cfg config.Config,
) (*intctl.Handler, error) {
	consumer, err := mq.NewKafkaConsumer(cfg.Kafka.Brokers, intctl.ConsumerGroup)
	if err != nil {
		return nil, err
	}
	h := intctl.New(consumer, cfg.Kafka.TopicPrefix)
	h.ConsumeAndHandle(ctx)
	return h, nil
}

func BuildWebhookHandler(
	ctx context.Context,
	cfg config.Config,
	db com.MainDB,
	webhookSender *sender.WebhookSender,
	catalogQuery catalog.QueryBus,
	customerQuery customering.QueryBus,
	inventoryQuery inventory.QueryBus,
	addressQuery addressing.QueryBus,
	locationQuery location.QueryBus,
	shipnowQuery shipnow.QueryBus,
) (*handler.Handler, error) {
	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := mq.NewKafkaConsumer(cfg.Kafka.Brokers, etophandler.ConsumerGroup, kafkaCfg)
	if err != nil {
		return nil, err
	}

	etopHandler := etophandler.New(db, webhookSender, catalogQuery, customerQuery, inventoryQuery, addressQuery, locationQuery, shipnowQuery)
	h := handler.New(consumer, cfg.Kafka)
	h.StartConsuming(ctx, etophandler.Topics(), etopHandler.TopicsAndHandlers())
	return h, nil
}

func BuildWaiters(
	intctlHandler *intctl.Handler,
	h *handler.Handler,
) (waiters []lifecycle.Waiter) {
	waiters = append(waiters, intctlHandler, h)
	return waiters
}

func BuildOneSignal(cfg cc.OnesignalConfig) (*notifier.Notifier, error) {
	if cfg.ApiKey != "" {
		return notifier.NewOneSignalNotifier(cfg)
	}

	if cmenv.IsDev() {
		ll.Warn("DEVELOPMENT. Skip connect to Onesignal")
	} else {
		ll.Fatal("Onesignal: No apikey")
	}
	return nil, nil
}

func BuildHandlers(
	ctx context.Context,
	cfg config.Config,
	db com.MainDB,
	notifierDB com.NotifierDB,
) ([]*handler.Handler, error) {
	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := mq.NewKafkaConsumer(cfg.Kafka.Brokers, etophandler.ConsumerGroup, kafkaCfg)
	if err != nil {
		return nil, err
	}

	hMain, hNotifier := notihandler.New(db, notifierDB, consumer, cfg.Kafka)
	hMain.StartConsuming(ctx, etophandler.GetTopics(notihandler.TopicsAndHandlersEtop()), notihandler.TopicsAndHandlersEtop())
	hNotifier.StartConsuming(ctx, notihandler.GetTopics(notihandler.TopicsAndHandlerNotifier()), notihandler.TopicsAndHandlerNotifier())
	return []*handler.Handler{hMain, hNotifier}, nil
}
