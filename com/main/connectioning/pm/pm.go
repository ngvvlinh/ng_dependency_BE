package pm

import (
	"context"

	"o.o/api/main/connectioning"
	"o.o/api/top/types/etc/connection_type"
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
	// chỉ quan tâm tới connection direct (NVC trực tiếp tích hợp)
	if conn.ConnectionMethod != connection_type.ConnectionMethodDirect {
		return nil
	}

	queryListConn := &connectioning.ListConnectionsByOriginConnectionIDQuery{
		OriginConnectionID: conn.ID,
	}
	if err := m.connectionQuery.Dispatch(ctx, queryListConn); err != nil {
		return err
	}

	// Lấy tất cả các connection có origin_connection_id = conn.ID
	// thay đổi thông tin của các conn đó theo thông tin của connection gốc
	for _, _conn := range queryListConn.Result {
		update := &connectioning.UpdateConnectionCommand{
			ID:           _conn.ID,
			ImageURL:     conn.ImageURL,
			DriverConfig: conn.DriverConfig,
		}
		if err := m.connectionAggr.Dispatch(ctx, update); err != nil {
			return err
		}
	}
	return nil
}
