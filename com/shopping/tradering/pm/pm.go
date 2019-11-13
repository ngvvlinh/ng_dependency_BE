package pm

import (
	"context"

	"etop.vn/api/shopping/tradering"

	"etop.vn/backend/pkg/common/bus"
	"etop.vn/capi"
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
