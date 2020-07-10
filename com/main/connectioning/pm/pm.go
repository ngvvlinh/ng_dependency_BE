package pm

import (
	"context"

	"o.o/api/main/connectioning"
	"o.o/backend/pkg/common/bus"
)

type ProcessManager struct {
	connectionAggr  connectioning.CommandBus
	connectionQuery connectioning.QueryBus
}

func New(eventBus bus.EventRegistry, connA connectioning.CommandBus, connQ connectioning.QueryBus) *ProcessManager {
	p := &ProcessManager{
		connectionAggr:  connA,
		connectionQuery: connQ,
	}
	p.registerEventHandlers(eventBus)
	return p
}

func (m *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.ConnectionUpdated)
}

func (m *ProcessManager) ConnectionUpdated(ctx context.Context, event *connectioning.ConnectionUpdatedEvent) error {
	query := &connectioning.GetConnectionByIDQuery{
		ID: event.ConnectionID,
	}
	if err := m.connectionQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	conn := query.Result

	args := &connectioning.UpdateConnectionFromOriginCommand{
		ConnectionID: conn.ID,
	}
	if err := m.connectionAggr.Dispatch(ctx, args); err != nil {
		return err
	}
	return nil
}
