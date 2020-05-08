package partner

import (
	service "o.o/api/top/external/partner"
	"o.o/capi/httprpc"
)

// +gen:wrapper=o.o/api/top/external/partner
// +gen:wrapper:package=partner
// +gen:wrapper:prefix=ext

func NewPartnerServer(m httprpc.Muxer) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService.Clone)),
		service.NewShopServiceServer(WrapShopService(shopService.Clone)),
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
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
