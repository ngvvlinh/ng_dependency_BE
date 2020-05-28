// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"o.o/backend/cmd/etop-server/config"
	aggregate21 "o.o/backend/com/etc/logging/payment/aggregate"
	"o.o/backend/com/etc/logging/smslog/aggregate"
	"o.o/backend/com/external/payment/manager"
	aggregate19 "o.o/backend/com/external/payment/payment/aggregate"
	vtpay2 "o.o/backend/com/external/payment/vtpay"
	aggregate20 "o.o/backend/com/external/payment/vtpay/gateway/aggregate"
	"o.o/backend/com/external/payment/vtpay/gateway/server"
	"o.o/backend/com/main/address"
	aggregate3 "o.o/backend/com/main/authorization/aggregate"
	aggregate4 "o.o/backend/com/main/catalog/aggregate"
	query3 "o.o/backend/com/main/catalog/query"
	aggregate7 "o.o/backend/com/main/connectioning/aggregate"
	query13 "o.o/backend/com/main/connectioning/query"
	"o.o/backend/com/main/identity"
	aggregate5 "o.o/backend/com/main/inventory/aggregate"
	query8 "o.o/backend/com/main/inventory/query"
	aggregate2 "o.o/backend/com/main/invitation/aggregate"
	"o.o/backend/com/main/invitation/query"
	aggregate11 "o.o/backend/com/main/ledgering/aggregate"
	query16 "o.o/backend/com/main/ledgering/query"
	"o.o/backend/com/main/location"
	aggregate18 "o.o/backend/com/main/moneytx/aggregate"
	query18 "o.o/backend/com/main/moneytx/query"
	"o.o/backend/com/main/ordering"
	aggregate12 "o.o/backend/com/main/purchaseorder/aggregate"
	query10 "o.o/backend/com/main/purchaseorder/query"
	aggregate16 "o.o/backend/com/main/purchaserefund/aggregate"
	query12 "o.o/backend/com/main/purchaserefund/query"
	aggregate8 "o.o/backend/com/main/receipting/aggregate"
	query9 "o.o/backend/com/main/receipting/query"
	aggregate15 "o.o/backend/com/main/refund/aggregate"
	query11 "o.o/backend/com/main/refund/query"
	"o.o/backend/com/main/shipmentpricing/pricelist"
	"o.o/backend/com/main/shipmentpricing/shipmentprice"
	"o.o/backend/com/main/shipmentpricing/shipmentservice"
	"o.o/backend/com/main/shipnow"
	"o.o/backend/com/main/shipnow-carrier"
	aggregate14 "o.o/backend/com/main/shipping/aggregate"
	"o.o/backend/com/main/shipping/carrier"
	query14 "o.o/backend/com/main/shipping/query"
	aggregate13 "o.o/backend/com/main/stocktaking/aggregate"
	query7 "o.o/backend/com/main/stocktaking/query"
	affiliate2 "o.o/backend/com/services/affiliate"
	aggregate10 "o.o/backend/com/shopping/carrying/aggregate"
	query4 "o.o/backend/com/shopping/carrying/query"
	aggregate6 "o.o/backend/com/shopping/customering/aggregate"
	query2 "o.o/backend/com/shopping/customering/query"
	aggregate9 "o.o/backend/com/shopping/suppliering/aggregate"
	query5 "o.o/backend/com/shopping/suppliering/query"
	query6 "o.o/backend/com/shopping/tradering/query"
	"o.o/backend/com/subscripting/subscription"
	"o.o/backend/com/subscripting/subscriptionbill"
	"o.o/backend/com/subscripting/subscriptionplan"
	"o.o/backend/com/subscripting/subscriptionproduct"
	query15 "o.o/backend/com/summary/query"
	aggregate17 "o.o/backend/com/web/webserver/aggregate"
	query17 "o.o/backend/com/web/webserver/query"
	"o.o/backend/pkg/common/apifw/service"
	"o.o/backend/pkg/common/authorization/auth"
	"o.o/backend/pkg/common/extservice/telebot"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/admin"
	"o.o/backend/pkg/etop/api/affiliate"
	"o.o/backend/pkg/etop/api/integration"
	"o.o/backend/pkg/etop/api/sadmin"
	"o.o/backend/pkg/etop/api/shop"
	"o.o/backend/pkg/etop/api/shop/imports"
	"o.o/backend/pkg/etop/apix/partner"
	"o.o/backend/pkg/etop/apix/partnercarrier"
	"o.o/backend/pkg/etop/apix/partnerimport"
	"o.o/backend/pkg/etop/apix/shipping"
	"o.o/backend/pkg/etop/apix/shop"
	"o.o/backend/pkg/etop/apix/shopping"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/eventstream"
	"o.o/backend/pkg/etop/logic/money-transaction/ghtk-imcsv"
	"o.o/backend/pkg/etop/logic/money-transaction/imcsv"
	"o.o/backend/pkg/etop/logic/money-transaction/vtpost-imxlsx"
	"o.o/backend/pkg/etop/logic/orders"
	imcsv2 "o.o/backend/pkg/etop/logic/orders/imcsv"
	imcsv3 "o.o/backend/pkg/etop/logic/products/imcsv"
	"o.o/backend/pkg/etop/logic/shipping_provider"
	"o.o/backend/pkg/integration/payment/vtpay"
	client2 "o.o/backend/pkg/integration/payment/vtpay/client"
	"o.o/backend/pkg/integration/shipnow/ahamove"
	"o.o/backend/pkg/integration/shipnow/ahamove/client"
	"o.o/backend/pkg/integration/shipping/ghn"
	"o.o/backend/pkg/integration/shipping/ghtk"
	"o.o/backend/pkg/integration/shipping/vtpost"
	"o.o/backend/pkg/integration/sms"
	api2 "o.o/backend/pkg/services/affiliate/api"
	"o.o/capi"
)

// Injectors from wire.go:

func BuildServers(db2 *cmsql.Database, cfg2 config.Config, bot2 *telebot.Channel, sd cmService.Shutdowner, eventBus capi.EventBus, rd redis.Store, s auth.Generator, ss session.Session, authURL partner.AuthURL) ([]Server, error) {
	miscService := &api.MiscService{}
	locationQuery := location.New(db2)
	queryBus := location.QueryMessageBus(locationQuery)
	queryService := shipnow.NewQueryService(db2)
	shipnowQueryBus := shipnow.QueryServiceMessageBus(queryService)
	clientConfig := cfg2.Ahamove
	clientClient := client.New(clientConfig)
	identityQueryService := identity.NewQueryService(db2)
	identityQueryBus := identity.QueryServiceMessageBus(identityQueryService)
	ahamoveCarrier := ahamove.New(clientClient, queryBus, identityQueryBus)
	urlConfig := AhamoveConfig(cfg2)
	carrierAccount := ahamove.NewCarrierAccount(clientClient, urlConfig, identityQueryBus)
	v := SupportedShipnowCarriers(ahamoveCarrier, carrierAccount)
	shipnowManager := shipnow_carrier.NewManager(db2, queryBus, shipnowQueryBus, v)
	identityAggregate := identity.NewAggregate(db2, shipnowManager)
	commandBus := identity.AggregateMessageBus(identityAggregate)
	invitationQuery := query.NewInvitationQuery(db2)
	invitationQueryBus := query.InvitationQueryMessageBus(invitationQuery)
	smsConfig := cfg2.SMS
	v2 := SupportedSMSDrivers(cfg2, smsConfig)
	smsLogAggregate := aggregate.NewSmsLogAggregate(eventBus, db2)
	smslogCommandBus := aggregate.SmsLogAggregateMessageBus(smsLogAggregate)
	smsClient := sms.New(smsConfig, bot2, v2, smslogCommandBus)
	userService := &api.UserService{
		IdentityAggr:    commandBus,
		IdentityQuery:   identityQueryBus,
		InvitationQuery: invitationQueryBus,
		EventBus:        eventBus,
		AuthStore:       s,
		RedisStore:      rd,
		SMSClient:       smsClient,
	}
	accountService := &api.AccountService{}
	locationService := &api.LocationService{
		LocationQuery: queryBus,
	}
	bankService := &api.BankService{}
	addressService := &api.AddressService{}
	invitationConfig := cfg2.Invitation
	customerQuery := query2.NewCustomerQuery(db2)
	customeringQueryBus := query2.CustomerQueryMessageBus(customerQuery)
	invitationAggregate := aggregate2.NewInvitationAggregate(db2, invitationConfig, customeringQueryBus, identityQueryBus, eventBus, cfg2, smsClient)
	invitationCommandBus := aggregate2.InvitationAggregateMessageBus(invitationAggregate)
	authorizationAggregate := aggregate3.NewAuthorizationAggregate()
	authorizationCommandBus := aggregate3.AuthorizationAggregateMessageBus(authorizationAggregate)
	accountRelationshipService := &api.AccountRelationshipService{
		InvitationAggr:    invitationCommandBus,
		InvitationQuery:   invitationQueryBus,
		AuthorizationAggr: authorizationCommandBus,
	}
	userRelationshipService := &api.UserRelationshipService{
		InvitationAggr:         invitationCommandBus,
		InvitationQuery:        invitationQueryBus,
		AuthorizationAggregate: authorizationCommandBus,
	}
	ecomService := &api.EcomService{}
	emailConfig := cfg2.Email
	servers := api.NewServers(miscService, userService, accountService, locationService, bankService, addressService, accountRelationshipService, userRelationshipService, ecomService, rd, emailConfig, smsConfig, bot2, sd)
	shopMiscService := &shop.MiscService{}
	queryQueryService := query3.New(db2)
	catalogQueryBus := query3.QueryServiceMessageBus(queryQueryService)
	aggregateAggregate := aggregate4.New(eventBus, db2)
	catalogCommandBus := aggregate4.AggregateMessageBus(aggregateAggregate)
	brandService := &shop.BrandService{
		CatalogQuery: catalogQueryBus,
		CatalogAggr:  catalogCommandBus,
	}
	carrierQuery := query4.NewCarrierQuery(db2)
	carryingQueryBus := query4.CarrierQueryMessageBus(carrierQuery)
	supplierQuery := query5.NewSupplierQuery(db2)
	supplieringQueryBus := query5.SupplierQueryMessageBus(supplierQuery)
	traderQuery := query6.NewTraderQuery(db2, customeringQueryBus, carryingQueryBus, supplieringQueryBus)
	traderingQueryBus := query6.TraderQueryMessageBus(traderQuery)
	stocktakeQuery := query7.NewQueryStocktake(db2)
	stocktakingQueryBus := query7.StocktakeQueryMessageBus(stocktakeQuery)
	inventoryQueryService := query8.NewQueryInventory(stocktakingQueryBus, eventBus, db2)
	inventoryQueryBus := query8.InventoryQueryServiceMessageBus(inventoryQueryService)
	receiptQuery := query9.NewReceiptQuery(db2)
	receiptingQueryBus := query9.ReceiptQueryMessageBus(receiptQuery)
	purchaseOrderQuery := query10.NewPurchaseOrderQuery(db2, eventBus, supplieringQueryBus, inventoryQueryBus, receiptingQueryBus)
	purchaseorderQueryBus := query10.PurchaseOrderQueryMessageBus(purchaseOrderQuery)
	refundQueryService := query11.NewQueryRefund(eventBus, db2)
	refundQueryBus := query11.RefundQueryServiceMessageBus(refundQueryService)
	purchaseRefundQueryService := query12.NewQueryPurchasePurchaseRefund(eventBus, db2)
	purchaserefundQueryBus := query12.PurchaseRefundQueryServiceMessageBus(purchaseRefundQueryService)
	inventoryAggregate := aggregate5.NewAggregateInventory(eventBus, db2, traderingQueryBus, purchaseorderQueryBus, stocktakingQueryBus, refundQueryBus, purchaserefundQueryBus)
	inventoryCommandBus := aggregate5.InventoryAggregateMessageBus(inventoryAggregate)
	inventoryService := &shop.InventoryService{
		TraderQuery:    traderingQueryBus,
		InventoryAggr:  inventoryCommandBus,
		InventoryQuery: inventoryQueryBus,
	}
	addressQueryService := address.NewQueryService(db2)
	addressQueryBus := address.QueryServiceMessageBus(addressQueryService)
	shopAccountService := &shop.AccountService{
		IdentityAggr:  commandBus,
		IdentityQuery: identityQueryBus,
		AddressQuery:  addressQueryBus,
	}
	collectionService := &shop.CollectionService{
		CatalogQuery: catalogQueryBus,
		CatalogAggr:  catalogCommandBus,
	}
	customerAggregate := aggregate6.NewCustomerAggregate(eventBus, db2)
	customeringCommandBus := aggregate6.CustomerAggregateMessageBus(customerAggregate)
	addressAggregate := aggregate6.NewAddressAggregate(db2)
	addressingCommandBus := aggregate6.AddressAggregateMessageBus(addressAggregate)
	addressQuery := query2.NewAddressQuery(db2)
	addressingQueryBus := query2.AddressQueryMessageBus(addressQuery)
	orderingQueryService := ordering.NewQueryService(db2)
	orderingQueryBus := ordering.QueryServiceMessageBus(orderingQueryService)
	customerService := &shop.CustomerService{
		LocationQuery: queryBus,
		CustomerQuery: customeringQueryBus,
		CustomerAggr:  customeringCommandBus,
		AddressAggr:   addressingCommandBus,
		AddressQuery:  addressingQueryBus,
		OrderQuery:    orderingQueryBus,
		ReceiptQuery:  receiptingQueryBus,
	}
	customerGroupService := &shop.CustomerGroupService{
		CustomerAggr:  customeringCommandBus,
		CustomerQuery: customeringQueryBus,
	}
	productService := &shop.ProductService{
		CatalogQuery:   catalogQueryBus,
		CatalogAggr:    catalogCommandBus,
		InventoryQuery: inventoryQueryBus,
	}
	categoryService := &shop.CategoryService{
		CatalogQuery: catalogQueryBus,
		CatalogAggr:  catalogCommandBus,
	}
	productSourceService := &shop.ProductSourceService{}
	orderingAggregate := ordering.NewAggregate(eventBus, db2)
	orderingCommandBus := ordering.AggregateMessageBus(orderingAggregate)
	v3 := SupportedCarrierDrivers(cfg2, queryBus)
	carrierManager := shipping_provider.NewCtrl(eventBus, queryBus, v3)
	connectionQuery := query13.NewConnectionQuery(db2)
	connectioningQueryBus := query13.ConnectionQueryMessageBus(connectionQuery)
	connectionAggregate := aggregate7.NewConnectionAggregate(db2, eventBus)
	connectioningCommandBus := aggregate7.ConnectionAggregateMessageBus(connectionAggregate)
	shipmentserviceQueryService := shipmentservice.NewQueryService(db2, rd)
	shipmentserviceQueryBus := shipmentservice.QueryServiceMessageBus(shipmentserviceQueryService)
	pricelistQueryService := pricelist.NewQueryService(db2, rd)
	pricelistQueryBus := pricelist.QueryServiceMessageBus(pricelistQueryService)
	shipmentpriceQueryService := shipmentprice.NewQueryService(db2, rd, queryBus, pricelistQueryBus)
	shipmentpriceQueryBus := shipmentprice.QueryServiceMessageBus(shipmentpriceQueryService)
	flagApplyShipmentPrice := cfg2.FlagApplyShipmentPrice
	carrierConfig := SupportedShippingCarrierConfig(cfg2)
	shipmentManager, err := carrier.NewShipmentManager(eventBus, queryBus, connectioningQueryBus, connectioningCommandBus, rd, shipmentserviceQueryBus, shipmentpriceQueryBus, flagApplyShipmentPrice, carrierConfig)
	if err != nil {
		return nil, err
	}
	orderLogic := orderS.New(carrierManager, catalogQueryBus, orderingCommandBus, customeringCommandBus, customeringQueryBus, addressingCommandBus, addressingQueryBus, queryBus, eventBus, shipmentManager)
	orderService := &shop.OrderService{
		OrderAggr:     orderingCommandBus,
		CustomerQuery: customeringQueryBus,
		OrderQuery:    orderingQueryBus,
		ReceiptQuery:  receiptingQueryBus,
		OrderLogic:    orderLogic,
	}
	queryService2 := query14.NewQueryService(db2)
	shippingQueryBus := query14.QueryServiceMessageBus(queryService2)
	fulfillmentService := &shop.FulfillmentService{
		ShippingQuery: shippingQueryBus,
		ShippingCtrl:  carrierManager,
	}
	shipnowAggregate := shipnow.NewAggregate(eventBus, db2, queryBus, identityQueryBus, addressQueryBus, orderingQueryBus, shipnowManager)
	shipnowCommandBus := shipnow.AggregateMessageBus(shipnowAggregate)
	shipnowService := &shop.ShipnowService{
		ShipnowAggr:  shipnowCommandBus,
		ShipnowQuery: shipnowQueryBus,
	}
	historyService := &shop.HistoryService{}
	moneyTransactionService := &shop.MoneyTransactionService{}
	dashboardQuery := query15.NewDashboardQuery(db2, rd, queryBus)
	summaryQueryBus := query15.DashboardQueryMessageBus(dashboardQuery)
	summaryService := &shop.SummaryService{
		SummaryQuery: summaryQueryBus,
	}
	exportService := &shop.ExportService{}
	notificationService := &shop.NotificationService{}
	authorizeService := &shop.AuthorizeService{}
	tradingService := &shop.TradingService{
		EventBus:       eventBus,
		IdentityQuery:  identityQueryBus,
		CatalogQuery:   catalogQueryBus,
		OrderQuery:     orderingQueryBus,
		InventoryQuery: inventoryQueryBus,
		OrderLogic:     orderLogic,
	}
	config2 := cfg2.VTPay
	provider := vtpay.New(config2)
	v4 := SupportedPaymentProvider(provider)
	paymentManager := manager.NewManager(v4, orderingQueryBus)
	managerCommandBus := manager.ManagerMesssageBus(paymentManager)
	paymentService := &shop.PaymentService{
		PaymentAggr: managerCommandBus,
	}
	ledgerQuery := query16.NewLedgerQuery(db2)
	ledgeringQueryBus := query16.LedgerQueryMessageBus(ledgerQuery)
	receiptAggregate := aggregate8.NewReceiptAggregate(db2, eventBus, traderingQueryBus, ledgeringQueryBus, orderingQueryBus, customeringQueryBus, carryingQueryBus, supplieringQueryBus, purchaseorderQueryBus)
	receiptingCommandBus := aggregate8.ReceiptAggregateMessageBus(receiptAggregate)
	receiptService := &shop.ReceiptService{
		CarrierQuery:  carryingQueryBus,
		CustomerQuery: customeringQueryBus,
		LedgerQuery:   ledgeringQueryBus,
		ReceiptAggr:   receiptingCommandBus,
		ReceiptQuery:  receiptingQueryBus,
		SupplierQuery: supplieringQueryBus,
		TraderQuery:   traderingQueryBus,
	}
	supplierAggregate := aggregate9.NewSupplierAggregate(eventBus, db2)
	supplieringCommandBus := aggregate9.SupplierAggregateMessageBus(supplierAggregate)
	supplierService := &shop.SupplierService{
		CatalogQuery:       catalogQueryBus,
		PurchaseOrderQuery: purchaseorderQueryBus,
		ReceiptQuery:       receiptingQueryBus,
		SupplierAggr:       supplieringCommandBus,
		SupplierQuery:      supplieringQueryBus,
	}
	carrierAggregate := aggregate10.NewCarrierAggregate(eventBus, db2)
	carryingCommandBus := aggregate10.CarrierAggregateMessageBus(carrierAggregate)
	carrierService := &shop.CarrierService{
		CarrierAggr:  carryingCommandBus,
		CarrierQuery: carryingQueryBus,
	}
	ledgerAggregate := aggregate11.NewLedgerAggregate(db2, receiptingQueryBus)
	ledgeringCommandBus := aggregate11.LedgerAggregateMessageBus(ledgerAggregate)
	ledgerService := &shop.LedgerService{
		LedgerAggr:  ledgeringCommandBus,
		LedgerQuery: ledgeringQueryBus,
	}
	purchaseOrderAggregate := aggregate12.NewPurchaseOrderAggregate(db2, eventBus, catalogQueryBus, supplieringQueryBus, inventoryQueryBus)
	purchaseorderCommandBus := aggregate12.PurchaseOrderAggregateMessageBus(purchaseOrderAggregate)
	purchaseOrderService := &shop.PurchaseOrderService{
		PurchaseOrderAggr:  purchaseorderCommandBus,
		PurchaseOrderQuery: purchaseorderQueryBus,
	}
	stocktakeAggregate := aggregate13.NewAggregateStocktake(db2, eventBus)
	stocktakingCommandBus := aggregate13.StocktakeAggregateMessageBus(stocktakeAggregate)
	stocktakeService := &shop.StocktakeService{
		CatalogQuery:   catalogQueryBus,
		StocktakeAggr:  stocktakingCommandBus,
		StocktakeQuery: stocktakingQueryBus,
		InventoryQuery: inventoryQueryBus,
	}
	aggregate22 := aggregate14.NewAggregate(db2, eventBus, queryBus, bot2, orderingQueryBus, shipmentManager, connectioningQueryBus)
	shippingCommandBus := aggregate14.AggregateMessageBus(aggregate22)
	shipmentService := &shop.ShipmentService{
		ShipmentManager:   shipmentManager,
		ShippingAggregate: shippingCommandBus,
	}
	connectionService := &shop.ConnectionService{
		ShipmentManager: shipmentManager,
		ConnectionQuery: connectioningQueryBus,
		ConnectionAggr:  connectioningCommandBus,
	}
	refundAggregate := aggregate15.NewRefundAggregate(db2, eventBus)
	refundCommandBus := aggregate15.RefundAggregateMessageBus(refundAggregate)
	refundService := &shop.RefundService{
		CustomerQuery:  customeringQueryBus,
		InventoryQuery: inventoryQueryBus,
		ReceiptQuery:   receiptingQueryBus,
		RefundAggr:     refundCommandBus,
		RefundQuery:    refundQueryBus,
	}
	purchaseRefundAggregate := aggregate16.NewPurchaseRefundAggregate(db2, eventBus, purchaseorderQueryBus)
	purchaserefundCommandBus := aggregate16.PurchaseRefundAggregateMessageBus(purchaseRefundAggregate)
	purchaseRefundService := &shop.PurchaseRefundService{
		PurchaseRefundAggr:  purchaserefundCommandBus,
		PurchaseRefundQuery: purchaserefundQueryBus,
		SupplierQuery:       supplieringQueryBus,
		PurchaseOrderQuery:  purchaseorderQueryBus,
		InventoryQuery:      inventoryQueryBus,
	}
	webserverAggregate := aggregate17.New(eventBus, db2, catalogQueryBus)
	webserverCommandBus := aggregate17.WebserverAggregateMessageBus(webserverAggregate)
	webserverQueryService := query17.New(eventBus, db2, catalogQueryBus)
	webserverQueryBus := query17.WebserverQueryServiceMessageBus(webserverQueryService)
	webServerService := &shop.WebServerService{
		CatalogQuery:   catalogQueryBus,
		WebserverAggr:  webserverCommandBus,
		WebserverQuery: webserverQueryBus,
		InventoryQuery: inventoryQueryBus,
	}
	subrProductQuery := subscriptionproduct.NewSubrProductQuery(db2)
	subscriptionproductQueryBus := subscriptionproduct.SubrProductQueryMessageBus(subrProductQuery)
	subrPlanQuery := subscriptionplan.NewSubrPlanQuery(db2, subscriptionproductQueryBus)
	subscriptionplanQueryBus := subscriptionplan.SubrPlanQueryMessageBus(subrPlanQuery)
	subscriptionQuery := subscription.NewSubscriptionQuery(db2, subscriptionplanQueryBus, subscriptionproductQueryBus)
	subscriptionQueryBus := subscription.SubscriptionQueryMessageBus(subscriptionQuery)
	subscriptionService := &shop.SubscriptionService{
		SubscriptionQuery: subscriptionQueryBus,
	}
	shopServers := shop.NewServers(rd, shopMiscService, brandService, inventoryService, shopAccountService, collectionService, customerService, customerGroupService, productService, categoryService, productSourceService, orderService, fulfillmentService, shipnowService, historyService, moneyTransactionService, summaryService, exportService, notificationService, authorizeService, tradingService, paymentService, receiptService, supplierService, carrierService, ledgerService, purchaseOrderService, stocktakeService, shipmentService, connectionService, refundService, purchaseRefundService, webServerService, subscriptionService)
	adminMiscService := &admin.MiscService{}
	adminAccountService := &admin.AccountService{}
	adminOrderService := &admin.OrderService{}
	adminFulfillmentService := &admin.FulfillmentService{
		RedisStore: rd,
	}
	moneyTxQuery := query18.NewMoneyTxQuery(db2, shippingQueryBus)
	moneytxQueryBus := query18.MoneyTxQueryMessageBus(moneyTxQuery)
	moneyTxAggregate := aggregate18.NewMoneyTxAggregate(db2, shippingQueryBus, identityQueryBus, eventBus)
	moneytxCommandBus := aggregate18.MoneyTxAggregateMessageBus(moneyTxAggregate)
	adminMoneyTransactionService := &admin.MoneyTransactionService{
		MoneyTxQuery: moneytxQueryBus,
		MoneyTxAggr:  moneytxCommandBus,
	}
	shopService := &admin.ShopService{
		IdentityQuery: identityQueryBus,
	}
	creditService := &admin.CreditService{}
	adminNotificationService := &admin.NotificationService{}
	adminConnectionService := &admin.ConnectionService{
		ConnectionAggr:  connectioningCommandBus,
		ConnectionQuery: connectioningQueryBus,
	}
	shipmentpriceAggregate := shipmentprice.NewAggregate(db2, rd)
	shipmentpriceCommandBus := shipmentprice.AggregateMessageBus(shipmentpriceAggregate)
	shipmentserviceAggregate := shipmentservice.NewAggregate(db2, rd)
	shipmentserviceCommandBus := shipmentservice.AggregateMessageBus(shipmentserviceAggregate)
	pricelistAggregate := pricelist.NewAggregate(db2, eventBus)
	pricelistCommandBus := pricelist.AggregateMessageBus(pricelistAggregate)
	shipmentPriceService := &admin.ShipmentPriceService{
		ShipmentManager:        shipmentManager,
		ShipmentPriceAggr:      shipmentpriceCommandBus,
		ShipmentPriceQuery:     shipmentpriceQueryBus,
		ShipmentServiceQuery:   shipmentserviceQueryBus,
		ShipmentServiceAggr:    shipmentserviceCommandBus,
		ShipmentPriceListAggr:  pricelistCommandBus,
		ShipmentPriceListQuery: pricelistQueryBus,
	}
	locationAggregate := location.NewAggregate(db2)
	locationCommandBus := location.AggregateMessageBus(locationAggregate)
	adminLocationService := &admin.LocationService{
		LocationAggr:  locationCommandBus,
		LocationQuery: queryBus,
	}
	subrProductAggregate := subscriptionproduct.NewSubrProductAggregate(db2)
	subscriptionproductCommandBus := subscriptionproduct.SubrProductAggregateMessageBus(subrProductAggregate)
	subrPlanAggregate := subscriptionplan.NewSubrPlanAggregate(db2)
	subscriptionplanCommandBus := subscriptionplan.SubrPlanAggregateMessageBus(subrPlanAggregate)
	subscriptionAggregate := subscription.NewSubscriptionAggregate(db2)
	subscriptionCommandBus := subscription.SubscriptionAggregateMessageBus(subscriptionAggregate)
	aggregate23 := aggregate19.NewAggregate(db2)
	paymentCommandBus := aggregate19.AggregateMessageBus(aggregate23)
	subrBillAggregate := subscriptionbill.NewSubrBillAggregate(db2, eventBus, paymentCommandBus, subscriptionQueryBus, subscriptionplanQueryBus)
	subscriptionbillCommandBus := subscriptionbill.SubrBillAggregateMessageBus(subrBillAggregate)
	subrBillQuery := subscriptionbill.NewSubrBillQuery(db2)
	subscriptionbillQueryBus := subscriptionbill.SubrBillQueryMessageBus(subrBillQuery)
	adminSubscriptionService := &admin.SubscriptionService{
		SubrProductAggr:   subscriptionproductCommandBus,
		SubrProductQuery:  subscriptionproductQueryBus,
		SubrPlanAggr:      subscriptionplanCommandBus,
		SubrPlanQuery:     subscriptionplanQueryBus,
		SubscriptionQuery: subscriptionQueryBus,
		SubscriptionAggr:  subscriptionCommandBus,
		SubrBillAggr:      subscriptionbillCommandBus,
		SubrBillQuery:     subscriptionbillQueryBus,
	}
	adminServers := admin.NewServers(adminMiscService, adminAccountService, adminOrderService, adminFulfillmentService, adminMoneyTransactionService, shopService, creditService, adminNotificationService, adminConnectionService, shipmentPriceService, adminLocationService, adminSubscriptionService)
	sadminMiscService := &sadmin.MiscService{
		Session: ss,
	}
	sadminUserService := &sadmin.UserService{
		Session: ss,
	}
	sadminServers := sadmin.NewServers(sadminMiscService, sadminUserService)
	integrationMiscService := &integration.MiscService{}
	integrationService := &integration.IntegrationService{
		AuthStore: s,
		SMSClient: smsClient,
	}
	integrationServers := integration.NewIntegrationServer(sd, rd, integrationMiscService, integrationService)
	affiliateMiscService := affiliate.MiscService{}
	affiliateAccountService := affiliate.AccountService{
		IdentityAggr: commandBus,
	}
	affiliateServers := affiliate.NewServers(affiliateMiscService, affiliateAccountService)
	secretToken := cfg2.Secret
	affiliateAggregate := affiliate2.NewAggregate(db2, identityQueryBus, catalogQueryBus, orderingQueryBus)
	affiliateCommandBus := affiliate2.AggregateMessageBus(affiliateAggregate)
	apiUserService := &api2.UserService{
		AffiliateAggr: affiliateCommandBus,
	}
	affiliateQueryService := affiliate2.NewQuery(db2)
	affiliateQueryBus := affiliate2.QueryServiceMessageBus(affiliateQueryService)
	apiTradingService := &api2.TradingService{
		AffiliateAggr:  affiliateCommandBus,
		AffiliateQuery: affiliateQueryBus,
		CatalogQuery:   catalogQueryBus,
		InventoryQuery: inventoryQueryBus,
	}
	apiShopService := &api2.ShopService{
		CatalogQuery:   catalogQueryBus,
		InventoryQuery: inventoryQueryBus,
		AffiliateQuery: affiliateQueryBus,
	}
	affiliateService := &api2.AffiliateService{
		AffiliateAggr:  affiliateCommandBus,
		CatalogQuery:   catalogQueryBus,
		AffiliateQuery: affiliateQueryBus,
		IdentityQuery:  identityQueryBus,
	}
	apiServers := api2.NewServers(secretToken, apiUserService, apiTradingService, apiShopService, affiliateService)
	shippingShipping := shipping.New(queryBus, db2, shipmentManager, shippingCommandBus, shippingQueryBus, orderLogic)
	partnerMiscService := &partner.MiscService{
		Shipping: shippingShipping,
	}
	partnerShopService := &partner.ShopService{}
	webhookService := &partner.WebhookService{}
	partnerHistoryService := &partner.HistoryService{}
	shippingService := &partner.ShippingService{
		Shipping: shippingShipping,
	}
	partnerOrderService := &partner.OrderService{
		Shipping: shippingShipping,
	}
	partnerFulfillmentService := &partner.FulfillmentService{
		Shipping: shippingShipping,
	}
	shoppingShopping := &shopping.Shopping{
		LocationQuery:     queryBus,
		CustomerQuery:     customeringQueryBus,
		CustomerAggregate: customeringCommandBus,
		AddressQuery:      addressingQueryBus,
		AddressAggregate:  addressingCommandBus,
		InventoryQuery:    inventoryQueryBus,
		CatalogQuery:      catalogQueryBus,
		CatalogAggregate:  catalogCommandBus,
	}
	partnerCustomerService := &partner.CustomerService{
		Shopping: shoppingShopping,
	}
	customerAddressService := &partner.CustomerAddressService{
		Shopping: shoppingShopping,
	}
	partnerCustomerGroupService := &partner.CustomerGroupService{
		Shopping: shoppingShopping,
	}
	customerGroupRelationshipService := &partner.CustomerGroupRelationshipService{
		Shopping: shoppingShopping,
	}
	partnerInventoryService := &partner.InventoryService{
		Shopping: shoppingShopping,
	}
	variantService := &partner.VariantService{
		Shopping: shoppingShopping,
	}
	partnerProductService := &partner.ProductService{
		Shopping: shoppingShopping,
	}
	productCollectionService := &partner.ProductCollectionService{
		Shopping: shoppingShopping,
	}
	productCollectionRelationshipService := &partner.ProductCollectionRelationshipService{
		Shopping: shoppingShopping,
	}
	partnerServers := partner.NewServers(sd, rd, s, authURL, partnerMiscService, partnerShopService, webhookService, partnerHistoryService, shippingService, partnerOrderService, partnerFulfillmentService, partnerCustomerService, customerAddressService, partnerCustomerGroupService, customerGroupRelationshipService, partnerInventoryService, variantService, partnerProductService, productCollectionService, productCollectionRelationshipService)
	xshopMiscService := &xshop.MiscService{
		Shipping: shippingShipping,
	}
	xshopWebhookService := &xshop.WebhookService{}
	xshopHistoryService := &xshop.HistoryService{}
	xshopShippingService := &xshop.ShippingService{
		Shipping: shippingShipping,
	}
	xshopOrderService := &xshop.OrderService{
		Shipping: shippingShipping,
	}
	xshopFulfillmentService := &xshop.FulfillmentService{
		Shipping: shippingShipping,
	}
	xshopCustomerService := &xshop.CustomerService{
		Shopping: shoppingShopping,
	}
	xshopCustomerAddressService := &xshop.CustomerAddressService{
		Shopping: shoppingShopping,
	}
	xshopCustomerGroupService := &xshop.CustomerGroupService{
		Shopping: shoppingShopping,
	}
	xshopCustomerGroupRelationshipService := &xshop.CustomerGroupRelationshipService{
		Shopping: shoppingShopping,
	}
	xshopInventoryService := &xshop.InventoryService{
		Shopping: shoppingShopping,
	}
	xshopVariantService := &xshop.VariantService{
		Shopping: shoppingShopping,
	}
	xshopProductService := &xshop.ProductService{
		Shopping: shoppingShopping,
	}
	xshopProductCollectionService := &xshop.ProductCollectionService{
		Shopping: shoppingShopping,
	}
	xshopProductCollectionRelationshipService := &xshop.ProductCollectionRelationshipService{
		Shopping: shoppingShopping,
	}
	xshopServers := xshop.NewServers(sd, rd, shippingShipping, xshopMiscService, xshopWebhookService, xshopHistoryService, xshopShippingService, xshopOrderService, xshopFulfillmentService, xshopCustomerService, xshopCustomerAddressService, xshopCustomerGroupService, xshopCustomerGroupRelationshipService, xshopInventoryService, xshopVariantService, xshopProductService, xshopProductCollectionService, xshopProductCollectionRelationshipService)
	partnercarrierMiscService := &partnercarrier.MiscService{
		Session:  ss,
		Shipping: shippingShipping,
	}
	shipmentConnectionService := &partnercarrier.ShipmentConnectionService{
		Session:         ss,
		ConnectionQuery: connectioningQueryBus,
		ConnectionAggr:  connectioningCommandBus,
	}
	partnercarrierShipmentService := &partnercarrier.ShipmentService{
		Session:         ss,
		ConnectionQuery: connectioningQueryBus,
		ShippingAggr:    shippingCommandBus,
		ShippingQuery:   shippingQueryBus,
	}
	partnercarrierServers := partnercarrier.NewServers(sd, rd, partnercarrierMiscService, shipmentConnectionService, partnercarrierShipmentService)
	importService := partnerimport.New(db2, catalogCommandBus)
	partnerimportServers := partnerimport.NewServers(importService)
	etopHandlers := NewEtopHandlers(servers, shopServers, adminServers, sadminServers, integrationServers, affiliateServers, apiServers, partnerServers, xshopServers, partnercarrierServers, partnerimportServers)
	eventstreamEventStreamer := eventstream.NewEventStreamer(sd)
	client3 := client2.New(config2)
	vtpayAggregate := vtpay2.NewAggregate(db2, orderingQueryBus, orderingCommandBus, paymentCommandBus, client3)
	vtpayCommandBus := vtpay2.AggregateMessageBus(vtpayAggregate)
	aggregate24 := aggregate20.NewAggregate(orderingQueryBus, orderingCommandBus, vtpayCommandBus, client3)
	gatewayCommandBus := aggregate20.AggregateMessageBus(aggregate24)
	aggregate25 := aggregate21.New(db2)
	serverServer := server.New(gatewayCommandBus, aggregate25)
	imcsvImport := imcsv.Import{
		MoneyTxAggr: moneytxCommandBus,
	}
	ghtkimcsvImport := ghtkimcsv.Import{
		MoneyTxAggr: moneytxCommandBus,
	}
	vtpostimxlsxImport := vtpostimxlsx.Import{
		MoneyTxAggr: moneytxCommandBus,
	}
	uploader, err := NewUploader(cfg2)
	if err != nil {
		return nil, err
	}
	import2 := imcsv2.New(queryBus, sd, rd, uploader, db2)
	import3 := imcsv3.New(sd, rd, uploader, db2)
	shopImport := imports.New(import2, import3)
	etopServer := NewEtopServer(etopHandlers, eventstreamEventStreamer, serverServer, imcsvImport, ghtkimcsvImport, vtpostimxlsxImport, shopImport)
	webServer := NewWebServer(webserverQueryBus, catalogQueryBus, subscriptionQueryBus, rd, queryBus)
	ghnConfig := cfg2.GHN
	ghnCarrier := ghn.New(ghnConfig, queryBus)
	ghnWebhookServer := NewGHNWebhookServer(shipmentManager, ghnCarrier, identityQueryBus, shippingCommandBus)
	ghtkConfig := cfg2.GHTK
	ghtkCarrier := ghtk.New(ghtkConfig, queryBus)
	ghtkWebhookServer := NewGHTKWebhookServer(shipmentManager, ghtkCarrier, identityQueryBus, shippingCommandBus)
	vtpostConfig := cfg2.VTPost
	vtpostCarrier := vtpost.New(vtpostConfig, queryBus)
	vtPostWebhookServer := NewVTPostWebhookServer(shipmentManager, vtpostCarrier, identityQueryBus, shippingCommandBus)
	ahamoveVerificationFileServer := NewAhamoveVerificationFileServer()
	ahamoveWebhookServer := NewAhamoveWebhookServer(shipmentManager, ahamoveCarrier, identityQueryBus, shipnowQueryBus, shipnowCommandBus, orderingCommandBus, orderingQueryBus, ahamoveVerificationFileServer)
	v5 := NewServers(etopServer, webServer, ghnWebhookServer, ghtkWebhookServer, vtPostWebhookServer, ahamoveWebhookServer)
	return v5, nil
}
