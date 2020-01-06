package pm

import (
	"context"

	"etop.vn/api/main/moneytx"
	"etop.vn/api/main/shipping"
	"etop.vn/api/top/types/etc/status3"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/capi"
)

type ProcessManager struct {
	eventBus      capi.EventBus
	moneytxQuery  moneytx.QueryBus
	shippingQuery shipping.QueryBus
}

func New(eventB capi.EventBus, moneyTxQ moneytx.QueryBus, shippingQ shipping.QueryBus) *ProcessManager {
	return &ProcessManager{
		eventBus:      eventB,
		moneytxQuery:  moneyTxQ,
		shippingQuery: shippingQ,
	}
}

func (m *ProcessManager) RegisterEvenHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.FulfillmentUpdating)
	eventBus.AddEventListener(m.FulfillmentShippingFeeChanged)
}

func (m *ProcessManager) FulfillmentUpdating(ctx context.Context, event *shipping.FulfillmentUpdatingEvent) error {
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
		ID:     ffm.MoneyTransactionID,
		ShopID: ffm.ShopID,
	}
	if err := m.moneytxQuery.Dispatch(ctx, queryMoneyTx); err != nil {
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
	return nil
}
