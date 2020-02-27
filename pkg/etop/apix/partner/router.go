package partner

import (
	service "etop.vn/api/top/external/partner"
	"etop.vn/capi/httprpc"
)

// +gen:wrapper=etop.vn/api/top/external/partner
// +gen:wrapper:package=partner
// +gen:wrapper:prefix=ext

func NewPartnerServer(m httprpc.Muxer) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService)),
		service.NewShopServiceServer(WrapShopService(shopService)),
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
		service.NewShipmentConnectionServiceServer(WrapShipmentConnectionService(shipmentConnectionService)),
		service.NewShipmentServiceServer(WrapShipmentService(shipmentService)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
