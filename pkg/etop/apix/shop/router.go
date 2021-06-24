package xshop

import (
	"o.o/backend/pkg/common/apifw/idemp"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/apix/shipping"
	"o.o/capi/httprpc"
)

var idempgroup *idemp.RedisGroup

const PrefixIdempShopAPI = "IdempShopAPI"

type Servers []httprpc.Server

func NewServers(
	rd redis.Store,
	_ *shipping.Shipping,
	miscService *MiscService,
	webhookService *WebhookService,
	historyService *HistoryService,
	shippingService *ShippingService,
	orderService *OrderService,
	fulfillmentService *FulfillmentService,
	customerService *CustomerService,
	customerAddressService *CustomerAddressService,
	customerGroupService *CustomerGroupService,
	customerGroupRelationshipService *CustomerGroupRelationshipService,
	inventoryService *InventoryService,
	variantService *VariantService,
	productService *ProductService,
	productCollectionService *ProductCollectionService,
	productCollectionRelationshipService *ProductCollectionRelationshipService,
	shipnowService *ShipnowService,
	etelecomService *EtelecomService,
) (Servers, func()) {
	idempgroup = idemp.NewRedisGroup(rd, PrefixIdempShopAPI, 0)
	servers := httprpc.MustNewServers(
		miscService.Clone,
		webhookService.Clone,
		historyService.Clone,
		shippingService.Clone,
		orderService.Clone,
		fulfillmentService.Clone,
		customerService.Clone,
		customerAddressService.Clone,
		customerGroupService.Clone,
		customerGroupRelationshipService.Clone,
		inventoryService.Clone,
		variantService.Clone,
		productService.Clone,
		productCollectionService.Clone,
		productCollectionRelationshipService.Clone,
		shipnowService.Clone,
		etelecomService.Clone,
	)
	return servers, idempgroup.Shutdown
}
