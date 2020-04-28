package pm

import (
	"context"

	"o.o/api/main/moneytx"
	"o.o/api/main/shipping"
	"o.o/api/top/types/etc/status3"
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

func New(eventB capi.EventBus, moneyTxQ moneytx.QueryBus,
	moneyTxA moneytx.CommandBus, shippingQ shipping.QueryBus) *ProcessManager {
	return &ProcessManager{
		eventBus:      eventB,
		moneyTxQuery:  moneyTxQ,
		moneyTxAggr:   moneyTxA,
		shippingQuery: shippingQ,
	}
}

func (m *ProcessManager) RegisterEvenHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.FulfillmentUpdating)
	eventBus.AddEventListener(m.FulfillmentShippingFeeChanged)
	eventBus.AddEventListener(m.MoneyTxShippingExternalDeleted)
}

func (m *ProcessManager) FulfillmentUpdating(ctx context.Context, event *shipping.FulfillmentUpdatingEvent) error {
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
	if err := canUpdateFulfillment(ffm); err != nil {
		return err
	}

	if ffm.MoneyTransactionID == 0 {
		return nil
	}

	queryMoneyTx := &moneytx.GetMoneyTxShippingByIDQuery{
		MoneyTxShippingID: ffm.MoneyTransactionID,
		ShopID:            ffm.ShopID,
	}
	if err := m.moneyTxQuery.Dispatch(ctx, queryMoneyTx); err != nil {
		return err
	}
	moneyTx := queryMoneyTx.Result
	if moneyTx.Status == status3.P {
		return cm.Errorf(cm.FailedPrecondition, nil, "Đơn vận chuyển đã đối soát.").WithMetap("money_transaction_id", ffm.MoneyTransactionID)
	}
	return nil
}

func canUpdateFulfillment(ffm *shipping.Fulfillment) error {
	if !ffm.CODEtopTransferedAt.IsZero() {
		return cm.Errorf(cm.FailedPrecondition, nil, "Đơn vận chuyển đã đối soát").WithMetap("money_transaction_id", ffm.MoneyTransactionID)
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
