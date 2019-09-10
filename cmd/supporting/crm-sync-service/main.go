package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"etop.vn/backend/cmd/supporting/crm-sync-service/config"
	vhtaggregate "etop.vn/backend/com/supporting/crm/vht/aggregate"
	vhtquery "etop.vn/backend/com/supporting/crm/vht/query"
	vtigeraggregate "etop.vn/backend/com/supporting/crm/vtiger/aggregate"
	"etop.vn/backend/com/supporting/crm/vtiger/mapping"
	vtigerquery "etop.vn/backend/com/supporting/crm/vtiger/query"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/health"
	"etop.vn/backend/pkg/common/metrics"
	vhtclient "etop.vn/backend/pkg/integration/vht/client"
	vtigerclient "etop.vn/backend/pkg/integration/vtiger/client"
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

	<-ctx.Done()

	mux := http.NewServeMux()
	svr := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: mux,
	}
	metrics.RegisterHTTPHandler(mux)
	healthservice.RegisterHTTPHandler(mux)
	healthservice.MarkReady()

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
	err = json.Unmarshal(body, &configMap)
	return
}
