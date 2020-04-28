package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"o.o/backend/cmd/supporting/crm-sync-service/config"
	vhtaggregate "o.o/backend/com/supporting/crm/vht/aggregate"
	vhtquery "o.o/backend/com/supporting/crm/vht/query"
	vtigeraggregate "o.o/backend/com/supporting/crm/vtiger/aggregate"
	"o.o/backend/com/supporting/crm/vtiger/mapping"
	vtigerquery "o.o/backend/com/supporting/crm/vtiger/query"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/metrics"
	"o.o/backend/pkg/common/sql/cmsql"
	vhtclient "o.o/backend/pkg/integration/vht/client"
	vtigerclient "o.o/backend/pkg/integration/vtiger/client"
	"o.o/common/jsonx"
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
		ll.Fatal("error while loading config", l.Error(err))
	}

	cmenv.SetEnvironment(cfg.Env)
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

	db, err := cmsql.Connect(cfg.Postgres)
	if err != nil {
		ll.Fatal("error while connecting to postgres", l.Error(err))
	}

	configMap, err := ReadMappingFile(cfg.Vtiger.MappingFile)
	if err != nil {
		ll.Fatal("error while reading field map file", l.String("file", cfg.Vtiger.MappingFile), l.Error(err))
	}

	vtigerClient := vtigerclient.NewVigerClient(cfg.Vtiger.ServiceURL, cfg.Vtiger.Username, cfg.Vtiger.APIKey)
	vhtClient := vhtclient.NewClient(cfg.Vht.Username, cfg.Vht.Password)

	vhtAggregate := vhtaggregate.New(db, vhtClient).MessageBus()
	vhtQuery := vhtquery.New(db).MessageBus()
	vtigerAggregate := vtigeraggregate.New(db, configMap, vtigerClient).MessageBus()
	vtigerQuery := vtigerquery.New(db, configMap, vtigerClient).MessageBus()
	go func() {
		SyncCallHistoryVht(vhtAggregate, vhtQuery)
	}()

	ll.Info("Sync Vtiger Starting")
	go func() {
		SyncVtiger(vtigerAggregate, vtigerQuery)
	}()

	mux := http.NewServeMux()
	svr := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: mux,
	}
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

// ReadMappingFile read mapping json file for mapping fields between vtiger and etop
func ReadMappingFile(filename string) (configMap mapping.ConfigMap, _ error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = jsonx.Unmarshal(body, &configMap)
	return
}
