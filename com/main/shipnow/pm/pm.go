package pm

import (
	"context"
	"time"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/ordering"
	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/main/shipnow"
	"etop.vn/api/main/shipnow/carrier"
	shipnowcarrier "etop.vn/api/main/shipnow/carrier"
	"etop.vn/api/meta"
	etopconvert "etop.vn/backend/com/main/etop/convert"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/common/bus"
)

type ProcessManager struct {
	eventBus     meta.EventBus
	shipnowQuery shipnow.QueryBus
	shipnow      shipnow.CommandBus

	order          ordering.CommandBus
	carrierManager carrier.Manager
}

func New(
	eventBus meta.EventBus,
	shipnowQuery shipnow.QueryBus,
	shipnowAggrBus shipnow.CommandBus,
	orderAggrBus ordering.CommandBus,
	carrierManager carrier.Manager,
) *ProcessManager {
	return &ProcessManager{
		eventBus:       eventBus,
		shipnowQuery:   shipnowQuery,
		shipnow:        shipnowAggrBus,
		order:          orderAggrBus,
		carrierManager: carrierManager,
	}
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.ShipnowOrderReservation)
	eventBus.AddEventListener(m.ShipnowOrderChanged)
	eventBus.AddEventListener(m.ShipnowCancelled)
	eventBus.AddEventListener(m.ValidateConfirmed)
	eventBus.AddEventListener(m.ShipnowCreateExternal)
}

func (m *ProcessManager) ShipnowOrderReservation(ctx context.Context, event *shipnow.ShipnowOrderReservationEvent) error {
	// Call orderAggr for ReserveOrdersForFfm
	cmd := &ordering.ReserveOrdersForFfmCommand{
		OrderIDs:   event.OrderIds,
		Fulfill:    ordertypes.FulfillShipnow,
		FulfillIDs: []int64{event.ShipnowFulfillmentId},
	}
	if err := m.order.Dispatch(ctx, cmd); err != nil {
		return err
	}
	return nil
}

func (m *ProcessManager) ShipnowOrderChanged(ctx context.Context, event *shipnow.ShipnowOrderChangedEvent) error {
	// release old orderIDs and reserve new orderIDs
	cmd := &ordering.ReleaseOrdersForFfmCommand{
		OrderIDs: event.OldOrderIds,
	}
	if err := m.order.Dispatch(ctx, cmd); err != nil {
		return err
	}

	cmd2 := &ordering.ReserveOrdersForFfmCommand{
		OrderIDs:   event.OrderIds,
		Fulfill:    ordertypes.FulfillShipnow,
		FulfillIDs: []int64{event.ShipnowFulfillmentId},
	}
	if err := m.order.Dispatch(ctx, cmd2); err != nil {
		return err
	}
	return nil
}

func (m *ProcessManager) ShipnowCancelled(ctx context.Context, event *shipnow.ShipnowCancelledEvent) error {
	// release orderIDs
	cmd := &ordering.ReleaseOrdersForFfmCommand{
		OrderIDs: event.OrderIds,
	}
	if err := m.order.Dispatch(ctx, cmd); err != nil {
		return err
	}

	query := &shipnow.GetShipnowFulfillmentQuery{
		Id: event.ShipnowFulfillmentId,
	}
	if err := m.shipnowQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	ffm := query.Result.ShipnowFulfillment
	if ffm.ShippingCode != "" {
		cmd2 := &shipnowcarrier.CancelExternalShipnowCommand{
			ShopID:               ffm.ShopId,
			ShipnowFulfillmentID: ffm.Id,
			ExternalShipnowID:    ffm.ShippingCode,
			CarrierServiceCode:   ffm.ShippingServiceCode,
			CancelReason:         event.CancelReason,
			Carrier:              ffm.Carrier,
		}
		if err := m.carrierManager.CancelExternalShipping(ctx, cmd2); err != nil {
			return err
		}
	}

	return nil
}

func (m *ProcessManager) ValidateConfirmed(ctx context.Context, event *shipnow.ShipnowValidateConfirmedEvent) error {
	cmd := &ordering.ValidateOrdersForShippingCommand{
		OrderIDs: event.OrderIds,
	}
	if err := m.order.Dispatch(ctx, cmd); err != nil {
		return err
	}

	// update order confirm status
	cmd2 := &ordering.UpdateOrdersConfirmStatusCommand{
		IDs:           event.OrderIds,
		ShopConfirm:   etop.S3Positive,
		ConfirmStatus: etop.S3Positive,
	}
	if err := m.order.Dispatch(ctx, cmd2); err != nil {
		return err
	}
	return nil
}

func (m *ProcessManager) ShipnowCreateExternal(ctx context.Context, event *shipnow.ShipnowCreateExternalEvent) (_err error) {
	query := &shipnow.GetShipnowFulfillmentQuery{
		Id: event.ShipnowFulfillmentId,
	}
	if err := m.shipnowQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	ffm := query.Result.ShipnowFulfillment
	{
		// update sync status
		update := &shipnow.UpdateShipnowFulfillmentStateCommand{
			Id:         ffm.Id,
			SyncStatus: etop.S4SuperPos,
			SyncStates: &shipnow.SyncStates{
				TrySyncAt: time.Now(),
			},
		}
		if err := m.shipnow.Dispatch(ctx, update); err != nil {
			return err
		}
	}

	defer func() {
		if _err == nil {
			return
		}
		update := &shipnow.UpdateShipnowFulfillmentStateCommand{
			Id:         ffm.Id,
			SyncStatus: etop.S4Negative,
			SyncStates: &shipnow.SyncStates{
				TrySyncAt: time.Now(),
				Error:     etopconvert.Error(model.ToError(_err)),
			},
		}
		// Keep the original error
		_ = m.shipnow.Dispatch(ctx, update)
	}()

	cmd := &shipnowcarrier.CreateExternalShipnowCommand{
		ShopID:               ffm.ShopId,
		ShipnowFulfillmentID: ffm.Id,
		PickupAddress:        ffm.PickupAddress,
		DeliveryPoints:       ffm.DeliveryPoints,
		ShippingNote:         ffm.ShippingNote,
	}
	xShipnow, err := m.carrierManager.CreateExternalShipping(ctx, cmd)
	if err != nil {
		return err
	}

	cmd2 := &shipnow.UpdateShipnowFulfillmentCarrierInfoCommand{
		Id:                         ffm.Id,
		ShippingCode:               xShipnow.ID,
		ShippingState:              xShipnow.State,
		TotalFee:                   int32(xShipnow.TotalFee),
		FeeLines:                   xShipnow.FeeLines,
		CarrierFeeLines:            xShipnow.FeeLines,
		ShippingCreatedAt:          xShipnow.CreatedAt,
		ShippingServiceName:        xShipnow.Service.Name,
		ShippingServiceDescription: xShipnow.Service.Description,
		ShippingSharedLink:         xShipnow.SharedLink,
	}
	if err := m.shipnow.Dispatch(ctx, cmd2); err != nil {
		return nil
	}
	return nil
}
