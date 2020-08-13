package shop_min

import (
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/api/shop"
	"o.o/backend/pkg/etop/api/shop/account"
	"o.o/backend/pkg/etop/api/shop/authorize"
	"o.o/backend/pkg/etop/api/shop/brand"
	"o.o/backend/pkg/etop/api/shop/carrier"
	"o.o/backend/pkg/etop/api/shop/category"
	"o.o/backend/pkg/etop/api/shop/collection"
	"o.o/backend/pkg/etop/api/shop/connection"
	"o.o/backend/pkg/etop/api/shop/customer"
	"o.o/backend/pkg/etop/api/shop/customergroup"
	"o.o/backend/pkg/etop/api/shop/export"
	"o.o/backend/pkg/etop/api/shop/fulfillment"
	"o.o/backend/pkg/etop/api/shop/history"
	"o.o/backend/pkg/etop/api/shop/inventory"
	"o.o/backend/pkg/etop/api/shop/notification"
	"o.o/backend/pkg/etop/api/shop/order"
	"o.o/backend/pkg/etop/api/shop/product"
	"o.o/backend/pkg/etop/api/shop/shipment"
	"o.o/backend/pkg/etop/api/shop/stocktake"
	"o.o/backend/pkg/etop/api/shop/summary"
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
	orderService *order.OrderService,
	fulfillmentService *fulfillment.FulfillmentService,
	historyService *history.HistoryService,
	summaryService *summary.SummaryService,
	exportService *export.ExportService,
	notificationService *notification.NotificationService,
	authorizeService *authorize.AuthorizeService,
	carrierService *carrier.CarrierService,
	stocktakeService *stocktake.StocktakeService,
	shipmentService *shipment.ShipmentService,
	connectionService *connection.ConnectionService,
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
		miscService.Clone,
		notificationService.Clone,
		orderService.Clone,
		productService.Clone,
		shipmentService.Clone,
		stocktakeService.Clone,
		summaryService.Clone,
	)
	return servers
}