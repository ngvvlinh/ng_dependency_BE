package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"etop.vn/backend/cmd/etop-server/config"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/auth"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/captcha"
	"etop.vn/backend/pkg/common/cmsql"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/health"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/common/mq"
	"etop.vn/backend/pkg/common/redis"
	cmService "etop.vn/backend/pkg/common/service"
	"etop.vn/backend/pkg/common/telebot"
	"etop.vn/backend/pkg/etop-handler/intctl"
	"etop.vn/backend/pkg/etop/api"
	"etop.vn/backend/pkg/etop/api/integration"
	"etop.vn/backend/pkg/etop/api/shop"
	"etop.vn/backend/pkg/etop/apix/partner"
	xshipping "etop.vn/backend/pkg/etop/apix/shipping"
	xshop "etop.vn/backend/pkg/etop/apix/shop"
	"etop.vn/backend/pkg/etop/apix/webhook"
	"etop.vn/backend/pkg/etop/authorize/middleware"
	"etop.vn/backend/pkg/etop/authorize/tokens"
	"etop.vn/backend/pkg/etop/eventstream"
	ffmexport "etop.vn/backend/pkg/etop/logic/fulfillments/export"
	orderS "etop.vn/backend/pkg/etop/logic/orders"
	imcsvorder "etop.vn/backend/pkg/etop/logic/orders/imcsv"
	imcsvproduct "etop.vn/backend/pkg/etop/logic/products/imcsv"
	"etop.vn/backend/pkg/etop/logic/shipping_provider"
	"etop.vn/backend/pkg/etop/logic/summary"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/backend/pkg/etop/upload"
	"etop.vn/backend/pkg/integration/ahamove"
	"etop.vn/backend/pkg/integration/email"
	"etop.vn/backend/pkg/integration/ghn"
	"etop.vn/backend/pkg/integration/ghtk"
	"etop.vn/backend/pkg/integration/sms"
	"etop.vn/backend/pkg/integration/vtpost"
	"etop.vn/backend/pkg/services/address"
	"etop.vn/backend/pkg/services/identity"
	servicelocation "etop.vn/backend/pkg/services/location"
	"etop.vn/backend/pkg/services/ordering"
	orderingpm "etop.vn/backend/pkg/services/ordering/pm"
	ordersqlstore "etop.vn/backend/pkg/services/ordering/sqlstore"
	"etop.vn/backend/pkg/services/shipnow"
	shipnow_carrier "etop.vn/backend/pkg/services/shipnow-carrier"
	shipnowpm "etop.vn/backend/pkg/services/shipnow/pm"
	shipsqlstore "etop.vn/backend/pkg/services/shipping/sqlstore"
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

	eventStreamer  *eventstream.EventStreamer
	dbLogs         cmsql.Database
	ghnCarrier     *ghn.Carrier
	ghtkCarrier    *ghtk.Carrier
	vtpostCarrier  *vtpost.Carrier
	ahamoveCarrier *ahamove.Carrier
)

func main() {
	cc.InitFlags()
	flag.Parse()

	var err error
	cfg, err = config.Load(*flTest)
	if err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}

	cm.SetEnvironment(cfg.Env)
	cm.SetMainSiteBaseURL(cfg.URL.MainSite)

	ll.Info("Service started with config", l.String("commit", cm.Commit()))
	if cm.IsDev() {
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

	var uploader *upload.Uploader
	if !*flNoUpload {
		uploader, err = upload.NewUploader(map[string]string{
			string(model.ImportTypeShopOrder):   cfg.Upload.DirImportShopOrder,
			string(model.ImportTypeShopProduct): cfg.Upload.DirImportShopProduct,
		})
		if err != nil {
			ll.Fatal("Unable to init uploader", l.Error(err))
		}
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
	middleware.Init(cfg.SAdminToken)
	db, err := cmsql.Connect(cfg.Postgres)
	if err != nil {
		ll.Fatal("Unable to connect to Postgres", l.Error(err))
	}
	sqlstore.Init(db)
	summary.Init(db)

	dbLogs, err = cmsql.Connect(cfg.PostgresLogs)
	if err != nil {
		ll.Fatal("Unable to connect to Postgres (webhook_logs)", l.Error(err))
	}

	{
		// init database notifier
		dbNotifier, err := cmsql.Connect(cfg.PostgresNotifier)
		if err != nil {
			ll.Fatal("Unable to connect to Postgres Notifier", l.Error(err))
		}
		sqlstore.InitDBNotifier(dbNotifier)
	}

	locationBus := servicelocation.New().MessageBus()
	if cfg.GHN.AccountDefault.Token != "" {
		ghnCarrier = ghn.New(cfg.GHN, locationBus)
		if err := ghnCarrier.InitAllClients(ctx); err != nil {
			ll.Fatal("Unable to connect to GHN", l.Error(err))
		}
	} else {
		if cm.IsDev() {
			ll.Warn("DEVELOPMENT. Skip connecting to GHN")
		} else {
			ll.Fatal("GHN: No token")
		}
	}
	if cfg.GHTK.AccountDefault.Token != "" {
		ghtkCarrier = ghtk.New(cfg.GHTK, locationBus)
		if err := ghtkCarrier.InitAllClients(ctx); err != nil {
			ll.Fatal("Unable to connect to GHTK", l.Error(err))
		}
	} else {
		if cm.IsDev() {
			ll.Warn("DEVELOPMENT. Skip connecting to GHTK.")
		} else {
			ll.Fatal("GHTK: No token")
		}
	}

	if cfg.VTPost.AccountDefault.Username != "" {
		vtpostCarrier = vtpost.New(cfg.VTPost, locationBus)
		if err := vtpostCarrier.InitAllClients(ctx); err != nil {
			ll.Fatal("Unable to connect to VTPost", l.Error(err))
		}
	} else {
		if cm.IsDev() {
			ll.Warn("DEVELOPMENT. Skip connecting to VTPost.")
		} else {
			ll.Fatal("VTPost: No token")
		}
	}

	if cfg.Ahamove.AccountDefault.Token != "" {
		ahamoveCarrier = ahamove.New(cfg.Ahamove, locationBus)
		if err := ahamoveCarrier.InitAllClients(ctx); err != nil {
			ll.Fatal("Unable to connect to ahamove", l.Error(err))
		}
	} else {
		if cm.IsDev() {
			ll.Warn("DEVELOPMENT. Skip connecting to ahamove.")
		} else {
			ll.Fatal("ahamove: No token")
		}
	}

	shippingManager := shipping_provider.NewCtrl(locationBus, ghnCarrier, ghtkCarrier, vtpostCarrier)

	authStore := auth.NewGenerator(redisStore)
	api.Init(shutdowner, redisStore, authStore, cfg.Email, cfg.SMS)
	imcsvorder.Init(locationBus, shutdowner, redisStore, uploader)
	imcsvproduct.Init(shutdowner, redisStore, uploader, db)
	ffmexport.Init(shutdowner, redisStore, eventStreamer, ffmexport.Config{
		UrlPrefix: cfg.Export.URLPrefix,
		DirExport: cfg.Export.DirExport,
	})

	eventBus := bus.New()

	shipnowCarrierManager := shipnow_carrier.NewManager(db, locationBus, ahamoveCarrier)
	// create aggregate, query service
	identityQuery := identity.NewQueryService(db)
	addressQuery := address.NewQueryService(db)
	shipnowQuery := shipnow.NewQueryService(db)

	orderAggregate := ordering.NewAggregate(db)
	shipnowAggregate := shipnow.NewAggregate(eventBus, db, locationBus, identityQuery, addressQuery)

	orderingPM := orderingpm.New(orderAggregate, shipnowAggregate)
	shipnowPM := shipnowpm.New(eventBus, shipnowQuery.MessageBus(), orderAggregate, orderAggregate.MessageBus(), identityQuery, addressQuery, shipnowCarrierManager)
	shipnowPM.RegisterEventHandlers(eventBus)

	orderAggregate.WithPM(orderingPM)

	shop.Init(shipnowAggregate, shipnowQuery, shippingManager, shutdowner, redisStore)
	partner.Init(shutdowner, redisStore, authStore, cfg.URL.Auth)
	xshop.Init(shutdowner, redisStore, authStore)
	integration.Init(shutdowner, redisStore, authStore)
	webhook.Init(ctlProducer, redisStore)
	xshipping.Init(shippingManager, ordersqlstore.NewOrderStore(db), shipsqlstore.NewFulfillmentStore(db))
	orderS.Init(shippingManager)

	svrs := startServers()
	if bot != nil {
		bot.SendMessage("–––\n✨ Server started ✨\n" + cm.Commit())
		defer bot.SendMessage("👻 Server stopped 👻\n–––")
	}

	if cm.IsDev() {
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
	if cfg.SMS.Enabled {
		sms.New(cfg.SMS).Register(bus.Global())
		ll.Info("Enabled sending sms")
	} else {
		ll.Warn("Disabled sending sms")
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
		go func(svr *http.Server) {
			defer wg.Done()
			svr.Shutdown(context.Background())
		}(svr)
	}
	wg.Wait()

	ll.Info("Gracefully stopped!")
	ll.Sync()
}
