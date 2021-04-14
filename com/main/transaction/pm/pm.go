package pm

import (
	"context"

	"o.o/api/main/credit"
	"o.o/api/main/invoicing"
	"o.o/api/main/transaction"
	"o.o/api/top/types/etc/invoice_type"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/subject_referral"
	"o.o/api/top/types/etc/transaction_type"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/capi/dot"
)

type ProcessManager struct {
	trxnQS    transaction.QueryBus
	trxnAggr  transaction.CommandBus
	creditQS  credit.QueryBus
	invoiceQS invoicing.QueryBus
}

func New(
	eventBus bus.EventRegistry,
	trxnQ transaction.QueryBus,
	trxnA transaction.CommandBus,
	creditQ credit.QueryBus,
	invoiceQS invoicing.QueryBus,
) *ProcessManager {
	p := &ProcessManager{
		trxnAggr:  trxnA,
		trxnQS:    trxnQ,
		creditQS:  creditQ,
		invoiceQS: invoiceQS,
	}
	p.registerEventHandler(eventBus)
	return p
}

func (m *ProcessManager) registerEventHandler(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.InvoiceDeleted)
	eventBus.AddEventListener(m.InvoicePaid)
}

func (m *ProcessManager) InvoiceDeleted(ctx context.Context, event *invoicing.InvoiceDeletedEvent) error {
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

func (m *ProcessManager) InvoicePaid(ctx context.Context, event *invoicing.InvoicePaidEvent) error {
	getInvoiceQuery := &invoicing.GetInvoiceByIDQuery{
		ID:        event.ID,
		AccountID: event.AccountID,
	}
	if err := m.invoiceQS.Dispatch(ctx, getInvoiceQuery); err != nil {
		return err
	}
	inv := getInvoiceQuery.Result

	getTransactionQuery := &transaction.GetTransactionByReferralQuery{
		ReferralType: subject_referral.Invoice,
		ReferralID:   event.ID,
	}
	err := m.trxnQS.Dispatch(ctx, getTransactionQuery)
	switch cm.ErrorCode(err) {
	case cm.NotFound:
		cmd := &transaction.CreateTransactionCommand{
			ID:           cm.NewID(),
			Name:         inv.Description,
			Amount:       inv.TotalAmount,
			AccountID:    inv.AccountID,
			Status:       status3.P,
			Type:         transaction_type.Invoice,
			Classify:     inv.Classify,
			Note:         "",
			ReferralType: subject_referral.Invoice,
			ReferralIDs:  []dot.ID{event.ID},
		}
		if inv.Type == invoice_type.Out {
			cmd.Amount = -inv.TotalAmount
		}
		if cmd.Name == "" {
			switch inv.ReferralType {
			case subject_referral.Subscription:
				cmd.Name = "Thanh toán subscription"
			case subject_referral.Credit:
				cmd.Name = "Nạp tiền vào tài khoản"
				if inv.Type == invoice_type.Out {
					cmd.Name = "Trừ tiền tài khoản"
				}
			default:

			}
		}

		if err = m.trxnAggr.Dispatch(ctx, cmd); err != nil {
			return err
		}
	case cm.NoError:
	// no-op
	default:
		return err
	}
	return nil
}
