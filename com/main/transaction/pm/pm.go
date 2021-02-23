package pm

import (
	"context"

	"o.o/api/main/credit"
	"o.o/api/main/transaction"
	"o.o/api/subscripting/invoice"
	"o.o/api/top/types/etc/service_classify"
	"o.o/api/top/types/etc/subject_referral"
	"o.o/api/top/types/etc/transaction_type"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/capi/dot"
)

type ProcessManager struct {
	trxnQS   transaction.QueryBus
	trxnAggr transaction.CommandBus
	creditQS credit.QueryBus
}

func New(
	eventBus bus.EventRegistry,
	trxnQ transaction.QueryBus,
	trxnA transaction.CommandBus,
	creditQ credit.QueryBus,
) *ProcessManager {
	p := &ProcessManager{
		trxnAggr: trxnA,
		trxnQS:   trxnQ,
		creditQS: creditQ,
	}
	p.registerEventHandler(eventBus)
	return p
}

func (m *ProcessManager) registerEventHandler(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.InvoinceDeleted)
	eventBus.AddEventListener(m.CreditConfirmed)
}

func (m *ProcessManager) InvoinceDeleted(ctx context.Context, event *invoice.InvoiceDeletedEvent) error {
	query := &transaction.GetTransactionByReferralQuery{
		ReferralType: subject_referral.Invoice,
		ReferralID:   event.InvoinceID,
	}
	err := m.trxnQS.Dispatch(ctx, query)
	if err != nil && cm.ErrorCode(err) != cm.NotFound {
		return err
	}

	trxn := query.Result
	if trxn == nil {
		return nil
	}

	cmd := &transaction.DeleteTransactionCommand{
		TrxnID:    trxn.ID,
		AccountID: trxn.AccountID,
	}
	return m.trxnAggr.Dispatch(ctx, cmd)
}

func (m *ProcessManager) CreditConfirmed(ctx context.Context, event *credit.CreditConfirmedEvent) error {
	query := &credit.GetCreditQuery{
		ID:     event.CreditID,
		ShopID: event.ShopID,
	}
	if err := m.creditQS.Dispatch(ctx, query); err != nil {
		return err
	}
	_credit := query.Result

	cmd := &transaction.CreateTransactionCommand{
		ID:           _credit.ID,
		Name:         "",
		Amount:       _credit.Amount,
		AccountID:    _credit.ShopID,
		Status:       _credit.Status,
		Type:         transaction_type.Credit,
		Classify:     service_classify.ServiceClassify(_credit.Classify),
		Note:         "",
		ReferralType: subject_referral.Credit,
		ReferralIDs:  []dot.ID{_credit.ID},
	}
	if err := m.trxnAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	return nil
}
