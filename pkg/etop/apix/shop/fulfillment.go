package xshop

import (
	"context"

	"o.o/backend/pkg/etop/apix/shipping"
)

type FulfillmentService struct {
	Shipping *shipping.Shipping
}

func (s *FulfillmentService) Clone() *FulfillmentService { res := *s; return &res }

func (s *FulfillmentService) GetFulfillment(ctx context.Context, r *FulfillmentGetFulfillmentEndpoint) error {
	resp, err := s.Shipping.GetFulfillment(ctx, r.Context.Shop.ID, r.FulfillmentIDRequest)
	r.Result = resp
	return err
}

func (s *FulfillmentService) ListFulfillments(ctx context.Context, r *FulfillmentListFulfillmentsEndpoint) error {
	resp, err := s.Shipping.ListFulfillments(ctx, r.Context.Shop.ID, r.ListFulfillmentsRequest)
	r.Result = resp
	return err
}
