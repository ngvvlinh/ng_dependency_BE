package pm

import (
	"context"

	"o.o/api/main/identity"
	"o.o/api/supporting/ticket"
	"o.o/api/top/types/etc/ticket/ticket_type"
	"o.o/backend/pkg/common/bus"
	"o.o/capi"
)

type ProcessManager struct {
	eventBus capi.EventBus

	ticketLabelAggregate ticket.CommandBus
}

var defaultTicketLabels = []ticket.TicketLabel{
	{
		Name:  "Góp ý",
		Code:  "feedback",
		Color: "#9b59b6",
	},
	{
		Name:  "Tư vấn",
		Code:  "support",
		Color: "#1abc9c",
	},
	{
		Name:  "Khiếu nại",
		Code:  "complain",
		Color: "#e74c3c",
	},
}

func NewProcessManager(
	eventBus bus.EventRegistry,
	ticketLabelAggregate ticket.CommandBus,
) *ProcessManager {
	p := &ProcessManager{
		eventBus:             eventBus,
		ticketLabelAggregate: ticketLabelAggregate,
	}
	p.registerEventHandlers(eventBus)
	return p
}

func (m *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.AccountCreated)
}

func (m *ProcessManager) AccountCreated(ctx context.Context, event *identity.AccountCreatedEvent) error {
	for _, ticketLabel := range defaultTicketLabels {
		if err := m.ticketLabelAggregate.Dispatch(ctx, &ticket.CreateTicketLabelCommand{
			ShopID: event.ShopID,
			Type:   ticket_type.Internal,
			Name:   ticketLabel.Name,
			Code:   ticketLabel.Code,
			Color:  ticketLabel.Color,
		}); err != nil {
			return err
		}
	}
	return nil
}
