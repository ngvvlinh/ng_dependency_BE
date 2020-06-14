package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"o.o/backend/cmd/fabo-sync-service/config"
	servicefbmessaging "o.o/backend/com/fabo/main/fbmessaging"
	servicefbpaging "o.o/backend/com/fabo/main/fbpage"
	servicefbusering "o.o/backend/com/fabo/main/fbuser"
	"o.o/backend/com/fabo/pkg/fbclient"
	faboRedis "o.o/backend/com/fabo/pkg/redis"
	"o.o/backend/com/fabo/pkg/sync"
	customerquery "o.o/backend/com/shopping/customering/query"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/metrics"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/common/l"
)

var (
	ll  = l.New()
	ctx context.Context
	cfg config.Config

	ctxCancel     context.CancelFunc
	healthservice = health.New()
)

func main() {
	cc.InitFlags()
	cc.ParseFlags()

	var err error
	cfg, err = config.Load()
	if err != nil {
		ll.Fatal("Unable to load config", l.Error(err))
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
		ll.SendMessage("👻 shipping-sync-service stopped (forced) 👻\n–––")
		ll.Fatal("Force shutdown due to timeout!")
	}()

	cfg.TelegramBot.MustRegister()

	ll.SendMessage("–––\n✨ fabo-sync-service started ✨\n" + cm.CommitMessage())
	defer ll.SendMessage("👹 fabo-sync-service stopped 👹\n–––")

	redisStore := redis.ConnectWithStr(cfg.Redis.ConnectionString())

	db, err := cmsql.Connect(cfg.Postgres)
	if err != nil {
		ll.Fatal("Unable to connect to Postgres", l.Error(err))
	}
	eventBus := bus.New()

	fbClient := fbclient.New(cfg.FacebookApp)
	if err := fbClient.Ping(); err != nil {
		ll.Fatal("Error while connection Facebook", l.Error(err))
	}

	fbRedis := faboRedis.NewFaboRedis(redisStore)

	customerQuery := customerquery.CustomerQueryMessageBus(customerquery.NewCustomerQuery(db))
	fbPagingQuery := servicefbpaging.FbPageQueryMessageBus(servicefbpaging.NewFbPageQuery(db))
	fbPagingAggr := servicefbpaging.FbExternalPageAggregateMessageBus(servicefbpaging.NewFbPageAggregate(db))
	fbMessagingAggr := servicefbmessaging.FbExternalMessagingAggregateMessageBus(servicefbmessaging.NewFbExternalMessagingAggregate(db, eventBus))
	fbMessagingQuery := servicefbmessaging.FbMessagingQueryMessageBus(servicefbmessaging.NewFbMessagingQuery(db))
	fbUseringQuery := servicefbusering.FbUserQueryMessageBus(servicefbusering.NewFbUserQuery(db, customerQuery))
	fbUseringAggr := servicefbusering.FbUserAggregateMessageBus(servicefbusering.NewFbUserAggregate(db, fbPagingAggr, customerQuery))
	// fbMessagingPM
	_ = servicefbmessaging.NewProcessManager(
		eventBus,
		fbMessagingQuery, fbMessagingAggr,
		fbPagingQuery,
		fbUseringQuery, fbUseringAggr,
		fbRedis)
	synchronizer := sync.New(
		db,
		fbClient,
		fbMessagingAggr, fbMessagingQuery,
		fbUseringAggr, fbRedis,
		cfg.TimeLimit,
	)
	if err := synchronizer.Init(); err != nil {
		panic(err)
	}

	go func() {
		synchronizer.Start()
	}()

	mux := http.NewServeMux()
	l.RegisterHTTPHandler(mux)
	metrics.RegisterHTTPHandler(mux)
	healthservice.RegisterHTTPHandler(mux)

	svr := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: mux,
	}
	ll.S.Infof("HTTP server listening at %v", cfg.HTTP.Address())
	go func() {
		defer ctxCancel()
		err := svr.ListenAndServe()
		if err != http.ErrServerClosed {
			ll.Error("HTTP server", l.Error(err))
		}
		ll.Sync()
	}()

	healthservice.MarkReady()

	// Wait for OS signal or any error from services
	<-ctx.Done()
	ll.Info("Waiting for all requests to finish")

	ll.Info("Gracefully stopped!")
}
