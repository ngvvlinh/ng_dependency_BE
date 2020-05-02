package xshop

import (
	"context"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/apix/shipping"
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

func (s *MiscService) Clone() *MiscService                       { res := *s; return &res }
func (s *WebhookService) Clone() *WebhookService                 { res := *s; return &res }
func (s *HistoryService) Clone() *HistoryService                 { res := *s; return &res }
func (s *ShippingService) Clone() *ShippingService               { res := *s; return &res }
func (s *CustomerService) Clone() *CustomerService               { res := *s; return &res }
func (s *CustomerAddressService) Clone() *CustomerAddressService { res := *s; return &res }
func (s *CustomerGroupService) Clone() *CustomerGroupService     { res := *s; return &res }
func (s *CustomerGroupRelationshipService) Clone() *CustomerGroupRelationshipService {
	res := *s
	return &res
}
func (s *InventoryService) Clone() *InventoryService                 { res := *s; return &res }
func (s *OrderService) Clone() *OrderService                         { res := *s; return &res }
func (s *FulfillmentService) Clone() *FulfillmentService             { res := *s; return &res }
func (s *ProductService) Clone() *ProductService                     { res := *s; return &res }
func (s *ProductCollectionService) Clone() *ProductCollectionService { res := *s; return &res }
func (s *ProductCollectionRelationshipService) Clone() *ProductCollectionRelationshipService {
	res := *s
	return &res
}
func (s *VariantService) Clone() *VariantService { res := *s; return &res }

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
	userID := cm.CoalesceID(r.Context.UserID, r.Context.Shop.OwnerID)
	resp, err := shipping.CreateAndConfirmOrder(ctx, userID, r.Context.Shop.ID, &r.Context, r.CreateAndConfirmOrderRequest)
	r.Result = resp
	return err
}

func (s *ShippingService) CancelOrder(ctx context.Context, r *CancelOrderEndpoint) error {
	userID := cm.CoalesceID(r.Context.UserID, r.Context.Shop.OwnerID)
	resp, err := shipping.CancelOrder(ctx, userID, r.Context.Shop.ID, r.CancelOrderRequest)
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
