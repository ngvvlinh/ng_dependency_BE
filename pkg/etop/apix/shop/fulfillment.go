package xshop

import (
	"context"

	"etop.vn/backend/pkg/etop/apix/shipping"
)

func (s *FulfillmentService) GetFulfillment(ctx context.Context, r *FulfillmentGetFulfillmentEndpoint) error {
	resp, err := shipping.GetFulfillment(ctx, r.Context.Shop.ID, r.FulfillmentIDRequest)
	r.Result = resp
	return err
}

func (s *FulfillmentService) ListFulfillments(ctx context.Context, r *FulfillmentListFulfillmentsEndpoint) error {
	resp, err := shipping.ListFulfillments(ctx, r.Context.Shop.ID, r.ListFulfillmentsRequest)
	r.Result = resp
	return err
}
