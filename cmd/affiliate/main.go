package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"etop.vn/backend/pkg/etop/sqlstore"

	"etop.vn/backend/cmd/affiliate/config"
	querycatalog "etop.vn/backend/com/main/catalog/query"
	serviceidentity "etop.vn/backend/com/main/identity"
	serviceaffiliate "etop.vn/backend/com/services/affiliate"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/health"
	"etop.vn/backend/pkg/common/metrics"
	"etop.vn/backend/pkg/common/redis"
	"etop.vn/backend/pkg/etop/authorize/middleware"
	"etop.vn/backend/pkg/etop/authorize/tokens"
	apiaff "etop.vn/backend/pkg/services/affiliate/api"
	wrapaff "etop.vn/backend/wrapper/services/affiliate"
	"etop.vn/common/l"
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
	flag.Parse()

	var err error
	cfg, err = config.Load()
	if err != nil {
		ll.Fatal("error while loading config", l.Error(err))
	}

	cm.SetEnvironment(cfg.Env)
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

	redisStore := redis.Connect(cfg.Redis.ConnectionString())
	tokens.Init(redisStore)

	db, err := cmsql.Connect(cfg.Postgres)
	if err != nil {
		ll.Fatal("error while connecting to etop postgres", l.Error(err))
	}

	dbaff, err := cmsql.Connect(cfg.AffPostgres)
	if err != nil {
		ll.Fatal("error while connecting to affiliate postgres")
	}

	sqlstore.Init(db)

	affiliateCmd := serviceaffiliate.NewAggregate(dbaff).MessageBus()
	affiliateQuery := serviceaffiliate.NewQuery(dbaff).MessageBus()
	catalogCmd := querycatalog.New(db).MessageBus()
	identityQuery := serviceidentity.NewQueryService(db).MessageBus()

	apiaff.Init(affiliateCmd, affiliateQuery, catalogCmd, identityQuery)

	apiMux := http.NewServeMux()

	middleware.Init("", identityQuery)

	mux := http.NewServeMux()
	mux.Handle("/api/", middleware.ForwardHeaders(apiMux))
	svr := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: mux,
	}

	wrapaff.NewAffiliateServer(apiMux, nil, cfg.Secret)

	metrics.RegisterHTTPHandler(mux)
	healthservice.RegisterHTTPHandler(mux)

	healthservice.MarkReady()
	go func() {
		defer ctxCancel()
		err = svr.ListenAndServe()
		if err != http.ErrServerClosed {
			ll.Error("HTTP server", l.Error(err))
		}
		ll.Sync()
	}()

	ll.Info("Server started")

	<-ctx.Done()
	_ = svr.Shutdown(context.Background())
	ll.Info("Waiting for all requests to finish")
	ll.Info("Gracefully stopped!")
}
