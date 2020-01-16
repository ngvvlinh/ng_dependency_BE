package xshop

import (
	"context"

	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/apix/shipping"
)

func init() {
	bus.AddHandlers("apix",
		miscService.GetLocationList,
		shippingService.GetShippingServices,
		shippingService.CreateAndConfirmOrder,
		shippingService.CancelOrder,
		shippingService.GetOrder,
		shippingService.GetFulfillment,
	)
}

type MiscService struct{}
type WebhookService struct{}
type HistoryService struct{}
type ShippingService struct{}
type CustomerService struct{}
type CustomerAddressService struct{}
type CustomerGroupService struct{}
type CustomerGroupRelationshipService struct{}
type InventoryService struct{}
type OrderService struct{}
type FulfillmentService struct{}
type ProductService struct{}
type ProductCollectionService struct{}
type ProductCollectionRelationshipService struct{}
type VariantService struct{}

var miscService = &MiscService{}
var webhookService = &WebhookService{}
var historyService = &HistoryService{}
var shippingService = &ShippingService{}
var customerService = &CustomerService{}
var customerAddressService = &CustomerAddressService{}
var customerGroupService = &CustomerGroupService{}
var customerGroupRelationshipService = &CustomerGroupRelationshipService{}
var inventoryService = &InventoryService{}
var orderService = &OrderService{}
var fulfillmentService = &FulfillmentService{}
var productService = &ProductService{}
var productCollectionService = &ProductCollectionService{}
var productCollectionRelationshipService = &ProductCollectionRelationshipService{}
var variantService = &VariantService{}

func (s *MiscService) GetLocationList(ctx context.Context, r *GetLocationListEndpoint) error {
	resp, err := shipping.GetLocationList(ctx)
	r.Result = resp
	return err
}

func (s *ShippingService) GetShippingServices(ctx context.Context, r *GetShippingServicesEndpoint) error {
	resp, err := shipping.GetShippingServices(ctx, r.Context.Shop.ID, r.GetShippingServicesRequest)
	r.Result = resp
	return err
}

func (s *ShippingService) CreateAndConfirmOrder(ctx context.Context, r *CreateAndConfirmOrderEndpoint) error {
	resp, err := shipping.CreateAndConfirmOrder(ctx, r.Context.Shop.ID, &r.Context, r.CreateAndConfirmOrderRequest)
	r.Result = resp
	return err
}

func (s *ShippingService) CancelOrder(ctx context.Context, r *CancelOrderEndpoint) error {
	resp, err := shipping.CancelOrder(ctx, r.Context.Shop.ID, r.CancelOrderRequest)
	r.Result = resp
	return err
}

func (s *ShippingService) GetOrder(ctx context.Context, r *GetOrderEndpoint) error {
	resp, err := shipping.GetOrder(ctx, r.Context.Shop.ID, r.OrderIDRequest)
	r.Result = resp
	return err
}

func (s *ShippingService) GetFulfillment(ctx context.Context, r *GetFulfillmentEndpoint) error {
	resp, err := shipping.GetFulfillment(ctx, r.Context.Shop.ID, r.FulfillmentIDRequest)
	r.Result = resp
	return err
}
