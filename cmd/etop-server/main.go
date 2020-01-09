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

	"etop.vn/api/main/identity"
	"etop.vn/api/main/invitation"
	"etop.vn/api/main/moneytx"
	"etop.vn/api/main/ordering"
	"etop.vn/api/main/receipting"
	"etop.vn/api/main/shipnow"
	"etop.vn/api/top/types/etc/connection_type"
	"etop.vn/backend/cmd/etop-server/config"
	smsAgg "etop.vn/backend/com/etc/logging/smslog/aggregate"
	servicepaymentmanager "etop.vn/backend/com/external/payment/manager"
	"etop.vn/backend/com/handler/etop-handler/intctl"
	"etop.vn/backend/com/main/address"
	authorizationaggregate "etop.vn/backend/com/main/authorization/aggregate"
	authorizationquery "etop.vn/backend/com/main/authorization/query"
	catalogaggregate "etop.vn/backend/com/main/catalog/aggregate"
	catalogpm "etop.vn/backend/com/main/catalog/pm"
	catalogquery "etop.vn/backend/com/main/catalog/query"
	connectionaggregate "etop.vn/backend/com/main/connectioning/aggregate"
	connectionquery "etop.vn/backend/com/main/connectioning/query"
	serviceidentity "etop.vn/backend/com/main/identity"
	identitypm "etop.vn/backend/com/main/identity/pm"
	inventoryaggregate "etop.vn/backend/com/main/inventory/aggregate"
	inventorypm "etop.vn/backend/com/main/inventory/pm"
	inventoryquery "etop.vn/backend/com/main/inventory/query"
	invitationaggregate "etop.vn/backend/com/main/invitation/aggregate"
	invitationpm "etop.vn/backend/com/main/invitation/pm"
	invitationquery "etop.vn/backend/com/main/invitation/query"
	ledgeraggregate "etop.vn/backend/com/main/ledgering/aggregate"
	ledgerpm "etop.vn/backend/com/main/ledgering/pm"
	ledgerquery "etop.vn/backend/com/main/ledgering/query"
	servicelocation "etop.vn/backend/com/main/location"
	moneytxaggregate "etop.vn/backend/com/main/moneytx/aggregate"
	moneytxpm "etop.vn/backend/com/main/moneytx/pm"
	moneytxquery "etop.vn/backend/com/main/moneytx/query"
	serviceordering "etop.vn/backend/com/main/ordering"
	serviceorderingpm "etop.vn/backend/com/main/ordering/pm"
	ordersqlstore "etop.vn/backend/com/main/ordering/sqlstore"
	purchaseorderaggregate "etop.vn/backend/com/main/purchaseorder/aggregate"
	purchaseorderpm "etop.vn/backend/com/main/purchaseorder/pm"
	purchaseorderquery "etop.vn/backend/com/main/purchaseorder/query"
	purchaserefundaggregate "etop.vn/backend/com/main/purchaserefund/aggregate"
	purchaserefundpm "etop.vn/backend/com/main/purchaserefund/pm"
	purchaserefundquery "etop.vn/backend/com/main/purchaserefund/query"
	receiptaggregate "etop.vn/backend/com/main/receipting/aggregate"
	receiptpm "etop.vn/backend/com/main/receipting/pm"
	receiptquery "etop.vn/backend/com/main/receipting/query"
	refundaggregate "etop.vn/backend/com/main/refund/aggregate"
	refundpm "etop.vn/backend/com/main/refund/pm"
	refundquery "etop.vn/backend/com/main/refund/query"
	"etop.vn/backend/com/main/shipmentpricing/pricelist"
	pricelistpm "etop.vn/backend/com/main/shipmentpricing/pricelist/pm"
	"etop.vn/backend/com/main/shipmentpricing/shipmentprice"
	"etop.vn/backend/com/main/shipmentpricing/shipmentservice"
	serviceshipnow "etop.vn/backend/com/main/shipnow"
	shipnowcarrier "etop.vn/backend/com/main/shipnow-carrier"
	shipnowpm "etop.vn/backend/com/main/shipnow/pm"
	shippingaggregate "etop.vn/backend/com/main/shipping/aggregate"
	shippingcarrier "etop.vn/backend/com/main/shipping/carrier"
	shippingpm "etop.vn/backend/com/main/shipping/pm"
	shippingquery "etop.vn/backend/com/main/shipping/query"
	shipsqlstore "etop.vn/backend/com/main/shipping/sqlstore"
	stocktakeaggregate "etop.vn/backend/com/main/stocktaking/aggregate"
	stocktakequery "etop.vn/backend/com/main/stocktaking/query"
	serviceaffiliate "etop.vn/backend/com/services/affiliate"
	affiliatepm "etop.vn/backend/com/services/affiliate/pm"
	carrieraggregate "etop.vn/backend/com/shopping/carrying/aggregate"
	carrierquery "etop.vn/backend/com/shopping/carrying/query"
	customeraggregate "etop.vn/backend/com/shopping/customering/aggregate"
	customerquery "etop.vn/backend/com/shopping/customering/query"
	supplieraggregate "etop.vn/backend/com/shopping/suppliering/aggregate"
	supplierquery "etop.vn/backend/com/shopping/suppliering/query"
	traderAgg "etop.vn/backend/com/shopping/tradering/aggregate"
	traderpm "etop.vn/backend/com/shopping/tradering/pm"
	traderquery "etop.vn/backend/com/shopping/tradering/query"
	summaryquery "etop.vn/backend/com/summary/query"
	webserveraggregate "etop.vn/backend/com/web/webserver/aggregate"
	webserverquery "etop.vn/backend/com/web/webserver/query"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/captcha"
	"etop.vn/backend/pkg/common/apifw/health"
	cmService "etop.vn/backend/pkg/common/apifw/service"
	"etop.vn/backend/pkg/common/apifw/whitelabel/wl"
	"etop.vn/backend/pkg/common/authorization/auth"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmenv"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/extservice/telebot"
	"etop.vn/backend/pkg/common/mq"
	"etop.vn/backend/pkg/common/redis"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sqltrace"
	"etop.vn/backend/pkg/etop/api"
	"etop.vn/backend/pkg/etop/api/admin"
	"etop.vn/backend/pkg/etop/api/affiliate"
	"etop.vn/backend/pkg/etop/api/export"
	"etop.vn/backend/pkg/etop/api/integration"
	"etop.vn/backend/pkg/etop/api/shop"
	"etop.vn/backend/pkg/etop/apix/partner"
	xshipping "etop.vn/backend/pkg/etop/apix/shipping"
	xshop "etop.vn/backend/pkg/etop/apix/shop"
	xshopping "etop.vn/backend/pkg/etop/apix/shopping"
	"etop.vn/backend/pkg/etop/apix/webhook"
	whitelabelapix "etop.vn/backend/pkg/etop/apix/whitelabel"
	authorizeauth "etop.vn/backend/pkg/etop/authorize/auth"
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
	imgroupsms "etop.vn/backend/pkg/integration/sms/imgroup"
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
