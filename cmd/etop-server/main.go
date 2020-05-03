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
	"o.o/api/main/invitation"
	"o.o/api/main/moneytx"
	"o.o/api/main/ordering"
	"o.o/api/main/receipting"
	"o.o/api/main/shipnow"
	"o.o/api/top/types/etc/connection_type"
	"o.o/backend/cmd/etop-server/config"
	smsAgg "o.o/backend/com/etc/logging/smslog/aggregate"
	servicepaymentmanager "o.o/backend/com/external/payment/manager"
	"o.o/backend/com/handler/etop-handler/intctl"
	"o.o/backend/com/main/address"
	authorizationaggregate "o.o/backend/com/main/authorization/aggregate"
	authorizationquery "o.o/backend/com/main/authorization/query"
	catalogaggregate "o.o/backend/com/main/catalog/aggregate"
	catalogpm "o.o/backend/com/main/catalog/pm"
	catalogquery "o.o/backend/com/main/catalog/query"
	connectionaggregate "o.o/backend/com/main/connectioning/aggregate"
	connectionquery "o.o/backend/com/main/connectioning/query"
	serviceidentity "o.o/backend/com/main/identity"
	identitypm "o.o/backend/com/main/identity/pm"
	inventoryaggregate "o.o/backend/com/main/inventory/aggregate"
	inventorypm "o.o/backend/com/main/inventory/pm"
	inventoryquery "o.o/backend/com/main/inventory/query"
	invitationaggregate "o.o/backend/com/main/invitation/aggregate"
	invitationpm "o.o/backend/com/main/invitation/pm"
	invitationquery "o.o/backend/com/main/invitation/query"
	ledgeraggregate "o.o/backend/com/main/ledgering/aggregate"
	ledgerpm "o.o/backend/com/main/ledgering/pm"
	ledgerquery "o.o/backend/com/main/ledgering/query"
	servicelocation "o.o/backend/com/main/location"
	moneytxaggregate "o.o/backend/com/main/moneytx/aggregate"
	moneytxpm "o.o/backend/com/main/moneytx/pm"
	moneytxquery "o.o/backend/com/main/moneytx/query"
	serviceordering "o.o/backend/com/main/ordering"
	serviceorderingpm "o.o/backend/com/main/ordering/pm"
	ordersqlstore "o.o/backend/com/main/ordering/sqlstore"
	purchaseorderaggregate "o.o/backend/com/main/purchaseorder/aggregate"
	purchaseorderpm "o.o/backend/com/main/purchaseorder/pm"
	purchaseorderquery "o.o/backend/com/main/purchaseorder/query"
	purchaserefundaggregate "o.o/backend/com/main/purchaserefund/aggregate"
	purchaserefundpm "o.o/backend/com/main/purchaserefund/pm"
	purchaserefundquery "o.o/backend/com/main/purchaserefund/query"
	receiptaggregate "o.o/backend/com/main/receipting/aggregate"
	receiptpm "o.o/backend/com/main/receipting/pm"
	receiptquery "o.o/backend/com/main/receipting/query"
	refundaggregate "o.o/backend/com/main/refund/aggregate"
	refundpm "o.o/backend/com/main/refund/pm"
	refundquery "o.o/backend/com/main/refund/query"
	"o.o/backend/com/main/shipmentpricing/pricelist"
	pricelistpm "o.o/backend/com/main/shipmentpricing/pricelist/pm"
	"o.o/backend/com/main/shipmentpricing/shipmentprice"
	"o.o/backend/com/main/shipmentpricing/shipmentservice"
	serviceshipnow "o.o/backend/com/main/shipnow"
	shipnowcarrier "o.o/backend/com/main/shipnow-carrier"
	shipnowpm "o.o/backend/com/main/shipnow/pm"
	shippingaggregate "o.o/backend/com/main/shipping/aggregate"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
	shippingpm "o.o/backend/com/main/shipping/pm"
	shippingquery "o.o/backend/com/main/shipping/query"
	shipsqlstore "o.o/backend/com/main/shipping/sqlstore"
	stocktakeaggregate "o.o/backend/com/main/stocktaking/aggregate"
	stocktakequery "o.o/backend/com/main/stocktaking/query"
	serviceaffiliate "o.o/backend/com/services/affiliate"
	affiliatepm "o.o/backend/com/services/affiliate/pm"
	carrieraggregate "o.o/backend/com/shopping/carrying/aggregate"
	carrierquery "o.o/backend/com/shopping/carrying/query"
	customeraggregate "o.o/backend/com/shopping/customering/aggregate"
	customerquery "o.o/backend/com/shopping/customering/query"
	supplieraggregate "o.o/backend/com/shopping/suppliering/aggregate"
	supplierquery "o.o/backend/com/shopping/suppliering/query"
	traderAgg "o.o/backend/com/shopping/tradering/aggregate"
	traderpm "o.o/backend/com/shopping/tradering/pm"
	traderquery "o.o/backend/com/shopping/tradering/query"
	summaryquery "o.o/backend/com/summary/query"
	webserveraggregate "o.o/backend/com/web/webserver/aggregate"
	webserverquery "o.o/backend/com/web/webserver/query"
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
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/admin"
	"o.o/backend/pkg/etop/api/affiliate"
	"o.o/backend/pkg/etop/api/export"
	"o.o/backend/pkg/etop/api/integration"
	"o.o/backend/pkg/etop/api/shop"
	"o.o/backend/pkg/etop/apix/partner"
	xshipping "o.o/backend/pkg/etop/apix/shipping"
	xshop "o.o/backend/pkg/etop/apix/shop"
	xshopping "o.o/backend/pkg/etop/apix/shopping"
	"o.o/backend/pkg/etop/apix/webhook"
	whitelabelapix "o.o/backend/pkg/etop/apix/whitelabel"
	authorizeauth "o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/authorize/middleware"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/authorize/tokens"
	"o.o/backend/pkg/etop/eventstream"
	orderS "o.o/backend/pkg/etop/logic/orders"
	imcsvorder "o.o/backend/pkg/etop/logic/orders/imcsv"
	imcsvproduct "o.o/backend/pkg/etop/logic/products/imcsv"
	"o.o/backend/pkg/etop/logic/shipping_provider"
	"o.o/backend/pkg/etop/logic/summary"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/backend/pkg/etop/upload"
	"o.o/backend/pkg/integration/email"
	"o.o/backend/pkg/integration/payment/vtpay"
	vtpayclient "o.o/backend/pkg/integration/payment/vtpay/client"
	"o.o/backend/pkg/integration/shipnow/ahamove"
	"o.o/backend/pkg/integration/shipping/ghn"
	"o.o/backend/pkg/integration/shipping/ghtk"
	"o.o/backend/pkg/integration/shipping/vtpost"
	"o.o/backend/pkg/integration/sms"
	imgroupsms "o.o/backend/pkg/integration/sms/imgroup"
	apiaff "o.o/backend/pkg/services/affiliate/api"
	"o.o/backend/tools/pkg/acl"
	"o.o/capi/httprpc"
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

	eventStreamer         *eventstream.EventStreamer
	db                    *cmsql.Database
	dbWebServer           *cmsql.Database
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

	invitationQuery invitation.QueryBus

	vtpayClient *vtpayclient.Client

	receiptQuery    receipting.QueryBus
	shipmentManager *shippingcarrier.ShipmentManager
	moneyTxQuery    moneytx.QueryBus
	moneyTxAggr     moneytx.CommandBus

	ss    *session.Session
	hooks *httprpc.Hooks
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

	var uploader *upload.Uploader
	if !*flNoUpload {
		uploader, err = upload.NewUploader(map[string]string{
			model.ImportTypeShopOrder.String():   cfg.Upload.DirImportShopOrder,
			model.ImportTypeShopProduct.String(): cfg.Upload.DirImportShopProduct,
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

	eventBus := bus.New()
	sqlstore.Init(db)
	sqlstore.AddEventBus(eventBus)
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

	locationBus := servicelocation.New(db).MessageBus()
	locationAggr := servicelocation.NewAggregate(db).MessageBus()
	identityQuery = serviceidentity.NewQueryService(db).MessageBus()
	if cfg.GHN.AccountDefault.Token != "" {
		ghnCarrier = ghn.New(cfg.GHN, locationBus)
		if err := ghnCarrier.InitAllClients(ctx); err != nil {
			ll.Fatal("Unable to connect to GHN", l.Error(err))
		}
	} else {
		if cmenv.IsDev() {
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
		if cmenv.IsDev() {
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
		if cmenv.IsDev() {
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
		if cmenv.IsDev() {
			ll.Warn("DEVELOPMENT. Skip connecting to ahamove.")
		} else {
			ll.Fatal("ahamove: No token")
		}
	}

	shippingManager := shipping_provider.NewCtrl(locationBus, ghnCarrier, ghtkCarrier, vtpostCarrier)

	authStore := auth.NewGenerator(redisStore)
	imcsvorder.Init(locationBus, shutdowner, redisStore, uploader, db)
	imcsvproduct.Init(shutdowner, redisStore, uploader, db)
	export.Init(shutdowner, redisStore, eventStreamer, export.Config{
		UrlPrefix: cfg.Export.URLPrefix,
		DirExport: cfg.Export.DirExport,
	})

	// create aggregate, query service
	summaryQuery := summaryquery.NewDashboardQuery(db, redisStore, locationBus).MessageBus()

	identityQuery = serviceidentity.NewQueryService(db).MessageBus()
	identityPM := identitypm.New(identityQuery, &invitationQuery)
	identityPM.RegisterEventHandlers(eventBus)
	catalogQuery := catalogquery.New(db).MessageBus()
	catalogAggr := catalogaggregate.New(eventBus, db).MessageBus()
	catalogPm := catalogpm.New(eventBus, catalogQuery, catalogAggr)
	catalogPm.RegisterEventHandlers(eventBus)

	addressQuery := address.NewQueryService(db).MessageBus()
	shipnowQuery = serviceshipnow.NewQueryService(db).MessageBus()
	orderQuery = serviceordering.NewQueryService(db).MessageBus()

	orderAggr = serviceordering.NewAggregate(eventBus, db)
	shipnowCarrierManager := shipnowcarrier.NewManager(db, locationBus, &shipnowcarrier.Carrier{
		ShipnowCarrier:        ahamoveCarrier,
		ShipnowCarrierAccount: ahamoveCarrierAccount,
	}, shipnowQuery)
	identityAggr := serviceidentity.NewAggregate(db, shipnowCarrierManager).MessageBus()
	shipnowAggr = serviceshipnow.NewAggregate(eventBus, db, locationBus, identityQuery, addressQuery, orderQuery, shipnowCarrierManager).MessageBus()

	shipnowPM := shipnowpm.New(eventBus, shipnowQuery, shipnowAggr, orderAggr.MessageBus(), shipnowCarrierManager)
	shipnowPM.RegisterEventHandlers(eventBus)

	stocktakeAggr := stocktakeaggregate.NewAggregateStocktake(db, eventBus).MessageBus()
	stocktakeQuery := stocktakequery.NewQueryStocktake(db).MessageBus()
	customerAggr := customeraggregate.NewCustomerAggregate(eventBus, db).MessageBus()
	supplierAggr := supplieraggregate.NewSupplierAggregate(eventBus, db).MessageBus()
	carrierAggr := carrieraggregate.NewCarrierAggregate(eventBus, db).MessageBus()
	traderAddressAggr := customeraggregate.NewAddressAggregate(db).MessageBus()
	traderAgg := traderAgg.NewTraderAgg(db).MessageBus()
	traderPM := traderpm.New(eventBus, traderAgg)
	traderPM.RegisterEventHandlers(eventBus)
	customerQuery := customerquery.NewCustomerQuery(db).MessageBus()
	supplierQuery := supplierquery.NewSupplierQuery(db).MessageBus()
	carrierQuery := carrierquery.NewCarrierQuery(db).MessageBus()
	traderQuery := traderquery.NewTraderQuery(db, customerQuery, carrierQuery, supplierQuery).MessageBus()
	traderAddressQuery := customerquery.NewAddressQuery(db).MessageBus()
	affiliateCmd := serviceaffiliate.NewAggregate(dbaff, identityQuery, catalogQuery, orderQuery).MessageBus()
	affilateQuery := serviceaffiliate.NewQuery(dbaff).MessageBus()
	affiliatePM := affiliatepm.New(affiliateCmd)
	affiliatePM.RegisterEventHandlers(eventBus)

	ledgerAggr := ledgeraggregate.NewLedgerAggregate(db, &receiptQuery).MessageBus()
	ledgerQuery := ledgerquery.NewLedgerQuery(db).MessageBus()
	ledgerPM := ledgerpm.New(eventBus, ledgerAggr)
	ledgerPM.RegisterEventHandlers(eventBus)

	inventoryQuery := inventoryquery.NewQueryInventory(stocktakeQuery, eventBus, db).MessageBus()
	purchaseOrderAggr := purchaseorderaggregate.NewPurchaseOrderAggregate(db, eventBus, catalogQuery, supplierQuery, inventoryQuery).MessageBus()
	purchaseOrderQuery := purchaseorderquery.NewPurchaseOrderQuery(db, eventBus, supplierQuery, inventoryQuery, &receiptQuery).MessageBus()

	purchaseOrderPM := purchaseorderpm.New(&purchaseOrderQuery, &receiptQuery)
	purchaseOrderPM.RegisterEventHandlers(eventBus)
	refundAggr := refundaggregate.NewRefundAggregate(db, eventBus).MessageBus()
	refundQuery := refundquery.NewQueryRefund(eventBus, db).MessageBus()

	purchaseRefundAggr := purchaserefundaggregate.NewPurchaseRefundAggregate(db, eventBus, purchaseOrderQuery).MessageBus()
	purchaseRefundQuery := purchaserefundquery.NewQueryPurchasePurchaseRefund(eventBus, db).MessageBus()

	inventoryAggr := inventoryaggregate.NewAggregateInventory(eventBus, db, traderQuery, purchaseOrderQuery, stocktakeQuery, refundQuery, purchaseRefundQuery).MessageBus()
	inventoryPm := inventorypm.New(eventBus, catalogQuery, orderQuery, inventoryAggr)
	inventoryPm.RegisterEventHandlers(eventBus)

	receiptAggr := receiptaggregate.NewReceiptAggregate(db, eventBus, traderQuery, ledgerQuery, orderQuery, customerQuery, carrierQuery, supplierQuery, purchaseOrderQuery).MessageBus()
	receiptQuery = receiptquery.NewReceiptQuery(db).MessageBus()
	receiptPM := receiptpm.New(eventBus, receiptQuery, receiptAggr, ledgerQuery, ledgerAggr, identityQuery)
	receiptPM.RegisterEventHandlers(eventBus)

	// payment
	var vtpayProvider *vtpay.Provider
	if cfg.VTPay.MerchantCode != "" {
		vtpayClient = vtpayclient.New(cfg.VTPay)
		vtpayProvider = vtpay.New(cfg.VTPay)
	}
	paymentManager := servicepaymentmanager.NewManager(vtpayProvider, orderQuery).MesssageBus()
	orderPM := serviceorderingpm.New(orderAggr.MessageBus(), affiliateCmd, receiptQuery, inventoryAggr, orderQuery, customerQuery)
	orderPM.RegisterEventHandlers(eventBus)
	refundPm := refundpm.New(&refundQuery, &receiptQuery, &refundAggr)
	refundPm.RegisterEventHandlers(eventBus)
	invitationAggr := invitationaggregate.NewInvitationAggregate(db, cfg.Invitation.Secret, customerQuery, identityQuery, eventBus, cfg).MessageBus()
	invitationQuery = invitationquery.NewInvitationQuery(db).MessageBus()
	invitationPM := invitationpm.New(eventBus, invitationQuery, invitationAggr)
	invitationPM.RegisterEventHandlers(eventBus)
	purchaseRefundPM := purchaserefundpm.New(&purchaseRefundAggr, &purchaseRefundQuery, &receiptQuery)
	purchaseRefundPM.RegisterEventHandlers(eventBus)
	authorizationQuery := authorizationquery.NewAuthorizationQuery().MessageBus()
	authorizationAggregate := authorizationaggregate.NewAuthorizationAggregate().MessageBus()

	authorizeauth.SetMode(cfg.FlagEnablePermission)

	smsArg := smsAgg.NewSmsLogAggregate(eventBus, dbLogs).MessageBus()
	connectionQuery := connectionquery.NewConnectionQuery(db).MessageBus()
	connectionAggregate := connectionaggregate.NewConnectionAggregate(db, eventBus).MessageBus()
	shipmentServiceAggr := shipmentservice.NewAggregate(db, redisStore).MessageBus()
	shipmentServiceQuery := shipmentservice.NewQueryService(db, redisStore).MessageBus()
	shipmentPriceListAggr := pricelist.NewAggregate(db, eventBus).MessageBus()
	shipmentPriceListQuery := pricelist.NewQueryService(db, redisStore).MessageBus()
	shipmentPriceAggr := shipmentprice.NewAggregate(db, redisStore).MessageBus()
	shipmentPriceQuery := shipmentprice.NewQueryService(db, redisStore, locationBus, shipmentPriceListQuery).MessageBus()
	shipmentPriceListPM := pricelistpm.New(redisStore)
	shipmentPriceListPM.RegisterEventHandlers(eventBus)

	shipmentManager = shippingcarrier.NewShipmentManager(locationBus, connectionQuery, connectionAggregate, redisStore, shipmentServiceQuery, shipmentPriceQuery, cfg.FlagApplyShipmentPrice)
	shipmentManager.SetWebhookEndpoint(connection_type.ConnectionProviderGHN, cfg.GHNWebhook.Endpoint)
	shippingAggr := shippingaggregate.NewAggregate(db, locationBus, orderQuery, shipmentManager, connectionQuery, eventBus).MessageBus()
	shippingQuery := shippingquery.NewQueryService(db).MessageBus()
	shippingPM := shippingpm.New(eventBus, shippingQuery, shippingAggr, redisStore)
	shippingPM.RegisterEventHandlers(eventBus)

	moneyTxQuery = moneytxquery.NewMoneyTxQuery(db, shippingQuery).MessageBus()
	moneyTxAggr = moneytxaggregate.NewMoneyTxAggregate(db, shippingQuery, identityQuery, eventBus).MessageBus()

	dbWebServer, err = cmsql.Connect(cfg.PostgresWebServer)
	if err != nil {
		ll.Fatal("Unable to connect to Postgres", l.Error(err))
	}
	webServerAggregate := webserveraggregate.New(eventBus, dbWebServer, catalogQuery).MessageBus()
	webServerQuery := webserverquery.New(eventBus, dbWebServer, catalogQuery).MessageBus()

	moneyTxPM := moneytxpm.New(eventBus, moneyTxQuery, moneyTxAggr, shippingQuery)
	moneyTxPM.RegisterEvenHandlers(eventBus)

	whiteLabel := wl.Init(cmenv.Env())
	if err := whiteLabel.VerifyPartners(context.Background(), identityQuery); err != nil {
		ll.Fatal("error loading white label partners", l.Error(err))
	}

	middleware.Init(cfg.SAdminToken, identityQuery)
	sms.Init(smsArg)
	api.Init(
		eventBus,
		smsArg,
		identityAggr,
		identityQuery,
		invitationAggr,
		invitationQuery,
		authorizationQuery,
		authorizationAggregate,
		shutdowner,
		redisStore,
		authStore,
		cfg.Email,
		cfg.SMS,
		bot,
	)
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
		customerAggr,
		customerQuery,
		traderAddressAggr,
		traderAddressQuery,
		orderAggr.MessageBus(),
		orderQuery,
		paymentManager,
		supplierAggr,
		supplierQuery,
		carrierAggr,
		carrierQuery,
		traderQuery,
		eventBus,
		receiptAggr,
		receiptQuery,
		shutdowner,
		redisStore,
		inventoryAggr,
		inventoryQuery,
		ledgerAggr,
		ledgerQuery,
		purchaseOrderAggr,
		purchaseOrderQuery,
		summaryQuery,
		stocktakeQuery,
		stocktakeAggr,
		shipmentManager,
		shippingAggr,
		refundAggr,
		refundQuery,
		purchaseRefundAggr,
		purchaseRefundQuery,
		connectionQuery,
		connectionAggregate,
		shippingQuery,
		webServerAggregate,
		webServerQuery,
	)
	partner.Init(
		shutdowner,
		redisStore,
		authStore,
		cfg.URL.Auth,
		locationBus,
		&customerQuery,
		&customerAggr,
		&traderAddressQuery,
		&traderAddressAggr,
		&inventoryQuery,
		&catalogQuery,
		&catalogAggr,
		connectionQuery,
		connectionAggregate,
		shippingAggr,
	)
	whitelabelapix.Init(db, &catalogAggr)
	xshop.Init(
		shutdowner,
		redisStore,
		locationBus,
		&customerQuery,
		&customerAggr,
		&traderAddressQuery,
		&traderAddressAggr,
		&inventoryQuery,
		&catalogQuery,
		&catalogAggr)
	xshopping.Init(
		locationBus,
		&customerQuery,
		&customerAggr,
		&traderAddressQuery,
		&traderAddressAggr,
		&inventoryQuery,
		&catalogQuery,
		&catalogAggr)
	integration.Init(shutdowner, redisStore, authStore)
	webhook.Init(ctlProducer, redisStore)
	xshipping.Init(shippingManager, ordersqlstore.NewOrderStore(db), shipsqlstore.NewFulfillmentStore(db), shipmentManager, shippingAggr, shippingQuery, connectionQuery)
	orderS.Init(shippingManager, catalogQuery, orderAggr.MessageBus(),
		customerAggr, customerQuery, traderAddressAggr, traderAddressQuery, locationBus, eventBus, shipmentManager)
	affiliate.Init(identityAggr)
	apiaff.Init(affiliateCmd, affilateQuery, catalogQuery, identityQuery)
	admin.Init(eventBus, moneyTxQuery, moneyTxAggr, connectionAggregate, connectionQuery, identityQuery, shipmentPriceAggr, shipmentPriceQuery, shipmentServiceAggr, shipmentServiceQuery, shipmentPriceListAggr, shipmentPriceListQuery, locationAggr, locationBus, shipmentManager)

	err = db.GetSchemaErrors()
	if err != nil && cmenv.IsDev() {
		ll.Error("Fail to verify Database", l.Error(err))
	} else if err != nil {
		// should move struct `callback` out of etop/model before change to ll.Fatal
		ll.Error("Fail to verify Database", l.Error(err))
	}

	ss = session.New(
		session.OptValidator(tokens.NewTokenStore(redisStore)),
		session.OptSuperAdmin(cfg.SAdminToken),
	)
	hooks = session.NewHook(acl.GetACL()).Build()

	svrs := startServers(webServerQuery)
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

	if cfg.SMS.Enabled {
		smsBot := cfg.TelegramBot.MustConnectChannel(config.ChannelSMS)

		var imgroupSMSClient *imgroupsms.Client
		if cfg.WhiteLabel.IMGroup.SMS.APIKey != "" {
			imgroupSMSClient = imgroupsms.New(cfg.WhiteLabel.IMGroup.SMS)
		} else if !cmenv.IsDev() {
			ll.Panic("no sms config for whitelabel/imgroup")
		}

		sms.New(cfg.SMS, smsBot, imgroupSMSClient).Register(bus.Global())
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
