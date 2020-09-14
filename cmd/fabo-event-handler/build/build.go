package build

import (
	"context"
	"net/http"

	"github.com/Shopify/sarama"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/api/main/identity"
	"o.o/backend/cmd/fabo-event-handler/config"
	fabohandler "o.o/backend/com/eventhandler/fabo/handler"
	"o.o/backend/com/eventhandler/handler"
	handlerapi "o.o/backend/com/eventhandler/handler/api"
	"o.o/backend/com/eventhandler/handler/intctl"
	"o.o/backend/com/eventhandler/pgevent"
	"o.o/backend/com/eventhandler/webhook/sender"
	com "o.o/backend/com/main"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/bus"
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

type Output struct {
	Servers   []lifecycle.HTTPServer
	Waiters   []lifecycle.Waiter
	PgService *pgevent.Service
	WhSender  *sender.WebhookSender
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
	s, err := pgevent.NewService(ctx, dbdecl.DBMain, dbcfg, producer, cfg.Kafka.TopicPrefix, fabohandler.Topics())
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
	fbUserQuery fbusering.QueryBus,
	fbMessagingQuery fbmessaging.QueryBus,
	fbPageQuery fbpaging.QueryBus,
	identityQuery identity.QueryBus,
) (*handler.Handler, error) {
	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := mq.NewKafkaConsumer(cfg.Kafka.Brokers, "fabo-handler", kafkaCfg)
	if err != nil {
		return nil, err
	}
	producer, err := mq.NewKafkaProducer(ctx, cfg.Kafka.Brokers)
	if err != nil {
		return nil, err
	}

	faboHandler := fabohandler.New(db, producer, cfg.Kafka.TopicPrefix, fbUserQuery, fbMessagingQuery, fbPageQuery, identityQuery)
	h := handler.New(consumer, cfg.Kafka)
	h.StartConsuming(ctx, fabohandler.Topics(), faboHandler.TopicsAndHandlers())
	return h, nil
}

func BuildWaiters(
	intctlHandler *intctl.Handler,
	h *handler.Handler,
) (waiters []lifecycle.Waiter) {
	waiters = append(waiters, intctlHandler, h)
	return waiters
}
