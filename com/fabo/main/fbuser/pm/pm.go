package pm

import (
	"context"

	"o.o/api/fabo/fbusering"
	"o.o/api/shopping/tradering"
	"o.o/backend/pkg/common/bus"
	"o.o/capi"
)

type ProcessManager struct {
	eventBus   capi.EventBus
	fbUserAggr fbusering.CommandBus
}

func New(
	eventBusArgs capi.EventBus,
	fbUserA fbusering.CommandBus,
) *ProcessManager {
	return &ProcessManager{
		eventBus:   eventBusArgs,
		fbUserAggr: fbUserA,
	}
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.ShopCustomerDeletedEvent)
}

func (m *ProcessManager) ShopCustomerDeletedEvent(ctx context.Context, event *tradering.TraderDeletedEvent) error {
	if event.TradingType != tradering.CustomerType {
		return nil
	}
	cmd := &fbusering.DeleteFbExternalUserShopCustomerCommand{
		ShopID:     event.ShopID,
		CustomerID: event.TraderID.Wrap(),
	}
	err := m.fbUserAggr.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
}
