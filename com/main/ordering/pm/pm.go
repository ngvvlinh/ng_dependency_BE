package pm

import (
	"context"

	"etop.vn/api/main/ordering"
	ordertrading "etop.vn/api/main/ordering/trading"
	"etop.vn/api/main/receipting"
	"etop.vn/api/services/affiliate"
	"etop.vn/backend/com/main/ordering/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/common/l"
)

type ProcessManager struct {
	order        ordering.CommandBus
	orderQS      ordering.QueryBus
	affiliate    affiliate.CommandBus
	receiptQuery receipting.QueryBus
}

var (
	ll = l.New()
)

func New(
	orderAggr ordering.CommandBus,
	orderQuery ordering.QueryBus,
	affiliateAggr affiliate.CommandBus,
	receiptQs receipting.QueryBus,
) *ProcessManager {
	return &ProcessManager{
		order:        orderAggr,
		orderQS:      orderQuery,
		affiliate:    affiliateAggr,
		receiptQuery: receiptQs,
	}
}

func (p *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(p.CheckTradingOrderValid)
	eventBus.AddEventListener(p.TradingOrderCreated)
	eventBus.AddEventListener(p.ReceiptConfirmedOrCancelled)
}

func (p *ProcessManager) CheckTradingOrderValid(ctx context.Context, event *ordertrading.TradingOrderCreatingEvent) error {
	checkCmd := &affiliate.TradingOrderCreatingCommand{
		ReferralCode: event.ReferralCode,
		UserID:       event.UserID,
	}
	if err := p.affiliate.Dispatch(ctx, checkCmd); err != nil {
		return err
	}
	return nil
}

func (p *ProcessManager) TradingOrderCreated(ctx context.Context, event *ordertrading.TradingOrderCreatedEvent) error {
	orderCreatedNotifyCmd := &affiliate.OnTradingOrderCreatedCommand{
		OrderID:      event.OrderID,
		ReferralCode: event.ReferralCode,
	}
	if err := p.affiliate.Dispatch(ctx, orderCreatedNotifyCmd); err != nil {
		return err
	}
	return nil
}

func (p *ProcessManager) ReceiptConfirmedOrCancelled(ctx context.Context, event *receipting.ReceiptConfirmedOrCancelledEvent) error {
	var orderIDs []int64
	mapOrderIDAndReceivedAmount := make(map[int64]int32)

	getReceiptByIDQuery := &receipting.GetReceiptByIDQuery{
		ID:     event.ReceiptID,
		ShopID: event.ShopID,
	}
	if err := p.receiptQuery.Dispatch(ctx, getReceiptByIDQuery); err != nil {
		return cm.MapError(err).
			Wrapf(cm.NotFound, "receipt %v not found", event.ReceiptID).
			Throw()
	}
	if len(getReceiptByIDQuery.Result.RefIDs) == 0 {
		return nil
	}
	if getReceiptByIDQuery.Result.RefType != receipting.ReceiptRefTypeOrder {
		return nil
	}

	for _, orderID := range getReceiptByIDQuery.Result.RefIDs {
		mapOrderIDAndReceivedAmount[orderID] = 0
		orderIDs = append(orderIDs, orderID)
	}
	listReceiptsByRefIDsAndStatusQuery := &receipting.ListReceiptsByRefIDsAndStatusQuery{
		ShopID: event.ShopID,
		RefIDs: orderIDs,
		Status: int32(model.S3Positive),
	}
	if err := p.receiptQuery.Dispatch(ctx, listReceiptsByRefIDsAndStatusQuery); err != nil {
		return err
	}

	orders, err := p.validateTotalAmountAndReceivedAmount(listReceiptsByRefIDsAndStatusQuery.Result.Receipts, mapOrderIDAndReceivedAmount, event, orderIDs, ctx)
	if err != nil {
		return err
	}

	if err := p.updatePaymentStatus(orders, mapOrderIDAndReceivedAmount, event, ctx); err != nil {
		return err
	}

	return nil
}

func (p *ProcessManager) updatePaymentStatus(orders []*ordering.Order, mapOrderIDAndReceivedAmount map[int64]int32, event *receipting.ReceiptConfirmedOrCancelledEvent, ctx context.Context) error {
	for _, order := range orders {
		if int(order.PaymentStatus) == int(model.S4Negative) || int(order.PaymentStatus) == int(model.S4SuperPos) {
			continue
		}
		var status *model.Status3
		if int32(order.TotalAmount) == mapOrderIDAndReceivedAmount[order.ID] {
			status = model.S3Positive.P()
		} else {
			status = model.S3Zero.P()
		}

		updateOrderPaymentStatus := &modelx.UpdateOrderPaymentStatusCommand{
			ShopID:  event.ShopID,
			OrderID: order.ID,
			Status:  status,
		}
		if err := bus.Dispatch(ctx, updateOrderPaymentStatus); err != nil {
			return err
		}
	}
	return nil
}

func (p *ProcessManager) validateTotalAmountAndReceivedAmount(receipts []*receipting.Receipt, mapOrderIDAndReceivedAmount map[int64]int32, event *receipting.ReceiptConfirmedOrCancelledEvent, orderIDs []int64, ctx context.Context) ([]*ordering.Order, error) {
	for _, receipt := range receipts {
		for _, receiptLine := range receipt.Lines {
			if receiptLine.RefID == 0 {
				continue
			}
			if _, ok := mapOrderIDAndReceivedAmount[receiptLine.RefID]; ok {
				switch receipt.Type {
				case receipting.ReceiptTypeReceipt:
					mapOrderIDAndReceivedAmount[receiptLine.RefID] += receiptLine.Amount
				case receipting.ReceiptTypePayment:
					mapOrderIDAndReceivedAmount[receiptLine.RefID] -= receiptLine.Amount
				}
			}
		}
	}
	listOrdersByIDsQuery := &ordering.GetOrdersQuery{
		ShopID: event.ShopID,
		IDs:    orderIDs,
	}
	if err := p.orderQS.Dispatch(ctx, listOrdersByIDsQuery); err != nil {
		return nil, err
	}
	if len(listOrdersByIDsQuery.Result.Orders) == 0 {
		return nil, cm.Errorf(cm.NotFound, nil, "Orders %v not found", orderIDs)
	}
	orders := listOrdersByIDsQuery.Result.Orders
	if len(orders) != len(orderIDs) {
		for _, order := range orders {
			var isSame bool
			for _, orderID := range orderIDs {
				if order.ID == orderID {
					isSame = true
					break
				}
			}
			if !isSame {
				return nil, cm.Errorf(cm.NotFound, nil, "Order %v not found", order.ID)
			}
		}
	}
	return orders, nil
}
