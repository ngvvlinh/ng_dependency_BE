package pm

import (
	"context"

	"etop.vn/api/main/ordering"
	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/main/shipnow"
	"etop.vn/api/main/shipnow/carrier"
	"etop.vn/api/meta"
	"etop.vn/backend/pkg/common/bus"
)

type ProcessManager struct {
	eventBus     meta.EventBus
	shipnowQuery shipnow.QueryBus

	order          ordering.CommandBus
	carrierManager carrier.Manager
}

func New(
	eventBus meta.EventBus,
	shipnowBus shipnow.QueryBus,
	orderAggrBus ordering.CommandBus,
	carrierManager carrier.Manager,
) *ProcessManager {
	return &ProcessManager{
		eventBus:       eventBus,
		shipnowQuery:   shipnowBus,
		order:          orderAggrBus,
		carrierManager: carrierManager,
	}
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.ShipnowOrderReservation)
	eventBus.AddEventListener(m.ShipnowOrderChanged)
	eventBus.AddEventListener(m.ShipnowCancelled)
	eventBus.AddEventListener(m.ValidateConfirmed)
}

func (m *ProcessManager) ShipnowOrderReservation(ctx context.Context, event *shipnow.ShipnowOrderReservationEvent) error {
	// Call orderAggr for ReserveOrdersForFfm
	cmd := &ordering.ReserveOrdersForFfmCommand{
		OrderIDs:   event.OrderIds,
		Fulfill:    ordertypes.FulfillShipnowFulfillment,
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
		Fulfill:    ordertypes.FulfillShipnowFulfillment,
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
	return nil
}

func (m *ProcessManager) ValidateConfirmed(ctx context.Context, event *shipnow.ShipnowValidatedEvent) error {
	cmd := &ordering.ValidateOrdersForShippingCommand{
		OrderIDs: event.OrderIds,
	}
	if err := m.order.Dispatch(ctx, cmd); err != nil {
		return err
	}
	return nil
}
