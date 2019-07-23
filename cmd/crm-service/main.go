package main

import (
	"context"
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"etop.vn/backend/cmd/crm-service/config"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/health"
	"etop.vn/backend/pkg/common/metrics"
	"etop.vn/backend/pkg/common/redis"
	"etop.vn/backend/pkg/crm-service/mapping"
	servicecrm "etop.vn/backend/pkg/crm-service/service"
	vs "etop.vn/backend/pkg/crm-service/vtiger-service"
	"etop.vn/backend/pkg/etop/authorize/middleware"
	"etop.vn/backend/pkg/etop/authorize/tokens"
	"etop.vn/backend/pkg/etop/sqlstore"
	wrapcrm "etop.vn/backend/wrapper/services/crmservice"
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
		ll.Fatal("error while connecting to postgres", l.Error(err))
	}

	vtigerConfig := vs.VtigerConfig{
		VtigerService:   cfg.VtigerService,
		VtigerUsername:  cfg.VtigerUsername,
		VtigerAccesskey: cfg.VtigerAccesskey,
	}

	f, err := vs.ReadFileConfig()
	var configMap *mapping.ConfigMap
	{
		err := json.Unmarshal([]byte(f), &configMap)
		if err != nil {
			ll.Error("Fail to Unmarshal field_mapping.json to  mappping.ConfigMap")
		}
	}
	if err != nil {
		ll.Fatal("err while read field map file", l.Error(err))
	}
	s := servicecrm.NewService(db, vtigerConfig, configMap)
	s.Register()

	apiMux := http.NewServeMux()
	wrapcrm.NewCrmserviceServer(apiMux, nil)

	mux := http.NewServeMux()
	mux.Handle("/api/", middleware.ForwardHeaders(apiMux))
	svr := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: mux,
	}
	sqlstore.Init(db)

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

	<-ctx.Done()
	_ = svr.Shutdown(context.Background())
	ll.Info("Waiting for all requests to finish")
	ll.Info("Gracefully stopped!")
}
