package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"o.o/api/main/identity"
	"o.o/backend/cmd/etop-server/config"
	"o.o/backend/com/handler/etop-handler/intctl"
	serviceidentity "o.o/backend/com/main/identity"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/captcha"
	"o.o/backend/pkg/common/apifw/health"
	cmService "o.o/backend/pkg/common/apifw/service"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/authorization/auth"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/extservice/telebot"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqltrace"
	"o.o/backend/pkg/etop/api/export"
	"o.o/backend/pkg/etop/apix/partner"
	"o.o/backend/pkg/etop/apix/webhook"
	"o.o/backend/pkg/etop/authorize/middleware"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/authorize/tokens"
	"o.o/backend/pkg/etop/eventstream"
	"o.o/backend/pkg/etop/logic/summary"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/backend/pkg/integration/email"
	"o.o/common/l"
)

var (
	flTest     = flag.Bool("test", false, "Start services with default config for testing")
	flDocOnly  = flag.Bool("doc-only", false, "Serve API documentation only")
	flNoUpload = flag.Bool("no-upload", false, "Don't upload file")

	ll  = l.New()
	ctx context.Context
	cfg config.Config
	bot *telebot.Channel

	ctxCancel     context.CancelFunc
	healthservice = health.New()

	eventStreamer *eventstream.EventStreamer
	db            *cmsql.Database
	dbWebServer   *cmsql.Database
	dbLogs        *cmsql.Database
	identityQuery identity.QueryBus
)

func main() {
	cc.InitFlags()
	cc.ParseFlags()

	var err error
	cfg, err = config.Load(*flTest)
	if err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}

	cmenv.SetEnvironment(cfg.Env)
	cm.SetMainSiteBaseURL(cfg.URL.MainSite)
	sqltrace.Init()

	ll.Info("Service started with config", l.String("commit", cm.CommitMessage()))
	if cmenv.IsDev() {
		ll.Info("config", l.Object("cfg", cfg))
	}
	if *flTest {
		ll.Warn("TESTING MODE ENABLED: Use default config for testing")
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
		if bot != nil {
			bot.SendMessage("👻 Server stopped (forced) 👻\n–––")
		}
		ll.Fatal("Force shutdown due to timeout!")
	}()

	bot, err = cfg.TelegramBot.ConnectDefault()
	if err != nil {
		ll.Fatal("Unable to connect to Telegram", l.Error(err))
	}

	var producer *mq.KafkaProducer
	if cfg.Kafka.Enabled {
		producer, err = mq.NewKafkaProducer(ctx, cfg.Kafka.Brokers)
		if err != nil {
			ll.Fatal("Error while connecting to Kafka", l.Error(err))
		}
		ll.Info("Connected to Kafka")
	} else {
		ll.Warn("Disabled sending events to Kafka")
	}

	model.GetShippingServiceRegistry().Initialize()
	ctlProducer := producer.WithTopic(intctl.Topic(cfg.Kafka.TopicPrefix))

	shutdowner := cmService.NewShutdowner()
	eventStreamer = eventstream.NewEventStreamer(shutdowner)
	go eventStreamer.RunForwarder()

	redisStore := redis.Connect(cfg.Redis.ConnectionString())
	tokens.Init(redisStore)
	db, err = cmsql.Connect(cfg.Postgres)
	if err != nil {
		ll.Fatal("Unable to connect to Postgres", l.Error(err))
	}

	eventBus := bus.New()
	sqlstore.Init(db)
	sqlstore.AddEventBus(eventBus)
	summary.Init(db)

	dbLogs, err = cmsql.Connect(cfg.PostgresLogs)
	if err != nil {
		ll.Fatal("Unable to connect to Postgres (webhook_logs)", l.Error(err))
	}

	// MUSTDO(qv): database dependencies (db affiliate)
	{
		// init database notifier
		dbNotifier, err := cmsql.Connect(cfg.PostgresNotifier)
		if err != nil {
			ll.Fatal("Unable to connect to Postgres Notifier", l.Error(err))
		}
		sqlstore.InitDBNotifier(dbNotifier)
	}

	identityQuery = serviceidentity.QueryServiceMessageBus(serviceidentity.NewQueryService(db))

	authStore := auth.NewGenerator(redisStore)
	export.Init(shutdowner, redisStore, eventStreamer, export.Config{
		UrlPrefix: cfg.Export.URLPrefix,
		DirExport: cfg.Export.DirExport,
	})

	dbWebServer, err = cmsql.Connect(cfg.PostgresWebServer)
	if err != nil {
		ll.Fatal("Unable to connect to Postgres", l.Error(err))
	}

	whiteLabel := wl.Init(cmenv.Env())
	if err := whiteLabel.VerifyPartners(context.Background(), identityQuery); err != nil {
		ll.Fatal("error loading white label partners", l.Error(err))
	}

	middleware.Init(cfg.SAdminToken, identityQuery)

	webhook.Init(ctlProducer, redisStore)

	err = db.GetSchemaErrors()
	if err != nil && cmenv.IsDev() {
		ll.Error("Fail to verify Database", l.Error(err))
	} else if err != nil {
		// should move struct `callback` out of etop/model before change to ll.Fatal
		ll.Error("Fail to verify Database", l.Error(err))
	}

	ss := session.New(
		session.OptValidator(tokens.NewTokenStore(redisStore)),
		session.OptSuperAdmin(cfg.SAdminToken),
	)

	svrs, err := BuildServers(
		db,
		cfg,
		bot,
		shutdowner,
		eventBus,
		redisStore,
		authStore,
		ss,
		partner.AuthURL(cfg.URL.Auth),
	)
	if err != nil {
		ll.Fatal("can not init server", l.Error(err))
	}
	for _, svr := range svrs {
		svr := svr // https://golang.org/doc/faq#closures_and_goroutines
		ll.S.Infof("HTTP server %v listening on %v", svr.Name, svr.Addr)
		go func() {
			defer ctxCancel()
			err := svr.ListenAndServe()
			if err != http.ErrServerClosed {
				ll.Error("server", l.Error(err))
			}
		}()
	}

	if bot != nil {
		bot.SendMessage(fmt.Sprintf("–––\n✨ Server started on %v✨\n%v", cmenv.Env(), cm.CommitMessage()))
		defer bot.SendMessage("👻 Server stopped 👻\n–––")
	}

	if cmenv.IsDev() {
		ll.Warn("DEVELOPMENT MODE ENABLED")
	}

	if cfg.Email.Enabled {
		emailClient := email.New(cfg.SMTP)
		emailClient.Register(bus.Global())
		if err := emailClient.Ping(); err != nil {
			ll.Fatal("Unable to connect to email server", l.Error(err))
		}
		ll.Info("Enabled sending email")

	} else {
		ll.Warn("Disabled sending email")
	}

	captcha.Init(cfg.Captcha)

	healthservice.MarkReady()

	// Wait for OS signal or any error from services
	<-ctx.Done()
	ll.Info("Waiting for all requests to finish")
	shutdowner.ShutdownAll()

	// Graceful stop
	var wg sync.WaitGroup
	wg.Add(len(svrs))
	for _, svr := range svrs {
		svr := svr // https://golang.org/doc/faq#closures_and_goroutines
		go func() {
			defer wg.Done()
			svr.Shutdown(context.Background())
		}()
	}
	wg.Wait()

	ll.Info("Gracefully stopped!")
	ll.Sync()
}
