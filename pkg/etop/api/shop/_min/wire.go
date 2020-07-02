package shop_min

import (
	"github.com/google/wire"

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
	receiptService *shop.ReceiptService,
	carrierService *shop.CarrierService,
	ledgerService *shop.LedgerService,
	purchaseOrderService *shop.PurchaseOrderService,
	stocktakeService *shop.StocktakeService,
	shipmentService *shop.ShipmentService,
	connectionService *shop.ConnectionService,
	refundService *shop.RefundService,
	purchaseRefundService *shop.PurchaseRefundService,
) shop.Servers {

	shop.InitIdemp(rd)
	shop.ProductServiceImpl = productService
	shop.StocktakeServiceImpl = stocktakeService
	shop.InventoryServiceImpl = inventoryService

	servers := httprpc.MustNewServers(
		accountService.Clone,
		authorizeService.Clone,
		brandService.Clone,
		carrierService.Clone,
		categoryService.Clone,
		collectionService.Clone,
		connectionService.Clone,
		customerGroupService.Clone,
		customerService.Clone,
		exportService.Clone,
		fulfillmentService.Clone,
		historyService.Clone,
		inventoryService.Clone,
		ledgerService.Clone,
		miscService.Clone,
		moneyTransactionService.Clone,
		notificationService.Clone,
		orderService.Clone,
		productService.Clone,
		productSourceService.Clone,
		purchaseOrderService.Clone,
		purchaseRefundService.Clone,
		receiptService.Clone,
		refundService.Clone,
		shipmentService.Clone,
		stocktakeService.Clone,
		summaryService.Clone,
	)
	return servers
}
