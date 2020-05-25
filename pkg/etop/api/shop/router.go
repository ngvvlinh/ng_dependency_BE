package shop

import (
	service "o.o/api/top/int/shop"
	"o.o/backend/pkg/common/apifw/idemp"
	"o.o/backend/pkg/common/redis"
	"o.o/capi/httprpc"
	"o.o/common/l"
)

var ll = l.New()
var idempgroup *idemp.RedisGroup

// +gen:wrapper=o.o/api/top/int/shop
// +gen:wrapper:package=shop

type Servers []httprpc.Server

// hack for imcsv
// TODO: remove this
var ProductServiceImpl *ProductService
var StocktakeServiceImpl *StocktakeService
var InventoryServiceImpl *InventoryService

func NewServers(
	rd redis.Store,
	miscService *MiscService,
	brandService *BrandService,
	inventoryService *InventoryService,
	accountService *AccountService,
	collectionService *CollectionService,
	customerService *CustomerService,
	customerGroupService *CustomerGroupService,
	productService *ProductService,
	categoryService *CategoryService,
	productSourceService *ProductSourceService,
	orderService *OrderService,
	fulfillmentService *FulfillmentService,
	shipnowService *ShipnowService,
	historyService *HistoryService,
	moneyTransactionService *MoneyTransactionService,
	summaryService *SummaryService,
	exportService *ExportService,
	notificationService *NotificationService,
	authorizeService *AuthorizeService,
	tradingService *TradingService,
	paymentService *PaymentService,
	receiptService *ReceiptService,
	supplierService *SupplierService,
	carrierService *CarrierService,
	ledgerService *LedgerService,
	purchaseOrderService *PurchaseOrderService,
	stocktakeService *StocktakeService,
	shipmentService *ShipmentService,
	connectionService *ConnectionService,
	refundService *RefundService,
	purchaseRefundService *PurchaseRefundService,
	webServerService *WebServerService,
	subscriptionService *SubscriptionService,
) Servers {

	idempgroup = idemp.NewRedisGroup(rd, "idemp_shop", 30)
	ProductServiceImpl = productService
	StocktakeServiceImpl = stocktakeService

	servers := []httprpc.Server{
		service.NewAccountServiceServer(WrapAccountService(accountService.Clone)),
		service.NewAuthorizeServiceServer(WrapAuthorizeService(authorizeService.Clone)),
		service.NewBrandServiceServer(WrapBrandService(brandService.Clone)),
		service.NewCarrierServiceServer(WrapCarrierService(carrierService.Clone)),
		service.NewCategoryServiceServer(WrapCategoryService(categoryService.Clone)),
		service.NewCollectionServiceServer(WrapCollectionService(collectionService.Clone)),
		service.NewConnectionServiceServer(WrapConnectionService(connectionService.Clone)),
		service.NewCustomerGroupServiceServer(WrapCustomerGroupService(customerGroupService.Clone)),
		service.NewCustomerServiceServer(WrapCustomerService(customerService.Clone)),
		service.NewExportServiceServer(WrapExportService(exportService.Clone)),
		service.NewFulfillmentServiceServer(WrapFulfillmentService(fulfillmentService.Clone)),
		service.NewHistoryServiceServer(WrapHistoryService(historyService.Clone)),
		service.NewInventoryServiceServer(WrapInventoryService(inventoryService.Clone)),
		service.NewLedgerServiceServer(WrapLedgerService(ledgerService.Clone)),
		service.NewMiscServiceServer(WrapMiscService(miscService.Clone)),
		service.NewMoneyTransactionServiceServer(WrapMoneyTransactionService(moneyTransactionService.Clone)),
		service.NewNotificationServiceServer(WrapNotificationService(notificationService.Clone)),
		service.NewOrderServiceServer(WrapOrderService(orderService.Clone)),
		service.NewPaymentServiceServer(WrapPaymentService(paymentService.Clone)),
		service.NewProductServiceServer(WrapProductService(productService.Clone)),
		service.NewProductSourceServiceServer(WrapProductSourceService(productSourceService.Clone)),
		service.NewPurchaseOrderServiceServer(WrapPurchaseOrderService(purchaseOrderService.Clone)),
		service.NewPurchaseRefundServiceServer(WrapPurchaseRefundService(purchaseRefundService.Clone)),
		service.NewReceiptServiceServer(WrapReceiptService(receiptService.Clone)),
		service.NewRefundServiceServer(WrapRefundService(refundService.Clone)),
		service.NewShipmentServiceServer(WrapShipmentService(shipmentService.Clone)),
		service.NewShipnowServiceServer(WrapShipnowService(shipnowService.Clone)),
		service.NewStocktakeServiceServer(WrapStocktakeService(stocktakeService.Clone)),
		service.NewSubscriptionServiceServer(WrapSubscriptionService(subscriptionService.Clone)),
		service.NewSummaryServiceServer(WrapSummaryService(summaryService.Clone)),
		service.NewSupplierServiceServer(WrapSupplierService(supplierService.Clone)),
		service.NewTradingServiceServer(WrapTradingService(tradingService.Clone)),
		service.NewWebServerServiceServer(WrapWebServerService(webServerService.Clone)),
	}
	return servers
}
