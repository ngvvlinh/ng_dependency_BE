package pm

import (
	"context"

	"etop.vn/api/main/identity"
	"etop.vn/api/shopping/customering"
	"etop.vn/api/shopping/customering/customer_type"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/capi"
)

type ProcessManager struct {
	eventBus capi.EventBus

	customerAggregate customering.CommandBus
}

func New(
	eventBus capi.EventBus,
	customerAggregate customering.CommandBus,
) *ProcessManager {
	return &ProcessManager{
		eventBus:          eventBus,
		customerAggregate: customerAggregate,
	}
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.AccountCreated)
}

func (m *ProcessManager) AccountCreated(ctx context.Context, event *identity.AccountCreatedEvent) error {
	cmd := &customering.CreateCustomerCommand{
		ShopID:   event.ShopID,
		FullName: "Khách lẻ",
		Type:     customer_type.Independent,
	}
	if err := m.customerAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	return nil
}
