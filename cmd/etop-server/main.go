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
	shippingcore "o.o/api/main/shipping"
	subscriptioncore "o.o/api/subscripting/subscription"
	"o.o/api/top/types/etc/connection_type"
	"o.o/backend/cmd/etop-server/config"
	smsAgg "o.o/backend/com/etc/logging/smslog/aggregate"
	servicepaymentmanager "o.o/backend/com/external/payment/manager"
	paymentaggregate "o.o/backend/com/external/payment/payment/aggregate"
	"o.o/backend/com/handler/etop-handler/intctl"
	"o.o/backend/com/main/address"
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
	stocktakequery "o.o/backend/com/main/stocktaking/query"
	serviceaffiliate "o.o/backend/com/services/affiliate"
	affiliatepm "o.o/backend/com/services/affiliate/pm"
	carrierquery "o.o/backend/com/shopping/carrying/query"
	customeraggregate "o.o/backend/com/shopping/customering/aggregate"
	customerquery "o.o/backend/com/shopping/customering/query"
	supplierquery "o.o/backend/com/shopping/suppliering/query"
	traderAgg "o.o/backend/com/shopping/tradering/aggregate"
	traderpm "o.o/backend/com/shopping/tradering/pm"
	traderquery "o.o/backend/com/shopping/tradering/query"
	"o.o/backend/com/subscripting/subscription"
	subscriptionpm "o.o/backend/com/subscripting/subscription/pm"
	"o.o/backend/com/subscripting/subscriptionbill"
	"o.o/backend/com/subscripting/subscriptionplan"
	"o.o/backend/com/subscripting/subscriptionproduct"
	webserveraggregate "o.o/backend/com/web/webserver/aggregate"
	webserverpm "o.o/backend/com/web/webserver/pm"
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
	"o.o/backend/pkg/etop/api/affiliate"
	"o.o/backend/pkg/etop/api/export"
	"o.o/backend/pkg/etop/api/integration"
	"o.o/backend/pkg/etop/apix/partner"
	"o.o/backend/pkg/etop/apix/partnercarrier"
	xshopping "o.o/backend/pkg/etop/apix/shopping"
	"o.o/backend/pkg/etop/apix/webhook"
	whitelabelapix "o.o/backend/pkg/etop/apix/whitelabel"
	"o.o/backend/pkg/etop/authorize/middleware"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/authorize/tokens"
	"o.o/backend/pkg/etop/eventstream"
	imcsvorder "o.o/backend/pkg/etop/logic/orders/imcsv"
	imcsvproduct "o.o/backend/pkg/etop/logic/products/imcsv"
	"o.o/backend/pkg/etop/logic/summary"
	"o.o/backend/pkg/etop/middlewares"
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
	shippingAggr shippingcore.CommandBus

	identityQuery identity.QueryBus

	invitationQuery invitation.QueryBus

	vtpayClient *vtpayclient.Client

	receiptQuery      receipting.QueryBus
	shipmentManager   *shippingcarrier.ShipmentManager
	moneyTxQuery      moneytx.QueryBus
	moneyTxAggr       moneytx.CommandBus
	subscriptionQuery subscriptioncore.QueryBus

	ss    *session.Session
	hooks httprpc.HooksBuilder

	servers    []httprpc.Server
	serversExt []httprpc.Server
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

	locationBus := servicelocation.QueryMessageBus(servicelocation.New(db))
	identityQuery = serviceidentity.QueryServiceMessageBus(serviceidentity.NewQueryService(db))
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

	authStore := auth.NewGenerator(redisStore)
	imcsvorder.Init(locationBus, shutdowner, redisStore, uploader, db)
	imcsvproduct.Init(shutdowner, redisStore, uploader, db)
	export.Init(shutdowner, redisStore, eventStreamer, export.Config{
		UrlPrefix: cfg.Export.URLPrefix,
		DirExport: cfg.Export.DirExport,
	})

	// create aggregate, query service
	identityQuery = serviceidentity.QueryServiceMessageBus(serviceidentity.NewQueryService(db))
	identityPM := identitypm.New(identityQuery, invitationQuery)
	identityPM.RegisterEventHandlers(eventBus)
	catalogQuery := catalogquery.QueryServiceMessageBus(catalogquery.New(db))
	catalogAggr := catalogaggregate.AggregateMessageBus(catalogaggregate.New(eventBus, db))
	catalogPm := catalogpm.New(eventBus, catalogQuery, catalogAggr)
	catalogPm.RegisterEventHandlers(eventBus)

	addressQuery := address.QueryServiceMessageBus(address.NewQueryService(db))
	shipnowQuery = serviceshipnow.QueryServiceMessageBus(serviceshipnow.NewQueryService(db))
	orderQuery = serviceordering.QueryServiceMessageBus(serviceordering.NewQueryService(db))

	orderAggr = serviceordering.NewAggregate(eventBus, db)
	shipnowCarrierManager := shipnowcarrier.NewManager(db, locationBus, &shipnowcarrier.Carrier{
		ShipnowCarrier:        ahamoveCarrier,
		ShipnowCarrierAccount: ahamoveCarrierAccount,
	}, shipnowQuery)
	identityAggr := serviceidentity.AggregateMessageBus(serviceidentity.NewAggregate(db, shipnowCarrierManager))
	shipnowAggr = serviceshipnow.AggregateMessageBus(serviceshipnow.NewAggregate(eventBus, db, locationBus, identityQuery, addressQuery, orderQuery, shipnowCarrierManager))

	shipnowPM := shipnowpm.New(eventBus, shipnowQuery, shipnowAggr, serviceordering.AggregateMessageBus(orderAggr), shipnowCarrierManager)
	shipnowPM.RegisterEventHandlers(eventBus)

	stocktakeQuery := stocktakequery.StocktakeQueryMessageBus(stocktakequery.NewQueryStocktake(db))
	customerAggr := customeraggregate.CustomerAggregateMessageBus(customeraggregate.NewCustomerAggregate(eventBus, db))
	traderAddressAggr := customeraggregate.AddressAggregateMessageBus(customeraggregate.NewAddressAggregate(db))
	traderAgg := traderAgg.TraderAggMessageBus(traderAgg.NewTraderAgg(db))
	traderPM := traderpm.New(eventBus, traderAgg)
	traderPM.RegisterEventHandlers(eventBus)
	customerQuery := customerquery.CustomerQueryMessageBus(customerquery.NewCustomerQuery(db))
	supplierQuery := supplierquery.SupplierQueryMessageBus(supplierquery.NewSupplierQuery(db))
	carrierQuery := carrierquery.CarrierQueryMessageBus(carrierquery.NewCarrierQuery(db))
	traderQuery := traderquery.TraderQueryMessageBus(traderquery.NewTraderQuery(db, customerQuery, carrierQuery, supplierQuery))
	traderAddressQuery := customerquery.AddressQueryMessageBus(customerquery.NewAddressQuery(db))
	affiliateCmd := serviceaffiliate.AggregateMessageBus(serviceaffiliate.NewAggregate(dbaff, identityQuery, catalogQuery, orderQuery))
	affiliatePM := affiliatepm.New(affiliateCmd)
	affiliatePM.RegisterEventHandlers(eventBus)

	ledgerAggr := ledgeraggregate.LedgerAggregateMessageBus(ledgeraggregate.NewLedgerAggregate(db, receiptQuery))
	ledgerQuery := ledgerquery.LedgerQueryMessageBus(ledgerquery.NewLedgerQuery(db))
	ledgerPM := ledgerpm.New(eventBus, ledgerAggr)
	ledgerPM.RegisterEventHandlers(eventBus)

	inventoryQuery := inventoryquery.InventoryQueryServiceMessageBus(inventoryquery.NewQueryInventory(stocktakeQuery, eventBus, db))
	purchaseOrderQuery := purchaseorderquery.PurchaseOrderQueryMessageBus(purchaseorderquery.NewPurchaseOrderQuery(db, eventBus, supplierQuery, inventoryQuery, receiptQuery))

	purchaseOrderPM := purchaseorderpm.New(purchaseOrderQuery, receiptQuery)
	purchaseOrderPM.RegisterEventHandlers(eventBus)
	refundAggr := refundaggregate.RefundAggregateMessageBus(refundaggregate.NewRefundAggregate(db, eventBus))
	refundQuery := refundquery.RefundQueryServiceMessageBus(refundquery.NewQueryRefund(eventBus, db))

	purchaseRefundAggr := purchaserefundaggregate.PurchaseRefundAggregateMessageBus(purchaserefundaggregate.NewPurchaseRefundAggregate(db, eventBus, purchaseOrderQuery))
	purchaseRefundQuery := purchaserefundquery.PurchaseRefundQueryServiceMessageBus(purchaserefundquery.NewQueryPurchasePurchaseRefund(eventBus, db))

	inventoryAggr := inventoryaggregate.InventoryAggregateMessageBus(inventoryaggregate.NewAggregateInventory(eventBus, db, traderQuery, purchaseOrderQuery, stocktakeQuery, refundQuery, purchaseRefundQuery))
	inventoryPm := inventorypm.New(eventBus, catalogQuery, orderQuery, inventoryAggr)
	inventoryPm.RegisterEventHandlers(eventBus)

	receiptAggr := receiptaggregate.ReceiptAggregateMessageBus(receiptaggregate.NewReceiptAggregate(db, eventBus, traderQuery, ledgerQuery, orderQuery, customerQuery, carrierQuery, supplierQuery, purchaseOrderQuery))
	receiptQuery = receiptquery.ReceiptQueryMessageBus(receiptquery.NewReceiptQuery(db))
	receiptPM := receiptpm.New(eventBus, receiptQuery, receiptAggr, ledgerQuery, ledgerAggr, identityQuery)
	receiptPM.RegisterEventHandlers(eventBus)

	// payment
	var vtpayProvider *vtpay.Provider
	if cfg.VTPay.MerchantCode != "" {
		vtpayClient = vtpayclient.New(cfg.VTPay)
		vtpayProvider = vtpay.New(cfg.VTPay)
	}
	paymentAggr := paymentaggregate.AggregateMessageBus(paymentaggregate.NewAggregate(db))
	paymentManager := servicepaymentmanager.NewManager(vtpayProvider, orderQuery).MesssageBus()
	orderPM := serviceorderingpm.New(serviceordering.AggregateMessageBus(orderAggr), affiliateCmd, receiptQuery, inventoryAggr, orderQuery, customerQuery)
	orderPM.RegisterEventHandlers(eventBus)
	refundPm := refundpm.New(refundQuery, receiptQuery, refundAggr)
	refundPm.RegisterEventHandlers(eventBus)
	invitationAggr := invitationaggregate.InvitationAggregateMessageBus(invitationaggregate.NewInvitationAggregate(db, cfg.Invitation, customerQuery, identityQuery, eventBus, cfg))
	invitationQuery = invitationquery.InvitationQueryMessageBus(invitationquery.NewInvitationQuery(db))
	invitationPM := invitationpm.New(eventBus, invitationQuery, invitationAggr)
	invitationPM.RegisterEventHandlers(eventBus)
	purchaseRefundPM := purchaserefundpm.New(purchaseRefundAggr, purchaseRefundQuery, receiptQuery)
	purchaseRefundPM.RegisterEventHandlers(eventBus)

	smsArg := smsAgg.SmsLogAggregateMessageBus(smsAgg.NewSmsLogAggregate(eventBus, dbLogs))
	connectionQuery := connectionquery.ConnectionQueryMessageBus(connectionquery.NewConnectionQuery(db))
	connectionAggregate := connectionaggregate.ConnectionAggregateMessageBus(connectionaggregate.NewConnectionAggregate(db, eventBus))
	shipmentServiceQuery := shipmentservice.QueryServiceMessageBus(shipmentservice.NewQueryService(db, redisStore))
	shipmentPriceListQuery := pricelist.QueryServiceMessageBus(pricelist.NewQueryService(db, redisStore))
	shipmentPriceQuery := shipmentprice.QueryServiceMessageBus(shipmentprice.NewQueryService(db, redisStore, locationBus, shipmentPriceListQuery))
	shipmentPriceListPM := pricelistpm.New(redisStore)
	shipmentPriceListPM.RegisterEventHandlers(eventBus)

	shipmentManager = shippingcarrier.NewShipmentManager(eventBus, locationBus, connectionQuery, connectionAggregate, redisStore, shipmentServiceQuery, shipmentPriceQuery, cfg.FlagApplyShipmentPrice)
	shipmentManager.SetWebhookEndpoint(connection_type.ConnectionProviderGHN, cfg.GHNWebhook.Endpoint)

	botCarrier := cfg.TelegramBot.MustConnectChannel(config.ChannelShipmentCarrier)
	shippingQuery := shippingquery.QueryServiceMessageBus(shippingquery.NewQueryService(db))
	shippingAggr := shippingaggregate.AggregateMessageBus(shippingaggregate.NewAggregate(db, eventBus, locationBus, botCarrier, orderQuery, shipmentManager, connectionQuery))
	shippingPM := shippingpm.New(eventBus, shippingQuery, shippingAggr, redisStore)
	shippingPM.RegisterEventHandlers(eventBus)

	moneyTxQuery = moneytxquery.MoneyTxQueryMessageBus(moneytxquery.NewMoneyTxQuery(db, shippingQuery))
	moneyTxAggr = moneytxaggregate.MoneyTxAggregateMessageBus(moneytxaggregate.NewMoneyTxAggregate(db, shippingQuery, identityQuery, eventBus))

	dbWebServer, err = cmsql.Connect(cfg.PostgresWebServer)
	if err != nil {
		ll.Fatal("Unable to connect to Postgres", l.Error(err))
	}
	webServerAggregate := webserveraggregate.WebserverAggregateMessageBus(webserveraggregate.New(eventBus, dbWebServer, catalogQuery))
	webServerQuery := webserverquery.WebserverQueryServiceMessageBus(webserverquery.New(eventBus, dbWebServer, catalogQuery))
	webserverPm := webserverpm.New(eventBus, webServerAggregate, webServerQuery)
	webserverPm.RegisterEventHandlers(eventBus)

	moneyTxPM := moneytxpm.New(eventBus, moneyTxQuery, moneyTxAggr, shippingQuery)
	moneyTxPM.RegisterEvenHandlers(eventBus)

	whiteLabel := wl.Init(cmenv.Env())
	if err := whiteLabel.VerifyPartners(context.Background(), identityQuery); err != nil {
		ll.Fatal("error loading white label partners", l.Error(err))
	}

	subrProductQuery := subscriptionproduct.SubrProductQueryMessageBus(subscriptionproduct.NewSubrProductQuery(db))
	subrPlanQuery := subscriptionplan.SubrPlanQueryMessageBus(subscriptionplan.NewSubrPlanQuery(db, subrProductQuery))
	subscriptionQuery = subscription.SubscriptionQueryMessageBus(subscription.NewSubscriptionQuery(db, subrPlanQuery, subrProductQuery))
	subscriptionAggr := subscription.SubscriptionAggregateMessageBus(subscription.NewSubscriptionAggregate(db))
	subrBillAggr := subscriptionbill.SubrBillAggregateMessageBus(subscriptionbill.NewSubrBillAggregate(db, eventBus, paymentAggr, subscriptionQuery, subrPlanQuery))
	subrBillQuery := subscriptionbill.SubrBillQueryMessageBus(subscriptionbill.NewSubrBillQuery(db))
	subscriptionPM := subscriptionpm.New(
		subrBillQuery, subrBillAggr,
		subscriptionQuery, subscriptionAggr,
		subrPlanQuery,
		identityQuery,
	)
	subscriptionPM.RegisterEventHandlers(eventBus)

	middleware.Init(cfg.SAdminToken, identityQuery)
	sms.Init(smsArg)
	servers = append(servers, BuildServers(
		db,
		cfg,
		bot,
		shutdowner,
		eventBus,
		redisStore,
		authStore,
		ss,
		shipnowCarrierManager,
		paymentManager,
		partner.AuthURL(cfg.URL.Auth),
	)...)
	partnercarrier.Init(
		shutdowner,
		redisStore,
		connectionQuery,
		connectionAggregate,
		shippingQuery,
		shippingAggr,
	)
	whitelabelapix.Init(db, catalogAggr)
	xshopping.Init(
		locationBus,
		customerQuery,
		customerAggr,
		traderAddressQuery,
		traderAddressAggr,
		inventoryQuery,
		catalogQuery,
		catalogAggr)
	integration.Init(shutdowner, redisStore, authStore)
	webhook.Init(ctlProducer, redisStore)
	affiliate.Init(identityAggr)

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
	hooks = httprpc.ChainHooks(
		middlewares.NewLogging(),
		session.NewHook(acl.GetACL()),
	)

	svrs := startServers(webServerQuery, catalogQuery, redisStore, locationBus)
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
