package shop

import (
	service "etop.vn/api/top/int/shop"
	"etop.vn/capi/httprpc"
)

// +gen:wrapper=etop.vn/api/top/int/shop
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
		service.NewMiscServiceServer(WrapMiscService(miscService)),
		service.NewBrandServiceServer(WrapBrandService(brandService)),
		service.NewInventoryServiceServer(WrapInventoryService(inventoryService)),
		service.NewAccountServiceServer(WrapAccountService(accountService)),
		service.NewExternalAccountServiceServer(WrapExternalAccountService(externalAccountService)),
		service.NewCollectionServiceServer(WrapCollectionService(collectionService)),
		service.NewCustomerServiceServer(WrapCustomerService(customerService)),
		service.NewCustomerGroupServiceServer(WrapCustomerGroupService(customerGroupService)),
		service.NewProductServiceServer(WrapProductService(ProductServiceImpl)),
		service.NewCategoryServiceServer(WrapCategoryService(categoryService)),
		service.NewProductSourceServiceServer(WrapProductSourceService(productSourceService)),
		service.NewOrderServiceServer(WrapOrderService(orderService)),
		service.NewFulfillmentServiceServer(WrapFulfillmentService(fulfillmentService)),
		service.NewShipnowServiceServer(WrapShipnowService(shipnowService)),
		service.NewHistoryServiceServer(WrapHistoryService(historyService)),
		service.NewMoneyTransactionServiceServer(WrapMoneyTransactionService(moneyTransactionService)),
		service.NewSummaryServiceServer(WrapSummaryService(summaryService)),
		service.NewExportServiceServer(WrapExportService(exportService)),
		service.NewNotificationServiceServer(WrapNotificationService(notificationService)),
		service.NewAuthorizeServiceServer(WrapAuthorizeService(authorizeService)),
		service.NewTradingServiceServer(WrapTradingService(tradingService)),
		service.NewPaymentServiceServer(WrapPaymentService(paymentService)),
		service.NewReceiptServiceServer(WrapReceiptService(receiptService)),
		service.NewSupplierServiceServer(WrapSupplierService(supplierService)),
		service.NewCarrierServiceServer(WrapCarrierService(carrierService)),
		service.NewLedgerServiceServer(WrapLedgerService(ledgerService)),
		service.NewPurchaseOrderServiceServer(WrapPurchaseOrderService(purchaseOrderService)),
		service.NewStocktakeServiceServer(WrapStocktakeService(stocktakeService)),
		service.NewInvitationServiceServer(WrapInvitationService(invitationService)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
