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
	eventBus bus.EventRegistry,
	traderAggregate tradering.CommandBus,
) *ProcessManager {
	p := &ProcessManager{
		eventBus:        eventBus,
		traderAggregate: traderAggregate,
	}
	p.registerEventHandlers(eventBus)
	return p
}

func (m *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
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
