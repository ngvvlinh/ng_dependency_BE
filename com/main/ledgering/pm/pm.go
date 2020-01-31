package pm

import (
	"context"

	"etop.vn/api/main/identity"
	"etop.vn/api/main/ledgering"
	"etop.vn/api/top/types/etc/ledger_type"
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
