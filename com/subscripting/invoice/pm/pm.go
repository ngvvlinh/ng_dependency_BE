package pm

import (
	"context"

	"o.o/api/main/transaction"
	"o.o/api/subscripting/invoice"
	"o.o/api/top/types/etc/payment_method"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/common/xerrors"
)

type ProcessManager struct {
	transactionQuery transaction.QueryBus
}

func New(
	eventBus bus.EventRegistry,
	transactionQ transaction.QueryBus,
) *ProcessManager {
	p := &ProcessManager{
		transactionQuery: transactionQ,
	}
	p.registerEventHandlers(eventBus)
	return p
}

func (m *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.InvoicePaying)
}

func (m *ProcessManager) InvoicePaying(ctx context.Context, event *invoice.InvoicePayingEvent) error {
	if event.PaymentMethod != payment_method.Balance {
		return nil
	}
	if !event.ServiceClassify.Valid {
		return wrapError(cm.Internal, nil, "Missing credit classify")
	}
	if event.OwnerID == 0 {
		return wrapError(cm.Internal, nil, "Missing ownerID")
	}
	var balance int

	query := &transaction.GetBalanceUserQuery{
		UserID:   event.OwnerID,
		Classify: event.ServiceClassify.Enum,
	}
	if err := m.transactionQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	balance = query.Result.ActualBalance

	if balance < event.TotalAmount {
		return cm.Errorf(cm.FailedPrecondition, nil, "Số dư không đủ để thanh toán")
	}
	return nil
}

func wrapError(code xerrors.Code, err error, msg string) error {
	if err != nil {
		return cm.Errorf(cm.ErrorCode(err), err, "Lỗi từ invoice process manager: %v", err.Error())
	}
	return cm.Errorf(code, err, "Lỗi từ invoice process manager: %v", msg)
}
