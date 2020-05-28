package pm

import (
	"context"
	"time"

	"o.o/api/main/ordering"
	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/main/shipnow"
	"o.o/api/main/shipnow/carrier"
	shipnowcarrier "o.o/api/main/shipnow/carrier"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	etopconvert "o.o/backend/com/main/etop/convert"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/model"
	"o.o/capi"
	"o.o/capi/dot"
)

type ProcessManager struct {
	eventBus     capi.EventBus
	shipnowQuery shipnow.QueryBus
	shipnow      shipnow.CommandBus

	order          ordering.CommandBus
	carrierManager carrier.Manager
}

func New(
	eventBus bus.EventRegistry,
	shipnowQuery shipnow.QueryBus,
	shipnowAggrBus shipnow.CommandBus,
	orderAggrBus ordering.CommandBus,
	carrierManager carrier.Manager,
) *ProcessManager {
	p := &ProcessManager{
		eventBus:       eventBus,
		shipnowQuery:   shipnowQuery,
		shipnow:        shipnowAggrBus,
		order:          orderAggrBus,
		carrierManager: carrierManager,
	}
	p.registerEventHandlers(eventBus)
	return p
}

func (m *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.ShipnowOrderReservation)
	eventBus.AddEventListener(m.ShipnowOrderChanged)
	eventBus.AddEventListener(m.ShipnowCancelled)
	eventBus.AddEventListener(m.ValidateConfirmed)
	eventBus.AddEventListener(m.ShipnowExternalCreated)
}

func (m *ProcessManager) ShipnowOrderReservation(ctx context.Context, event *shipnow.ShipnowOrderReservationEvent) error {
	// Call orderAggr for ReserveOrdersForFfm
	cmd := &ordering.ReserveOrdersForFfmCommand{
		OrderIDs:   event.OrderIds,
		Fulfill:    ordertypes.ShippingTypeShipnow,
		FulfillIDs: []dot.ID{event.ShipnowFulfillmentId},
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
		Fulfill:    ordertypes.ShippingTypeShipnow,
		FulfillIDs: []dot.ID{event.ShipnowFulfillmentId},
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
		ShopConfirm:   status3.P,
		ConfirmStatus: status3.P,
	}
	if err := m.order.Dispatch(ctx, cmd2); err != nil {
		return err
	}
	return nil
}

func (m *ProcessManager) ShipnowExternalCreated(ctx context.Context, event *shipnow.ShipnowExternalCreatedEvent) (_err error) {
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
			SyncStatus: status4.S,
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
			SyncStatus: status4.N,
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
		TotalFee:                   xShipnow.TotalFee,
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
