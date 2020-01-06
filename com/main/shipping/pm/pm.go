package pm

import (
	"context"

	"etop.vn/api/main/moneytx"
	"etop.vn/api/main/shipping"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/capi"
)

type ProcessManager struct {
	eventBus  capi.EventBus
	shippingA shipping.CommandBus
}

func New(eventBus capi.EventBus, shippingAggregate shipping.CommandBus) *ProcessManager {
	return &ProcessManager{
		eventBus:  eventBus,
		shippingA: shippingAggregate,
	}
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.MoneyTxShippingExternalCreated)
}

func (m *ProcessManager) MoneyTxShippingExternalCreated(ctx context.Context, event *moneytx.MoneyTransactionShippingExternalCreatedEvent) error {
	if event.MoneyTxShippingExternalID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Event MoneyTransactionShippingExternalCreated missing ID")
	}
	if len(event.FulfillementIDs) == 0 {
		return nil
	}
	cmd := &shipping.UpdateFulfillmentsMoneyTxShippingExternalIDCommand{
		FulfillmentIDs:            event.FulfillementIDs,
		MoneyTxShippingExternalID: event.MoneyTxShippingExternalID,
	}
	if err := m.shippingA.Dispatch(ctx, cmd); err != nil {
		return err
	}
	return nil
}
