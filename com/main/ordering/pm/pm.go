package pm

import (
	"context"

	"o.o/api/main/inventory"
	"o.o/api/main/ordering"
	ordertrading "o.o/api/main/ordering/trading"
	"o.o/api/main/receipting"
	"o.o/api/main/shipnow"
	"o.o/api/main/shipping"
	"o.o/api/services/affiliate"
	"o.o/api/shopping/customering"
	"o.o/api/top/types/etc/receipt_ref"
	"o.o/api/top/types/etc/receipt_type"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/capi/dot"
	"o.o/common/l"
)

type ProcessManager struct {
	order         ordering.CommandBus
	affiliate     affiliate.CommandBus
	receiptQuery  receipting.QueryBus
	inventoryAgg  inventory.CommandBus
	orderQuery    ordering.QueryBus
	customerQuery customering.QueryBus
	shipnowQuery  shipnow.QueryBus
}

var ll = l.New()

func New(
	eventBus bus.EventRegistry,
	orderAggr ordering.CommandBus,
	affiliateAggr affiliate.CommandBus,
	receiptQs receipting.QueryBus,
	inventoryAgg inventory.CommandBus,
	orderQ ordering.QueryBus,
	customerQ customering.QueryBus,
	shipnowQ shipnow.QueryBus,
) *ProcessManager {
	p := &ProcessManager{
		order:         orderAggr,
		affiliate:     affiliateAggr,
		receiptQuery:  receiptQs,
		inventoryAgg:  inventoryAgg,
		orderQuery:    orderQ,
		customerQuery: customerQ,
		shipnowQuery:  shipnowQ,
	}
	p.registerEventHandlers(eventBus)
	return p
}

func (p *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(p.CheckTradingOrderValid)
	eventBus.AddEventListener(p.TradingOrderCreated)
	eventBus.AddEventListener(p.ReceiptConfirmed)
	eventBus.AddEventListener(p.ReceiptCancelled)
	eventBus.AddEventListener(p.ReceiptCreating)
	eventBus.AddEventListener(p.FulfillmentsCreatingEvent)
	eventBus.AddEventListener(p.FulfillmentsCreatedEvent)
	eventBus.AddEventListener(p.ReceiptConfirming)
	eventBus.AddEventListener(p.FulfillmentUpdatedInfoEvent)
	eventBus.AddEventListener(p.HandleShipnowFulfillmentCreatedEvent)
}

func (p *ProcessManager) ReceiptConfirming(ctx context.Context, event *receipting.ReceiptConfirmingEvent) error {
	queryReceipt := &receipting.GetReceiptByIDQuery{
		ID:     event.ReceiptID,
		ShopID: event.ShopID,
	}
	if err := p.receiptQuery.Dispatch(ctx, queryReceipt); err != nil {
		return err
	}
	receipt := queryReceipt.Result
	if receipt.RefType != receipt_ref.Order {
		return nil
	}
	// Ki???m tra receipt ???????c kh???i t???o t??? ?????ng t??? vi???c h???y order hay kh??ng (t???o ra khi h???y order ???? c?? s???n m???t receipt ???? ???????c x??c nh???n, receipt_type = payment).
	// L??c n??y kh??ng c???n ki???m tra tr???ng th??i c???a ????n h??ng n???a.
	// C??c tr?????ng h???p kh??c, kh??ng ???????c confirm receipt c???a ????n h??ng ???? h???y.
	if event.ReceiptType == receipt_type.Payment {
		return nil
	}
	queryOrder := &ordering.GetOrdersQuery{
		ShopID: event.ShopID,
		IDs:    receipt.RefIDs,
	}
	if err := p.orderQuery.Dispatch(ctx, queryOrder); err != nil {
		return err
	}
	for _, order := range queryOrder.Result.Orders {
		if order.Status == status5.N {
			return cm.Errorf(cm.InvalidArgument, nil, "????n h??ng %v ???? b??? h???y, kh??ng th??? x??c nh???n phi???u %v", order.Code, receipt.Type.GetLabelRefName())
		}
	}
	return nil
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

func (p *ProcessManager) handleReceiptConfirmedOrCancelled(ctx context.Context, receiptID, shopID dot.ID) error {
	var orderIDs []dot.ID
	mapOrderIDAndReceivedAmount := make(map[dot.ID]int)
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
	if getReceiptByIDQuery.Result.RefType != receipt_ref.Order {
		return nil
	}
	for _, orderID := range getReceiptByIDQuery.Result.RefIDs {
		mapOrderIDAndReceivedAmount[orderID] = 0
		orderIDs = append(orderIDs, orderID)
	}
	listReceiptsByRefIDsAndStatusQuery := &receipting.ListReceiptsByRefsAndStatusQuery{
		ShopID:  shopID,
		RefIDs:  orderIDs,
		RefType: receipt_ref.Order,
		Status:  int(status3.P),
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
	ctx context.Context, shopID dot.ID,
	orders []*ordering.Order, mapOrderIDAndReceivedAmount map[dot.ID]int,
) error {
	for _, order := range orders {
		var status status4.NullStatus
		receivedAmount := mapOrderIDAndReceivedAmount[order.ID]
		if receivedAmount >= order.TotalAmount {
			status = status4.P.Wrap()
		} else if receivedAmount > 0 {
			status = status4.S.Wrap()
		} else {
			status = status4.Z.Wrap()
		}

		cmd := &ordering.UpdateOrderPaymentStatusCommand{
			OrderID:       order.ID,
			ShopID:        shopID,
			PaymentStatus: status,
		}
		if err := p.order.Dispatch(ctx, cmd); err != nil {
			return err
		}
	}
	return nil
}

func (p *ProcessManager) validateTotalAmountAndReceivedAmount(
	ctx context.Context, shopID dot.ID, orderIDs []dot.ID,
	receipts []*receipting.Receipt, mapOrderIDAndReceivedAmount map[dot.ID]int,
) ([]*ordering.Order, error) {
	for _, receipt := range receipts {
		for _, receiptLine := range receipt.Lines {
			if receiptLine.RefID == 0 {
				continue
			}
			if _, ok := mapOrderIDAndReceivedAmount[receiptLine.RefID]; ok {
				switch receipt.Type {
				case receipt_type.Receipt:
					mapOrderIDAndReceivedAmount[receiptLine.RefID] += receiptLine.Amount
				case receipt_type.Payment:
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
	mOrder := make(map[dot.ID]*ordering.Order)
	receipt := event.Receipt
	refIDs := event.RefIDs
	mapRefIDAmount := event.MapRefIDAmount
	if receipt.RefType != receipt_ref.Order {
		return nil
	}
	if receipt.TraderID != 0 {
		getCustomerQuery := &customering.GetCustomerByIDQuery{
			ID:             receipt.TraderID,
			ShopID:         receipt.ShopID,
			IncludeDeleted: true,
		}
		if err := p.customerQuery.Dispatch(ctx, getCustomerQuery); err != nil {
			return err
		}
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
		if receipt.TraderID == customering.CustomerAnonymous && order.CustomerID == 0 {
			continue
		}
		if order.CustomerID != receipt.TraderID {
			return cm.Errorf(cm.FailedPrecondition, nil, "????n h??ng %v kh??ng thu???c ?????i t??c ???? ch???n", order.Code)
		}
	}

	// Check refIDs and orderIDs (result of query above)
	if len(refIDs) != len(orders) {
		for _, refID := range refIDs {
			if _, ok := mOrder[refID]; !ok {
				return cm.Errorf(cm.FailedPrecondition, nil, "ref_id %d kh??ng t??m th???y", refID)
			}
		}
	}

	// IGNORE: check received_amount of receipt type payment
	if receipt.Type == receipt_type.Payment {
		return nil
	}

	// List receipts by refIDs
	listReceiptsQuery := &receipting.ListReceiptsByRefsAndStatusQuery{
		ShopID:  receipt.ShopID,
		RefIDs:  refIDs,
		RefType: receipt_ref.Order,
		Status:  int(status3.P),
	}
	if err := p.receiptQuery.Dispatch(ctx, listReceiptsQuery); err != nil {
		return err
	}
	receipts := listReceiptsQuery.Result.Receipts

	// Get total amount each orderID
	// Map of [ orderId ] amount of receiptLines (current receipts into DB)
	mapRefIDAmountOld := make(map[dot.ID]int)
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
				case receipt_type.Receipt:
					mapRefIDAmountOld[receiptLine.RefID] += receiptLine.Amount
				case receipt_type.Payment:
					mapRefIDAmountOld[receiptLine.RefID] -= receiptLine.Amount
				}
			}
		}
	}
	// Check each amount of receiptLine (param) with (total amount of old receiptLines + total amount of order)
	for key, value := range mapRefIDAmount {
		if value > (mOrder[key].TotalAmount)-mapRefIDAmountOld[key] {
			return cm.Errorf(cm.InvalidArgument, nil, "Gi?? tr??? c???a ????n h??ng kh??ng h???p l???, Vui l??ng t???i l???i trang v?? th??? l???i")
		}
	}
	return nil
}

func (p *ProcessManager) FulfillmentsCreatingEvent(ctx context.Context, event *shipping.FulfillmentsCreatingEvent) error {
	// update order status to processing
	cmd := &ordering.UpdateOrderStatusCommand{
		OrderID: event.OrderID,
		ShopID:  event.ShopID,
		Status:  status5.S,
	}
	return p.order.Dispatch(ctx, cmd)
}

func (p *ProcessManager) FulfillmentsCreatedEvent(ctx context.Context, event *shipping.FulfillmentsCreatedEvent) error {
	// Update order: fulfillmentIDs & fulfillmentType (shippingType)
	cmd := &ordering.ReserveOrdersForFfmCommand{
		OrderIDs:   []dot.ID{event.OrderID},
		Fulfill:    event.ShippingType,
		FulfillIDs: event.FulfillmentIDs,
	}
	return p.order.Dispatch(ctx, cmd)
}

func (p *ProcessManager) FulfillmentUpdatedInfoEvent(ctx context.Context, event *shipping.FulfillmentUpdatedInfoEvent) error {
	if event.OrderID == 0 {
		return nil
	}
	cmd := &ordering.UpdateOrderCustomerInfoCommand{
		ID:       event.OrderID,
		FullName: event.FullName,
		Phone:    event.Phone,
	}
	return p.order.Dispatch(ctx, cmd)
}

func (p *ProcessManager) HandleShipnowFulfillmentCreatedEvent(ctx context.Context, event *shipnow.ShipnowCreatedEvent) error {
	if event.ShipnowFulfillmentID == 0 {
		return nil
	}

	query := &shipnow.GetShipnowFulfillmentQuery{
		ID: event.ShipnowFulfillmentID,
	}
	if err := p.shipnowQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	ffm := query.Result.ShipnowFulfillment

	cmd := &ordering.UpdateOrdersFulfillmentShippingCodesCommand{
		OrderIDs:                 ffm.OrderIDs,
		FulfillmentShippingCodes: []string{ffm.ShippingCode},
	}
	return p.order.Dispatch(ctx, cmd)
}
