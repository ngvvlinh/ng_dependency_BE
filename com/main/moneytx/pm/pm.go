package pm

import (
	"context"

	"o.o/api/main/moneytx"
	"o.o/api/main/shipping"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/capi"
)

type ProcessManager struct {
	eventBus      capi.EventBus
	moneyTxQuery  moneytx.QueryBus
	moneyTxAggr   moneytx.CommandBus
	shippingQuery shipping.QueryBus
}

func New(eventBus bus.EventRegistry, moneyTxQ moneytx.QueryBus,
	moneyTxA moneytx.CommandBus, shippingQ shipping.QueryBus) *ProcessManager {
	p := &ProcessManager{
		eventBus:      eventBus,
		moneyTxQuery:  moneyTxQ,
		moneyTxAggr:   moneyTxA,
		shippingQuery: shippingQ,
	}
	p.registerEvenHandlers(eventBus)
	return p
}

func (m *ProcessManager) registerEvenHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.FulfillmentUpdated)
	eventBus.AddEventListener(m.FulfillmentShippingFeeChanged)
	eventBus.AddEventListener(m.MoneyTxShippingExternalDeleted)
}

func (m *ProcessManager) FulfillmentUpdated(ctx context.Context, event *shipping.FulfillmentUpdatedEvent) error {
	if event.MoneyTxShippingID == 0 {
		return nil
	}
	query := &shipping.GetFulfillmentByIDOrShippingCodeQuery{
		ID: event.FulfillmentID,
	}
	if err := m.shippingQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	ffm := query.Result

	if ffm.MoneyTransactionID == 0 {
		return nil
	}

	// Chỉ tính lại phiên khi trạng thái ffm nằm trong các trạng thái cho phép
	// Các trường hợp khác, giữ nguyên phiên. Khi admin xác nhận phiên gặp lỗi sẽ xử lý sau
	if cm.StringsContain(moneytx.ShippingAcceptStates, ffm.ShippingState.String()) {
		updateMoneyTx := &moneytx.ReCalcMoneyTxShippingCommand{
			MoneyTxShippingID: ffm.MoneyTransactionID,
		}
		if err := m.moneyTxAggr.Dispatch(ctx, updateMoneyTx); err != nil {
			return err
		}
	}
	return nil
}

func (m *ProcessManager) FulfillmentShippingFeeChanged(ctx context.Context, event *shipping.FulfillmentShippingFeeChangedEvent) error {
	if event.MoneyTxShippingID == 0 {
		return nil
	}
	query := &shipping.GetFulfillmentByIDOrShippingCodeQuery{
		ID: event.FulfillmentID,
	}
	if err := m.shippingQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	ffm := query.Result
	if ffm.MoneyTransactionID == 0 {
		return nil
	}
	// TODO: change total_cod, total_amount ... if needed
	return nil
}

func (m *ProcessManager) MoneyTxShippingExternalDeleted(ctx context.Context, event *moneytx.MoneyTxShippingExternalDeletedEvent) error {
	if event.MoneyTxShippingExternalID == 0 {
		return nil
	}
	cmd := &moneytx.DeleteMoneyTxShippingExternalLinesCommand{
		MoneyTxShippingExternalID: event.MoneyTxShippingExternalID,
	}
	return m.moneyTxAggr.Dispatch(ctx, cmd)
}
