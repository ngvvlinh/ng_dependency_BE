package shop

import (
	service "o.o/api/top/int/shop"
	"o.o/capi/httprpc"
)

// +gen:wrapper=o.o/api/top/int/shop
// +gen:wrapper:package=shop

// hack for imcsv
// TODO: remove this
var ProductServiceImpl *ProductService
var StocktakeServiceImpl *StocktakeService
var InventoryServiceImpl *InventoryService

func NewShopServer(m httprpc.Muxer) {
	ProductServiceImpl = productService
	StocktakeServiceImpl = stocktakeService
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService.Clone)),
		service.NewBrandServiceServer(WrapBrandService(brandService.Clone)),
		service.NewInventoryServiceServer(WrapInventoryService(inventoryService.Clone)),
		service.NewAccountServiceServer(WrapAccountService(accountService.Clone)),
		service.NewCollectionServiceServer(WrapCollectionService(collectionService.Clone)),
		service.NewCustomerServiceServer(WrapCustomerService(customerService.Clone)),
		service.NewCustomerGroupServiceServer(WrapCustomerGroupService(customerGroupService.Clone)),
		service.NewProductServiceServer(WrapProductService(ProductServiceImpl.Clone)),
		service.NewCategoryServiceServer(WrapCategoryService(categoryService.Clone)),
		service.NewProductSourceServiceServer(WrapProductSourceService(productSourceService.Clone)),
		service.NewOrderServiceServer(WrapOrderService(orderService.Clone)),
		service.NewFulfillmentServiceServer(WrapFulfillmentService(fulfillmentService.Clone)),
		service.NewShipnowServiceServer(WrapShipnowService(shipnowService.Clone)),
		service.NewHistoryServiceServer(WrapHistoryService(historyService.Clone)),
		service.NewMoneyTransactionServiceServer(WrapMoneyTransactionService(moneyTransactionService.Clone)),
		service.NewSummaryServiceServer(WrapSummaryService(summaryService.Clone)),
		service.NewExportServiceServer(WrapExportService(exportService.Clone)),
		service.NewNotificationServiceServer(WrapNotificationService(notificationService.Clone)),
		service.NewAuthorizeServiceServer(WrapAuthorizeService(authorizeService.Clone)),
		service.NewTradingServiceServer(WrapTradingService(tradingService.Clone)),
		service.NewPaymentServiceServer(WrapPaymentService(paymentService.Clone)),
		service.NewReceiptServiceServer(WrapReceiptService(receiptService.Clone)),
		service.NewSupplierServiceServer(WrapSupplierService(supplierService.Clone)),
		service.NewCarrierServiceServer(WrapCarrierService(carrierService.Clone)),
		service.NewLedgerServiceServer(WrapLedgerService(ledgerService.Clone)),
		service.NewPurchaseOrderServiceServer(WrapPurchaseOrderService(purchaseOrderService.Clone)),
		service.NewStocktakeServiceServer(WrapStocktakeService(stocktakeService.Clone)),
		service.NewShipmentServiceServer(WrapShipmentService(shipmentService.Clone)),
		service.NewConnectionServiceServer(WrapConnectionService(connectionService.Clone)),
		service.NewRefundServiceServer(WrapRefundService(refundService.Clone)),
		service.NewPurchaseRefundServiceServer(WrapPurchaseRefundService(purchaseRefundService.Clone)),
		service.NewWebServerServiceServer(WrapWebServerService(webServerService.Clone)),
		service.NewSubscriptionServiceServer(WrapSubscriptionService(subscriptionService.Clone)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
