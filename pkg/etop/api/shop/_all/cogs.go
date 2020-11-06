package shop_all

import (
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/api/shop"
	"o.o/backend/pkg/etop/api/shop/account"
	"o.o/backend/pkg/etop/api/shop/accountshipnow"
	"o.o/backend/pkg/etop/api/shop/authorize"
	"o.o/backend/pkg/etop/api/shop/brand"
	"o.o/backend/pkg/etop/api/shop/carrier"
	"o.o/backend/pkg/etop/api/shop/category"
	"o.o/backend/pkg/etop/api/shop/collection"
	"o.o/backend/pkg/etop/api/shop/connection"
	"o.o/backend/pkg/etop/api/shop/contact"
	"o.o/backend/pkg/etop/api/shop/customer"
	"o.o/backend/pkg/etop/api/shop/customergroup"
	"o.o/backend/pkg/etop/api/shop/export"
	"o.o/backend/pkg/etop/api/shop/fulfillment"
	"o.o/backend/pkg/etop/api/shop/history"
	"o.o/backend/pkg/etop/api/shop/inventory"
	"o.o/backend/pkg/etop/api/shop/ledger"
	"o.o/backend/pkg/etop/api/shop/money_transaction"
	"o.o/backend/pkg/etop/api/shop/notification"
	"o.o/backend/pkg/etop/api/shop/order"
	"o.o/backend/pkg/etop/api/shop/payment"
	"o.o/backend/pkg/etop/api/shop/product"
	"o.o/backend/pkg/etop/api/shop/product_source"
	"o.o/backend/pkg/etop/api/shop/purchase_order"
	"o.o/backend/pkg/etop/api/shop/purchase_refund"
	"o.o/backend/pkg/etop/api/shop/receipt"
	"o.o/backend/pkg/etop/api/shop/refund"
	"o.o/backend/pkg/etop/api/shop/shipment"
	"o.o/backend/pkg/etop/api/shop/shipnow"
	"o.o/backend/pkg/etop/api/shop/stocktake"
	"o.o/backend/pkg/etop/api/shop/subscription"
	"o.o/backend/pkg/etop/api/shop/summary"
	"o.o/backend/pkg/etop/api/shop/supplier"
	"o.o/backend/pkg/etop/api/shop/ticket"
	"o.o/backend/pkg/etop/api/shop/trading"
	"o.o/backend/pkg/etop/api/shop/ws"
	"o.o/capi/httprpc"
)

func NewServers(
	rd redis.Store,
	miscService *shop.MiscService,
	brandService *brand.BrandService,
	inventoryService *inventory.InventoryService,
	accountService *account.AccountService,
	collectionService *collection.CollectionService,
	customerService *customer.CustomerService,
	customerGroupService *customergroup.CustomerGroupService,
	productService *product.ProductService,
	categoryService *category.CategoryService,
	productSourceService *product_source.ProductSourceService,
	orderService *order.OrderService,
	fulfillmentService *fulfillment.FulfillmentService,
	shipnowService *shipnow.ShipnowService,
	historyService *history.HistoryService,
	moneyTransactionService *money_transaction.MoneyTransactionService,
	summaryService *summary.SummaryService,
	exportService *export.ExportService,
	notificationService *notification.NotificationService,
	authorizeService *authorize.AuthorizeService,
	tradingService *trading.TradingService,
	paymentService *payment.PaymentService,
	receiptService *receipt.ReceiptService,
	supplierService *supplier.SupplierService,
	carrierService *carrier.CarrierService,
	ledgerService *ledger.LedgerService,
	purchaseOrderService *purchase_order.PurchaseOrderService,
	stocktakeService *stocktake.StocktakeService,
	shipmentService *shipment.ShipmentService,
	connectionService *connection.ConnectionService,
	refundService *refund.RefundService,
	purchaseRefundService *purchase_refund.PurchaseRefundService,
	webServerService *ws.WebServerService,
	subscriptionService *subscription.SubscriptionService,
	ticketService *ticket.TicketService,
	accountshipnowService *accountshipnow.AccountShipnowService,
	contactService *contact.ContactService,
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
		paymentService.Clone,
		productService.Clone,
		productSourceService.Clone,
		purchaseOrderService.Clone,
		purchaseRefundService.Clone,
		receiptService.Clone,
		refundService.Clone,
		shipmentService.Clone,
		shipnowService.Clone,
		stocktakeService.Clone,
		subscriptionService.Clone,
		summaryService.Clone,
		supplierService.Clone,
		tradingService.Clone,
		webServerService.Clone,
		ticketService.Clone,
		accountshipnowService.Clone,
		contactService.Clone,
	)
	return servers
}
