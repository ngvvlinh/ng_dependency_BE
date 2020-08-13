package xshop

import (
	"context"

	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/apix/shipping"
)

type ShipnowService struct {
	Shipping *shipping.Shipping
}

func (s *ShipnowService) Clone() *ShipnowService {
	res := *s
	return &res
}

func (s *ShipnowService) GetShipnowServices(ctx context.Context, r *GetShipnowServicesEndpoint) error {
	resp, err := s.Shipping.GetShipnowServices(ctx, r.Context.Shop.ID, r.GetShipnowServicesRequest)
	r.Result = resp
	return err
}

func (s *ShipnowService) CreateShipnowFulfillment(ctx context.Context, r *CreateShipnowFulfillmentEndpoint) error {
	resp, err := s.Shipping.CreateShipnowFulfillment(ctx, &r.Context, r.CreateShipnowFulfillmentRequest)
	r.Result = resp
	return err
}

func (s *ShipnowService) CancelShipnowFulfillment(ctx context.Context, r *CancelShipnowFulfillmentEndpoint) error {
	err := s.Shipping.CancelShipnowFulfillment(ctx, r.Context.Shop.ID, r.CancelShipnowFulfillmentRequest)
	if err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: 1}
	return nil
}

func (s *ShipnowService) GetShipnowFulfillment(ctx context.Context, r *GetShipnowFulfillmentEndpoint) error {
	res, err := s.Shipping.GetShipnowFulfillment(ctx, r.Context.Shop.ID, r.FulfillmentIDRequest)
	r.Result = res
	return err
}
