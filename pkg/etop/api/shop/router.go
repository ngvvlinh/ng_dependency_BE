package shop

import (
	"etop.vn/backend/pkg/common/httprpc"
	service "etop.vn/backend/zexp/api/root/int/shop"
)

// +gen:wrapper=etop.vn/backend/pb/etop/shop
// +gen:wrapper:package=shop

// hack for imcsv
// TODO: remove this
var ProductServiceImpl *ProductService

func NewShopServer(m httprpc.Muxer) {
	ProductServiceImpl = productService
	servers := []httprpc.Server{
		service.NewMiscServiceServer(NewMiscService(miscService)),
		service.NewBrandServiceServer(NewBrandService(brandService)),
		service.NewInventoryServiceServer(NewInventoryService(inventoryService)),
		service.NewAccountServiceServer(NewAccountService(accountService)),
		service.NewExternalAccountServiceServer(NewExternalAccountService(externalAccountService)),
		service.NewCollectionServiceServer(NewCollectionService(collectionService)),
		service.NewCustomerServiceServer(NewCustomerService(customerService)),
		service.NewCustomerGroupServiceServer(NewCustomerGroupService(customerGroupService)),
		service.NewProductServiceServer(NewProductService(ProductServiceImpl)),
		service.NewCategoryServiceServer(NewCategoryService(categoryService)),
		service.NewProductSourceServiceServer(NewProductSourceService(productSourceService)),
		service.NewOrderServiceServer(NewOrderService(orderService)),
		service.NewFulfillmentServiceServer(NewFulfillmentService(fulfillmentService)),
		service.NewShipnowServiceServer(NewShipnowService(shipnowService)),
		service.NewHistoryServiceServer(NewHistoryService(historyService)),
		service.NewMoneyTransactionServiceServer(NewMoneyTransactionService(moneyTransactionService)),
		service.NewSummaryServiceServer(NewSummaryService(summaryService)),
		service.NewExportServiceServer(NewExportService(exportService)),
		service.NewNotificationServiceServer(NewNotificationService(notificationService)),
		service.NewAuthorizeServiceServer(NewAuthorizeService(authorizeService)),
		service.NewTradingServiceServer(NewTradingService(tradingService)),
		service.NewPaymentServiceServer(NewPaymentService(paymentService)),
		service.NewReceiptServiceServer(NewReceiptService(receiptService)),
		service.NewSupplierServiceServer(NewSupplierService(supplierService)),
		service.NewCarrierServiceServer(NewCarrierService(carrierService)),
		service.NewLedgerServiceServer(NewLedgerService(ledgerService)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
