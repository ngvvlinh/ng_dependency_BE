package pm

import (
	"context"

	"etop.vn/api/shopping/customering"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/inventory"
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
	order         ordering.CommandBus
	affiliate     affiliate.CommandBus
	receiptQuery  receipting.QueryBus
	inventoryAgg  inventory.CommandBus
	orderQuery    ordering.QueryBus
	customerQuery customering.QueryBus
}

var (
	ll = l.New()
)

func New(
	orderAggr ordering.CommandBus,
	affiliateAggr affiliate.CommandBus,
	receiptQs receipting.QueryBus,
	inventoryAgg inventory.CommandBus,
	orderQ ordering.QueryBus,
	customerQ customering.QueryBus,
) *ProcessManager {
	return &ProcessManager{
		order:         orderAggr,
		affiliate:     affiliateAggr,
		receiptQuery:  receiptQs,
		inventoryAgg:  inventoryAgg,
		orderQuery:    orderQ,
		customerQuery: customerQ,
	}
}

func (p *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(p.CheckTradingOrderValid)
	eventBus.AddEventListener(p.TradingOrderCreated)
	eventBus.AddEventListener(p.ReceiptConfirmed)
	eventBus.AddEventListener(p.ReceiptCancelled)
	eventBus.AddEventListener(p.ReceiptCreating)
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

func (p *ProcessManager) ReceiptConfirmed(ctx context.Context, event *receipting.ReceiptConfirmedEvent) error {
	if err := p.handleReceiptConfirmedOrCancelled(ctx, event.ReceiptID, event.ShopID); err != nil {
		return err
	}
	return nil
}

func (p *ProcessManager) ReceiptCancelled(ctx context.Context, event *receipting.ReceiptCancelledEvent) error {
	if err := p.handleReceiptConfirmedOrCancelled(ctx, event.ReceiptID, event.ShopID); err != nil {
		return err
	}
	return nil
}

func (p *ProcessManager) handleReceiptConfirmedOrCancelled(ctx context.Context, receiptID, shopID int64) error {
	var orderIDs []int64
	mapOrderIDAndReceivedAmount := make(map[int64]int32)
	getReceiptByIDQuery := &receipting.GetReceiptByIDQuery{
		ID:     receiptID,
		ShopID: shopID,
	}
	if err := p.receiptQuery.Dispatch(ctx, getReceiptByIDQuery); err != nil {
		return cm.MapError(err).
			Wrapf(cm.NotFound, "receipt %v not found", receiptID).
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
	listReceiptsByRefIDsAndStatusQuery := &receipting.ListReceiptsByRefsAndStatusQuery{
		ShopID:  shopID,
		RefIDs:  orderIDs,
		RefType: receipting.ReceiptRefTypeOrder,
		Status:  int32(model.S3Positive),
	}
	if err := p.receiptQuery.Dispatch(ctx, listReceiptsByRefIDsAndStatusQuery); err != nil {
		return err
	}
	orders, err := p.validateTotalAmountAndReceivedAmount(ctx, shopID, orderIDs, listReceiptsByRefIDsAndStatusQuery.Result.Receipts, mapOrderIDAndReceivedAmount)
	if err != nil {
		return err
	}
	if err := p.updatePaymentStatus(ctx, shopID, orders, mapOrderIDAndReceivedAmount); err != nil {
		return err
	}
	return nil
}

func (p *ProcessManager) updatePaymentStatus(
	ctx context.Context, shopID int64,
	orders []*ordering.Order, mapOrderIDAndReceivedAmount map[int64]int32,
) error {
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
			ShopID:  shopID,
			OrderID: order.ID,
			Status:  status,
		}
		if err := bus.Dispatch(ctx, updateOrderPaymentStatus); err != nil {
			return err
		}
	}
	return nil
}

func (p *ProcessManager) validateTotalAmountAndReceivedAmount(
	ctx context.Context, shopID int64, orderIDs []int64,
	receipts []*receipting.Receipt, mapOrderIDAndReceivedAmount map[int64]int32,
) ([]*ordering.Order, error) {
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
		ShopID: shopID,
		IDs:    orderIDs,
	}
	if err := p.orderQuery.Dispatch(ctx, listOrdersByIDsQuery); err != nil {
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

func (p *ProcessManager) ReceiptCreating(ctx context.Context, event *receipting.ReceiptCreatingEvent) error {
	var orders []*ordering.Order
	var isIndependentCustomer bool
	mOrder := make(map[int64]*ordering.Order)
	receipt := event.Receipt
	refIDs := event.RefIDs
	mapRefIDAmount := event.MapRefIDAmount
	if receipt.RefType != receipting.ReceiptRefTypeOrder {
		return nil
	}

	getCustomerQuery := &customering.GetCustomerByIDQuery{
		ID:     receipt.TraderID,
		ShopID: receipt.ShopID,
	}
	if err := p.customerQuery.Dispatch(ctx, getCustomerQuery); err != nil {
		return err
	}
	if getCustomerQuery.Result.Type == customering.CustomerTypeIndependent {
		isIndependentCustomer = true
	}

	// List orders depend on refIDs
	query := &ordering.GetOrdersQuery{
		ShopID: receipt.ShopID,
		IDs:    refIDs,
	}
	if err := p.orderQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	orders = query.Result.Orders
	for _, order := range orders {
		mOrder[order.ID] = order
		if isIndependentCustomer && order.CustomerID == 0 {
			continue
		}
		if order.CustomerID != receipt.TraderID {
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn nhập hàng %v không thuộc đối tác đã chọn", order.Code)
		}
	}

	// Check refIDs and orderIDs (result of query above)
	if len(refIDs) != len(orders) {
		for _, refID := range refIDs {
			if _, ok := mOrder[refID]; !ok {
				return cm.Errorf(cm.FailedPrecondition, nil, "ref_id %d không tìm thấy", refID)
			}
		}
	}

	// IGNORE: check received_amount of receipt type payment
	if receipt.Type == receipting.ReceiptTypePayment {
		return nil
	}

	// List receipts by refIDs
	listReceiptsQuery := &receipting.ListReceiptsByRefsAndStatusQuery{
		ShopID:  receipt.ShopID,
		RefIDs:  refIDs,
		RefType: receipting.ReceiptRefTypeOrder,
		Status:  int32(etop.S3Positive),
	}
	if err := p.receiptQuery.Dispatch(ctx, listReceiptsQuery); err != nil {
		return err
	}
	receipts := listReceiptsQuery.Result.Receipts

	// Get total amount each orderID
	// Map of [ orderId ] amount of receiptLines (current receipts into DB)
	mapRefIDAmountOld := make(map[int64]int32)
	for _, receiptElem := range receipts {
		// Ignore current receipt when updating
		if receiptElem.ID == receipt.ID {
			continue
		}
		for _, receiptLine := range receiptElem.Lines {
			if receiptLine.RefID == 0 {
				continue
			}
			if _, has := mapRefIDAmount[receiptLine.RefID]; has {
				switch receipt.Type {
				case receipting.ReceiptTypeReceipt:
					mapRefIDAmountOld[receiptLine.RefID] += receiptLine.Amount
				case receipting.ReceiptTypePayment:
					mapRefIDAmountOld[receiptLine.RefID] -= receiptLine.Amount
				}
			}
		}
	}
	// Check each amount of receiptLine (param) with (total amount of old receiptLines + total amount of order)
	for key, value := range mapRefIDAmount {
		if value > int32(mOrder[key].TotalAmount)-mapRefIDAmountOld[key] {
			return cm.Errorf(cm.InvalidArgument, nil, "Giá trị của đơn hàng không hợp lệ, Vui lòng tải lại trang và thử lại")
		}
	}
	return nil
}
