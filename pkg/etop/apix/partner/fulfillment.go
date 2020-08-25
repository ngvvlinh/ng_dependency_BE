package partner

import (
	"context"

	api "o.o/api/top/external/partner"
	externaltypes "o.o/api/top/external/types"
	pbcm "o.o/api/top/types/common"
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

func (s *FulfillmentService) CreateFulfillment(ctx context.Context, r *externaltypes.CreateFulfillmentRequest) (*externaltypes.Fulfillment, error) {
	resp, err := s.Shipping.CreateFulfillment(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *FulfillmentService) CancelFulfillment(ctx context.Context, r *externaltypes.CancelFulfillmentRequest) (*pbcm.Empty, error) {
	err := s.Shipping.CancelFulfillment(ctx, r.FulfillmentID, r.CancelReason)
	return &pbcm.Empty{}, err
}
