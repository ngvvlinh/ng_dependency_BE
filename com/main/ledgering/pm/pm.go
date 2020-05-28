package pm

import (
	"context"

	"o.o/api/main/identity"
	"o.o/api/main/ledgering"
	"o.o/api/top/types/etc/ledger_type"
	"o.o/backend/pkg/common/bus"
	"o.o/capi"
)

type ProcessManager struct {
	eventBus capi.EventBus

	ledgerAggregate ledgering.CommandBus
}

func New(
	eventBus bus.EventRegistry,
	ledgerAggregate ledgering.CommandBus,
) *ProcessManager {
	p := &ProcessManager{
		eventBus:        eventBus,
		ledgerAggregate: ledgerAggregate,
	}
	p.registerEventHandlers(eventBus)
	return p
}

func (m *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.AccountCreated)
}

func (m *ProcessManager) AccountCreated(ctx context.Context, event *identity.AccountCreatedEvent) error {
	cmd := &ledgering.CreateLedgerCommand{
		ShopID:      event.ShopID,
		Name:        "Tiền mặt",
		BankAccount: nil,
		Note:        "Tài khoản thanh toán mặc định",
		Type:        ledger_type.LedgerTypeCash,
		CreatedBy:   event.UserID,
	}
	if err := m.ledgerAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	return nil
}
