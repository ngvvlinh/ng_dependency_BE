package shop

import (
	"context"
	"fmt"
	"time"

	"o.o/api/main/shipping"
	shippingtypes "o.o/api/main/shipping/types"
	"o.o/api/top/int/shop"
	"o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
	shipmodelx "o.o/backend/com/main/shipping/modelx"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

type ShipmentService struct {
	ShipmentManager   *shippingcarrier.ShipmentManager
	ShippingAggregate shipping.CommandBus
}

func (s *ShipmentService) Clone() *ShipmentService { res := *s; return &res }

func (s *ShipmentService) GetShippingServices(ctx context.Context, q *GetShippingServicesEndpoint) error {
	shopID := q.Context.Shop.ID
	args, err := s.ShipmentManager.PrepareDataGetShippingServices(ctx, q.GetShippingServicesRequest)
	args.AccountID = shopID
	if err != nil {
		return err
	}
	resp, err := s.ShipmentManager.GetShippingServices(ctx, args)
	if err != nil {
		return err
	}
	q.Result = &types.GetShippingServicesResponse{
		Services: convertpb.PbAvailableShippingServices(resp),
	}
	return nil
}

func (s *ShipmentService) CreateFulfillments(ctx context.Context, q *CreateFulfillmentsEndpoint) error {
	key := fmt.Sprintf("CreateFulfillments %v-%v", q.Context.Shop.ID, q.OrderID)
	res, _, err := idempgroup.DoAndWrap(
		ctx, key, 10*time.Second, "tạo đơn giao hàng",
		func() (interface{}, error) { return s.createFulfillments(ctx, q) })

	if err != nil {
		return err
	}
	q.Result = res.(*CreateFulfillmentsEndpoint).Result
	return nil
}

func (s *ShipmentService) createFulfillments(ctx context.Context, q *CreateFulfillmentsEndpoint) (_ *CreateFulfillmentsEndpoint, _err error) {
	shopID := q.Context.Shop.ID
	args := &shipping.CreateFulfillmentsCommand{
		ShopID:              shopID,
		OrderID:             q.OrderID,
		PickupAddress:       convertpb.Convert_api_OrderAddress_To_core_OrderAddress(q.PickupAddress),
		ShippingAddress:     convertpb.Convert_api_OrderAddress_To_core_OrderAddress(q.ShippingAddress),
		ReturnAddress:       convertpb.Convert_api_OrderAddress_To_core_OrderAddress(q.ReturnAddress),
		ShippingServiceCode: q.ShippingServiceCode,
		ShippingServiceFee:  q.ShippingServiceFee,
		ShippingServiceName: q.ShippingServiceName,
		WeightInfo: shippingtypes.WeightInfo{
			GrossWeight:      q.GrossWeight,
			ChargeableWeight: q.ChargeableWeight,
			Length:           q.Length,
			Width:            q.Width,
			Height:           q.Height,
		},
		ValueInfo: shippingtypes.ValueInfo{
			CODAmount:        q.CODAmount,
			IncludeInsurance: q.IncludeInsurance,
		},
		TryOn:         q.TryOn,
		ShippingNote:  q.ShippingNote,
		ShippingType:  q.ShippingType,
		ConnectionID:  q.ConnectionID,
		ShopCarrierID: q.ShopCarrierID,
	}
	if err := s.ShippingAggregate.Dispatch(ctx, args); err != nil {
		return nil, err
	}
	query := &shipmodelx.GetFulfillmentExtendedsQuery{
		ShopIDs: []dot.ID{shopID},
		IDs:     args.Result,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	ffms := convertpb.PbFulfillmentExtendeds(query.Result.Fulfillments, model.TagShop)
	res := &CreateFulfillmentsEndpoint{
		Result: &shop.CreateFulfillmentsResponse{
			Fulfillments: ffms,
		},
	}
	return res, nil
}

func (s *ShipmentService) CancelFulfillment(ctx context.Context, q *CancelFulfillmentEndpoint) error {
	key := fmt.Sprintf("CancelFulfillment %v-%v", q.Context.Shop.ID, q.FulfillmentID)
	res, _, err := idempgroup.DoAndWrap(
		ctx, key, 10*time.Second, "huỷ đơn giao hàng",
		func() (interface{}, error) { return s.cancelFulfillment(ctx, q) })

	if err != nil {
		return err
	}
	q.Result = res.(*CancelFulfillmentEndpoint).Result
	return nil
}

func (s *ShipmentService) cancelFulfillment(ctx context.Context, q *CancelFulfillmentEndpoint) (*CancelFulfillmentEndpoint, error) {
	cmd := &shipping.CancelFulfillmentCommand{
		FulfillmentID: q.FulfillmentID,
		CancelReason:  q.CancelReason,
	}
	if err := s.ShippingAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &CancelFulfillmentEndpoint{
		Result: &pbcm.UpdatedResponse{Updated: 1},
	}, nil
}
