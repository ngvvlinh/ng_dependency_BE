package partner

import (
	"context"

	"o.o/api/top/types/common"
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

func (s *FulfillmentService) CreateFulfillment(ctx context.Context, r *FulfillmentCreateFulfillmentEndpoint) error {
	resp, err := s.Shipping.CreateFulfillment(ctx, r.Context.Shop.ID, r.CreateFulfillmentRequest)
	r.Result = resp
	return err
}

func (s *FulfillmentService) CancelFulfillment(ctx context.Context, r *FulfillmentCancelFulfillmentEndpoint) error {
	err := s.Shipping.CancelFulfillment(ctx, r.FulfillmentID, r.CancelReason)
	r.Result = &common.Empty{}
	return err
}
