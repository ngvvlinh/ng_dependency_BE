// +build wireinject

package fabo

import (
	"github.com/google/wire"

	"o.o/backend/pkg/etop/api/shop"
	"o.o/backend/pkg/etop/api/shop/account"
	"o.o/backend/pkg/etop/api/shop/authorize"
	"o.o/backend/pkg/etop/api/shop/category"
	"o.o/backend/pkg/etop/api/shop/collection"
	"o.o/backend/pkg/etop/api/shop/connection"
	"o.o/backend/pkg/etop/api/shop/customer"
	"o.o/backend/pkg/etop/api/shop/customergroup"
	"o.o/backend/pkg/etop/api/shop/export"
	"o.o/backend/pkg/etop/api/shop/fulfillment"
	"o.o/backend/pkg/etop/api/shop/history"
	"o.o/backend/pkg/etop/api/shop/notification"
	"o.o/backend/pkg/etop/api/shop/order"
	"o.o/backend/pkg/etop/api/shop/product"
	"o.o/backend/pkg/etop/api/shop/setting"
	"o.o/backend/pkg/etop/api/shop/shipment"
)

// TODO(ngoc): remove unnecessary services
var WireSet = wire.NewSet(
	wire.Struct(new(account.AccountService), "*"),
	wire.Struct(new(authorize.AuthorizeService), "*"),
	wire.Struct(new(category.CategoryService), "*"),
	wire.Struct(new(collection.CollectionService), "*"),
	wire.Struct(new(connection.ConnectionService), "*"),
	wire.Struct(new(customergroup.CustomerGroupService), "*"),
	wire.Struct(new(customer.CustomerService), "*"),
	wire.Struct(new(export.ExportService), "*"),
	wire.Struct(new(fulfillment.FulfillmentService), "*"),
	wire.Struct(new(history.HistoryService), "*"),
	wire.Struct(new(shop.MiscService), "*"),
	wire.Struct(new(notification.NotificationService), "*"),
	wire.Struct(new(order.OrderService), "*"),
	wire.Struct(new(product.ProductService), "*"),
	wire.Struct(new(shipment.ShipmentService), "*"),
	wire.Struct(new(setting.SettingService), "*"),
)
