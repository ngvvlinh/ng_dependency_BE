package pm

import (
	"context"

	"o.o/api/main/bankstatement"
	"o.o/api/main/credit"
	"o.o/api/top/types/etc/credit_type"
	"o.o/backend/pkg/common/bus"
)

type ProcessManager struct {
	creditAggr     credit.CommandBus
	creditQ        credit.QueryBus
	bankStatementQ bankstatement.QueryBus
}

func New(eventBus bus.EventRegistry, creditAggr credit.CommandBus, creditQ credit.QueryBus, bankStatementQ bankstatement.QueryBus) *ProcessManager {
	p := &ProcessManager{
		creditAggr:     creditAggr,
		creditQ:        creditQ,
		bankStatementQ: bankStatementQ,
	}
	p.registerEventHandlers(eventBus)
	return p
}

func (m *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.BankStatementCreated)
}

func (m *ProcessManager) BankStatementCreated(ctx context.Context, event *bankstatement.BankStatementCreatedEvent) error {
	query := &bankstatement.GetBankStatementQuery{
		ID: event.ID,
	}
	if err := m.bankStatementQ.Dispatch(ctx, query); err != nil {
		return err
	}
	bankStatement := query.Result

	// create credit
	creditArgs := &credit.CreateCreditCommand{
		Amount:          bankStatement.Amount,
		ShopID:          bankStatement.AccountID,
		Type:            credit_type.Shop,
		PaidAt:          bankStatement.TransferedAt,
		Classify:        credit_type.CreditClassifyShipping,
		BankStatementID: bankStatement.ID,
	}
	if err := m.creditAggr.Dispatch(ctx, creditArgs); err != nil {
		return err
	}

	// confirm credit
	creditConfirmArgs := &credit.ConfirmCreditCommand{
		ID:     creditArgs.Result.ID,
		ShopID: creditArgs.Result.Shop.ID,
	}
	if err := m.creditAggr.Dispatch(ctx, creditConfirmArgs); err != nil {
		return err
	}
	return nil
}
