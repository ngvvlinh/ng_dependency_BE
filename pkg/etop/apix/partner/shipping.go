package partner

import (
	"context"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/apix/shipping"
)

type ShippingService struct {
	Shipping *shipping.Shipping
}

func (s *ShippingService) Clone() *ShippingService { res := *s; return &res }

func (s *ShippingService) GetShippingServices(ctx context.Context, r *GetShippingServicesEndpoint) error {
	resp, err := s.Shipping.GetShippingServices(ctx, r.Context.Shop.ID, r.GetShippingServicesRequest)
	r.Result = resp
	return err
}

func (s *ShippingService) CreateAndConfirmOrder(ctx context.Context, r *CreateAndConfirmOrderEndpoint) error {
	userID := cm.CoalesceID(r.Context.UserID, r.Context.Shop.OwnerID)
	resp, err := s.Shipping.CreateAndConfirmOrder(ctx, userID, r.Context.Shop.ID, &r.Context, r.CreateAndConfirmOrderRequest)
	r.Result = resp
	return err
}

func (s *ShippingService) CancelOrder(ctx context.Context, r *CancelOrderEndpoint) error {
	userID := cm.CoalesceID(r.Context.UserID, r.Context.Shop.OwnerID)
	resp, err := s.Shipping.CancelOrder(ctx, userID, r.Context.Shop.ID, r.CancelOrderRequest)
	r.Result = resp
	return err
}

func (s *ShippingService) GetOrder(ctx context.Context, r *GetOrderEndpoint) error {
	resp, err := s.Shipping.GetOrder(ctx, r.Context.Shop.ID, r.OrderIDRequest)
	r.Result = resp
	return err
}

func (s *ShippingService) GetFulfillment(ctx context.Context, r *GetFulfillmentEndpoint) error {
	resp, err := s.Shipping.GetFulfillment(ctx, r.Context.Shop.ID, r.FulfillmentIDRequest)
	r.Result = resp
	return err
}
