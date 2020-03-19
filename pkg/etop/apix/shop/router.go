package xshop

import (
	service "etop.vn/api/top/external/shop"
	"etop.vn/capi/httprpc"
)

// +gen:wrapper=etop.vn/api/top/external/shop
// +gen:wrapper:package=shop
// +gen:wrapper:prefix=ext

func NewShopServer(m httprpc.Muxer) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService)),
		service.NewWebhookServiceServer(WrapWebhookService(webhookService)),
		service.NewHistoryServiceServer(WrapHistoryService(historyService)),
		service.NewShippingServiceServer(WrapShippingService(shippingService)),
		service.NewOrderServiceServer(WrapOrderService(orderService)),
		service.NewFulfillmentServiceServer(WrapFulfillmentService(fulfillmentService)),
		service.NewCustomerServiceServer(WrapCustomerService(customerService)),
		service.NewCustomerAddressServiceServer(WrapCustomerAddressService(customerAddressService)),
		service.NewCustomerGroupServiceServer(WrapCustomerGroupService(customerGroupService)),
		service.NewCustomerGroupRelationshipServiceServer(WrapCustomerGroupRelationshipService(customerGroupRelationshipService)),
		service.NewInventoryServiceServer(WrapInventoryService(inventoryService)),
		service.NewVariantServiceServer(WrapVariantService(variantService)),
		service.NewProductServiceServer(WrapProductService(productService)),
		service.NewProductCollectionServiceServer(WrapProductCollectionService(productCollectionService)),
		service.NewProductCollectionRelationshipServiceServer(WrapProductCollectionRelationshipService(productCollectionRelationshipService)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
