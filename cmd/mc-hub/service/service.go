package service

import (
	"context"

	"o.o/api/top/external/mc/vnp"
	externaltypes "o.o/api/top/external/types"
	pbcm "o.o/api/top/types/common"
)

var _ vnp.ShipnowService = &MCShipnowService{}

type MCShipnowService struct {
	Client vnp.ShipnowService
}

func (s *MCShipnowService) Clone() vnp.ShipnowService {
	res := *s
	return &res
}

func (s *MCShipnowService) Ping(ctx context.Context, r *pbcm.Empty) (*pbcm.Empty, error) {
	return s.Client.Ping(ctx, r)
}

func (s *MCShipnowService) GetShipnowServices(ctx context.Context, r *externaltypes.GetShipnowServicesRequest) (*externaltypes.GetShipnowServicesResponse, error) {
	return s.Client.GetShipnowServices(ctx, r)
}

func (s *MCShipnowService) CreateShipnowFulfillment(ctx context.Context, r *externaltypes.CreateShipnowFulfillmentRequest) (*externaltypes.ShipnowFulfillment, error) {
	return s.Client.CreateShipnowFulfillment(ctx, r)
}

func (s *MCShipnowService) CancelShipnowFulfillment(ctx context.Context, r *externaltypes.CancelShipnowFulfillmentRequest) (*pbcm.UpdatedResponse, error) {
	return s.Client.CancelShipnowFulfillment(ctx, r)
}

func (s *MCShipnowService) GetShipnowFulfillment(ctx context.Context, r *externaltypes.FulfillmentIDRequest) (*externaltypes.ShipnowFulfillment, error) {
	return s.Client.GetShipnowFulfillment(ctx, r)
}
