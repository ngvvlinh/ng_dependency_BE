package shop

import (
	"context"
	"time"

	"o.o/api/main/shipnow"
	carriertypes "o.o/api/main/shipnow/carrier/types"
	shippingtypes "o.o/api/main/shipping/types"
	apitypes "o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/convertpb"
)

type ShipnowService struct {
	ShipnowAggr  shipnow.CommandBus
	ShipnowQuery shipnow.QueryBus
}

func (s *ShipnowService) Clone() *ShipnowService { res := *s; return &res }

func (s *ShipnowService) GetShipnowFulfillment(ctx context.Context, q *GetShipnowFulfillmentEndpoint) error {
	query := &shipnow.GetShipnowFulfillmentQuery{
		Id:     q.Id,
		ShopId: q.Context.Shop.ID,
		Result: nil,
	}
	if err := s.ShipnowQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = convertpb.Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(query.Result.ShipnowFulfillment)
	return nil
}

func (s *ShipnowService) GetShipnowFulfillments(ctx context.Context, q *GetShipnowFulfillmentsEndpoint) error {
	shopIDs, err := api.MixAccount(q.Context.Claim, q.Mixed)
	if err != nil {
		return err
	}
	paging := cmapi.CMPaging(q.Paging)

	query := &shipnow.GetShipnowFulfillmentsQuery{
		ShopIds: shopIDs,
		Paging:  paging,
		Filters: cmapi.ToFiltersPtr(q.Filters),
	}
	if err := s.ShipnowQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &apitypes.ShipnowFulfillments{
		ShipnowFulfillments: convertpb.Convert_core_ShipnowFulfillments_To_api_ShipnowFulfillments(query.Result.ShipnowFulfillments),
		Paging:              cmapi.PbPageInfo(paging),
	}
	return nil
}

func (s *ShipnowService) CreateShipnowFulfillment(ctx context.Context, q *CreateShipnowFulfillmentEndpoint) error {
	pickupAddress, err := convertpb.OrderAddressFulfilled(q.PickupAddress)
	if err != nil {
		return err
	}
	_carrier, _ := carriertypes.ParseCarrier(q.Carrier)
	cmd := &shipnow.CreateShipnowFulfillmentCommand{
		OrderIds:            q.OrderIds,
		Carrier:             _carrier,
		ShopId:              q.Context.Shop.ID,
		ShippingServiceCode: q.ShippingServiceCode,
		ShippingServiceFee:  q.ShippingServiceFee,
		ShippingNote:        q.ShippingNote,
		RequestPickupAt:     time.Time{},
		PickupAddress:       convertpb.Convert_api_OrderAddress_To_core_OrderAddress(pickupAddress),
	}
	if err := s.ShipnowAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(cmd.Result)
	return nil
}

func (s *ShipnowService) CreateShipnowFulfillmentV2(ctx context.Context, q *CreateShipnowFulfillmentV2Endpoint) error {
	pickupAddress, err := convertpb.OrderAddressFulfilled(q.PickupAddress)
	if err != nil {
		return err
	}

	var deliveryPoints []*shipnow.OrderShippingInfo
	for _, point := range q.DeliveryPoints {
		shippingAddress, err := convertpb.OrderAddressFulfilled(point.ShippingAddress)
		if err != nil {
			return err
		}
		p := &shipnow.OrderShippingInfo{
			OrderID:         point.OrderID,
			ShippingAddress: convertpb.Convert_api_OrderAddress_To_core_OrderAddress(shippingAddress),
			ShippingNote:    point.ShippingNote,
			WeightInfo: shippingtypes.WeightInfo{
				GrossWeight:      point.GrossWeight,
				ChargeableWeight: point.ChargeableWeight,
			},
			ValueInfo: shippingtypes.ValueInfo{
				CODAmount: point.CODAmount,
			},
		}
		deliveryPoints = append(deliveryPoints, p)
	}
	cmd := &shipnow.CreateShipnowFulfillmentV2Command{
		DeliveryPoints:      deliveryPoints,
		Carrier:             q.Carrier,
		ShopID:              q.Context.Shop.ID,
		ShippingServiceCode: q.ShippingServiceCode,
		ShippingServiceFee:  q.ShippingServiceFee,
		ShippingNote:        q.ShippingNote,
		PickupAddress:       convertpb.Convert_api_OrderAddress_To_core_OrderAddress(pickupAddress),
	}
	if err := s.ShipnowAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(cmd.Result)
	return nil
}

func (s *ShipnowService) ConfirmShipnowFulfillment(ctx context.Context, q *ConfirmShipnowFulfillmentEndpoint) error {
	cmd := &shipnow.ConfirmShipnowFulfillmentCommand{
		Id:     q.Id,
		ShopId: q.Context.Shop.ID,
	}
	if err := s.ShipnowAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(cmd.Result)
	return nil
}

func (s *ShipnowService) UpdateShipnowFulfillment(ctx context.Context, q *UpdateShipnowFulfillmentEndpoint) error {
	pickupAddress, err := convertpb.OrderAddressFulfilled(q.PickupAddress)
	if err != nil {
		return err
	}
	_carrier, _ := carriertypes.ParseCarrier(q.Carrier)
	cmd := &shipnow.UpdateShipnowFulfillmentCommand{
		Id:                  q.Id,
		OrderIds:            q.OrderIds,
		Carrier:             _carrier,
		ShopId:              q.Context.Shop.ID,
		ShippingServiceCode: q.ShippingServiceCode,
		ShippingServiceFee:  q.ShippingServiceFee,
		ShippingNote:        q.ShippingNote,
		RequestPickupAt:     time.Time{},
		PickupAddress:       convertpb.Convert_api_OrderAddress_To_core_OrderAddress(pickupAddress),
	}
	if err := s.ShipnowAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(cmd.Result)
	return nil
}

func (s *ShipnowService) CancelShipnowFulfillment(ctx context.Context, q *CancelShipnowFulfillmentEndpoint) error {
	cmd := &shipnow.CancelShipnowFulfillmentCommand{
		Id:           q.Id,
		ShopId:       q.Context.Shop.ID,
		CancelReason: q.CancelReason,
	}
	if err := s.ShipnowAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &pbcm.UpdatedResponse{
		Updated: 1,
	}
	return nil
}

func (s *ShipnowService) GetShipnowServices(ctx context.Context, q *GetShipnowServicesEndpoint) error {
	pickupAddress, err := convertpb.OrderAddressFulfilled(q.PickupAddress)
	if err != nil {
		return err
	}
	var points []*shipnow.DeliveryPoint
	if len(q.DeliveryPoints) > 0 {
		for _, p := range q.DeliveryPoints {
			addr, err := convertpb.OrderAddressFulfilled(p.ShippingAddress)
			if err != nil {
				return err
			}
			points = append(points, &shipnow.DeliveryPoint{
				ShippingAddress: convertpb.Convert_api_OrderAddress_To_core_OrderAddress(addr),
				ValueInfo: shippingtypes.ValueInfo{
					CODAmount: p.CodAmount,
				},
			})
		}
	}

	cmd := &shipnow.GetShipnowServicesCommand{
		ShopId:         q.Context.Shop.ID,
		OrderIds:       q.OrderIds,
		PickupAddress:  convertpb.Convert_api_OrderAddress_To_core_OrderAddress(pickupAddress),
		DeliveryPoints: points,
	}
	if err := s.ShipnowAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &apitypes.GetShipnowServicesResponse{
		Services: convertpb.Convert_core_ShipnowServices_To_api_ShipnowServices(cmd.Result.Services),
	}
	return nil
}
