package shipnow

import (
	"context"
	"time"

	"o.o/api/main/shipnow"
	carriertypes "o.o/api/main/shipnow/carrier/types"
	shippingtypes "o.o/api/main/shipping/types"
	"o.o/api/top/int/shop"
	"o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type ShipnowService struct {
	session.Session

	ShipnowAggr  shipnow.CommandBus
	ShipnowQuery shipnow.QueryBus
}

func (s *ShipnowService) Clone() shop.ShipnowService { res := *s; return &res }

func (s *ShipnowService) GetShipnowFulfillment(ctx context.Context, q *pbcm.IDRequest) (*types.ShipnowFulfillment, error) {
	query := &shipnow.GetShipnowFulfillmentQuery{
		ID:     q.Id,
		ShopID: s.SS.Shop().ID,
		Result: nil,
	}
	if err := s.ShipnowQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	result := convertpb.Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(query.Result.ShipnowFulfillment)
	return result, nil
}

func (s *ShipnowService) GetShipnowFulfillments(ctx context.Context, q *types.GetShipnowFulfillmentsRequest) (*types.ShipnowFulfillments, error) {
	shopIDs, err := api.MixAccount(s.SS.Claim(), q.Mixed)
	if err != nil {
		return nil, err
	}
	paging := cmapi.CMPaging(q.Paging)

	query := &shipnow.GetShipnowFulfillmentsQuery{
		ShopIds: shopIDs,
		Paging:  paging,
		Filters: cmapi.ToFiltersPtr(q.Filters),
	}
	if err := s.ShipnowQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &types.ShipnowFulfillments{
		ShipnowFulfillments: convertpb.Convert_core_ShipnowFulfillments_To_api_ShipnowFulfillments(query.Result.ShipnowFulfillments),
		Paging:              cmapi.PbPageInfo(paging),
	}
	return result, nil
}

func (s *ShipnowService) CreateShipnowFulfillmentV2(ctx context.Context, q *types.CreateShipnowFulfillmentV2Request) (*types.ShipnowFulfillment, error) {
	pickupAddress, err := convertpb.OrderAddressFulfilled(q.PickupAddress)
	if err != nil {
		return nil, err
	}

	var deliveryPoints []*shipnow.OrderShippingInfo
	for _, point := range q.DeliveryPoints {
		shippingAddress, err := convertpb.OrderAddressFulfilled(point.ShippingAddress)
		if err != nil {
			return nil, err
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
	cmd := &shipnow.CreateShipnowFulfillmentCommand{
		DeliveryPoints:      deliveryPoints,
		Carrier:             q.Carrier,
		ShopID:              s.SS.Shop().ID,
		ShippingServiceCode: q.ShippingServiceCode,
		ShippingServiceFee:  q.ShippingServiceFee,
		ShippingNote:        q.ShippingNote,
		PickupAddress:       convertpb.Convert_api_OrderAddress_To_core_OrderAddress(pickupAddress),
		ConnectionID:        q.ConnectionID,
	}
	if err := s.ShipnowAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(cmd.Result)
	return result, nil
}

func (s *ShipnowService) ConfirmShipnowFulfillment(ctx context.Context, q *pbcm.IDRequest) (*types.ShipnowFulfillment, error) {
	cmd := &shipnow.ConfirmShipnowFulfillmentCommand{
		ID:     q.Id,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.ShipnowAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(cmd.Result)
	return result, nil
}

func (s *ShipnowService) UpdateShipnowFulfillment(ctx context.Context, q *types.UpdateShipnowFulfillmentRequest) (*types.ShipnowFulfillment, error) {
	pickupAddress, err := convertpb.OrderAddressFulfilled(q.PickupAddress)
	if err != nil {
		return nil, err
	}
	_carrier, _ := carriertypes.ParseShipnowCarrier(q.Carrier)

	var deliveryPoints []*shipnow.OrderShippingInfo
	for _, point := range q.DeliveryPoints {
		shippingAddress, err := convertpb.OrderAddressFulfilled(point.ShippingAddress)
		if err != nil {
			return nil, err
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

	cmd := &shipnow.UpdateShipnowFulfillmentCommand{
		ID:                  q.Id,
		DeliveryPoints:      deliveryPoints,
		Carrier:             _carrier,
		ShopID:              s.SS.Shop().ID,
		ShippingServiceCode: q.ShippingServiceCode,
		ShippingServiceFee:  q.ShippingServiceFee,
		ShippingNote:        q.ShippingNote,
		RequestPickupAt:     time.Time{},
		PickupAddress:       convertpb.Convert_api_OrderAddress_To_core_OrderAddress(pickupAddress),
	}
	if err := s.ShipnowAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(cmd.Result)
	return result, nil
}

func (s *ShipnowService) CancelShipnowFulfillment(ctx context.Context, q *types.CancelShipnowFulfillmentRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &shipnow.CancelShipnowFulfillmentCommand{
		ID:           q.Id,
		ShopID:       s.SS.Shop().ID,
		CancelReason: q.CancelReason,
	}
	if err := s.ShipnowAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	result := &pbcm.UpdatedResponse{
		Updated: 1,
	}
	return result, nil
}

func (s *ShipnowService) GetShipnowServices(ctx context.Context, q *types.GetShipnowServicesRequest) (*types.GetShipnowServicesResponse, error) {
	pickupAddress, err := convertpb.OrderAddressFulfilled(q.PickupAddress)
	if err != nil {
		return nil, err
	}
	var points []*shipnow.DeliveryPoint
	if len(q.DeliveryPoints) > 0 {
		for _, p := range q.DeliveryPoints {
			addr, err := convertpb.OrderAddressFulfilled(p.ShippingAddress)
			if err != nil {
				return nil, err
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
		ShopId:         s.SS.Shop().ID,
		OrderIds:       q.OrderIds,
		PickupAddress:  convertpb.Convert_api_OrderAddress_To_core_OrderAddress(pickupAddress),
		DeliveryPoints: points,
	}
	if err := s.ShipnowAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &types.GetShipnowServicesResponse{
		Services: convertpb.Convert_core_ShipnowServices_To_api_ShipnowServices(cmd.Result.Services),
	}
	return result, nil
}
