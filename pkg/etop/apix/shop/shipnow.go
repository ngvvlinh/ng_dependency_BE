package xshop

import (
	"context"

	api "o.o/api/top/external/shop"
	externaltypes "o.o/api/top/external/types"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/apix/shipping"
	"o.o/backend/pkg/etop/authorize/session"
)

type ShipnowService struct {
	session.Session

	Shipping *shipping.Shipping
}

func (s *ShipnowService) Clone() api.ShipnowService {
	res := *s
	return &res
}

func (s *ShipnowService) GetShipnowServices(ctx context.Context, r *externaltypes.GetShipnowServicesRequest) (*externaltypes.GetShipnowServicesResponse, error) {
	resp, err := s.Shipping.GetShipnowServices(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *ShipnowService) CreateShipnowFulfillment(ctx context.Context, r *externaltypes.CreateShipnowFulfillmentRequest) (*externaltypes.ShipnowFulfillment, error) {
	resp, err := s.Shipping.CreateShipnowFulfillment(ctx, s.SS.Claim().UserID, s.SS.Shop(), s.SS.Partner(), r)
	return resp, err
}

func (s *ShipnowService) CancelShipnowFulfillment(ctx context.Context, r *externaltypes.CancelShipnowFulfillmentRequest) (*pbcm.UpdatedResponse, error) {
	err := s.Shipping.CancelShipnowFulfillment(ctx, s.SS.Shop().ID, r)
	if err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: 1}
	return result, nil
}

func (s *ShipnowService) GetShipnowFulfillment(ctx context.Context, r *externaltypes.FulfillmentIDRequest) (*externaltypes.ShipnowFulfillment, error) {
	res, err := s.Shipping.GetShipnowFulfillment(ctx, s.SS.Shop().ID, r)
	return res, err
}
