package pm

import (
	"context"

	"etop.vn/api/main/address"
	etoptypes "etop.vn/api/main/etop"
	"etop.vn/api/main/identity"
	"etop.vn/api/main/ordering"
	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/main/shipnow"
	"etop.vn/api/main/shipnow/carrier"
	shipnowtypes "etop.vn/api/main/shipnow/types"
	"etop.vn/api/meta"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	shipnowconvert "etop.vn/backend/pkg/services/shipnow/convert"
)

type ProcessManager struct {
	eventBus meta.EventBus

	shipnowQuery shipnow.QueryBus

	orderAggr    ordering.Aggregate
	orderAggrBus ordering.CommandBus

	identityQuery  identity.QueryService
	addressQuery   address.QueryService
	carrierManager carrier.Manager
}

func New(
	eventBus meta.EventBus,
	shipnowBus shipnow.QueryBus,
	orderAggr ordering.Aggregate,
	orderAggrBus ordering.CommandBus,
	identityQuery identity.QueryService,
	addressQuery address.QueryService,
	carrierManager carrier.Manager,
) *ProcessManager {
	return &ProcessManager{
		eventBus:       eventBus,
		shipnowQuery:   shipnowBus,
		orderAggr:      orderAggr,
		orderAggrBus:   orderAggrBus,
		identityQuery:  identityQuery,
		addressQuery:   addressQuery,
		carrierManager: carrierManager,
	}
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.ShipnowOrderReservation)
	// TODO: add more
}

func (m *ProcessManager) ShipnowOrderReservation(ctx context.Context, event *shipnow.ShipnowOrderReservationEvent) error {
	// Call orderAggr for ReserveOrdersForFfm
	cmd := &ordering.ReserveOrdersForFfmCommand{
		OrderIDs:   event.OrderIds,
		Fulfill:    ordertypes.FulfillShipnowFulfillment,
		FulfillIDs: []int64{event.ShipnowFulfillmentId},
	}
	if err := m.orderAggrBus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	// args := &ordering.ReserveOrdersForFfmArgs{
	// 	OrderIDs: event.OrderIds,
	// }
	// result, err := m.orderAggr.ReserveOrdersForFfm(ctx, args)
	// if err != nil {
	// 	return nil, err
	// }

	return nil
}

func (m *ProcessManager) HandleShipnowCancellation(ctx context.Context, cmd *shipnow.CancelShipnowFulfillmentArgs) error {
	queryFfm := &shipnow.GetShipnowFulfillmentQuery{
		Id:     cmd.Id,
		ShopId: cmd.ShopId,
	}
	if err := m.shipnowQuery.Dispatch(ctx, queryFfm); err != nil {
		return err
	}
	ffm := queryFfm.Result.ShipnowFulfillment

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

	// if err := m.carrierManager.CancelExternalShipping(ctx, nil); err != nil {
	// 	return err
	// }

	// updateArgs := sqlstore.UpdateSyncStateArgs{
	// 	ID:         ffm.Id,
	// 	SyncStatus: etoptypes.S4Negative,
	// 	State:      shipnowtypes.StateCancelled,
	// 	Status:     etoptypes.S5Negative,
	// 	SyncStates: &model.FulfillmentSyncStates{
	// 		TrySyncAt:         time.Now(),
	// 		NextShippingState: model.StateCreated,
	// 	},
	// }
	// ffm, err = m.s.WithContext(ctx).UpdateSyncState(updateArgs)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (m *ProcessManager) validateOrders(ctx context.Context, orderIDs []int64, shipnowFfmID int64) error {
	if len(orderIDs) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Vui lòng chọn đơn hàng")
	}
	cmd := &ordering.ValidateOrdersForShippingArgs{
		OrderIDs: orderIDs,
	}
	if _, err := m.orderAggr.ValidateOrders(ctx, cmd); err != nil {
		return err
	}
	// for _, id := range orderIDs {
	// 	cmd1 := &shipnowmodelx.GetActiveShipnowFulfillmentsByOrderIDArgs{
	// 		OrderID:                     id,
	// 		ExcludeShipnowFulfillmentID: shipnowFfmID,
	// 	}
	// 	ffms, err := m.s.WithContext(ctx).GetActiveShipnowFulfillmentsByOrderID(cmd1)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if len(ffms) > 0 {
	// 		return cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng %v đã thuộc đơn vận chuyển khác (%v)", id, ffms[0].Id)
	// 	}
	// }
	return nil
}

func (m *ProcessManager) ValidateConfirm(ctx context.Context, ffm *shipnow.ShipnowFulfillment) error {
	switch ffm.ConfirmStatus {
	case etoptypes.S3Negative:
		return cm.Errorf(cm.FailedPrecondition, nil, "Đơn giao hàng đã hủy")
	case etoptypes.S3Positive:
		return cm.Errorf(cm.FailedPrecondition, nil, "Đơn giao hàng đã xác nhận")
	}
	if ffm.Status == etoptypes.S5Negative || ffm.Status == etoptypes.S5Positive {
		return cm.Errorf(cm.FailedPrecondition, nil, "Không thể xác nhận đơn giao hàng này")
	}

	if len(ffm.DeliveryPoints) == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "Số điểm giao hàng không hợp lệ.")
	}
	var orderIDs []int64
	for _, point := range ffm.DeliveryPoints {
		if point.OrderId == 0 {
			continue
		}
		orderIDs = append(orderIDs, point.OrderId)
	}
	if err := m.validateOrders(ctx, orderIDs, ffm.Id); err != nil {
		return err
	}
	return nil
}

func (m *ProcessManager) HandleUpdate(ctx context.Context, cmd *shipnow.UpdateShipnowFulfillmentArgs) (shipnowFfm *shipnow.ShipnowFulfillment, err error) {
	queryFfm := &shipnow.GetShipnowFulfillmentQuery{
		Id:     cmd.Id,
		ShopId: cmd.ShopId,
	}
	if err := m.shipnowQuery.Dispatch(ctx, queryFfm); err != nil {
		return nil, err
	}
	ffm := queryFfm.Result.ShipnowFulfillment
	if ffm.ConfirmStatus != etoptypes.S3Zero || ffm.ShippingCode != "" {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Không thể cập nhật đơn giao hàng này.")
	}

	shipnowFfm = &shipnow.ShipnowFulfillment{
		Id:                  cmd.Id,
		ShopId:              cmd.ShopId,
		PickupAddress:       cmd.PickupAddress,
		Carrier:             cmd.Carrier,
		ShippingServiceCode: cmd.ShippingServiceCode,
		ShippingServiceFee:  cmd.ShippingServiceFee,
		ShippingNote:        cmd.ShippingNote,
		RequestPickupAt:     nil,
	}

	if len(cmd.OrderIds) > 0 {
		if err := m.validateOrders(ctx, cmd.OrderIds, cmd.Id); err != nil {
			return nil, err
		}
		cmd3 := &ordering.GetOrdersArgs{
			ShopID: cmd.ShopId,
			IDs:    cmd.OrderIds,
		}
		orders, err := m.orderAggr.GetOrders(ctx, cmd3)
		if err != nil {
			return nil, err
		}
		shipnowFfm.WeightInfo = shipnowconvert.GetWeightInfo(orders.Orders)
		shipnowFfm.ValueInfo = shipnowconvert.GetValueInfo(orders.Orders)
		var deliveryPoints []*shipnow.DeliveryPoint
		for _, order := range orders.Orders {
			deliveryPoints = append(deliveryPoints, shipnowconvert.OrderToDeliveryPoint(order))
		}
		shipnowFfm.DeliveryPoints = deliveryPoints
	}
	return shipnowFfm, nil
}
