package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	pgeventapi "o.o/backend/cmd/pgevent-forwarder/api"
	"o.o/backend/cmd/pgevent-forwarder/config"
	"o.o/backend/com/handler/pgevent"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/etop/model"
	"o.o/common/l"
)

var (
	ll  = l.New()
	cfg config.Config
	ctx context.Context

	ctxCancel     context.CancelFunc
	healthservice = health.New()

	flPrintTopics = flag.Bool("print-topics", false, "Print all topics then exit")
)

func main() {
	cc.InitFlags()
	cc.ParseFlags()

	if *flPrintTopics {
		printAllTopics()
		os.Exit(0)
	}

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

	if err := cfg.Postgres.RegisterCloudSQL(); err != nil {
		ll.Fatal("Error while registering cloudsql", l.Error(err))
	}

	producer, err := mq.NewKafkaProducer(ctx, cfg.Kafka.Brokers)
	if err != nil {
		ll.Fatal("Error while connecting to Kafka", l.Error(err))
	}

	sMain, err := pgevent.NewService(ctx, model.DBMain, cfg.Postgres, producer, cfg.Kafka.TopicPrefix)
	if err != nil {
		ll.Fatal("Error while listening to Postgres")
	}

	pgeventapi.Init(&sMain)

	sNotifier, err := pgevent.NewService(ctx, model.DBNotifier, cfg.PostgresNotifier, producer, cfg.Kafka.TopicPrefix)
	if err != nil {
		ll.Fatal("Error while listening to Postgres")
	}

	apiMux := http.NewServeMux()
	apiMux.Handle("/api/", http.NotFoundHandler())
	pgeventapi.NewPgeventServer(apiMux, cfg.Secret)

	mux := http.NewServeMux()
	mux.Handle("/api/", headers.ForwardHeaders(apiMux))
	svr := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: mux,
	}

	go func() {
		defer ctxCancel()
		err := svr.ListenAndServe()
		if err != http.ErrServerClosed {
			ll.Error("HTTP server", l.Error(err))
		}
		ll.Sync()
	}()
	healthservice.RegisterHTTPHandler(mux)
	healthservice.MarkReady()

	ll.Info("Start forwarding events from Postgres to Kafka")
	pgevent.StartForwardings(ctx, []pgevent.Service{sMain, sNotifier})
	// s.StartForwarding(ctx)
	// Wait for OS signal or any error from services
	<-ctx.Done()
	_ = svr.Shutdown(context.Background())

	ll.Info("Gracefully stopped!")
}

func printAllTopics() {
	for _, d := range pgevent.Topics {
		fmt.Printf("\t%3v %v\n", d.Partitions, d.Name)
	}
}
