package fabo

import (
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/api/shop"
	"o.o/backend/pkg/etop/api/shop/account"
	"o.o/backend/pkg/etop/api/shop/authorize"
	"o.o/backend/pkg/etop/api/shop/category"
	"o.o/backend/pkg/etop/api/shop/collection"
	"o.o/backend/pkg/etop/api/shop/connection"
	"o.o/backend/pkg/etop/api/shop/customer"
	"o.o/backend/pkg/etop/api/shop/customergroup"
	"o.o/backend/pkg/etop/api/shop/fulfillment"
	"o.o/backend/pkg/etop/api/shop/history"
	"o.o/backend/pkg/etop/api/shop/notification"
	"o.o/backend/pkg/etop/api/shop/order"
	"o.o/backend/pkg/etop/api/shop/product"
	"o.o/backend/pkg/etop/api/shop/setting"
	"o.o/backend/pkg/etop/api/shop/shipment"
	"o.o/capi/httprpc"
)

func NewServers(
	rd redis.Store,
	miscService *shop.MiscService,
	accountService *account.AccountService,
	collectionService *collection.CollectionService,
	customerService *customer.CustomerService,
	customerGroupService *customergroup.CustomerGroupService,
	productService *product.ProductService,
	categoryService *category.CategoryService,
	orderService *order.OrderService,
	fulfillmentService *fulfillment.FulfillmentService,
	historyService *history.HistoryService,
	// exportService *export.ExportService,
	notificationService *notification.NotificationService,
	authorizeService *authorize.AuthorizeService,
	shipmentService *shipment.ShipmentService,
	settingService *setting.SettingService,
	connectionService *connection.ConnectionService,
) shop.Servers {

	shop.InitIdemp(rd)
	shop.ProductServiceImpl = productService

	servers := httprpc.MustNewServers(
		accountService.Clone,
		authorizeService.Clone,
		categoryService.Clone,
		collectionService.Clone,
		connectionService.Clone,
		customerGroupService.Clone,
		customerService.Clone,
		// exportService.Clone,
		fulfillmentService.Clone,
		historyService.Clone,
		miscService.Clone,
		notificationService.Clone,
		orderService.Clone,
		productService.Clone,
		shipmentService.Clone,
		settingService.Close,
	)
	return servers
}
