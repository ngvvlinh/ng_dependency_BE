package pm

import (
	"context"

	"o.o/api/shopping/tradering"
	"o.o/backend/pkg/common/bus"
	"o.o/capi"
)

type ProcessManager struct {
	eventBus capi.EventBus

	traderAggregate tradering.CommandBus
}

func New(
	eventBus capi.EventBus,
	traderAggregate tradering.CommandBus,
) *ProcessManager {
	return &ProcessManager{
		eventBus:        eventBus,
		traderAggregate: traderAggregate,
	}
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.DeleteTraderEvent)
}

func (m *ProcessManager) DeleteTraderEvent(ctx context.Context, event *tradering.TraderDeletedEvent) error {
	cmd := &tradering.DeleteTraderCommand{
		ShopID: event.ShopID,
		ID:     event.TraderID,
	}
	if err := m.traderAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	return nil
}
