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

	"etop.vn/api/main/identity"
	"etop.vn/api/main/ordering"
	"etop.vn/api/main/shipnow"
	"etop.vn/backend/cmd/etop-server/config"
	haravanidentity "etop.vn/backend/com/external/haravan/identity"
	servicepaymentmanager "etop.vn/backend/com/external/payment/manager"
	"etop.vn/backend/com/handler/etop-handler/intctl"
	"etop.vn/backend/com/main/address"
	catalogaggregate "etop.vn/backend/com/main/catalog/aggregate"
	catalogquery "etop.vn/backend/com/main/catalog/query"
	serviceidentity "etop.vn/backend/com/main/identity"
	inventoryaggregate "etop.vn/backend/com/main/inventory/aggregate"
	inventorypm "etop.vn/backend/com/main/inventory/pm"
	inventoryquery "etop.vn/backend/com/main/inventory/query"
	servicelocation "etop.vn/backend/com/main/location"
	serviceordering "etop.vn/backend/com/main/ordering"
	serviceorderingpm "etop.vn/backend/com/main/ordering/pm"
	ordersqlstore "etop.vn/backend/com/main/ordering/sqlstore"
	receiptaggregate "etop.vn/backend/com/main/receipting/aggregate"
	receiptquery "etop.vn/backend/com/main/receipting/query"
	serviceshipnow "etop.vn/backend/com/main/shipnow"
	shipnowcarrier "etop.vn/backend/com/main/shipnow-carrier"
	shipnowpm "etop.vn/backend/com/main/shipnow/pm"
	shipsqlstore "etop.vn/backend/com/main/shipping/sqlstore"
	serviceaffiliate "etop.vn/backend/com/services/affiliate"
	affiliatepm "etop.vn/backend/com/services/affiliate/pm"
	carrieraggregate "etop.vn/backend/com/shopping/carrying/aggregate"
	carrierquery "etop.vn/backend/com/shopping/carrying/query"
	customeraggregate "etop.vn/backend/com/shopping/customering/aggregate"
	customerquery "etop.vn/backend/com/shopping/customering/query"
	vendorpm "etop.vn/backend/com/shopping/pm"
	traderquery "etop.vn/backend/com/shopping/tradering/query"
	vendoraggregate "etop.vn/backend/com/shopping/vendoring/aggregate"
	vendorquery "etop.vn/backend/com/shopping/vendoring/query"
	summaryquery "etop.vn/backend/com/summary/query"
	vhtaggregate "etop.vn/backend/com/supporting/crm/vht/aggregate"
	vhtquery "etop.vn/backend/com/supporting/crm/vht/query"
	vtigeraggregate "etop.vn/backend/com/supporting/crm/vtiger/aggregate"
	vtigerquery "etop.vn/backend/com/supporting/crm/vtiger/query"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/auth"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/captcha"
	"etop.vn/backend/pkg/common/cmsql"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/health"
	"etop.vn/backend/pkg/common/mq"
	"etop.vn/backend/pkg/common/redis"
	cmService "etop.vn/backend/pkg/common/service"
	"etop.vn/backend/pkg/common/telebot"
	"etop.vn/backend/pkg/etop/api"
	"etop.vn/backend/pkg/etop/api/affiliate"
	"etop.vn/backend/pkg/etop/api/crm"
	"etop.vn/backend/pkg/etop/api/export"
	"etop.vn/backend/pkg/etop/api/integration"
	"etop.vn/backend/pkg/etop/api/shop"
	"etop.vn/backend/pkg/etop/apix/partner"
	xshipping "etop.vn/backend/pkg/etop/apix/shipping"
	xshop "etop.vn/backend/pkg/etop/apix/shop"
	"etop.vn/backend/pkg/etop/apix/webhook"
	"etop.vn/backend/pkg/etop/authorize/middleware"
	"etop.vn/backend/pkg/etop/authorize/tokens"
	"etop.vn/backend/pkg/etop/eventstream"
	orderS "etop.vn/backend/pkg/etop/logic/orders"
	imcsvorder "etop.vn/backend/pkg/etop/logic/orders/imcsv"
	imcsvproduct "etop.vn/backend/pkg/etop/logic/products/imcsv"
	"etop.vn/backend/pkg/etop/logic/shipping_provider"
	"etop.vn/backend/pkg/etop/logic/summary"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/backend/pkg/etop/upload"
	"etop.vn/backend/pkg/integration/email"
	"etop.vn/backend/pkg/integration/payment/vtpay"
	vtpayclient "etop.vn/backend/pkg/integration/payment/vtpay/client"
	"etop.vn/backend/pkg/integration/shipnow/ahamove"
	"etop.vn/backend/pkg/integration/shipping/ghn"
	"etop.vn/backend/pkg/integration/shipping/ghtk"
	"etop.vn/backend/pkg/integration/shipping/vtpost"
	"etop.vn/backend/pkg/integration/sms"
	vtigerclient "etop.vn/backend/pkg/integration/vtiger/client"
	apiaff "etop.vn/backend/pkg/services/affiliate/api"
	"etop.vn/common/l"
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

	eventStreamer         *eventstream.EventStreamer
	db                    *cmsql.Database
	dbLogs                *cmsql.Database
	ghnCarrier            *ghn.Carrier
	ghtkCarrier           *ghtk.Carrier
	vtpostCarrier         *vtpost.Carrier
	ahamoveCarrier        *ahamove.Carrier
	ahamoveCarrierAccount *ahamove.CarrierAccount

	shipnowQuery shipnow.QueryBus
	shipnowAggr  shipnow.CommandBus
	orderAggr    *serviceordering.Aggregate
	orderQuery   ordering.QueryBus

	identityQuery identity.QueryBus

	vtpayClient *vtpayclient.Client
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

	ll.Info("Service started with config", l.String("commit", cm.CommitMessage()))
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
	db, err = cmsql.Connect(cfg.Postgres)
	if err != nil {
		ll.Fatal("Unable to connect to Postgres", l.Error(err))
	}
	sqlstore.Init(db)
	summary.Init(db)

	dbLogs, err = cmsql.Connect(cfg.PostgresLogs)
	if err != nil {
		ll.Fatal("Unable to connect to Postgres (webhook_logs)", l.Error(err))
	}

	dbaff, err := cmsql.Connect(cfg.PostgresAffiliate)
	if err != nil {
		ll.Fatal("error while connecting to affiliate postgres")
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
	identityQuery = serviceidentity.NewQueryService(db).MessageBus()
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
	if cfg.Ahamove.ApiKey != "" {
		ahamoveCarrier, ahamoveCarrierAccount = ahamove.New(
			cfg.Ahamove, locationBus, identityQuery,
			ahamove.URLConfig{
				ThirdPartyHost:       cfg.ThirdPartyHost,
				PathUserVerification: config.PathAhamoveUserVerification,
			})
		if err := ahamoveCarrier.InitClient(ctx); err != nil {
			ll.Fatal("Unable to connect to ahamove", l.Error(err))
		}
	} else {
		if cm.IsDev() {
			ll.Warn("DEVELOPMENT. Skip connecting to ahamove.")
		} else {
			ll.Fatal("ahamove: No token")
		}
	}
	configMap, err := config.ReadMappingFile(cfg.Vtiger.MappingFile)
	if err != nil {
		ll.Fatal("error while reading field map file", l.String("file", cfg.Vtiger.MappingFile), l.Error(err))
	}

	shippingManager := shipping_provider.NewCtrl(locationBus, ghnCarrier, ghtkCarrier, vtpostCarrier)

	authStore := auth.NewGenerator(redisStore)
	imcsvorder.Init(locationBus, shutdowner, redisStore, uploader, db)
	imcsvproduct.Init(shutdowner, redisStore, uploader, db)
	export.Init(shutdowner, redisStore, eventStreamer, export.Config{
		UrlPrefix: cfg.Export.URLPrefix,
		DirExport: cfg.Export.DirExport,
	})

	eventBus := bus.New()
	crmDB, err := cmsql.Connect(cfg.PostgresCRM)
	if err != nil {
		ll.Fatal("Unable to connect to Postgres", l.Error(err))
	}

	vtigerClient := vtigerclient.NewVigerClient(cfg.Vtiger.ServiceURL, cfg.Vtiger.Username, cfg.Vtiger.APIKey)
	// create aggregate, query service
	summaryQuery := summaryquery.NewDashboardQuery(db, redisStore).MessageBus()
	inventoryaggregate := inventoryaggregate.NewAggregateInventory(eventBus, db).MessageBus()
	inventoryQuery := inventoryquery.NewQueryInventory(eventBus, db).MessageBus()
	vhtQuery := vhtquery.New(crmDB).MessageBus()
	vhtAggregate := vhtaggregate.New(crmDB, nil).MessageBus()
	vtigerQuery := vtigerquery.New(crmDB, configMap, vtigerClient).MessageBus()
	vtigerAggregate := vtigeraggregate.New(crmDB, configMap, vtigerClient).MessageBus()

	identityQuery = serviceidentity.NewQueryService(db).MessageBus()
	catalogQuery := catalogquery.New(db).MessageBus()
	catalogAggr := catalogaggregate.New(eventBus, db).MessageBus()
	addressQuery := address.NewQueryService(db).MessageBus()
	shipnowQuery = serviceshipnow.NewQueryService(db).MessageBus()
	orderQuery = serviceordering.NewQueryService(db).MessageBus()
	haravanIdentityAggr := haravanidentity.NewAggregate(db, cfg.ThirdPartyHost, cfg.Haravan).MessageBus()
	haravanIdentityQuery := haravanidentity.NewQueryService(db).MessageBus()

	orderAggr = serviceordering.NewAggregate(eventBus, db)
	shipnowCarrierManager := shipnowcarrier.NewManager(db, locationBus, &shipnowcarrier.Carrier{
		ShipnowCarrier:        ahamoveCarrier,
		ShipnowCarrierAccount: ahamoveCarrierAccount,
	}, shipnowQuery)
	identityAggr := serviceidentity.NewAggregate(db, shipnowCarrierManager).MessageBus()
	shipnowAggr = serviceshipnow.NewAggregate(eventBus, db, locationBus, identityQuery, addressQuery, orderQuery, shipnowCarrierManager).MessageBus()

	inventoryPm := inventorypm.New(eventBus, catalogQuery, orderQuery, inventoryaggregate)
	inventoryPm.RegisterEventHandlers(eventBus)
	shipnowPM := shipnowpm.New(eventBus, shipnowQuery, shipnowAggr, orderAggr.MessageBus(), shipnowCarrierManager)
	shipnowPM.RegisterEventHandlers(eventBus)

	customerAggr := customeraggregate.NewCustomerAggregate(db).MessageBus()
	vendorAggr := vendoraggregate.NewVendorAggregate(db).MessageBus()
	carrierAggr := carrieraggregate.NewCarrierAggregate(db).MessageBus()
	traderAddressAggr := customeraggregate.NewAddressAggregate(db).MessageBus()
	customerQuery := customerquery.NewCustomerQuery(db).MessageBus()
	vendorQuery := vendorquery.NewVendorQuery(db).MessageBus()
	carrierQuery := carrierquery.NewCarrierQuery(db).MessageBus()
	traderQuery := traderquery.NewTraderQuery(db).MessageBus()
	traderAddressQuery := customerquery.NewAddressQuery(db).MessageBus()
	affiliateCmd := serviceaffiliate.NewAggregate(dbaff, identityQuery, catalogQuery, orderQuery).MessageBus()
	affilateQuery := serviceaffiliate.NewQuery(dbaff).MessageBus()
	affiliatePM := affiliatepm.New(affiliateCmd)
	affiliatePM.RegisterEventHandlers(eventBus)

	receiptAggr := receiptaggregate.NewReceiptAggregate(db).MessageBus()
	receiptQuery := receiptquery.NewReceiptQuery(db).MessageBus()

	vendorPM := vendorpm.New(eventBus, vendorQuery)
	vendorPM.RegisterEventHandlers(eventBus)
	// payment
	var vtpayProvider *vtpay.Provider
	if cfg.VTPay.MerchantCode != "" {
		vtpayClient = vtpayclient.New(cfg.VTPay)
		vtpayProvider = vtpay.New(cfg.VTPay)
	}
	paymentManager := servicepaymentmanager.NewManager(vtpayProvider, orderQuery).MesssageBus()
	orderPM := serviceorderingpm.New(orderAggr.MessageBus(), orderQuery, affiliateCmd)
	orderPM.RegisterEventHandlers(eventBus)

	middleware.Init(cfg.SAdminToken, identityQuery)
	api.Init(identityAggr, identityQuery, shutdowner, redisStore, authStore, cfg.Email, cfg.SMS)
	shop.Init(
		locationBus,
		catalogQuery,
		catalogAggr,
		shipnowAggr,
		shipnowQuery,
		identityAggr,
		identityQuery,
		addressQuery,
		shippingManager,
		haravanIdentityAggr,
		haravanIdentityQuery,
		customerAggr,
		customerQuery,
		traderAddressAggr,
		traderAddressQuery,
		orderAggr.MessageBus(),
		orderQuery,
		paymentManager,
		vendorAggr,
		vendorQuery,
		carrierAggr,
		carrierQuery,
		traderQuery,
		eventBus,
		receiptAggr,
		receiptQuery,
		shutdowner,
		redisStore,
		inventoryaggregate,
		inventoryQuery,
		summaryQuery,
	)
	partner.Init(shutdowner, redisStore, authStore, cfg.URL.Auth)
	xshop.Init(shutdowner, redisStore, authStore)
	integration.Init(shutdowner, redisStore, authStore)
	webhook.Init(ctlProducer, redisStore)
	xshipping.Init(shippingManager, ordersqlstore.NewOrderStore(db), shipsqlstore.NewFulfillmentStore(db))
	orderS.Init(shippingManager, catalogQuery, orderAggr.MessageBus(),
		customerAggr, customerQuery, traderAddressAggr, traderAddressQuery, locationBus)
	crm.Init(ghnCarrier, vtigerQuery, vtigerAggregate, vhtQuery, vhtAggregate)
	affiliate.Init(identityAggr)
	apiaff.Init(affiliateCmd, affilateQuery, catalogQuery, identityQuery, orderQuery)

	err = db.GetSchemaErrors()
	if err != nil && cm.IsDev() {
		ll.Error("Fail to verify Database", l.Error(err))
	} else if err != nil {
		ll.Fatal("Fail to verify Database", l.Error(err))
	}

	svrs := startServers()
	if bot != nil {
		bot.SendMessage("–––\n✨ Server started ✨\n" + cm.CommitMessage())
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
		smsBot := cfg.TelegramBot.MustConnectChannel(config.ChannelSMS)
		sms.New(cfg.SMS, smsBot).Register(bus.Global())
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
