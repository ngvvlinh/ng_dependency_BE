package shipnow

import (
	"context"
	"time"

	"o.o/api/main/address"
	"o.o/api/main/identity"
	"o.o/api/main/location"
	"o.o/api/main/ordering"
	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/main/shipnow"
	"o.o/api/main/shipnow/carrier"
	shipnowtypes "o.o/api/main/shipnow/types"
	shippingtypes "o.o/api/main/shipping/types"
	"o.o/api/meta"
	"o.o/api/top/types/etc/shipnow_state"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status5"
	"o.o/backend/com/main/shipnow/convert"
	"o.o/backend/com/main/shipnow/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
	"o.o/capi/dot"
)

var _ shipnow.Aggregate = &Aggregate{}

type Aggregate struct {
	location      location.QueryBus
	identityQuery identity.QueryBus
	addressQuery  address.QueryBus
	order         ordering.QueryBus

	db             cmsql.Transactioner
	store          sqlstore.ShipnowStoreFactory
	eventBus       capi.EventBus
	carrierManager carrier.Manager
}

func NewAggregate(eventBus capi.EventBus, db *cmsql.Database, location location.QueryBus, identityQuery identity.QueryBus, addressQuery address.QueryBus, order ordering.QueryBus, carrierManager carrier.Manager) *Aggregate {
	return &Aggregate{
		db:       db,
		store:    sqlstore.NewShipnowStore(db),
		eventBus: eventBus,

		location:      location,
		identityQuery: identityQuery,
		addressQuery:  addressQuery,
		order:         order,

		carrierManager: carrierManager,
	}
}

func AggregateMessageBus(a *Aggregate) shipnow.CommandBus {
	b := bus.New()
	return shipnow.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreateShipnowFulfillment(ctx context.Context, cmd *shipnow.CreateShipnowFulfillmentArgs) (_result *shipnow.ShipnowFulfillment, _ error) {
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		ffmID := cm.NewID()
		// ShipnowOrderReservationEvent
		event := &shipnow.ShipnowOrderReservationEvent{
			EventMeta:            meta.NewEvent(),
			OrderIds:             cmd.OrderIds,
			ShipnowFulfillmentId: ffmID,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}

		pickupAddress, err := a.PreparePickupAddress(ctx, cmd.ShopId, cmd.PickupAddress)
		if err != nil {
			return err
		}

		points, weightInfo, valueInfo, err := a.PrepareDeliveryPoints(ctx, cmd.OrderIds)
		if err != nil {
			return err
		}
		shipnowFfm := &shipnow.ShipnowFulfillment{
			Id:                  ffmID,
			ShopId:              cmd.ShopId,
			PickupAddress:       pickupAddress,
			DeliveryPoints:      points,
			Carrier:             cmd.Carrier,
			ShippingServiceCode: cmd.ShippingServiceCode,
			ShippingServiceFee:  cmd.ShippingServiceFee,
			WeightInfo:          weightInfo,
			ValueInfo:           valueInfo,
			ShippingNote:        cmd.ShippingNote,
			RequestPickupAt:     time.Time{},
		}

		if err := a.store(ctx).Create(shipnowFfm); err != nil {
			return err
		}
		_result = shipnowFfm
		return nil
	})
	return _result, err
}

func (a *Aggregate) CreateShipnowFulfillmentV2(ctx context.Context, args *shipnow.CreateShipnowFulfillmentV2Args) (_result *shipnow.ShipnowFulfillment, _ error) {
	orderIDs := make([]dot.ID, len(args.DeliveryPoints))
	for i, point := range args.DeliveryPoints {
		if cm.IDsContain(orderIDs, point.OrderID) {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Một đơn hàng không được chọn nhiều lần.")
		}
		orderIDs[i] = point.OrderID
	}
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		ffmID := cm.NewID()
		// ShipnowOrderReservationEvent
		event := &shipnow.ShipnowOrderReservationEvent{
			EventMeta:            meta.NewEvent(),
			OrderIds:             orderIDs,
			ShipnowFulfillmentId: ffmID,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}
		pickupAddress, err := a.PreparePickupAddress(ctx, args.ShopID, args.PickupAddress)
		if err != nil {
			return err
		}
		points, weightInfo, valueInfo, err := a.PrepareDeliveryPointsV2(ctx, args.DeliveryPoints)
		if err != nil {
			return err
		}
		shipnowFfm := &shipnow.ShipnowFulfillment{
			Id:                  ffmID,
			ShopId:              args.ShopID,
			PickupAddress:       pickupAddress,
			DeliveryPoints:      points,
			Carrier:             args.Carrier,
			ShippingServiceCode: args.ShippingServiceCode,
			ShippingServiceFee:  args.ShippingServiceFee,
			WeightInfo:          weightInfo,
			ValueInfo:           valueInfo,
			ShippingNote:        args.ShippingNote,
		}

		if err := a.store(ctx).Create(shipnowFfm); err != nil {
			return err
		}
		_result = shipnowFfm
		return nil
	})
	return _result, err
}

func (a *Aggregate) UpdateShipnowFulfillment(ctx context.Context, cmd *shipnow.UpdateShipnowFulfillmentArgs) (_result *shipnow.ShipnowFulfillment, _ error) {
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		ffm, err := a.store(ctx).ID(cmd.Id).ShopID(cmd.ShopId).GetShipnow()
		if err != nil {
			return err
		}
		if ffm.ConfirmStatus != status3.Z || ffm.ShippingCode != "" {
			return cm.Errorf(cm.FailedPrecondition, nil, "Không thể cập nhật đơn giao hàng này.")
		}

		updateArgs := sqlstore.UpdateInfoArgs{
			ID:                  cmd.Id,
			PickupAddress:       cmd.PickupAddress,
			Carrier:             cmd.Carrier,
			ShippingServiceCode: cmd.ShippingServiceCode,
			ShippingServiceFee:  cmd.ShippingServiceFee,
			ShippingNote:        cmd.ShippingNote,
			RequestPickupAt:     cmd.RequestPickupAt,
		}

		if len(cmd.OrderIds) > 0 {
			// ShipnowOrderChangedEvent
			event := &shipnow.ShipnowOrderChangedEvent{
				EventMeta:            meta.NewEvent(),
				ShipnowFulfillmentId: ffm.Id,
				OldOrderIds:          ffm.OrderIds,
				OrderIds:             cmd.OrderIds,
			}
			if err := a.eventBus.Publish(ctx, event); err != nil {
				return nil
			}
			points, weightInfo, valueInfo, err := a.PrepareDeliveryPoints(ctx, cmd.OrderIds)
			if err != nil {
				return err
			}
			updateArgs.DeliveryPoints = points
			updateArgs.WeightInfo = weightInfo
			updateArgs.ValueInfo = valueInfo
		}

		result, err := a.store(ctx).UpdateInfo(updateArgs)
		if err != nil {
			return err
		}
		_result = result
		return nil
	})
	return _result, err
}

func (a *Aggregate) CancelShipnowFulfillment(ctx context.Context, cmd *shipnow.CancelShipnowFulfillmentArgs) (*meta.Empty, error) {
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		ffm, err := a.store(ctx).ID(cmd.Id).ShopID(cmd.ShopId).GetShipnow()
		if err != nil {
			return err
		}

		switch ffm.Status {
		case status5.P, status5.N, status5.NS:
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn vận chuyển không thể hủy")
		}

		switch ffm.ShippingState {
		case shipnow_state.StateCancelled:
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn vận chuyển đã bị hủy")
		case shipnow_state.StateDelivering:
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn vận chuyển đang giao. Không thể hủy đơn.")
		case shipnow_state.StateDelivered,
			shipnow_state.StateReturning, shipnow_state.StateReturned:
			return cm.Errorf(cm.FailedPrecondition, nil, "Không thể hủy đơn.")
		}

		event := &shipnow.ShipnowCancelledEvent{
			EventMeta:            meta.NewEvent(),
			ShipnowFulfillmentId: ffm.Id,
			OrderIds:             ffm.OrderIds,
			CancelReason:         cmd.CancelReason,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}

		updateArgs := sqlstore.UpdateCancelArgs{
			ID:            ffm.Id,
			ShippingState: shipnow_state.StateCancelled,
			Status:        status5.N,
			ConfirmStatus: status3.N,
			CancelReason:  cmd.CancelReason,
		}
		ffm, err = a.store(ctx).UpdateCancelled(updateArgs)
		if err != nil {
			return err
		}
		return nil
	})
	return &meta.Empty{}, err
}

func (a *Aggregate) ConfirmShipnowFulfillment(ctx context.Context, cmd *shipnow.ConfirmShipnowFulfillmentArgs) (_result *shipnow.ShipnowFulfillment, _err error) {
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		ffm, err := a.store(ctx).ID(cmd.Id).ShopID(cmd.ShopId).GetShipnow()
		if err != nil {
			return err
		}
		if err := ValidateConfirmFulfillment(ffm); err != nil {
			return err
		}

		event := &shipnow.ShipnowValidateConfirmedEvent{
			EventMeta:            meta.NewEvent(),
			ShipnowFulfillmentId: ffm.Id,
			OrderIds:             ffm.OrderIds,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}

		event2 := &shipnow.ShipnowExternalCreatedEvent{
			ShipnowFulfillmentId: ffm.Id,
		}
		if err := a.eventBus.Publish(ctx, event2); err != nil {
			return err
		}

		update := sqlstore.UpdateStateArgs{
			ID:             cmd.Id,
			ConfirmStatus:  status3.P,
			ShippingStatus: status5.S,
			Status:         status5.S,
		}
		shipnowFfm, err := a.store(ctx).UpdateSyncState(update)
		if err != nil {
			return err
		}
		_result = shipnowFfm
		return nil
	})
	return _result, err
}

func (a *Aggregate) PreparePickupAddress(ctx context.Context, shopID dot.ID, pickupAddress *ordertypes.Address) (*ordertypes.Address, error) {
	if pickupAddress != nil {
		return pickupAddress, nil
	}
	query := &identity.GetShopByIDQuery{ID: shopID}
	if err := a.identityQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	shop := query.Result
	shopAddressID := shop.ShipFromAddressID
	if shopAddressID == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Bán hàng: Cần cung cấp thông tin địa chỉ lấy hàng trong đơn hàng hoặc tại thông tin cửa hàng. Vui lòng cập nhật.", nil)
	}
	query2 := &address.GetAddressByIDQuery{
		ID: shopAddressID,
	}
	if err := a.addressQuery.Dispatch(ctx, query2); err != nil {
		return nil, err
	}
	shopAddress := query2.Result
	pickupAddress = shopAddress.ToOrderAddress()
	return pickupAddress, nil
}

func (a *Aggregate) PrepareDeliveryPoints(ctx context.Context, orderIDs []dot.ID) (points []*shipnow.DeliveryPoint, weightInfo shippingtypes.WeightInfo, valueinfo shippingtypes.ValueInfo, _err error) {
	query := &ordering.GetOrdersQuery{
		IDs: orderIDs,
	}
	if err := a.order.Dispatch(ctx, query); err != nil {
		_err = err
		return
	}
	orders := query.Result.Orders

	// Note: Không thay đổi thứ tự đơn hàng vì nó ảnh hưởng tới giá
	mapOrders := make(map[dot.ID]*ordering.Order)
	for _, order := range orders {
		if order.Shipping == nil {
			_err = cm.Errorf(cm.InvalidArgument, nil, "Đơn hàng thiếu thông tin giao hàng: khối lượng, COD,...")
			return
		}
		mapOrders[order.ID] = order
	}
	for _, orderID := range orderIDs {
		order := mapOrders[orderID]
		points = append(points, convert.OrderToDeliveryPoint(order))
	}
	weightInfo = convert.GetWeightInfo(orders)
	valueinfo = convert.GetValueInfo(orders)
	return
}

func (a *Aggregate) PrepareDeliveryPointsV2(ctx context.Context, points []*shipnow.OrderShippingInfo) (deliveryPoints []*shipnow.DeliveryPoint, weightInfo shippingtypes.WeightInfo, valueInfo shippingtypes.ValueInfo, _err error) {
	orderIDs := make([]dot.ID, len(points))
	for i, point := range points {
		orderIDs[i] = point.OrderID
	}
	query := &ordering.GetOrdersQuery{
		IDs: orderIDs,
	}
	if err := a.order.Dispatch(ctx, query); err != nil {
		_err = err
		return
	}
	orders := query.Result.Orders

	// Note: Không thay đổi thứ tự đơn hàng vì nó ảnh hưởng tới giá
	mapOrders := make(map[dot.ID]*ordering.Order)
	for _, order := range orders {
		mapOrders[order.ID] = order
	}

	for _, point := range points {
		if point.ShippingAddress == nil {
			_err = cm.Errorf(cm.InvalidArgument, nil, "Vui lòng cung cấp địa chỉ giao hàng")
		}
		order := mapOrders[point.OrderID]
		p := &shipnowtypes.DeliveryPoint{
			ShippingAddress: point.ShippingAddress,
			Lines:           order.Lines,
			ShippingNote:    cm.Coalesce(point.ShippingNote, order.ShippingNote),
			OrderId:         order.ID,
			OrderCode:       order.Code,
			WeightInfo: shippingtypes.WeightInfo{
				GrossWeight:      point.GrossWeight,
				ChargeableWeight: point.ChargeableWeight,
				Length:           point.Length,
				Width:            point.Width,
				Height:           point.Height,
			},
			ValueInfo: shippingtypes.ValueInfo{
				BasketValue:      cm.CoalesceInt(point.BasketValue, order.BasketValue),
				CODAmount:        point.CODAmount,
				IncludeInsurance: point.IncludeInsurance,
			},
			TryOn: point.TryOn,
		}

		deliveryPoints = append(deliveryPoints, p)
		weightInfo.ChargeableWeight += point.ChargeableWeight
		weightInfo.GrossWeight += point.GrossWeight
		valueInfo.BasketValue += cm.CoalesceInt(point.BasketValue, order.BasketValue)
		valueInfo.CODAmount += point.CODAmount
	}
	return
}

func ValidateConfirmFulfillment(ffm *shipnow.ShipnowFulfillment) error {
	switch ffm.ConfirmStatus {
	case status3.N:
		return cm.Errorf(cm.FailedPrecondition, nil, "Đơn giao hàng đã hủy")
	case status3.P:
		return cm.Errorf(cm.FailedPrecondition, nil, "Đơn giao hàng đã xác nhận")
	}
	if ffm.Status == status5.N || ffm.Status == status5.P {
		return cm.Errorf(cm.FailedPrecondition, nil, "Không thể xác nhận đơn giao hàng này")
	}

	if len(ffm.DeliveryPoints) == 0 || len(ffm.OrderIds) == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "Số điểm giao hàng không hợp lệ")
	}
	return nil
}

func (a *Aggregate) UpdateShipnowFulfillmentCarrierInfo(ctx context.Context, args *shipnow.UpdateShipnowFulfillmentCarrierInfoArgs) (*shipnow.ShipnowFulfillment, error) {
	updateArgs := sqlstore.UpdateCarrierInfoArgs{
		ID:                  args.Id,
		FeeLines:            args.FeeLines,
		CarrierFeeLines:     args.CarrierFeeLines,
		ShippingCode:        args.ShippingCode,
		ShippingCreatedAt:   args.ShippingCreatedAt,
		ShippingState:       args.ShippingState,
		ShippingStatus:      args.ShippingStatus,
		EtopPaymentStatus:   args.EtopPaymentStatus,
		CODEtopTransferedAt: args.CodEtopTransferedAt,
		Status:              args.Status,

		ShippingPickingAt:          args.ShippingPickingAt,
		ShippingDeliveringAt:       args.ShippingDeliveringAt,
		ShippingDeliveredAt:        args.ShippingDeliveredAt,
		ShippingCancelledAt:        args.ShippingCancelledAt,
		ShippingServiceName:        args.ShippingServiceName,
		ShippingServiceDescription: args.ShippingServiceDescription,
		CancelReason:               args.CancelReason,
		ShippingSharedLink:         args.ShippingSharedLink,
	}
	updateArgs.TotalFee = shippingtypes.TotalFee(args.FeeLines)
	ffm, err := a.store(ctx).UpdateCarrierInfo(updateArgs)
	return ffm, err
}

func (a *Aggregate) UpdateShipnowFulfillmentState(ctx context.Context, args *shipnow.UpdateShipnowFulfillmentStateArgs) (*shipnow.ShipnowFulfillment, error) {
	updateArgs := sqlstore.UpdateStateArgs{
		ID:             args.Id,
		SyncStatus:     args.SyncStatus,
		Status:         args.Status,
		ShippingState:  args.ShippingState,
		SyncStates:     args.SyncStates,
		ConfirmStatus:  args.ConfirmStatus,
		ShippingStatus: args.ShippingStatus,
	}
	ffm, err := a.store(ctx).UpdateSyncState(updateArgs)
	return ffm, err
}

func (a *Aggregate) GetShipnowServices(ctx context.Context, args *shipnow.GetShipnowServicesArgs) (*shipnow.GetShipnowServicesResult, error) {
	if len(args.OrderIds) == 0 && len(args.DeliveryPoints) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng cung cấp địa chỉ giao hàng")
	}

	pickupAddress, err := a.PreparePickupAddress(ctx, args.ShopId, args.PickupAddress)
	if err != nil {
		return nil, err
	}

	var points = args.DeliveryPoints
	if len(args.OrderIds) > 0 {
		points, _, _, err = a.PrepareDeliveryPoints(ctx, args.OrderIds)
		if err != nil {
			return nil, err
		}
	}

	cmd := &carrier.GetExternalShipnowServicesCommand{
		ShopID:         args.ShopId,
		PickupAddress:  pickupAddress,
		DeliveryPoints: points,
	}

	services, err := a.carrierManager.GetExternalShippingServices(ctx, cmd)
	return &shipnow.GetShipnowServicesResult{
		Services: services,
	}, nil
}
