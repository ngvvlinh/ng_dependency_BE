package shipping

import (
	"context"
	"fmt"

	"o.o/api/main/connectioning"
	"o.o/api/main/shipnow"
	shipnowtypes "o.o/api/main/shipnow/types"
	shippingtypes "o.o/api/main/shipping/types"
	typesx "o.o/api/top/external/types"
	typesint "o.o/api/top/int/types"
	"o.o/api/top/types/etc/inventory_auto"
	pbsource "o.o/api/top/types/etc/order_source"
	"o.o/api/top/types/etc/payment_method"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	convertx "o.o/backend/pkg/etop/apix/convertpb"
	"o.o/backend/pkg/etop/authorize/claims"
	"o.o/capi/dot"
	"o.o/common/l"
)

func (s *Shipping) GetShipnowServices(ctx context.Context, accountID dot.ID, r *typesx.GetShipnowServicesRequest) (*typesx.GetShipnowServicesResponse, error) {
	if r.PickupAddress == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp địa chỉ lấy hàng")
	}
	if len(r.DeliveryPoints) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp địa chỉ giao hàng")
	}
	var points []*shipnow.DeliveryPoint
	for _, p := range r.DeliveryPoints {
		point := &shipnow.DeliveryPoint{
			ShippingAddress: convertx.Convert_apix_ShipnowLocationAddressShortVersion_To_core_OrderAddress(p.ShippingAddress),
			ValueInfo: shippingtypes.ValueInfo{
				CODAmount: p.CODAmount.Int(),
			},
		}
		points = append(points, point)
	}
	cmd := &shipnow.GetShipnowServicesCommand{
		ShopId:         accountID,
		PickupAddress:  convertx.Convert_apix_ShipnowLocationAddressShortVersion_To_core_OrderAddress(r.PickupAddress),
		DeliveryPoints: points,
	}
	if err := s.ShipnowAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	services := cmd.Result.Services
	if err := s.buildCodeForShipnowServices(ctx, services); err != nil {
		return nil, err
	}
	res := &typesx.GetShipnowServicesResponse{
		Services: convertx.Convert_core_ShipnowServices_To_apix_ShipnowServices(services),
	}
	return res, nil
}

func (s *Shipping) buildCodeForShipnowServices(ctx context.Context, services []*shipnowtypes.ShipnowService) error {
	// add connection code to service code to identify which connects
	// code format: XXXXYYYYYYYY (12 characters)
	for _, srv := range services {
		if srv.ConnectionInfo == nil {
			continue
		}
		connection, err := s.ShipmentManager.ConnectionManager.GetConnectionByID(ctx, srv.ConnectionInfo.ID)
		if err != nil {
			return err
		}
		srv.Code = connection.Code + srv.Code
	}
	return nil
}

func (s *Shipping) CreateShipnowFulfillment(ctx context.Context, shopClaim *claims.ShopClaim, r *typesx.CreateShipnowFulfillmentRequest) (_ *typesx.ShipnowFulfillment, _err error) {
	var partner *identitymodel.Partner
	if shopClaim.AuthPartnerID != 0 {
		queryPartner := &identitymodelx.GetPartner{
			PartnerID: shopClaim.AuthPartnerID,
		}
		if err := bus.Dispatch(ctx, queryPartner); err != nil {
			return nil, err
		}
		partner = queryPartner.Result.Partner
	}

	conn, serviceCode, err := s.parseServiceCode(ctx, r.ShippingServiceCode)
	if err != nil {
		return nil, err
	}

	if len(r.DeliveryPoints) > 10 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Không thể có quá 10 điểm giao hàng")
	}

	var deliveryPoints []*shipnow.OrderShippingInfo
	// Cần giữ nguyên thứ tự các điểm giao hàng
	for _, point := range r.DeliveryPoints {
		lines, err := convertx.OrderLinesToCreateOrderLines(point.Lines)
		if err != nil {
			return nil, err
		}
		totalItems := 0
		for _, line := range lines {
			totalItems += line.Quantity
		}

		shippingAddress := point.ShippingAddress
		orderAddress := convertx.Convert_apix_ShipnowAddress_To_api_OrderAddress(shippingAddress)
		weight := cm.CoalesceInt(point.ChargeableWeight.Int(), point.GrossWeight.Int())
		args := &typesint.CreateOrderRequest{
			Source:        pbsource.API,
			PaymentMethod: payment_method.Other,
			Customer: &typesint.OrderCustomer{
				FullName: shippingAddress.FullName,
				Email:    shippingAddress.Email,
				Phone:    shippingAddress.Phone,
			},
			CustomerAddress: orderAddress,
			BillingAddress:  orderAddress,
			ShippingAddress: orderAddress,
			Lines:           lines,
			TotalItems:      totalItems,
			BasketValue:     point.BasketValue.Int(),
			TotalWeight:     weight,
			TotalAmount:     point.BasketValue.Int(),
			OrderNote:       point.ShippingNote,
			ShippingNote:    point.ShippingNote,
		}
		respOrder, err := s.OrderLogic.CreateOrder(ctx, shopClaim.Shop, partner, args, nil, 0)
		if err != nil {
			return nil, err
		}

		deliveryPoint := &shipnow.OrderShippingInfo{
			OrderID:         respOrder.Id,
			ShippingAddress: convertx.Convert_apix_ShipnowAddress_To_core_OrderAddress(shippingAddress),
			ShippingNote:    point.ShippingNote,
			WeightInfo: shippingtypes.WeightInfo{
				GrossWeight:      point.GrossWeight.Int(),
				ChargeableWeight: weight,
			},
			ValueInfo: shippingtypes.ValueInfo{
				BasketValue: point.BasketValue.Int(),
				CODAmount:   point.CODAmount.Int(),
			},
			TryOn: 0,
		}
		deliveryPoints = append(deliveryPoints, deliveryPoint)
	}

	var shipnowFulfillment *shipnow.ShipnowFulfillment
	defer func() {
		if _err == nil {
			return
		}
		// cancel shipnow ffm if exist
		if shipnowFulfillment != nil {
			cmd := &shipnow.CancelShipnowFulfillmentCommand{
				ID:           shipnowFulfillment.ID,
				ShopID:       shopClaim.AccountID,
				CancelReason: fmt.Sprintf("Tạo đơn shipnow không thành công: %v", _err.Error()),
			}
			s.ShipnowAggr.Dispatch(ctx, cmd)
		}

		// always cancel order if cannot create shipnow ffm
		for _, point := range deliveryPoints {
			_, err = s.OrderLogic.CancelOrder(ctx, shopClaim.UserID, shopClaim.AccountID, shopClaim.AuthPartnerID, point.OrderID, fmt.Sprintf("Tạo đơn shipnow không thành công: %v", _err.Error()), inventory_auto.Unknown)
		}
		if err != nil {
			ll.Error("cancelling order", l.Error(err))
		}
		return
	}()

	// Create shipnow fulfillment
	createCmd := &shipnow.CreateShipnowFulfillmentCommand{
		DeliveryPoints:      deliveryPoints,
		ShopID:              shopClaim.AccountID,
		ShippingServiceCode: serviceCode,
		ShippingServiceFee:  r.ShippingServiceFee.Int(),
		ShippingNote:        r.ShippingNote,
		PickupAddress:       convertx.Convert_apix_ShipnowAddress_To_core_OrderAddress(r.PickupAddress),
		ConnectionID:        conn.ID,
		ExternalID:          r.ExternalID,
	}
	if err := s.ShipnowAggr.Dispatch(ctx, createCmd); err != nil {
		return nil, err
	}
	shipnowFulfillment = createCmd.Result
	// confirm shipnow fulfillment
	confirmCmd := &shipnow.ConfirmShipnowFulfillmentCommand{
		ID:     shipnowFulfillment.ID,
		ShopID: shopClaim.AccountID,
	}
	if err := s.ShipnowAggr.Dispatch(ctx, confirmCmd); err != nil {
		return nil, err
	}

	return s.GetShipnowFulfillment(ctx, shopClaim.AccountID, &typesx.FulfillmentIDRequest{
		Id: confirmCmd.ID,
	})
}

func (s *Shipping) CancelShipnowFulfillment(ctx context.Context, accountID dot.ID, r *typesx.CancelShipnowFulfillmentRequest) error {
	cmd := &shipnow.CancelShipnowFulfillmentCommand{
		ID:           r.ID,
		ShippingCode: r.ShippingCode,
		ShopID:       accountID,
		CancelReason: r.CancelReason,
	}
	return s.ShipnowAggr.Dispatch(ctx, cmd)
}

func (s *Shipping) GetShipnowFulfillment(ctx context.Context, accountID dot.ID, r *typesx.FulfillmentIDRequest) (*typesx.ShipnowFulfillment, error) {
	query := &shipnow.GetShipnowFulfillmentQuery{
		ID:           r.Id,
		ShippingCode: r.ShippingCode,
		ShopID:       accountID,
	}
	if err := s.ShipnowQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	shipnowFfm := query.Result.ShipnowFulfillment
	var conn *connectioning.Connection
	if shipnowFfm.ConnectionID != 0 {
		conn, _ = s.ShipmentManager.ConnectionManager.GetConnectionByID(ctx, shipnowFfm.ConnectionID)
	}
	return convertx.Convert_core_ShipnowFulfillment_To_apix_ShipnowFulfillment(query.Result.ShipnowFulfillment, conn), nil
}
