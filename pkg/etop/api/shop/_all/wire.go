package shop_all

import (
	"github.com/google/wire"

	service "o.o/api/top/int/shop"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/api/shop"
	"o.o/capi/httprpc"
)

var WireSet = wire.NewSet(
	NewServers,
)

func NewServers(
	rd redis.Store,
	miscService *shop.MiscService,
	brandService *shop.BrandService,
	inventoryService *shop.InventoryService,
	accountService *shop.AccountService,
	collectionService *shop.CollectionService,
	customerService *shop.CustomerService,
	customerGroupService *shop.CustomerGroupService,
	productService *shop.ProductService,
	categoryService *shop.CategoryService,
	productSourceService *shop.ProductSourceService,
	orderService *shop.OrderService,
	fulfillmentService *shop.FulfillmentService,
	shipnowService *shop.ShipnowService,
	historyService *shop.HistoryService,
	moneyTransactionService *shop.MoneyTransactionService,
	summaryService *shop.SummaryService,
	exportService *shop.ExportService,
	notificationService *shop.NotificationService,
	authorizeService *shop.AuthorizeService,
	tradingService *shop.TradingService,
	paymentService *shop.PaymentService,
	receiptService *shop.ReceiptService,
	supplierService *shop.SupplierService,
	carrierService *shop.CarrierService,
	ledgerService *shop.LedgerService,
	purchaseOrderService *shop.PurchaseOrderService,
	stocktakeService *shop.StocktakeService,
	shipmentService *shop.ShipmentService,
	connectionService *shop.ConnectionService,
	refundService *shop.RefundService,
	purchaseRefundService *shop.PurchaseRefundService,
	webServerService *shop.WebServerService,
	subscriptionService *shop.SubscriptionService,
) shop.Servers {

	shop.InitIdemp(rd)
	shop.ProductServiceImpl = productService
	shop.StocktakeServiceImpl = stocktakeService
	shop.InventoryServiceImpl = inventoryService

	servers := []httprpc.Server{
		service.NewAccountServiceServer(shop.WrapAccountService(accountService.Clone)),
		service.NewAuthorizeServiceServer(shop.WrapAuthorizeService(authorizeService.Clone)),
		service.NewBrandServiceServer(shop.WrapBrandService(brandService.Clone)),
		service.NewCarrierServiceServer(shop.WrapCarrierService(carrierService.Clone)),
		service.NewCategoryServiceServer(shop.WrapCategoryService(categoryService.Clone)),
		service.NewCollectionServiceServer(shop.WrapCollectionService(collectionService.Clone)),
		service.NewConnectionServiceServer(shop.WrapConnectionService(connectionService.Clone)),
		service.NewCustomerGroupServiceServer(shop.WrapCustomerGroupService(customerGroupService.Clone)),
		service.NewCustomerServiceServer(shop.WrapCustomerService(customerService.Clone)),
		service.NewExportServiceServer(shop.WrapExportService(exportService.Clone)),
		service.NewFulfillmentServiceServer(shop.WrapFulfillmentService(fulfillmentService.Clone)),
		service.NewHistoryServiceServer(shop.WrapHistoryService(historyService.Clone)),
		service.NewInventoryServiceServer(shop.WrapInventoryService(inventoryService.Clone)),
		service.NewLedgerServiceServer(shop.WrapLedgerService(ledgerService.Clone)),
		service.NewMiscServiceServer(shop.WrapMiscService(miscService.Clone)),
		service.NewMoneyTransactionServiceServer(shop.WrapMoneyTransactionService(moneyTransactionService.Clone)),
		service.NewNotificationServiceServer(shop.WrapNotificationService(notificationService.Clone)),
		service.NewOrderServiceServer(shop.WrapOrderService(orderService.Clone)),
		service.NewPaymentServiceServer(shop.WrapPaymentService(paymentService.Clone)),
		service.NewProductServiceServer(shop.WrapProductService(productService.Clone)),
		service.NewProductSourceServiceServer(shop.WrapProductSourceService(productSourceService.Clone)),
		service.NewPurchaseOrderServiceServer(shop.WrapPurchaseOrderService(purchaseOrderService.Clone)),
		service.NewPurchaseRefundServiceServer(shop.WrapPurchaseRefundService(purchaseRefundService.Clone)),
		service.NewReceiptServiceServer(shop.WrapReceiptService(receiptService.Clone)),
		service.NewRefundServiceServer(shop.WrapRefundService(refundService.Clone)),
		service.NewShipmentServiceServer(shop.WrapShipmentService(shipmentService.Clone)),
		service.NewShipnowServiceServer(shop.WrapShipnowService(shipnowService.Clone)),
		service.NewStocktakeServiceServer(shop.WrapStocktakeService(stocktakeService.Clone)),
		service.NewSubscriptionServiceServer(shop.WrapSubscriptionService(subscriptionService.Clone)),
		service.NewSummaryServiceServer(shop.WrapSummaryService(summaryService.Clone)),
		service.NewSupplierServiceServer(shop.WrapSupplierService(supplierService.Clone)),
		service.NewTradingServiceServer(shop.WrapTradingService(tradingService.Clone)),
		service.NewWebServerServiceServer(shop.WrapWebServerService(webServerService.Clone)),
	}
	return servers
}
