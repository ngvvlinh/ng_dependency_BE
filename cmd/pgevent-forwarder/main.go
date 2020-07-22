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
	"o.o/backend/com/eventhandler"
	etophandler "o.o/backend/com/eventhandler/etop/handler"
	fabohandler "o.o/backend/com/eventhandler/fabo/handler"
	"o.o/backend/com/eventhandler/pgevent"
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
	ll = l.New()

	flPrintTopics = flag.Bool("print-topics", false, "Print all topics then exit")
)

func main() {
	cc.InitFlags()
	cc.ParseFlags()

	cfg, err := config.Load()
	ll.Must(err, "can not load config", l.Error(err))

	cmenv.SetEnvironment("pgevent-forwarder", cfg.Env)
	ll.Info("Service started with config", l.String("commit", cm.CommitMessage()))
	if cmenv.IsDev() {
		ll.Info("config", l.Object("cfg", cfg))
	}

	ctx, ctxCancel := context.WithCancel(context.Background())
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

	ll.Must(cfg.Postgres.RegisterCloudSQL(), "can not register cloudsql", l.Error(err))

	producer, err := mq.NewKafkaProducer(ctx, cfg.Kafka.Brokers)
	ll.Must(err, "can not connect to Kafka", l.Error(err))

	topics := []eventhandler.TopicDef{}
	topics = append(topics, etophandler.Topics()...)
	topics = append(topics, fabohandler.Topics()...)
	if *flPrintTopics {
		printAllTopics(topics)
		os.Exit(0)
	}

	sMain, err := pgevent.NewService(ctx, model.DBMain, cfg.Postgres, producer, cfg.Kafka.TopicPrefix, topics)
	ll.Must(err, "Error while listening to Postgres")

	pgeventapi.Init(&sMain)

	sNotifier, err := pgevent.NewService(ctx, model.DBNotifier, cfg.PostgresNotifier, producer, cfg.Kafka.TopicPrefix, topics)
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
	healthService := health.New(nil)
	healthService.RegisterHTTPHandler(mux)
	defer healthService.Shutdown()
	healthService.MarkReady()

	ll.Info("Start forwarding events from Postgres to Kafka")
	pgevent.StartForwardings(ctx, []pgevent.Service{sMain, sNotifier})

	// Wait for OS signal or any error from services
	<-ctx.Done()
	_ = svr.Shutdown(context.Background())
}

func printAllTopics(topics []eventhandler.TopicDef) {
	for _, d := range topics {
		fmt.Printf("\t%3v %v\n", d.Partitions, d.Name)
	}
}
