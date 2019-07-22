package shipnow

import (
	"context"

	"etop.vn/api/main/address"
	etoptypes "etop.vn/api/main/etop"
	"etop.vn/api/main/identity"
	"etop.vn/api/main/location"
	"etop.vn/api/main/ordering"
	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/main/shipnow"
	"etop.vn/api/main/shipnow/carrier"
	shipnowtypes "etop.vn/api/main/shipnow/types"
	shippingtypes "etop.vn/api/main/shipping/types"
	"etop.vn/api/meta"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	shipnowconvert "etop.vn/backend/pkg/services/shipnow/convert"
	"etop.vn/backend/pkg/services/shipnow/sqlstore"
	"etop.vn/common/bus"
)

var _ shipnow.Aggregate = &Aggregate{}

type Aggregate struct {
	location      location.QueryBus
	identityQuery identity.QueryBus
	addressQuery  address.QueryBus
	order         ordering.QueryBus

	db             cmsql.Transactioner
	store          sqlstore.ShipnowStoreFactory
	eventBus       meta.EventBus
	carrierManager carrier.Manager
}

func NewAggregate(eventBus meta.EventBus, db cmsql.Database, location location.QueryBus, identityQuery identity.QueryBus, addressQuery address.QueryBus, order ordering.QueryBus, carrierManager carrier.Manager) *Aggregate {
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

func (a *Aggregate) MessageBus() shipnow.CommandBus {
	b := bus.New()
	return shipnow.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreateShipnowFulfillment(ctx context.Context, cmd *shipnow.CreateShipnowFulfillmentArgs) (_result *shipnow.ShipnowFulfillment, _ error) {
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		ffmID := cm.NewID()
		// ShipnowOrderReservationEvent
		event := &shipnow.ShipnowOrderReservationEvent{
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
			RequestPickupAt:     nil,
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
		if ffm.ConfirmStatus != etoptypes.S3Zero || ffm.ShippingCode != "" {
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
		case etoptypes.S5Positive, etoptypes.S5Negative, etoptypes.S5NegSuper:
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn vận chuyển không thể hủy")
		}

		switch ffm.ShippingState {
		case shipnowtypes.StateCancelled:
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn vận chuyển đã bị hủy")
		case shipnowtypes.StateDelivering:
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn vận chuyển đang giao. Không thể hủy đơn.")
		case shipnowtypes.StateDelivered,
			shipnowtypes.StateReturning, shipnowtypes.StateReturned:
			return cm.Errorf(cm.FailedPrecondition, nil, "Không thể hủy đơn.")
		}

		event := &shipnow.ShipnowCancelledEvent{
			ShipnowFulfillmentId: ffm.Id,
			OrderIds:             ffm.OrderIds,
			CancelReason:         cmd.CancelReason,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}

		updateArgs := sqlstore.UpdateCancelArgs{
			ID:            ffm.Id,
			ShippingState: shipnowtypes.StateCancelled,
			Status:        etoptypes.S5Negative,
			ConfirmStatus: etoptypes.S3Negative,
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
			ShipnowFulfillmentId: ffm.Id,
			OrderIds:             ffm.OrderIds,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}

		event2 := &shipnow.ShipnowCreateExternalEvent{
			ShipnowFulfillmentId: ffm.Id,
		}
		if err := a.eventBus.Publish(ctx, event2); err != nil {
			return err
		}

		update := sqlstore.UpdateStateArgs{
			ID:             cmd.Id,
			ConfirmStatus:  etoptypes.S3Positive,
			ShippingStatus: etoptypes.S5SuperPos,
			Status:         etoptypes.S5SuperPos,
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

func (a *Aggregate) PreparePickupAddress(ctx context.Context, shopID int64, pickupAddress *ordertypes.Address) (*ordertypes.Address, error) {
	if pickupAddress != nil {
		return pickupAddress, nil
	}
	query := &identity.GetShopByIDQuery{ID: shopID}
	if err := a.identityQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	shop := query.Result.Shop
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

func (a *Aggregate) PrepareDeliveryPoints(ctx context.Context, orderIDs []int64) (points []*shipnow.DeliveryPoint, weightInfo shippingtypes.WeightInfo, valueinfo shippingtypes.ValueInfo, _err error) {
	query := &ordering.GetOrdersQuery{
		IDs: orderIDs,
	}
	if err := a.order.Dispatch(ctx, query); err != nil {
		_err = err
		return
	}
	orders := query.Result.Orders

	// Note: Không thay đổi thứ tự đơn hàng vì nó ảnh hưởng tới giá
	mapOrders := make(map[int64]*ordering.Order)
	for _, order := range orders {
		if order.Shipping == nil {
			_err = cm.Errorf(cm.InvalidArgument, nil, "Đơn hàng thiếu thông tin giao hàng: khối lượng, COD,...")
			return
		}
		mapOrders[order.ID] = order
	}
	for _, orderID := range orderIDs {
		order := mapOrders[orderID]
		points = append(points, shipnowconvert.OrderToDeliveryPoint(order))
	}
	weightInfo = shipnowconvert.GetWeightInfo(orders)
	valueinfo = shipnowconvert.GetValueInfo(orders)
	return
}

func ValidateConfirmFulfillment(ffm *shipnow.ShipnowFulfillment) error {
	switch ffm.ConfirmStatus {
	case etoptypes.S3Negative:
		return cm.Errorf(cm.FailedPrecondition, nil, "Đơn giao hàng đã hủy")
	case etoptypes.S3Positive:
		return cm.Errorf(cm.FailedPrecondition, nil, "Đơn giao hàng đã xác nhận")
	}
	if ffm.Status == etoptypes.S5Negative || ffm.Status == etoptypes.S5Positive {
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
		ShippingCreatedAt:   args.ShippingCreatedAt.ToTime(),
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
