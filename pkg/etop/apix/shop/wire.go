package xshop

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	wire.Struct(new(MiscService), "*"),
	wire.Struct(new(WebhookService), "*"),
	wire.Struct(new(HistoryService), "*"),
	wire.Struct(new(ShippingService), "*"),
	wire.Struct(new(ShipnowService), "*"),
	wire.Struct(new(OrderService), "*"),
	wire.Struct(new(FulfillmentService), "*"),
	wire.Struct(new(CustomerService), "*"),
	wire.Struct(new(CustomerAddressService), "*"),
	wire.Struct(new(CustomerGroupService), "*"),
	wire.Struct(new(CustomerGroupRelationshipService), "*"),
	wire.Struct(new(InventoryService), "*"),
	wire.Struct(new(VariantService), "*"),
	wire.Struct(new(ProductService), "*"),
	wire.Struct(new(ProductCollectionService), "*"),
	wire.Struct(new(ProductCollectionRelationshipService), "*"),
	NewServers,
)
