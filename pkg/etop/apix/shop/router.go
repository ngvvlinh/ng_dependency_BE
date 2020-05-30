package xshop

import (
	service "o.o/api/top/external/shop"
	"o.o/backend/pkg/common/apifw/idemp"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/apix/shipping"
	"o.o/capi/httprpc"
)

// +gen:wrapper=o.o/api/top/external/shop
// +gen:wrapper:package=shop
// +gen:wrapper:prefix=ext

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
) (Servers, func()) {
	idempgroup = idemp.NewRedisGroup(rd, PrefixIdempShopAPI, 0)
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService.Clone)),
		service.NewWebhookServiceServer(WrapWebhookService(webhookService.Clone)),
		service.NewHistoryServiceServer(WrapHistoryService(historyService.Clone)),
		service.NewShippingServiceServer(WrapShippingService(shippingService.Clone)),
		service.NewOrderServiceServer(WrapOrderService(orderService.Clone)),
		service.NewFulfillmentServiceServer(WrapFulfillmentService(fulfillmentService.Clone)),
		service.NewCustomerServiceServer(WrapCustomerService(customerService.Clone)),
		service.NewCustomerAddressServiceServer(WrapCustomerAddressService(customerAddressService.Clone)),
		service.NewCustomerGroupServiceServer(WrapCustomerGroupService(customerGroupService.Clone)),
		service.NewCustomerGroupRelationshipServiceServer(WrapCustomerGroupRelationshipService(customerGroupRelationshipService.Clone)),
		service.NewInventoryServiceServer(WrapInventoryService(inventoryService.Clone)),
		service.NewVariantServiceServer(WrapVariantService(variantService.Clone)),
		service.NewProductServiceServer(WrapProductService(productService.Clone)),
		service.NewProductCollectionServiceServer(WrapProductCollectionService(productCollectionService.Clone)),
		service.NewProductCollectionRelationshipServiceServer(WrapProductCollectionRelationshipService(productCollectionRelationshipService.Clone)),
	}
	return servers, idempgroup.Shutdown
}
