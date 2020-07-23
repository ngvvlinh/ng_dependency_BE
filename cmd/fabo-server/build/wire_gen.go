// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package build

import (
	"context"
	"o.o/api/services/affiliate"
	"o.o/backend/cmd/fabo-server/config"
	"o.o/backend/cogs/config/_server"
	"o.o/backend/cogs/database/_min"
	"o.o/backend/cogs/server/admin"
	"o.o/backend/cogs/server/shop"
	"o.o/backend/cogs/shipment/_fabo"
	v2_2 "o.o/backend/cogs/shipment/ghn/v2"
	"o.o/backend/cogs/shipment/ghtk"
	"o.o/backend/cogs/shipment/vtpost"
	"o.o/backend/cogs/sms/_min"
	"o.o/backend/cogs/uploader"
	"o.o/backend/com/etc/logging/smslog/aggregate"
	"o.o/backend/com/eventhandler/fabo/publisher"
	"o.o/backend/com/eventhandler/handler"
	"o.o/backend/com/fabo/main/fbmessaging"
	"o.o/backend/com/fabo/main/fbpage"
	"o.o/backend/com/fabo/main/fbuser"
	pm11 "o.o/backend/com/fabo/main/fbuser/pm"
	"o.o/backend/com/fabo/pkg/fbclient"
	redis2 "o.o/backend/com/fabo/pkg/redis"
	webhook3 "o.o/backend/com/fabo/pkg/webhook"
	"o.o/backend/com/main/address"
	aggregate3 "o.o/backend/com/main/authorization/aggregate"
	aggregate4 "o.o/backend/com/main/catalog/aggregate"
	"o.o/backend/com/main/catalog/pm"
	query3 "o.o/backend/com/main/catalog/query"
	aggregate7 "o.o/backend/com/main/connectioning/aggregate"
	query13 "o.o/backend/com/main/connectioning/query"
	"o.o/backend/com/main/credit"
	"o.o/backend/com/main/identity"
	pm2 "o.o/backend/com/main/identity/pm"
	aggregate5 "o.o/backend/com/main/inventory/aggregate"
	pm3 "o.o/backend/com/main/inventory/pm"
	query8 "o.o/backend/com/main/inventory/query"
	aggregate2 "o.o/backend/com/main/invitation/aggregate"
	pm4 "o.o/backend/com/main/invitation/pm"
	"o.o/backend/com/main/invitation/query"
	aggregate10 "o.o/backend/com/main/ledgering/aggregate"
	pm5 "o.o/backend/com/main/ledgering/pm"
	query17 "o.o/backend/com/main/ledgering/query"
	"o.o/backend/com/main/location"
	aggregate16 "o.o/backend/com/main/moneytx/aggregate"
	pm6 "o.o/backend/com/main/moneytx/pm"
	query15 "o.o/backend/com/main/moneytx/query"
	"o.o/backend/com/main/ordering"
	pm7 "o.o/backend/com/main/ordering/pm"
	aggregate11 "o.o/backend/com/main/purchaseorder/aggregate"
	query10 "o.o/backend/com/main/purchaseorder/query"
	aggregate15 "o.o/backend/com/main/purchaserefund/aggregate"
	query12 "o.o/backend/com/main/purchaserefund/query"
	aggregate8 "o.o/backend/com/main/receipting/aggregate"
	pm8 "o.o/backend/com/main/receipting/pm"
	query9 "o.o/backend/com/main/receipting/query"
	aggregate14 "o.o/backend/com/main/refund/aggregate"
	pm9 "o.o/backend/com/main/refund/pm"
	query11 "o.o/backend/com/main/refund/query"
	"o.o/backend/com/main/shipmentpricing/pricelist"
	"o.o/backend/com/main/shipmentpricing/shipmentprice"
	"o.o/backend/com/main/shipmentpricing/shipmentservice"
	"o.o/backend/com/main/shipmentpricing/shopshipmentpricelist"
	"o.o/backend/com/main/shipnow"
	aggregate13 "o.o/backend/com/main/shipping/aggregate"
	"o.o/backend/com/main/shipping/carrier"
	pm10 "o.o/backend/com/main/shipping/pm"
	query14 "o.o/backend/com/main/shipping/query"
	aggregate12 "o.o/backend/com/main/stocktaking/aggregate"
	query7 "o.o/backend/com/main/stocktaking/query"
	aggregate9 "o.o/backend/com/shopping/carrying/aggregate"
	query4 "o.o/backend/com/shopping/carrying/query"
	aggregate6 "o.o/backend/com/shopping/customering/aggregate"
	query2 "o.o/backend/com/shopping/customering/query"
	query5 "o.o/backend/com/shopping/suppliering/query"
	query6 "o.o/backend/com/shopping/tradering/query"
	query16 "o.o/backend/com/summary/query"
	"o.o/backend/pkg/common/apifw/captcha"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/authorization/auth"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/admin"
	"o.o/backend/pkg/etop/api/admin/_min"
	"o.o/backend/pkg/etop/api/export"
	"o.o/backend/pkg/etop/api/sadmin"
	"o.o/backend/pkg/etop/api/shop"
	"o.o/backend/pkg/etop/api/shop/_min"
	"o.o/backend/pkg/etop/authorize/middleware"
	"o.o/backend/pkg/etop/authorize/tokens"
	"o.o/backend/pkg/etop/eventstream"
	"o.o/backend/pkg/etop/logic/hotfix"
	"o.o/backend/pkg/etop/logic/money-transaction/ghnimport"
	"o.o/backend/pkg/etop/logic/money-transaction/ghtkimport"
	"o.o/backend/pkg/etop/logic/money-transaction/handlers"
	"o.o/backend/pkg/etop/logic/money-transaction/jtexpressimport"
	"o.o/backend/pkg/etop/logic/money-transaction/vtpostimport"
	"o.o/backend/pkg/etop/logic/orders"
	"o.o/backend/pkg/etop/logic/orders/imcsv"
	imcsv2 "o.o/backend/pkg/etop/logic/products/imcsv"
	"o.o/backend/pkg/etop/logic/shipping_provider"
	"o.o/backend/pkg/etop/logic/summary"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/backend/pkg/fabo"
	"o.o/backend/pkg/fabo/faboinfo"
	"o.o/backend/pkg/integration/email"
	"o.o/backend/pkg/integration/shipping/ghn"
	"o.o/backend/pkg/integration/shipping/ghn/webhook/v2"
	"o.o/backend/pkg/integration/shipping/ghtk"
	"o.o/backend/pkg/integration/shipping/ghtk/webhook"
	"o.o/backend/pkg/integration/shipping/vtpost"
	webhook2 "o.o/backend/pkg/integration/shipping/vtpost/webhook"
	"o.o/backend/pkg/integration/sms"
)

// Injectors from wire.go:

func Build(ctx context.Context, cfg config.Config, consumer mq.KafkaConsumer) (Output, func(), error) {
	redisRedis := cfg.Redis
	store := redis.Connect(redisRedis)
	service := health.New(store)
	miscService := &api.MiscService{}
	sharedConfig := cfg.SharedConfig
	session := config_server.NewSession(sharedConfig, store)
	database_minConfig := cfg.Databases
	databases, err := database_min.BuildDatabases(database_minConfig)
	if err != nil {
		return Output{}, nil, err
	}
	mainDB := databases.Main
	manager := SupportedShipnowManager()
	identityAggregate := identity.NewAggregate(mainDB, manager)
	commandBus := identity.AggregateMessageBus(identityAggregate)
	queryService := identity.NewQueryService(mainDB)
	queryBus := identity.QueryServiceMessageBus(queryService)
	invitationQuery := query.NewInvitationQuery(mainDB)
	invitationQueryBus := query.InvitationQueryMessageBus(invitationQuery)
	busBus := bus.New()
	generator := auth.NewGenerator(store)
	tokenStore := tokens.NewTokenStore(store)
	smsConfig := cfg.SMS
	v := sms_min.SupportedSMSDrivers(smsConfig)
	logDB := databases.Log
	smsLogAggregate := aggregate.NewSmsLogAggregate(busBus, logDB)
	smslogCommandBus := aggregate.SmsLogAggregateMessageBus(smsLogAggregate)
	client := sms.New(smsConfig, v, smslogCommandBus)
	smtpConfig := cfg.SMTP
	emailClient := email.New(smtpConfig)
	userService := &api.UserService{
		Session:         session,
		IdentityAggr:    commandBus,
		IdentityQuery:   queryBus,
		InvitationQuery: invitationQueryBus,
		EventBus:        busBus,
		AuthStore:       generator,
		TokenStore:      tokenStore,
		RedisStore:      store,
		SMSClient:       client,
		EmailClient:     emailClient,
	}
	accountService := &api.AccountService{
		Session: session,
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
	addressService := &api.AddressService{
		Session: session,
	}
	invitationConfig := cfg.Invitation
	customerQuery := query2.NewCustomerQuery(mainDB)
	customeringQueryBus := query2.CustomerQueryMessageBus(customerQuery)
	secretToken := cfg.Secret
	flagEnableNewLinkInvitation := cfg.FlagEnableNewLinkInvitation
	invitationAggregate := aggregate2.NewInvitationAggregate(mainDB, invitationConfig, customeringQueryBus, queryBus, busBus, client, emailClient, secretToken, flagEnableNewLinkInvitation)
	invitationCommandBus := aggregate2.InvitationAggregateMessageBus(invitationAggregate)
	authorizationAggregate := aggregate3.NewAuthorizationAggregate()
	authorizationCommandBus := aggregate3.AuthorizationAggregateMessageBus(authorizationAggregate)
	accountRelationshipService := &api.AccountRelationshipService{
		Session:           session,
		InvitationAggr:    invitationCommandBus,
		InvitationQuery:   invitationQueryBus,
		AuthorizationAggr: authorizationCommandBus,
	}
	userRelationshipService := &api.UserRelationshipService{
		Session:                session,
		InvitationAggr:         invitationCommandBus,
		InvitationQuery:        invitationQueryBus,
		AuthorizationAggregate: authorizationCommandBus,
	}
	ecomService := &api.EcomService{
		Session: session,
	}
	emailConfig := cfg.Email
	servers, cleanup := api.NewServers(miscService, userService, accountService, locationService, bankService, addressService, accountRelationshipService, userRelationshipService, ecomService, store, emailConfig, smsConfig)
	shopMiscService := &shop.MiscService{
		Session: session,
	}
	queryQueryService := query3.New(mainDB)
	catalogQueryBus := query3.QueryServiceMessageBus(queryQueryService)
	aggregateAggregate := aggregate4.New(busBus, mainDB)
	catalogCommandBus := aggregate4.AggregateMessageBus(aggregateAggregate)
	brandService := &shop.BrandService{
		Session:      session,
		CatalogQuery: catalogQueryBus,
		CatalogAggr:  catalogCommandBus,
	}
	carrierQuery := query4.NewCarrierQuery(mainDB)
	carryingQueryBus := query4.CarrierQueryMessageBus(carrierQuery)
	supplierQuery := query5.NewSupplierQuery(mainDB)
	supplieringQueryBus := query5.SupplierQueryMessageBus(supplierQuery)
	traderQuery := query6.NewTraderQuery(mainDB, customeringQueryBus, carryingQueryBus, supplieringQueryBus)
	traderingQueryBus := query6.TraderQueryMessageBus(traderQuery)
	stocktakeQuery := query7.NewQueryStocktake(mainDB)
	stocktakingQueryBus := query7.StocktakeQueryMessageBus(stocktakeQuery)
	inventoryQueryService := query8.NewQueryInventory(stocktakingQueryBus, busBus, mainDB)
	inventoryQueryBus := query8.InventoryQueryServiceMessageBus(inventoryQueryService)
	receiptQuery := query9.NewReceiptQuery(mainDB)
	receiptingQueryBus := query9.ReceiptQueryMessageBus(receiptQuery)
	purchaseOrderQuery := query10.NewPurchaseOrderQuery(mainDB, busBus, supplieringQueryBus, inventoryQueryBus, receiptingQueryBus)
	purchaseorderQueryBus := query10.PurchaseOrderQueryMessageBus(purchaseOrderQuery)
	refundQueryService := query11.NewQueryRefund(busBus, mainDB)
	refundQueryBus := query11.RefundQueryServiceMessageBus(refundQueryService)
	purchaseRefundQueryService := query12.NewQueryPurchasePurchaseRefund(busBus, mainDB)
	purchaserefundQueryBus := query12.PurchaseRefundQueryServiceMessageBus(purchaseRefundQueryService)
	inventoryAggregate := aggregate5.NewAggregateInventory(busBus, mainDB, traderingQueryBus, purchaseorderQueryBus, stocktakingQueryBus, refundQueryBus, purchaserefundQueryBus, catalogQueryBus)
	inventoryCommandBus := aggregate5.InventoryAggregateMessageBus(inventoryAggregate)
	inventoryService := &shop.InventoryService{
		Session:        session,
		TraderQuery:    traderingQueryBus,
		InventoryAggr:  inventoryCommandBus,
		InventoryQuery: inventoryQueryBus,
	}
	addressQueryService := address.NewQueryService(mainDB)
	addressQueryBus := address.QueryServiceMessageBus(addressQueryService)
	shopAccountService := &shop.AccountService{
		Session:       session,
		IdentityAggr:  commandBus,
		IdentityQuery: queryBus,
		AddressQuery:  addressQueryBus,
	}
	collectionService := &shop.CollectionService{
		Session:      session,
		CatalogQuery: catalogQueryBus,
		CatalogAggr:  catalogCommandBus,
	}
	customerAggregate := aggregate6.NewCustomerAggregate(busBus, mainDB)
	customeringCommandBus := aggregate6.CustomerAggregateMessageBus(customerAggregate)
	addressAggregate := aggregate6.NewAddressAggregate(mainDB)
	addressingCommandBus := aggregate6.AddressAggregateMessageBus(addressAggregate)
	addressQuery := query2.NewAddressQuery(mainDB)
	addressingQueryBus := query2.AddressQueryMessageBus(addressQuery)
	orderingQueryService := ordering.NewQueryService(mainDB)
	orderingQueryBus := ordering.QueryServiceMessageBus(orderingQueryService)
	customerService := &shop.CustomerService{
		Session:       session,
		LocationQuery: locationQueryBus,
		CustomerQuery: customeringQueryBus,
		CustomerAggr:  customeringCommandBus,
		AddressAggr:   addressingCommandBus,
		AddressQuery:  addressingQueryBus,
		OrderQuery:    orderingQueryBus,
		ReceiptQuery:  receiptingQueryBus,
	}
	customerGroupService := &shop.CustomerGroupService{
		Session:       session,
		CustomerAggr:  customeringCommandBus,
		CustomerQuery: customeringQueryBus,
	}
	productService := &shop.ProductService{
		Session:        session,
		CatalogQuery:   catalogQueryBus,
		CatalogAggr:    catalogCommandBus,
		InventoryQuery: inventoryQueryBus,
	}
	categoryService := &shop.CategoryService{
		Session:      session,
		CatalogQuery: catalogQueryBus,
		CatalogAggr:  catalogCommandBus,
	}
	productSourceService := &shop.ProductSourceService{
		Session: session,
	}
	orderingAggregate := ordering.NewAggregate(busBus, mainDB)
	orderingCommandBus := ordering.AggregateMessageBus(orderingAggregate)
	notifierDB := databases.Notifier
	sqlstoreStore := sqlstore.New(mainDB, notifierDB, locationQueryBus, busBus)
	shipment_allConfig := cfg.Shipment
	v3 := shipment_all.SupportedCarrierDrivers(ctx, sqlstoreStore, shipment_allConfig, locationQueryBus)
	carrierManager := shipping_provider.NewCtrl(busBus, locationQueryBus, v3)
	flagFaboOrderAutoConfirmPaymentStatus := cfg.FlagFaboOrderAutoConfirmPaymentStatus
	connectionQuery := query13.NewConnectionQuery(mainDB)
	connectioningQueryBus := query13.ConnectionQueryMessageBus(connectionQuery)
	connectionAggregate := aggregate7.NewConnectionAggregate(mainDB, busBus)
	connectioningCommandBus := aggregate7.ConnectionAggregateMessageBus(connectionAggregate)
	shipmentserviceQueryService := shipmentservice.NewQueryService(mainDB, store)
	shipmentserviceQueryBus := shipmentservice.QueryServiceMessageBus(shipmentserviceQueryService)
	pricelistQueryService := pricelist.NewQueryService(mainDB, store)
	pricelistQueryBus := pricelist.QueryServiceMessageBus(pricelistQueryService)
	shopshipmentpricelistQueryService := shopshipmentpricelist.NewQueryService(mainDB, store)
	shopshipmentpricelistQueryBus := shopshipmentpricelist.QueryServiceMessageBus(shopshipmentpricelistQueryService)
	shipmentpriceQueryService := shipmentprice.NewQueryService(mainDB, store, locationQueryBus, pricelistQueryBus, shopshipmentpricelistQueryBus)
	shipmentpriceQueryBus := shipmentprice.QueryServiceMessageBus(shipmentpriceQueryService)
	flagApplyShipmentPrice := cfg.FlagApplyShipmentPrice
	carrierConfig := shipment_all.SupportedShippingCarrierConfig(shipment_allConfig)
	shipmentManager, err := carrier.NewShipmentManager(busBus, locationQueryBus, connectioningQueryBus, connectioningCommandBus, store, shipmentserviceQueryBus, shipmentpriceQueryBus, flagApplyShipmentPrice, carrierConfig)
	if err != nil {
		cleanup()
		return Output{}, nil, err
	}
	orderLogic := orderS.New(carrierManager, catalogQueryBus, orderingCommandBus, customeringCommandBus, customeringQueryBus, addressingCommandBus, addressingQueryBus, locationQueryBus, busBus, flagFaboOrderAutoConfirmPaymentStatus, shipmentManager)
	orderService := &shop.OrderService{
		Session:       session,
		OrderAggr:     orderingCommandBus,
		CustomerQuery: customeringQueryBus,
		OrderQuery:    orderingQueryBus,
		ReceiptQuery:  receiptingQueryBus,
		OrderLogic:    orderLogic,
	}
	queryService2 := query14.NewQueryService(mainDB)
	shippingQueryBus := query14.QueryServiceMessageBus(queryService2)
	fulfillmentService := &shop.FulfillmentService{
		Session:         session,
		ShipmentManager: shipmentManager,
		ShippingQuery:   shippingQueryBus,
		ShippingCtrl:    carrierManager,
	}
	shipnowAggregate := shipnow.NewAggregate(busBus, mainDB, locationQueryBus, queryBus, addressQueryBus, orderingQueryBus, manager)
	shipnowCommandBus := shipnow.AggregateMessageBus(shipnowAggregate)
	shipnowQueryService := shipnow.NewQueryService(mainDB)
	shipnowQueryBus := shipnow.QueryServiceMessageBus(shipnowQueryService)
	shipnowService := &shop.ShipnowService{
		Session:      session,
		ShipnowAggr:  shipnowCommandBus,
		ShipnowQuery: shipnowQueryBus,
	}
	historyService := &shop.HistoryService{
		Session: session,
	}
	moneyTxQuery := query15.NewMoneyTxQuery(mainDB, shippingQueryBus)
	moneytxQueryBus := query15.MoneyTxQueryMessageBus(moneyTxQuery)
	moneyTransactionService := &shop.MoneyTransactionService{
		Session:      session,
		MoneyTxQuery: moneytxQueryBus,
	}
	dashboardQuery := query16.NewDashboardQuery(mainDB, store, locationQueryBus)
	summaryQueryBus := query16.DashboardQueryMessageBus(dashboardQuery)
	summarySummary := summary.New(mainDB)
	summaryService := &shop.SummaryService{
		Session:      session,
		SummaryQuery: summaryQueryBus,
		SummaryOld:   summarySummary,
	}
	eventStream := eventstream.New(ctx)
	exportConfig := cfg.Export
	exportService, cleanup2 := export.New(store, eventStream, exportConfig)
	shopExportService := &shop.ExportService{
		Session:     session,
		ExportInner: exportService,
	}
	notificationService := &shop.NotificationService{
		Session: session,
	}
	authorizeService := &shop.AuthorizeService{
		Session: session,
	}
	ledgerQuery := query17.NewLedgerQuery(mainDB)
	ledgeringQueryBus := query17.LedgerQueryMessageBus(ledgerQuery)
	receiptAggregate := aggregate8.NewReceiptAggregate(mainDB, busBus, traderingQueryBus, ledgeringQueryBus, orderingQueryBus, customeringQueryBus, carryingQueryBus, supplieringQueryBus, purchaseorderQueryBus)
	receiptingCommandBus := aggregate8.ReceiptAggregateMessageBus(receiptAggregate)
	receiptService := &shop.ReceiptService{
		Session:       session,
		CarrierQuery:  carryingQueryBus,
		CustomerQuery: customeringQueryBus,
		LedgerQuery:   ledgeringQueryBus,
		ReceiptAggr:   receiptingCommandBus,
		ReceiptQuery:  receiptingQueryBus,
		SupplierQuery: supplieringQueryBus,
		TraderQuery:   traderingQueryBus,
	}
	carrierAggregate := aggregate9.NewCarrierAggregate(busBus, mainDB)
	carryingCommandBus := aggregate9.CarrierAggregateMessageBus(carrierAggregate)
	carrierService := &shop.CarrierService{
		Session:      session,
		CarrierAggr:  carryingCommandBus,
		CarrierQuery: carryingQueryBus,
	}
	ledgerAggregate := aggregate10.NewLedgerAggregate(mainDB, receiptingQueryBus)
	ledgeringCommandBus := aggregate10.LedgerAggregateMessageBus(ledgerAggregate)
	ledgerService := &shop.LedgerService{
		Session:     session,
		LedgerAggr:  ledgeringCommandBus,
		LedgerQuery: ledgeringQueryBus,
	}
	purchaseOrderAggregate := aggregate11.NewPurchaseOrderAggregate(mainDB, busBus, catalogQueryBus, supplieringQueryBus, inventoryQueryBus)
	purchaseorderCommandBus := aggregate11.PurchaseOrderAggregateMessageBus(purchaseOrderAggregate)
	purchaseOrderService := &shop.PurchaseOrderService{
		Session:            session,
		PurchaseOrderAggr:  purchaseorderCommandBus,
		PurchaseOrderQuery: purchaseorderQueryBus,
	}
	stocktakeAggregate := aggregate12.NewAggregateStocktake(mainDB, busBus, store)
	stocktakingCommandBus := aggregate12.StocktakeAggregateMessageBus(stocktakeAggregate)
	stocktakeService := &shop.StocktakeService{
		Session:        session,
		CatalogQuery:   catalogQueryBus,
		StocktakeAggr:  stocktakingCommandBus,
		StocktakeQuery: stocktakingQueryBus,
		InventoryQuery: inventoryQueryBus,
	}
	aggregate17 := aggregate13.NewAggregate(mainDB, busBus, locationQueryBus, orderingQueryBus, shipmentManager, connectioningQueryBus)
	shippingCommandBus := aggregate13.AggregateMessageBus(aggregate17)
	shipmentService := &shop.ShipmentService{
		Session:           session,
		ShipmentManager:   shipmentManager,
		ShippingAggregate: shippingCommandBus,
	}
	connectionService := &shop.ConnectionService{
		Session:         session,
		ShipmentManager: shipmentManager,
		ConnectionQuery: connectioningQueryBus,
		ConnectionAggr:  connectioningCommandBus,
	}
	refundAggregate := aggregate14.NewRefundAggregate(mainDB, busBus)
	refundCommandBus := aggregate14.RefundAggregateMessageBus(refundAggregate)
	refundService := &shop.RefundService{
		Session:        session,
		CustomerQuery:  customeringQueryBus,
		InventoryQuery: inventoryQueryBus,
		ReceiptQuery:   receiptingQueryBus,
		RefundAggr:     refundCommandBus,
		RefundQuery:    refundQueryBus,
	}
	purchaseRefundAggregate := aggregate15.NewPurchaseRefundAggregate(mainDB, busBus, purchaseorderQueryBus)
	purchaserefundCommandBus := aggregate15.PurchaseRefundAggregateMessageBus(purchaseRefundAggregate)
	purchaseRefundService := &shop.PurchaseRefundService{
		Session:             session,
		PurchaseRefundAggr:  purchaserefundCommandBus,
		PurchaseRefundQuery: purchaserefundQueryBus,
		SupplierQuery:       supplieringQueryBus,
		PurchaseOrderQuery:  purchaseorderQueryBus,
		InventoryQuery:      inventoryQueryBus,
	}
	shopServers := shop_min.NewServers(store, shopMiscService, brandService, inventoryService, shopAccountService, collectionService, customerService, customerGroupService, productService, categoryService, productSourceService, orderService, fulfillmentService, shipnowService, historyService, moneyTransactionService, summaryService, shopExportService, notificationService, authorizeService, receiptService, carrierService, ledgerService, purchaseOrderService, stocktakeService, shipmentService, connectionService, refundService, purchaseRefundService)
	adminMiscService := admin.MiscService{
		Session: session,
	}
	adminAccountService := admin.AccountService{
		Session: session,
	}
	adminOrderService := admin.OrderService{
		Session: session,
	}
	adminFulfillmentService := admin.FulfillmentService{
		Session:       session,
		EventBus:      busBus,
		IdentityQuery: queryBus,
		RedisStore:    store,
		ShippingAggr:  shippingCommandBus,
		ShippingQuery: shippingQueryBus,
	}
	moneyTxAggregate := aggregate16.NewMoneyTxAggregate(mainDB, shippingQueryBus, queryBus, busBus)
	moneytxCommandBus := aggregate16.MoneyTxAggregateMessageBus(moneyTxAggregate)
	adminMoneyTransactionService := admin.MoneyTransactionService{
		Session:      session,
		MoneyTxQuery: moneytxQueryBus,
		MoneyTxAggr:  moneytxCommandBus,
	}
	shopService := admin.ShopService{
		Session:       session,
		IdentityQuery: queryBus,
	}
	creditAggregate := credit.NewAggregateCredit(busBus, mainDB, queryBus)
	creditCommandBus := credit.CreditAggregateMessageBus(creditAggregate)
	creditQueryService := credit.NewQueryCredit(busBus, mainDB, queryBus)
	creditQueryBus := credit.CreditQueryServiceMessageBus(creditQueryService)
	creditService := admin.CreditService{
		Session:     session,
		CreditAggr:  creditCommandBus,
		CreditQuery: creditQueryBus,
	}
	adminNotificationService := admin.NotificationService{
		Session: session,
	}
	adminConnectionService := admin.ConnectionService{
		Session:         session,
		ConnectionAggr:  connectioningCommandBus,
		ConnectionQuery: connectioningQueryBus,
	}
	shipmentpriceAggregate := shipmentprice.NewAggregate(mainDB, store, pricelistQueryBus, shipmentserviceQueryBus)
	shipmentpriceCommandBus := shipmentprice.AggregateMessageBus(shipmentpriceAggregate)
	shipmentserviceAggregate := shipmentservice.NewAggregate(mainDB, store)
	shipmentserviceCommandBus := shipmentservice.AggregateMessageBus(shipmentserviceAggregate)
	pricelistAggregate := pricelist.NewAggregate(mainDB, busBus, shopshipmentpricelistQueryBus)
	pricelistCommandBus := pricelist.AggregateMessageBus(pricelistAggregate)
	shopshipmentpricelistAggregate := shopshipmentpricelist.NewAggregate(mainDB, pricelistQueryBus)
	shopshipmentpricelistCommandBus := shopshipmentpricelist.AggregateMessageBus(shopshipmentpricelistAggregate)
	shipmentPriceService := admin.ShipmentPriceService{
		Session:                    session,
		ShipmentManager:            shipmentManager,
		ShipmentPriceAggr:          shipmentpriceCommandBus,
		ShipmentPriceQuery:         shipmentpriceQueryBus,
		ShipmentServiceQuery:       shipmentserviceQueryBus,
		ShipmentServiceAggr:        shipmentserviceCommandBus,
		ShipmentPriceListAggr:      pricelistCommandBus,
		ShipmentPriceListQuery:     pricelistQueryBus,
		ShopShipmentPriceListQuery: shopshipmentpricelistQueryBus,
		ShopShipmentPriceListAggr:  shopshipmentpricelistCommandBus,
	}
	locationAggregate := location.NewAggregate(mainDB)
	locationCommandBus := location.AggregateMessageBus(locationAggregate)
	adminLocationService := admin.LocationService{
		Session:       session,
		LocationAggr:  locationCommandBus,
		LocationQuery: locationQueryBus,
	}
	adminServers := admin_min.NewServers(store, adminMiscService, adminAccountService, adminOrderService, adminFulfillmentService, adminMoneyTransactionService, shopService, creditService, adminNotificationService, adminConnectionService, shipmentPriceService, adminLocationService)
	sadminMiscService := &sadmin.MiscService{
		Session: session,
	}
	sadminUserService := &sadmin.UserService{
		Session: session,
	}
	webhookCallbackService := sadmin.NewWebhookCallbackService(store)
	webhookService := &sadmin.WebhookService{
		Session:                session,
		WebhookCallbackService: webhookCallbackService,
	}
	sadminServers := sadmin.NewServers(sadminMiscService, sadminUserService, webhookService)
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
	faboServers := fabo.NewServers(pageService, customerConversationService, faboCustomerService, store)
	intHandlers := BuildIntHandlers(servers, shopServers, adminServers, sadminServers, faboServers)
	ghnimportImport := ghnimport.Import{
		MoneyTxAggr: moneytxCommandBus,
	}
	ghtkimportImport := ghtkimport.Import{
		MoneyTxAggr:   moneytxCommandBus,
		ShippingAggr:  shippingCommandBus,
		ShippingQuery: shippingQueryBus,
	}
	vtpostimportImport := vtpostimport.Import{
		MoneyTxAggr: moneytxCommandBus,
	}
	vtPostImporter := &vtpostimport.VTPostImporter{}
	ghtkImporter := &ghtkimport.GHTKImporter{
		ShippingAggr:  shippingCommandBus,
		ShippingQuery: shippingQueryBus,
	}
	ghnImporter := &ghnimport.GHNImporter{}
	jtImporter := &jtexpressimport.JTImporter{}
	importService := handlers.ImportService{
		MoneyTxAggr:       moneytxCommandBus,
		ConnectionQuery:   connectioningQueryBus,
		VTPostImporter:    vtPostImporter,
		GHTKImporter:      ghtkImporter,
		GHNImporter:       ghnImporter,
		JTExpressImporter: jtImporter,
	}
	hotFixMoneyTxService := hotfix.New(mainDB)
	importServer := server_admin.BuildImportHandlers(ghnimportImport, ghtkimportImport, vtpostimportImport, importService, hotFixMoneyTxService, session)
	uploadConfig := cfg.Upload
	uploader, err := _uploader.NewUploader(uploadConfig)
	if err != nil {
		cleanup2()
		cleanup()
		return Output{}, nil, err
	}
	imcsvImport, cleanup3 := imcsv.New(locationQueryBus, store, uploader, mainDB)
	import2, cleanup4 := imcsv2.New(store, uploader, mainDB)
	importHandler := server_shop.BuildImportHandler(imcsvImport, import2, session)
	eventStreamHandler := server_shop.BuildEventStreamHandler(eventStream, session)
	downloadHandler := server_shop.BuildDownloadHandler()
	mainServer := BuildMainServer(service, intHandlers, sharedConfig, importServer, importHandler, eventStreamHandler, downloadHandler)
	webhookConfig := shipment_allConfig.GHNWebhook
	ghnConfig := shipment_allConfig.GHN
	ghnCarrier := ghn.New(ghnConfig, locationQueryBus)
	v2Webhook := v2.New(mainDB, logDB, ghnCarrier, shipmentManager, queryBus, shippingCommandBus)
	ghnWebhookServer := v2_2.NewGHNWebhookServer(webhookConfig, shipmentManager, ghnCarrier, queryBus, shippingCommandBus, v2Webhook)
	_ghtkWebhookConfig := shipment_allConfig.GHTKWebhook
	ghtkConfig := shipment_allConfig.GHTK
	ghtkCarrier := ghtk.New(ghtkConfig, locationQueryBus)
	webhookWebhook := webhook.New(mainDB, logDB, ghtkCarrier, shipmentManager, queryBus, shippingCommandBus)
	ghtkWebhookServer := _ghtk.NewGHTKWebhookServer(_ghtkWebhookConfig, shipmentManager, ghtkCarrier, queryBus, shippingCommandBus, webhookWebhook)
	_vtpostWebhookConfig := shipment_allConfig.VTPostWebhook
	vtpostConfig := shipment_allConfig.VTPost
	vtpostCarrier := vtpost.New(vtpostConfig, locationQueryBus)
	webhook4 := webhook2.New(mainDB, logDB, vtpostCarrier, shipmentManager, queryBus, shippingCommandBus)
	vtPostWebhookServer := _vtpost.NewVTPostWebhookServer(_vtpostWebhookConfig, shipmentManager, vtpostCarrier, queryBus, shippingCommandBus, webhook4)
	configWebhookConfig := cfg.Webhook
	faboRedis := redis2.NewFaboRedis(store)
	webhook5 := webhook3.New(mainDB, logDB, store, configWebhookConfig, faboRedis, fbClient, fbmessagingQueryBus, fbmessagingCommandBus, fbpagingQueryBus)
	fbWebhookServer := BuildWebhookServer(configWebhookConfig, webhook5)
	v4 := BuildServers(mainServer, ghnWebhookServer, ghtkWebhookServer, vtPostWebhookServer, fbWebhookServer)
	kafka := cfg.Kafka
	handlerHandler := handler.New(consumer, kafka)
	publisherPublisher := publisher.New(consumer, eventStream)
	processManager := pm.New(busBus, catalogQueryBus, catalogCommandBus)
	pmProcessManager := pm2.New(busBus, queryBus, invitationQueryBus)
	processManager2 := pm3.New(busBus, catalogQueryBus, orderingQueryBus, inventoryCommandBus)
	processManager3 := pm4.New(busBus, invitationQueryBus, invitationCommandBus)
	processManager4 := pm5.New(busBus, ledgeringCommandBus)
	processManager5 := pm6.New(busBus, moneytxQueryBus, moneytxCommandBus, shippingQueryBus)
	affiliateCommandBus := _wireCommandBusValue
	processManager6 := pm7.New(busBus, orderingCommandBus, affiliateCommandBus, receiptingQueryBus, inventoryCommandBus, orderingQueryBus, customeringQueryBus)
	processManager7 := pm8.New(busBus, receiptingQueryBus, receiptingCommandBus, ledgeringQueryBus, ledgeringCommandBus, queryBus)
	processManager8 := pm9.New(busBus, refundQueryBus, receiptingQueryBus, refundCommandBus)
	processManager9 := pm10.New(busBus, shippingQueryBus, shippingCommandBus, store, connectioningQueryBus)
	processManager10 := pm11.New(busBus, fbuseringCommandBus)
	fbmessagingProcessManager := fbmessaging.NewProcessManager(busBus, fbmessagingQueryBus, fbmessagingCommandBus, fbpagingQueryBus, fbuseringQueryBus, fbuseringCommandBus, faboRedis)
	sAdminToken := config_server.WireSAdminToken(sharedConfig)
	middlewareMiddleware := middleware.New(sAdminToken, tokenStore, queryBus)
	captchaConfig := cfg.Captcha
	captchaCaptcha := captcha.New(captchaConfig)
	output := Output{
		Servers:        v4,
		EventStream:    eventStream,
		Health:         service,
		Handler:        handlerHandler,
		Publisher:      publisherPublisher,
		_catalogPM:     processManager,
		_identityPM:    pmProcessManager,
		_inventoryPM:   processManager2,
		_invitationPM:  processManager3,
		_ledgerPM:      processManager4,
		_moneytxPM:     processManager5,
		_orderPM:       processManager6,
		_receiptPM:     processManager7,
		_refundPM:      processManager8,
		_shippingPM:    processManager9,
		_fbuserPM:      processManager10,
		_fbMessagingPM: fbmessagingProcessManager,
		_s:             sqlstoreStore,
		_m:             middlewareMiddleware,
		_c:             captchaCaptcha,
	}
	return output, func() {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

var (
	_wireCommandBusValue = affiliate.CommandBus{}
)
