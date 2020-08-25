package xshop

import (
	"context"

	api "o.o/api/top/external/shop"
	externaltypes "o.o/api/top/external/types"
	"o.o/backend/pkg/etop/apix/shipping"
	"o.o/backend/pkg/etop/authorize/session"
)

type FulfillmentService struct {
	session.Session

	Shipping *shipping.Shipping
}

func (s *FulfillmentService) Clone() api.FulfillmentService { res := *s; return &res }

func (s *FulfillmentService) GetFulfillment(ctx context.Context, r *externaltypes.FulfillmentIDRequest) (*externaltypes.Fulfillment, error) {
	resp, err := s.Shipping.GetFulfillment(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *FulfillmentService) ListFulfillments(ctx context.Context, r *externaltypes.ListFulfillmentsRequest) (*externaltypes.FulfillmentsResponse, error) {
	resp, err := s.Shipping.ListFulfillments(ctx, s.SS.Shop().ID, r)
	return resp, err
}
