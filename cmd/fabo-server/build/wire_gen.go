// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package build

import (
	"context"
	"o.o/api/main/accountshipnow"
	"o.o/api/services/affiliate"
	"o.o/api/shopping/tradering"
	"o.o/backend/cmd/fabo-server/config"
	"o.o/backend/cogs/config/_server"
	"o.o/backend/cogs/database/_min"
	fabo2 "o.o/backend/cogs/server/fabo"
	"o.o/backend/cogs/server/shop"
	"o.o/backend/cogs/shipment/_fabo"
	v2_2 "o.o/backend/cogs/shipment/ghn/v2"
	"o.o/backend/cogs/sms/_min"
	"o.o/backend/cogs/storage/_all"
	"o.o/backend/cogs/uploader"
	"o.o/backend/com/etc/logging/shippingwebhook"
	"o.o/backend/com/etc/logging/smslog/aggregate"
	"o.o/backend/com/eventhandler/fabo/publisher"
	"o.o/backend/com/eventhandler/handler"
	"o.o/backend/com/eventhandler/notifier"
	sqlstore2 "o.o/backend/com/eventhandler/notifier/sqlstore"
	"o.o/backend/com/fabo/main/fbmessaging"
	"o.o/backend/com/fabo/main/fbpage"
	"o.o/backend/com/fabo/main/fbuser"
	pm5 "o.o/backend/com/fabo/main/fbuser/pm"
	"o.o/backend/com/fabo/pkg/fbclient"
	redis2 "o.o/backend/com/fabo/pkg/redis"
	"o.o/backend/com/fabo/pkg/webhook"
	"o.o/backend/com/main/address"
	aggregate3 "o.o/backend/com/main/authorization/aggregate"
	aggregate4 "o.o/backend/com/main/catalog/aggregate"
	"o.o/backend/com/main/catalog/pm"
	query4 "o.o/backend/com/main/catalog/query"
	aggregate6 "o.o/backend/com/main/connectioning/aggregate"
	"o.o/backend/com/main/connectioning/manager"
	query8 "o.o/backend/com/main/connectioning/query"
	"o.o/backend/com/main/identity"
	pm2 "o.o/backend/com/main/identity/pm"
	"o.o/backend/com/main/inventory/aggregatex"
	query6 "o.o/backend/com/main/inventory/query"
	aggregate2 "o.o/backend/com/main/invitation/aggregate"
	"o.o/backend/com/main/invitation/query"
	"o.o/backend/com/main/location"
	"o.o/backend/com/main/ordering"
	pm3 "o.o/backend/com/main/ordering/pm"
	query7 "o.o/backend/com/main/receipting/query"
	"o.o/backend/com/main/shipmentpricing/pricelist"
	"o.o/backend/com/main/shipmentpricing/pricelistpromotion"
	"o.o/backend/com/main/shipmentpricing/shipmentprice"
	"o.o/backend/com/main/shipmentpricing/shipmentservice"
	"o.o/backend/com/main/shipmentpricing/shopshipmentpricelist"
	aggregate9 "o.o/backend/com/main/shipping/aggregate"
	"o.o/backend/com/main/shipping/carrier"
	pm4 "o.o/backend/com/main/shipping/pm"
	query10 "o.o/backend/com/main/shipping/query"
	query9 "o.o/backend/com/main/shippingcode/query"
	aggregate8 "o.o/backend/com/main/stocktaking/aggregate"
	query5 "o.o/backend/com/main/stocktaking/query"
	aggregate7 "o.o/backend/com/shopping/carrying/aggregate"
	query12 "o.o/backend/com/shopping/carrying/query"
	aggregate5 "o.o/backend/com/shopping/customering/aggregate"
	query2 "o.o/backend/com/shopping/customering/query"
	query11 "o.o/backend/com/summary/query"
	query3 "o.o/backend/com/supporting/ticket/query"
	"o.o/backend/pkg/common/apifw/captcha"
	"o.o/backend/pkg/common/apifw/health"
	auth2 "o.o/backend/pkg/common/authorization/auth"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/export"
	"o.o/backend/pkg/etop/api/sadmin"
	"o.o/backend/pkg/etop/api/sadmin/_fabo"
	"o.o/backend/pkg/etop/api/shop"
	"o.o/backend/pkg/etop/api/shop/_min"
	"o.o/backend/pkg/etop/api/shop/account"
	"o.o/backend/pkg/etop/api/shop/authorize"
	"o.o/backend/pkg/etop/api/shop/brand"
	carrier2 "o.o/backend/pkg/etop/api/shop/carrier"
	"o.o/backend/pkg/etop/api/shop/category"
	"o.o/backend/pkg/etop/api/shop/collection"
	"o.o/backend/pkg/etop/api/shop/connection"
	"o.o/backend/pkg/etop/api/shop/customer"
	"o.o/backend/pkg/etop/api/shop/customergroup"
	export2 "o.o/backend/pkg/etop/api/shop/export"
	"o.o/backend/pkg/etop/api/shop/fulfillment"
	"o.o/backend/pkg/etop/api/shop/history"
	"o.o/backend/pkg/etop/api/shop/inventory"
	"o.o/backend/pkg/etop/api/shop/notification"
	"o.o/backend/pkg/etop/api/shop/order"
	"o.o/backend/pkg/etop/api/shop/product"
	"o.o/backend/pkg/etop/api/shop/shipment"
	"o.o/backend/pkg/etop/api/shop/stocktake"
	summary2 "o.o/backend/pkg/etop/api/shop/summary"
	"o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/authorize/middleware"
	"o.o/backend/pkg/etop/authorize/tokens"
	"o.o/backend/pkg/etop/eventstream"
	imcsv3 "o.o/backend/pkg/etop/logic/fulfillments/imcsv"
	"o.o/backend/pkg/etop/logic/orders"
	"o.o/backend/pkg/etop/logic/orders/imcsv"
	imcsv2 "o.o/backend/pkg/etop/logic/products/imcsv"
	"o.o/backend/pkg/etop/logic/summary"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/backend/pkg/fabo"
	"o.o/backend/pkg/fabo/faboinfo"
	"o.o/backend/pkg/integration/email"
	"o.o/backend/pkg/integration/shipping/ghn/webhook/v2"
	"o.o/backend/pkg/integration/sms"
)

// Injectors from wire.go:

func Build(ctx context.Context, cfg config.Config, consumer mq.KafkaConsumer) (Output, func(), error) {
	redisRedis := cfg.Redis
	store := redis.Connect(redisRedis)
	service := health.New(store)
	miscService := &api.MiscService{}
	policy := ProvidePolicy()
	authorizer := auth.New(policy)
	sharedConfig := cfg.SharedConfig
	sAdminToken := config_server.WireSAdminToken(sharedConfig)
	database_minConfig := cfg.Databases
	databases, err := database_min.BuildDatabases(database_minConfig)
	if err != nil {
		return Output{}, nil, err
	}
	mainDB := databases.Main
	queryService := identity.NewQueryService(mainDB)
	queryBus := identity.QueryServiceMessageBus(queryService)
	tokenStore := tokens.NewTokenStore(store)
	busBus := bus.New()
	addressStore := &sqlstore.AddressStore{
		DB: mainDB,
	}
	userStore := sqlstore.BuildUserStore(mainDB)
	userStoreInterface := sqlstore.BindUserStore(userStore)
	accountStore := &sqlstore.AccountStore{
		DB:           mainDB,
		EventBus:     busBus,
		AddressStore: addressStore,
		UserStore:    userStoreInterface,
	}
	accountStoreInterface := sqlstore.BindAccountStore(accountStore)
	partnerStore := sqlstore.BuildPartnerStore(mainDB)
	partnerStoreInterface := sqlstore.BindPartnerStore(partnerStore)
	accountUserStore := &sqlstore.AccountUserStore{
		DB: mainDB,
	}
	accountUserStoreInterface := sqlstore.BindAccountUserStore(accountUserStore)
	shopStore := &sqlstore.ShopStore{
		DB: mainDB,
	}
	shopStoreInterface := sqlstore.BindShopStore(shopStore)
	sessionStarter := &middleware.SessionStarter{
		SAdminToken:      sAdminToken,
		IdentityQuery:    queryBus,
		TokenStore:       tokenStore,
		AccountStore:     accountStoreInterface,
		UserStore:        userStoreInterface,
		PartnerStore:     partnerStoreInterface,
		AccountUserStore: accountUserStoreInterface,
		ShopStore:        shopStoreInterface,
	}
	session := config_server.NewSession(authorizer, sessionStarter, userStoreInterface, accountUserStoreInterface, sharedConfig, store)
	identityAggregate := identity.NewAggregate(mainDB, busBus)
	commandBus := identity.AggregateMessageBus(identityAggregate)
	flagEnableNewLinkInvitation := cfg.FlagEnableNewLinkInvitation
	invitationQuery := query.NewInvitationQuery(mainDB, flagEnableNewLinkInvitation)
	invitationQueryBus := query.InvitationQueryMessageBus(invitationQuery)
	notifierQueryService := notifier.NewQueryService(mainDB)
	notifyQueryBus := notifier.QueryServiceNotifyBus(notifierQueryService)
	carrierManager := SupportedShipnowManager()
	notifierAggregate := notifier.NewNotifyAggregate(mainDB, carrierManager)
	notifyCommandBus := notifier.NewNotifyAggregateMessageBus(notifierAggregate)
	generator := auth2.NewGenerator(store)
	smsConfig := cfg.SMS
	v := sms_min.SupportedSMSDrivers(smsConfig)
	logDB := databases.Log
	smsLogAggregate := aggregate.NewSmsLogAggregate(busBus, logDB)
	smslogCommandBus := aggregate.SmsLogAggregateMessageBus(smsLogAggregate)
	client := sms.New(smsConfig, v, smslogCommandBus)
	smtpConfig := cfg.SMTP
	emailClient := email.New(smtpConfig)
	userStoreFactory := sqlstore.NewUserStore(mainDB)
	login := &sqlstore.Login{
		UserStore: userStoreInterface,
	}
	loginInterface := sqlstore.BindLogin(login)
	userService := &api.UserService{
		Session:          session,
		IdentityAggr:     commandBus,
		IdentityQuery:    queryBus,
		InvitationQuery:  invitationQueryBus,
		NotifyQuery:      notifyQueryBus,
		NotifyAggr:       notifyCommandBus,
		EventBus:         busBus,
		AuthStore:        generator,
		TokenStore:       tokenStore,
		RedisStore:       store,
		SMSClient:        client,
		EmailClient:      emailClient,
		UserStore:        userStoreFactory,
		UserStoreIface:   userStoreInterface,
		ShopStore:        shopStoreInterface,
		AccountUserStore: accountUserStoreInterface,
		LoginIface:       loginInterface,
	}
	partnerStoreFactory := sqlstore.NewPartnerStore(mainDB)
	accountService := &api.AccountService{
		Session:           session,
		PartnerStore:      partnerStoreFactory,
		AccountStore:      accountStoreInterface,
		AccountUserStore:  accountUserStoreInterface,
		PartnerStoreIface: partnerStoreInterface,
	}
	locationQuery := location.New(mainDB)
	locationQueryBus := location.QueryMessageBus(locationQuery)
	locationService := &api.LocationService{
		Session:       session,
		LocationQuery: locationQueryBus,
	}
	bankService := &api.BankService{
		Session: session,
	}
	addressAggregate := address.NewAggregateAddress(busBus, mainDB, locationQueryBus)
	addressCommandBus := address.AddressAggregateMessageBus(addressAggregate)
	addressQueryService := address.NewQueryAddress(mainDB, busBus)
	addressQueryBus := address.QueryServiceMessageBus(addressQueryService)
	addressService := &api.AddressService{
		Session:     session,
		AddressAggr: addressCommandBus,
		AddressQS:   addressQueryBus,
	}
	invitationConfig := cfg.Invitation
	customerQuery := query2.NewCustomerQuery(mainDB)
	customeringQueryBus := query2.CustomerQueryMessageBus(customerQuery)
	invitationAggregate := aggregate2.NewInvitationAggregate(mainDB, invitationConfig, customeringQueryBus, queryBus, busBus, client, emailClient, flagEnableNewLinkInvitation, accountUserStoreInterface, shopStoreInterface, userStoreInterface)
	invitationCommandBus := aggregate2.InvitationAggregateMessageBus(invitationAggregate)
	authorizationAggregate := &aggregate3.AuthorizationAggregate{
		Auth:             authorizer,
		AccountUserStore: accountUserStoreInterface,
		ShopStore:        shopStoreInterface,
	}
	authorizationCommandBus := aggregate3.AuthorizationAggregateMessageBus(authorizationAggregate)
	accountRelationshipService := &api.AccountRelationshipService{
		Session:           session,
		InvitationAggr:    invitationCommandBus,
		InvitationQuery:   invitationQueryBus,
		AuthorizationAggr: authorizationCommandBus,
		UserStore:         userStoreFactory,
		AccountUserStore:  accountUserStoreInterface,
	}
	userRelationshipService := &api.UserRelationshipService{
		Session:                session,
		InvitationAggr:         invitationCommandBus,
		InvitationQuery:        invitationQueryBus,
		AuthorizationAggregate: authorizationCommandBus,
		ShopStore:              shopStoreInterface,
		UserStore:              userStoreInterface,
	}
	ticketQuery := query3.NewTicketQuery(store, busBus, mainDB)
	ticketQueryBus := query3.TicketQueryMessageBus(ticketQuery)
	ticketService := &api.TicketService{
		Session:     session,
		TicketQuery: ticketQueryBus,
	}
	ecomService := &api.EcomService{
		Session:        session,
		SessionStarter: sessionStarter,
	}
	emailConfig := cfg.Email
	servers, cleanup := api.NewServers(miscService, userService, accountService, locationService, bankService, addressService, accountRelationshipService, userRelationshipService, ticketService, ecomService, store, emailConfig, smsConfig)
	shopMiscService := &shop.MiscService{
		Session: session,
	}
	queryQueryService := query4.New(mainDB)
	catalogQueryBus := query4.QueryServiceMessageBus(queryQueryService)
	aggregateAggregate := aggregate4.New(busBus, mainDB)
	catalogCommandBus := aggregate4.AggregateMessageBus(aggregateAggregate)
	brandService := &brand.BrandService{
		Session:      session,
		CatalogQuery: catalogQueryBus,
		CatalogAggr:  catalogCommandBus,
	}
	traderingQueryBus := _wireQueryBusValue
	stocktakeQuery := query5.NewQueryStocktake(mainDB)
	stocktakingQueryBus := query5.StocktakeQueryMessageBus(stocktakeQuery)
	orderStore := &sqlstore.OrderStore{
		DB:           mainDB,
		LocationBus:  locationQueryBus,
		AccountStore: accountStoreInterface,
		ShopStore:    shopStoreInterface,
	}
	orderStoreInterface := sqlstore.BindOrderStore(orderStore)
	inventoryAggregate := aggregatex.NewAggregateInventory(busBus, mainDB, stocktakingQueryBus, catalogQueryBus, orderStoreInterface)
	inventoryCommandBus := aggregatex.InventoryAggregateMessageBus(inventoryAggregate)
	inventoryQueryService := query6.NewQueryInventory(stocktakingQueryBus, busBus, mainDB)
	inventoryQueryBus := query6.InventoryQueryServiceMessageBus(inventoryQueryService)
	inventoryService := &inventory.InventoryService{
		Session:        session,
		TraderQuery:    traderingQueryBus,
		InventoryAggr:  inventoryCommandBus,
		InventoryQuery: inventoryQueryBus,
	}
	accountshipnowQueryBus := _wireAccountshipnowQueryBusValue
	accountshipnowCommandBus := _wireCommandBusValue
	accountAccountService := &account.AccountService{
		Session:             session,
		IdentityAggr:        commandBus,
		IdentityQuery:       queryBus,
		AddressQuery:        addressQueryBus,
		AddressAggr:         addressCommandBus,
		UserStore:           userStoreFactory,
		AccountStore:        accountStoreInterface,
		UserStoreIface:      userStoreInterface,
		AccountshipnowQuery: accountshipnowQueryBus,
		AccountshipnowAggr:  accountshipnowCommandBus,
	}
	collectionService := &collection.CollectionService{
		Session:      session,
		CatalogQuery: catalogQueryBus,
		CatalogAggr:  catalogCommandBus,
	}
	customerAggregate := aggregate5.NewCustomerAggregate(busBus, mainDB)
	customeringCommandBus := aggregate5.CustomerAggregateMessageBus(customerAggregate)
	aggregateAddressAggregate := aggregate5.NewAddressAggregate(mainDB)
	addressingCommandBus := aggregate5.AddressAggregateMessageBus(aggregateAddressAggregate)
	addressQuery := query2.NewAddressQuery(mainDB)
	addressingQueryBus := query2.AddressQueryMessageBus(addressQuery)
	orderingQueryService := ordering.NewQueryService(mainDB)
	orderingQueryBus := ordering.QueryServiceMessageBus(orderingQueryService)
	receiptQuery := query7.NewReceiptQuery(mainDB)
	receiptingQueryBus := query7.ReceiptQueryMessageBus(receiptQuery)
	customerService := &customer.CustomerService{
		Session:       session,
		LocationQuery: locationQueryBus,
		CustomerQuery: customeringQueryBus,
		CustomerAggr:  customeringCommandBus,
		AddressAggr:   addressingCommandBus,
		AddressQuery:  addressingQueryBus,
		OrderQuery:    orderingQueryBus,
		ReceiptQuery:  receiptingQueryBus,
	}
	customerGroupService := &customergroup.CustomerGroupService{
		Session:       session,
		CustomerAggr:  customeringCommandBus,
		CustomerQuery: customeringQueryBus,
	}
	shopVariantStore := &sqlstore.ShopVariantStore{
		DB: mainDB,
	}
	shopVariantStoreInterface := sqlstore.BindShopVariantStore(shopVariantStore)
	productService := &product.ProductService{
		Session:          session,
		CatalogQuery:     catalogQueryBus,
		CatalogAggr:      catalogCommandBus,
		InventoryQuery:   inventoryQueryBus,
		ShopVariantStore: shopVariantStoreInterface,
	}
	categoryService := &category.CategoryService{
		Session:      session,
		CatalogQuery: catalogQueryBus,
		CatalogAggr:  catalogCommandBus,
	}
	orderingAggregate := ordering.NewAggregate(busBus, mainDB)
	orderingCommandBus := ordering.AggregateMessageBus(orderingAggregate)
	mapShipmentServices := shipment_all.SupportedShipmentServices()
	connectionQuery := query8.NewConnectionQuery(mainDB, mapShipmentServices)
	connectioningQueryBus := query8.ConnectionQueryMessageBus(connectionQuery)
	connectionAggregate := aggregate6.NewConnectionAggregate(mainDB, busBus)
	connectioningCommandBus := aggregate6.ConnectionAggregateMessageBus(connectionAggregate)
	queryService2 := query9.NewQueryService(mainDB)
	shippingcodeQueryBus := query9.QueryServiceMessageBus(queryService2)
	shipmentserviceQueryService := shipmentservice.NewQueryService(mainDB, store)
	shipmentserviceQueryBus := shipmentservice.QueryServiceMessageBus(shipmentserviceQueryService)
	pricelistQueryService := pricelist.NewQueryService(mainDB, store)
	pricelistQueryBus := pricelist.QueryServiceMessageBus(pricelistQueryService)
	shopshipmentpricelistQueryService := shopshipmentpricelist.NewQueryService(mainDB, store)
	shopshipmentpricelistQueryBus := shopshipmentpricelist.QueryServiceMessageBus(shopshipmentpricelistQueryService)
	shipmentpriceQueryService := shipmentprice.NewQueryService(mainDB, store, locationQueryBus, pricelistQueryBus, shopshipmentpricelistQueryBus)
	shipmentpriceQueryBus := shipmentprice.QueryServiceMessageBus(shipmentpriceQueryService)
	pricelistpromotionQueryService := pricelistpromotion.NewQueryService(mainDB, store, locationQueryBus, queryBus, shopshipmentpricelistQueryBus, pricelistQueryBus)
	pricelistpromotionQueryBus := pricelistpromotion.QueryServiceMessageBus(pricelistpromotionQueryService)
	shipment_allConfig := cfg.Shipment
	typesConfig := shipment_all.SupportedShippingCarrierConfig(shipment_allConfig)
	driver := shipment_all.SupportedCarrierDriver()
	connectionManager := manager.NewConnectionManager(store, connectioningQueryBus)
	shipmentManager, err := carrier.NewShipmentManager(busBus, locationQueryBus, queryBus, connectioningQueryBus, connectioningCommandBus, shippingcodeQueryBus, shipmentserviceQueryBus, shipmentpriceQueryBus, pricelistpromotionQueryBus, typesConfig, driver, connectionManager, orderStoreInterface)
	if err != nil {
		cleanup()
		return Output{}, nil, err
	}
	addressStoreInterface := sqlstore.BindAddressStore(addressStore)
	flagFaboOrderAutoConfirmPaymentStatus := cfg.FlagFaboOrderAutoConfirmPaymentStatus
	orderLogic := &orderS.OrderLogic{
		CatalogQuery:                          catalogQueryBus,
		OrderAggr:                             orderingCommandBus,
		CustomerAggr:                          customeringCommandBus,
		CustomerQuery:                         customeringQueryBus,
		TraderAddressAggr:                     addressingCommandBus,
		TraderAddressQuery:                    addressingQueryBus,
		LocationQuery:                         locationQueryBus,
		EventBus:                              busBus,
		ShipmentManager:                       shipmentManager,
		AddressStore:                          addressStoreInterface,
		OrderStore:                            orderStoreInterface,
		FlagFaboOrderUpdatePaymentSatusConfig: flagFaboOrderAutoConfirmPaymentStatus,
	}
	orderService := &order.OrderService{
		Session:       session,
		OrderAggr:     orderingCommandBus,
		CustomerQuery: customeringQueryBus,
		OrderQuery:    orderingQueryBus,
		ReceiptQuery:  receiptingQueryBus,
		OrderLogic:    orderLogic,
		OrderStore:    orderStoreInterface,
	}
	queryService3 := query10.NewQueryService(mainDB)
	shippingQueryBus := query10.QueryServiceMessageBus(queryService3)
	fulfillmentService := &fulfillment.FulfillmentService{
		Session:         session,
		ShipmentManager: shipmentManager,
		ShippingQuery:   shippingQueryBus,
		OrderStore:      orderStoreInterface,
	}
	historyStore := &sqlstore.HistoryStore{
		DB: mainDB,
	}
	historyStoreInterface := sqlstore.BindHistoryStore(historyStore)
	historyService := &history.HistoryService{
		Session:      session,
		HistoryStore: historyStoreInterface,
	}
	dashboardQuery := query11.NewDashboardQuery(mainDB, store, locationQueryBus)
	summaryQueryBus := query11.DashboardQueryMessageBus(dashboardQuery)
	summarySummary := summary.New(mainDB)
	moneyTxStore := &sqlstore.MoneyTxStore{
		DB:               mainDB,
		EventBus:         busBus,
		AccountUserStore: accountUserStoreInterface,
		ShopStore:        shopStoreInterface,
		OrderStore:       orderStoreInterface,
	}
	moneyTxStoreInterface := sqlstore.BindMoneyTxStore(moneyTxStore)
	summaryService := &summary2.SummaryService{
		Session:      session,
		SummaryQuery: summaryQueryBus,
		SummaryOld:   summarySummary,
		MoneyTxStore: moneyTxStoreInterface,
	}
	eventStream := eventstream.New(ctx)
	configDirs := cfg.ExportDirs
	driverConfig := cfg.StorageDriver
	bucket, err := storage_all.Build(ctx, driverConfig)
	if err != nil {
		cleanup()
		return Output{}, nil, err
	}
	exportAttemptStoreFactory := sqlstore.NewExportAttemptStore(mainDB)
	exportService, cleanup2 := export.New(store, eventStream, configDirs, bucket, exportAttemptStoreFactory, orderStoreInterface)
	exportExportService := &export2.ExportService{
		Session:     session,
		Auth:        authorizer,
		ExportInner: exportService,
	}
	notifierDB := databases.Notifier
	notificationStore := sqlstore2.NewNotificationStore(notifierDB, accountUserStoreInterface)
	deviceStore := sqlstore2.NewDeviceStore(notifierDB)
	notificationService := &notification.NotificationService{
		Session:           session,
		NotificationStore: notificationStore,
		DeviceStore:       deviceStore,
	}
	authorizeService := &authorize.AuthorizeService{
		Session:      session,
		PartnerStore: partnerStoreInterface,
	}
	carrierAggregate := aggregate7.NewCarrierAggregate(busBus, mainDB)
	carryingCommandBus := aggregate7.CarrierAggregateMessageBus(carrierAggregate)
	carrierQuery := query12.NewCarrierQuery(mainDB)
	carryingQueryBus := query12.CarrierQueryMessageBus(carrierQuery)
	carrierService := &carrier2.CarrierService{
		Session:      session,
		CarrierAggr:  carryingCommandBus,
		CarrierQuery: carryingQueryBus,
	}
	stocktakeAggregate := aggregate8.NewAggregateStocktake(mainDB, busBus, store)
	stocktakingCommandBus := aggregate8.StocktakeAggregateMessageBus(stocktakeAggregate)
	stocktakeService := &stocktake.StocktakeService{
		Session:        session,
		CatalogQuery:   catalogQueryBus,
		StocktakeAggr:  stocktakingCommandBus,
		StocktakeQuery: stocktakingQueryBus,
		InventoryQuery: inventoryQueryBus,
	}
	aggregate10 := aggregate9.NewAggregate(mainDB, busBus, locationQueryBus, orderingQueryBus, shipmentManager, connectioningQueryBus, queryBus, addressQueryBus)
	shippingCommandBus := aggregate9.AggregateMessageBus(aggregate10)
	shipmentService := &shipment.ShipmentService{
		Session:           session,
		ShipmentManager:   shipmentManager,
		ShippingAggregate: shippingCommandBus,
		OrderStore:        orderStoreInterface,
	}
	connectionService := &connection.ConnectionService{
		Session:            session,
		ShipmentManager:    shipmentManager,
		ConnectionQuery:    connectioningQueryBus,
		ConnectionAggr:     connectioningCommandBus,
		IdentityQuery:      queryBus,
		AccountshipnowAggr: accountshipnowCommandBus,
	}
	shopServers := shop_min.NewServers(store, shopMiscService, brandService, inventoryService, accountAccountService, collectionService, customerService, customerGroupService, productService, categoryService, orderService, fulfillmentService, historyService, summaryService, exportExportService, notificationService, authorizeService, carrierService, stocktakeService, shipmentService, connectionService)
	fbPageQuery := fbpage.NewFbPageQuery(mainDB)
	fbpagingQueryBus := fbpage.FbPageQueryMessageBus(fbPageQuery)
	fbUserQuery := fbuser.NewFbUserQuery(mainDB, customeringQueryBus)
	fbuseringQueryBus := fbuser.FbUserQueryMessageBus(fbUserQuery)
	faboPagesKit := &faboinfo.FaboPagesKit{
		FBPageQuery: fbpagingQueryBus,
		FBUserQuery: fbuseringQueryBus,
	}
	fbExternalPageAggregate := fbpage.NewFbPageAggregate(mainDB)
	fbpagingCommandBus := fbpage.FbExternalPageAggregateMessageBus(fbExternalPageAggregate)
	fbUserAggregate := fbuser.NewFbUserAggregate(mainDB, fbpagingCommandBus, customeringQueryBus)
	fbuseringCommandBus := fbuser.FbUserAggregateMessageBus(fbUserAggregate)
	appConfig := cfg.FacebookApp
	fbClient := fbclient.New(appConfig)
	pageService := &fabo.PageService{
		Session:             session,
		FaboInfo:            faboPagesKit,
		FBExternalUserQuery: fbuseringQueryBus,
		FBExternalUserAggr:  fbuseringCommandBus,
		FBExternalPageQuery: fbpagingQueryBus,
		FBExternalPageAggr:  fbpagingCommandBus,
		FBClient:            fbClient,
	}
	fbMessagingQuery := fbmessaging.NewFbMessagingQuery(mainDB)
	fbmessagingQueryBus := fbmessaging.FbMessagingQueryMessageBus(fbMessagingQuery)
	fbExternalMessagingAggregate := fbmessaging.NewFbExternalMessagingAggregate(mainDB, busBus, fbClient)
	fbmessagingCommandBus := fbmessaging.FbExternalMessagingAggregateMessageBus(fbExternalMessagingAggregate)
	customerConversationService := &fabo.CustomerConversationService{
		Session:          session,
		FaboPagesKit:     faboPagesKit,
		FBMessagingQuery: fbmessagingQueryBus,
		FBMessagingAggr:  fbmessagingCommandBus,
		FBPagingQuery:    fbpagingQueryBus,
		FBClient:         fbClient,
		FBUserQuery:      fbuseringQueryBus,
	}
	faboCustomerService := &fabo.CustomerService{
		Session:        session,
		CustomerQuery:  customeringQueryBus,
		FBUseringQuery: fbuseringQueryBus,
		FBUseringAggr:  fbuseringCommandBus,
	}
	shopService := &fabo.ShopService{
		Session:             session,
		FBExternalUserQuery: fbuseringQueryBus,
		FBExternalUserAggr:  fbuseringCommandBus,
	}
	extraShipmentService := &fabo.ExtraShipmentService{
		Session:         session,
		ShipmentManager: shipmentManager,
		ConnectionQS:    connectioningQueryBus,
	}
	faboServers := fabo.NewServers(pageService, customerConversationService, faboCustomerService, shopService, extraShipmentService, store)
	webhookCallbackService := sadmin.NewWebhookCallbackService(store)
	webhookService := &sadmin.WebhookService{
		Session:                session,
		WebhookCallbackService: webhookCallbackService,
	}
	sadminServers := _fabo.NewServers(webhookService)
	captchaConfig := cfg.Captcha
	captchaCaptcha := captcha.New(captchaConfig)
	intHandlers, err := BuildIntHandlers(servers, shopServers, faboServers, sadminServers, captchaCaptcha)
	if err != nil {
		cleanup2()
		cleanup()
		return Output{}, nil, err
	}
	dirConfigs := cfg.UploadDirs
	uploader, err := _uploader.NewUploader(ctx, dirConfigs, bucket)
	if err != nil {
		cleanup2()
		cleanup()
		return Output{}, nil, err
	}
	exportAttemptStore := sqlstore.BuildExportAttemptStore(mainDB)
	exportAttemptStoreInterface := sqlstore.BindExportAttemptStore(exportAttemptStore)
	imcsvImport, cleanup3 := imcsv.New(authorizer, locationQueryBus, store, uploader, mainDB, orderStoreInterface, exportAttemptStoreInterface)
	categoryStore := &sqlstore.CategoryStore{
		DB: mainDB,
	}
	categoryStoreInterface := sqlstore.BindCategoryStore(categoryStore)
	import2, cleanup4 := imcsv2.New(store, uploader, mainDB, exportAttemptStoreInterface, categoryStoreInterface, shopStoreInterface)
	import3, cleanup5 := imcsv3.New(store, uploader, exportAttemptStoreInterface)
	importHandler := server_shop.BuildImportHandler(imcsvImport, import2, import3, session)
	eventStreamHandler := server_shop.BuildEventStreamHandler(eventStream, session)
	downloadHandler := server_shop.BuildDownloadHandler()
	faboImageHandler := fabo2.BuildFaboImageHandler()
	mainServer := BuildMainServer(service, intHandlers, sharedConfig, importHandler, eventStreamHandler, downloadHandler, faboImageHandler)
	webhookConfig := shipment_allConfig.GHNWebhook
	shippingwebhookAggregate := shippingwebhook.NewAggregate(logDB)
	v2Webhook := v2.New(mainDB, shipmentManager, queryBus, shippingCommandBus, shippingwebhookAggregate, orderStoreInterface)
	ghnWebhookServer := v2_2.NewGHNWebhookServer(webhookConfig, shipmentManager, queryBus, shippingCommandBus, v2Webhook)
	configWebhookConfig := cfg.Webhook
	faboRedis := redis2.NewFaboRedis(store)
	webhookWebhook := webhook.New(mainDB, logDB, store, configWebhookConfig, faboRedis, fbClient, fbmessagingQueryBus, fbmessagingCommandBus, fbpagingQueryBus)
	fbWebhookServer := BuildWebhookServer(configWebhookConfig, webhookWebhook)
	v3 := BuildServers(mainServer, ghnWebhookServer, fbWebhookServer)
	kafka := cfg.Kafka
	handlerHandler := handler.New(consumer, kafka)
	publisherPublisher := publisher.New(eventStream)
	processManager := pm.New(busBus, catalogQueryBus, catalogCommandBus)
	pmProcessManager := pm2.New(busBus, queryBus, commandBus, invitationQueryBus, addressQueryBus, addressCommandBus, accountUserStoreInterface)
	affiliateCommandBus := _wireAffiliateCommandBusValue
	processManager2 := pm3.New(busBus, orderingCommandBus, affiliateCommandBus, receiptingQueryBus, inventoryCommandBus, orderingQueryBus, customeringQueryBus)
	processManager3 := pm4.New(busBus, shippingQueryBus, shippingCommandBus, store, connectioningQueryBus, shopStoreInterface, moneyTxStoreInterface)
	processManager4 := pm5.New(busBus, fbuseringCommandBus)
	fbmessagingProcessManager := fbmessaging.NewProcessManager(busBus, fbmessagingQueryBus, fbmessagingCommandBus, fbpagingQueryBus, fbuseringQueryBus, fbuseringCommandBus, faboRedis)
	output := Output{
		Servers:        v3,
		EventStream:    eventStream,
		Health:         service,
		Handler:        handlerHandler,
		Publisher:      publisherPublisher,
		FbClient:       fbClient,
		_catalogPM:     processManager,
		_identityPM:    pmProcessManager,
		_orderPM:       processManager2,
		_shippingPM:    processManager3,
		_fbuserPM:      processManager4,
		_fbMessagingPM: fbmessagingProcessManager,
	}
	return output, func() {
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

var (
	_wireQueryBusValue               = tradering.QueryBus{}
	_wireAccountshipnowQueryBusValue = accountshipnow.QueryBus{}
	_wireCommandBusValue             = accountshipnow.CommandBus{}
	_wireAffiliateCommandBusValue    = affiliate.CommandBus{}
)
