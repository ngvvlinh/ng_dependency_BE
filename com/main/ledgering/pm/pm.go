package pm

import (
	"context"

	"etop.vn/api/main/ledgering"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/capi"
)

type ProcessManager struct {
	eventBus capi.EventBus

	ledgerAggregate ledgering.CommandBus
}

func New(
	eventBus capi.EventBus,
	ledgerAggregate ledgering.CommandBus,
) *ProcessManager {
	return &ProcessManager{
		eventBus:        eventBus,
		ledgerAggregate: ledgerAggregate,
	}
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.AccountCreated)
}

func (m *ProcessManager) AccountCreated(ctx context.Context, event *ledgering.AccountCreatedEvent) error {
	cmd := &ledgering.CreateLedgerCommand{
		ShopID:      event.ShopID,
		Name:        "Tiền mặt",
		BankAccount: nil,
		Note:        "Số quỹ mặc định",
		Type:        string(ledgering.LedgerTypeCash),
		CreatedBy:   event.UserID,
	}
	if err := m.ledgerAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	return nil
}
