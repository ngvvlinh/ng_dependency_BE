package shipnow

import (
	"context"

	shippingtypes "etop.vn/api/main/shipping/types"

	"etop.vn/api/main/ordering"

	"etop.vn/api/main/address"
	etoptypes "etop.vn/api/main/etop"
	"etop.vn/api/main/identity"
	"etop.vn/api/main/location"
	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/main/shipnow"
	shipnowtypes "etop.vn/api/main/shipnow/types"
	"etop.vn/api/meta"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	shipnowconvert "etop.vn/backend/pkg/services/shipnow/convert"
	shipnowmodelx "etop.vn/backend/pkg/services/shipnow/modelx"
	"etop.vn/backend/pkg/services/shipnow/sqlstore"
)

var _ shipnow.Aggregate = &Aggregate{}

type Aggregate struct {
	location      location.QueryBus
	identityQuery identity.QueryService
	addressQuery  address.QueryService
	order         ordering.QueryBus

	db       cmsql.Transactioner
	store    sqlstore.ShipnowStoreFactory
	eventBus meta.EventBus
}

func NewAggregate(eventBus meta.EventBus, db cmsql.Database, location location.QueryBus, identityQuery identity.QueryService, addressQuery address.QueryService, order ordering.QueryBus) *Aggregate {
	return &Aggregate{
		db:       db,
		store:    sqlstore.NewShipnowStore(db),
		eventBus: eventBus,

		location:      location,
		identityQuery: identityQuery,
		addressQuery:  addressQuery,
		order:         order,
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
		args := shipnowmodelx.GetByIDArgs{
			ID:     cmd.Id,
			ShopID: cmd.ShopId,
		}
		ffm, err := a.store(ctx).GetByID(args)
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
		args := shipnowmodelx.GetByIDArgs{
			ID:     cmd.Id,
			ShopID: cmd.ShopId,
		}
		ffm, err := a.store(ctx).GetByID(args)
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
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}

		// if err := m.carrierManager.CancelExternalShipping(ctx, nil); err != nil {
		// 	return err
		// }

		updateArgs := sqlstore.UpdateStateArgs{
			ID: ffm.Id,
			// SyncStatus: etoptypes.S4Negative,
			State:         shipnowtypes.StateCancelled,
			Status:        etoptypes.S5Negative,
			ConfirmStatus: etoptypes.S3Negative,
			// SyncStates: &model.FulfillmentSyncStates{
			// 	TrySyncAt:         time.Now(),
			// 	NextShippingState: model.StateCreated,
			// },
		}
		ffm, err = a.store(ctx).UpdateSyncState(updateArgs)
		if err != nil {
			return err
		}
		return nil
	})
	return &meta.Empty{}, err
}

func (a *Aggregate) ConfirmShipnowFulfillment(ctx context.Context, cmd *shipnow.ConfirmShipnowFulfillmentArgs) (_result *shipnow.ShipnowFulfillment, _ error) {
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		query := shipnowmodelx.GetByIDArgs{
			ID:     cmd.Id,
			ShopID: cmd.ShopId,
		}
		ffm, err := a.store(ctx).GetByID(query)
		if err != nil {
			return err
		}
		if err := ValidateConfirmFulfillment(ffm); err != nil {
			return err
		}

		event := &shipnow.ShipnowValidatedEvent{
			ShipnowFulfillmentId: ffm.Id,
			OrderIds:             ffm.OrderIds,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}

		update := sqlstore.UpdateStateArgs{
			ID:            cmd.Id,
			ConfirmStatus: etoptypes.S3Positive,
		}
		shipnowFfm, err := a.store(ctx).UpdateSyncState(update)
		if err != nil {
			return err
		}
		_result = shipnowFfm

		// if err := a.shipnowManagerCtrl.CreateExternalShipping(ctx, ffm); err != nil {
		// 	return &meta.Empty{}, err
		// }

		return nil
	})
	return _result, err
}

func (a *Aggregate) PreparePickupAddress(ctx context.Context, shopID int64, pickupAddress *ordertypes.Address) (*ordertypes.Address, error) {
	if pickupAddress != nil {
		return pickupAddress, nil
	}
	shopResult, err := a.identityQuery.GetShopByID(ctx, &identity.GetShopByIDQueryArgs{ID: shopID})
	if err != nil {
		return nil, err
	}
	shop := shopResult.Shop
	shopAddressID := shop.ShipFromAddressID
	if shopAddressID == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Bán hàng: Cần cung cấp thông tin địa chỉ lấy hàng trong đơn hàng hoặc tại thông tin cửa hàng. Vui lòng cập nhật.", nil)
	}
	shopAddress, err := a.addressQuery.GetAddressByID(ctx, &address.GetAddressByIDQueryArgs{ID: shopAddressID})
	if err != nil {
		return nil, err
	}
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
	for _, order := range orders {
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
		return cm.Errorf(cm.FailedPrecondition, nil, "Số điểm giao hàng không hợp lệ.")
	}
	return nil
}
