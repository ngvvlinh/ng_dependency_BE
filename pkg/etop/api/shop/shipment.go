package shop

import (
	"context"
	"fmt"
	"time"

	"o.o/api/main/shipping"
	shippingtypes "o.o/api/main/shipping/types"
	api "o.o/api/top/int/shop"
	inttypes "o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
	shipmodelx "o.o/backend/com/main/shipping/modelx"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

type ShipmentService struct {
	session.Session

	ShipmentManager   *shippingcarrier.ShipmentManager
	ShippingAggregate shipping.CommandBus
}

func (s *ShipmentService) Clone() api.ShipmentService { res := *s; return &res }

func (s *ShipmentService) GetShippingServices(ctx context.Context, q *inttypes.GetShippingServicesRequest) (*inttypes.GetShippingServicesResponse, error) {
	shopID := s.SS.Shop().ID
	args, err := s.ShipmentManager.PrepareDataGetShippingServices(ctx, q)
	if err != nil {
		return nil, err
	}
	args.AccountID = shopID
	resp, err := s.ShipmentManager.GetShippingServices(ctx, args)
	if err != nil {
		return nil, err
	}
	result := &inttypes.GetShippingServicesResponse{
		Services: convertpb.PbAvailableShippingServices(resp),
	}
	return result, nil
}

func (s *ShipmentService) CreateFulfillments(ctx context.Context, q *api.CreateFulfillmentsRequest) (*api.CreateFulfillmentsResponse, error) {
	key := fmt.Sprintf("CreateFulfillments %v-%v", s.SS.Shop().ID, q.OrderID)
	res, _, err := idempgroup.DoAndWrap(
		ctx, key, 10*time.Second, "tạo đơn giao hàng",
		func() (interface{}, error) { return s.createFulfillments(ctx, q) })

	if err != nil {
		return nil, err
	}
	result := res.(*api.CreateFulfillmentsResponse)
	return result, nil
}

func (s *ShipmentService) createFulfillments(ctx context.Context, q *api.CreateFulfillmentsRequest) (_ *api.CreateFulfillmentsResponse, _err error) {
	shopID := s.SS.Shop().ID
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
			InsuranceValue:   q.InsuranceValue,
		},
		TryOn:         q.TryOn,
		ShippingNote:  q.ShippingNote,
		ShippingType:  q.ShippingType,
		ConnectionID:  q.ConnectionID,
		ShopCarrierID: q.ShopCarrierID,
		Coupon:        q.Coupon,
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
	res := &api.CreateFulfillmentsResponse{
		Fulfillments: ffms,
	}
	return res, nil
}

func (s *ShipmentService) UpdateFulfillmentInfo(ctx context.Context, q *api.UpdateFulfillmentInfoRequest) (res *pbcm.UpdatedResponse, _ error) {
	updateFulfillmentInfo := &shipping.ShopUpdateFulfillmentInfoCommand{
		FulfillmentID:    q.FulfillmentID,
		AddressTo:        convertpb.Convert_api_OrderAddress_To_core_OrderAddress(q.PickupAddress),
		AddressFrom:      convertpb.Convert_api_OrderAddress_To_core_OrderAddress(q.ShippingAddress),
		IncludeInsurance: q.IncludeInsurance,
		InsuranceValue:   q.InsuranceValue,
		GrossWeight:      q.GrossWeight,
		TryOn:            q.TryOn,
		ShippingNote:     q.ShippingNote,
	}
	if err := s.ShippingAggregate.Dispatch(ctx, updateFulfillmentInfo); err != nil {
		return nil, err
	}
	res = &pbcm.UpdatedResponse{Updated: updateFulfillmentInfo.Result}
	return res, nil
}

func (s *ShipmentService) UpdateFulfillmentCOD(ctx context.Context, q *api.UpdateFulfillmentCODRequest) (*pbcm.UpdatedResponse, error) {
	key := fmt.Sprintf("UpdateFulfillmentCOD %v-%v", s.SS.Shop().ID, q.FulfillmentID)
	res, _, err := idempgroup.DoAndWrap(
		ctx, key, 10*time.Second, "cập nhật COD",
		func() (interface{}, error) { return s.updateFulfillmentCOD(ctx, q) })

	if err != nil {
		return nil, err
	}
	return res.(*pbcm.UpdatedResponse), nil
}

func (s *ShipmentService) updateFulfillmentCOD(ctx context.Context, q *api.UpdateFulfillmentCODRequest) (*pbcm.UpdatedResponse, error) {
	updateFulfillmentCODCmd := &shipping.ShopUpdateFulfillmentCODCommand{
		FulfillmentID:  q.FulfillmentID,
		TotalCODAmount: dot.Int(q.CODAmount),
		UpdatedBy:      s.SS.User().ID,
	}
	if err := s.ShippingAggregate.Dispatch(ctx, updateFulfillmentCODCmd); err != nil {
		return nil, err
	}
	return &pbcm.UpdatedResponse{
		Updated: 1,
	}, nil
}

func (s *ShipmentService) CancelFulfillment(ctx context.Context, q *api.CancelFulfillmentRequest) (*pbcm.UpdatedResponse, error) {
	key := fmt.Sprintf("CancelFulfillment %v-%v", s.SS.Shop().ID, q.FulfillmentID)
	res, _, err := idempgroup.DoAndWrap(
		ctx, key, 10*time.Second, "huỷ đơn giao hàng",
		func() (interface{}, error) { return s.cancelFulfillment(ctx, q) })

	if err != nil {
		return nil, err
	}
	result := res.(*pbcm.UpdatedResponse)
	return result, nil
}

func (s *ShipmentService) cancelFulfillment(ctx context.Context, q *api.CancelFulfillmentRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &shipping.CancelFulfillmentCommand{
		FulfillmentID: q.FulfillmentID,
		CancelReason:  q.CancelReason,
	}
	if err := s.ShippingAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &pbcm.UpdatedResponse{Updated: 1}, nil
}
