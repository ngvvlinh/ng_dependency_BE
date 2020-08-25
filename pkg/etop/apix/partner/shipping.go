package partner

import (
	"context"

	api "o.o/api/top/external/partner"
	externaltypes "o.o/api/top/external/types"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/apix/shipping"
	"o.o/backend/pkg/etop/authorize/session"
)

type ShippingService struct {
	session.Session

	Shipping *shipping.Shipping
}

func (s *ShippingService) Clone() api.ShippingService { res := *s; return &res }

func (s *ShippingService) GetShippingServices(ctx context.Context, r *externaltypes.GetShippingServicesRequest) (*externaltypes.GetShippingServicesResponse, error) {
	resp, err := s.Shipping.GetShippingServices(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *ShippingService) CreateAndConfirmOrder(ctx context.Context, r *externaltypes.CreateAndConfirmOrderRequest) (*externaltypes.OrderAndFulfillments, error) {
	userID := cm.CoalesceID(s.SS.Claim().UserID, s.SS.Shop().OwnerID)
	resp, err := s.Shipping.CreateAndConfirmOrder(ctx, userID, s.SS.Shop(), s.SS.Partner(), r)
	return resp, err
}

func (s *ShippingService) CancelOrder(ctx context.Context, r *externaltypes.CancelOrderRequest) (*externaltypes.OrderAndFulfillments, error) {
	userID := cm.CoalesceID(s.SS.Claim().UserID, s.SS.Shop().OwnerID)
	resp, err := s.Shipping.CancelOrder(ctx, userID, s.SS.Shop().ID, r)
	return resp, err
}

func (s *ShippingService) GetOrder(ctx context.Context, r *externaltypes.OrderIDRequest) (*externaltypes.OrderAndFulfillments, error) {
	resp, err := s.Shipping.GetOrder(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *ShippingService) GetFulfillment(ctx context.Context, r *externaltypes.FulfillmentIDRequest) (*externaltypes.Fulfillment, error) {
	resp, err := s.Shipping.GetFulfillment(ctx, s.SS.Shop().ID, r)
	return resp, err
}
