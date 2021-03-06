package pm

import (
	"context"
	"fmt"
	"time"

	"o.o/api/main/identity"
	identitytypes "o.o/api/main/identity/types"
	"o.o/api/main/ledgering"
	"o.o/api/main/moneytx"
	"o.o/api/main/ordering"
	"o.o/api/main/receipting"
	shippingcore "o.o/api/main/shipping"
	"o.o/api/shopping/customering"
	"o.o/api/top/types/etc/ledger_type"
	"o.o/api/top/types/etc/receipt_mode"
	"o.o/api/top/types/etc/receipt_ref"
	"o.o/api/top/types/etc/receipt_type"
	"o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/status3"
	ordermodelx "o.o/backend/com/main/ordering/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etc/idutil"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi"
	"o.o/capi/dot"
	"o.o/common/l"
)

type ProcessManager struct {
	eventBus capi.EventBus

	receiptQuery  receipting.QueryBus
	receiptAggr   receipting.CommandBus
	ledgerQuery   ledgering.QueryBus
	ledgerAggr    ledgering.CommandBus
	identityQuery identity.QueryBus

	moneyTxQuery  moneytx.QueryBus
	OrderStore    sqlstore.OrderStoreInterface
	shippingQuery shippingcore.QueryBus
}

func New(
	eventBus bus.EventRegistry,
	receiptQuery receipting.QueryBus,
	receiptAggregate receipting.CommandBus,
	ledgerQuery ledgering.QueryBus,
	ledgerAggregate ledgering.CommandBus,
	identityQuery identity.QueryBus,
	moneyTxQ moneytx.QueryBus,
	OrderStore sqlstore.OrderStoreInterface,
	shippingQ shippingcore.QueryBus,
) *ProcessManager {
	p := &ProcessManager{
		eventBus:      eventBus,
		receiptQuery:  receiptQuery,
		receiptAggr:   receiptAggregate,
		ledgerQuery:   ledgerQuery,
		ledgerAggr:    ledgerAggregate,
		identityQuery: identityQuery,
		moneyTxQuery:  moneyTxQ,
		OrderStore:    OrderStore,
		shippingQuery: shippingQ,
	}
	p.registerEventHandlers(eventBus)
	return p
}

var ll = l.New()

func (m *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.MoneyTransactionConfirmed)
	eventBus.AddEventListener(m.MoneyTxShippingEtopConfirmed)
	eventBus.AddEventListener(m.OrderCancelled)
}

func (m *ProcessManager) OrderCancelled(ctx context.Context, event *ordering.OrderCancelledEvent) error {
	cmd := &receipting.CancelReceiptByRefIDCommand{
		UpdatedBy: event.UpdatedBy,
		ShopID:    event.ShopID,
		RefID:     event.OrderID,
		RefType:   receipt_ref.Order,
	}
	err := m.receiptAggr.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (m *ProcessManager) MoneyTransactionConfirmed(ctx context.Context, event *moneytx.MoneyTxShippingConfirmedEvent) error {
	var (
		ledgerID         dot.ID
		totalShippingFee int
		fulfillments     []*shippingcore.FulfillmentExtended
		orderIDs         []dot.ID
	)
	mapOrderAndReceivedAmount := make(map[dot.ID]int)
	mapOrder := make(map[dot.ID]ordermodelx.OrderWithFulfillments)
	mapOrderFulfillment := make(map[dot.ID]*shippingcore.FulfillmentExtended)

	getMoneyTransaction := &moneytx.GetMoneyTxShippingByIDQuery{
		MoneyTxShippingID: event.MoneyTxShippingID,
		ShopID:            event.ShopID,
	}
	if err := m.moneyTxQuery.Dispatch(ctx, getMoneyTransaction); err != nil {
		return err
	}

	ffmQuery := &shippingcore.ListFulfillmentExtendedsByMoneyTxShippingIDQuery{
		ShopID:            event.ShopID,
		MoneyTxShippingID: event.MoneyTxShippingID,
	}
	if err := m.shippingQuery.Dispatch(ctx, ffmQuery); err != nil {
		return err
	}

	for _, fulfillment := range ffmQuery.Result {
		fulfillments = append(fulfillments, fulfillment)
		orderIDs = append(orderIDs, fulfillment.OrderID)
		totalShippingFee += fulfillment.ShippingFeeShop
		mapOrderFulfillment[fulfillment.OrderID] = fulfillment
	}

	if len(orderIDs) == 0 {
		return nil
	}
	// s??? ti???n th???c t??? c???a m???i ????n h??ng
	getOrdersQuery := &ordermodelx.GetOrdersQuery{
		IDs: orderIDs,
	}
	if err := m.OrderStore.GetOrders(ctx, getOrdersQuery); err != nil {
		return err
	}
	for _, order := range getOrdersQuery.Result.Orders {
		mapOrderAndReceivedAmount[order.ID] = 0
		mapOrder[order.ID] = order
	}
	// T??nh ReceivedAmount cho m???i Order
	getReceiptsByOrderIDs := &receipting.ListReceiptsByRefsAndStatusQuery{
		ShopID:  event.ShopID,
		RefIDs:  orderIDs,
		RefType: receipt_ref.Order,
		Status:  int(status3.P),
	}
	if err := m.receiptQuery.Dispatch(ctx, getReceiptsByOrderIDs); err != nil {
		return err
	}
	for _, receipt := range getReceiptsByOrderIDs.Result.Receipts {
		if receipt.RefType != receipt_ref.Order {
			continue
		}
		for _, receiptLine := range receipt.Lines {
			if receiptLine.RefID == 0 {
				continue
			}
			if _, ok := mapOrderAndReceivedAmount[receiptLine.RefID]; ok {
				switch receipt.Type {
				case receipt_type.Receipt:
					mapOrderAndReceivedAmount[receiptLine.RefID] += receiptLine.Amount
				case receipt_type.Payment:
					mapOrderAndReceivedAmount[receiptLine.RefID] -= receiptLine.Amount
				}
			}
		}
	}

	// Get bank_account
	bankAccount := getMoneyTransaction.Result.BankAccount
	if bankAccount == nil {
		// get shop bankAccount
		query := &identity.GetShopByIDQuery{
			ID: event.ShopID,
		}
		if err := m.identityQuery.Dispatch(ctx, query); err != nil {
			return err
		}
		bankAccount = query.Result.BankAccount
	}
	if bankAccount == nil {
		// B??? qua tr?????ng h???p kh??ng t??m th???y s??? qu???
		// M???t s??? shop n???p ti???n tr?????c (credit) ????? x??i n??n kh??ng c???p nh???t th??ng tin t??i kho???n ng??n h??ng
		// -> Gi???i ph??p t???m th???i: b??? qua, ko t???o receipt
		ll.Error("MoneyTxShippingConfirmedEvent failed: kh??ng t??m th???y t??i kho???n ng??n h??ng", l.ID("shop_id", event.ShopID), l.ID("money_transaction_id", event.MoneyTxShippingID))
		return nil
	}
	ledgerID, err := m.getOrCreateLedgerID(ctx, bankAccount, event.ShopID)
	if err != nil {
		return cm.Errorf(cm.NotFound, err, "Kh??ng t??m th???y s??? qu???").WithMetap("shop_id", event.ShopID).WithMetap("money_transaction_id", event.MoneyTxShippingID)
	}

	if err := m.createReceipts(ctx, mapOrderFulfillment, mapOrderAndReceivedAmount, mapOrder, event.ShopID, ledgerID); err != nil {
		return cm.Errorf(cm.FailedPrecondition, err, "T???o phi???u thu th???t b???i (%v)", err.Error()).WithMetap("shop_id", event.ShopID).WithMetap("money_transaction_id", event.MoneyTxShippingID)
	}

	// Create receipt type payment
	if err := m.createPayment(ctx, totalShippingFee, fulfillments, event.ShopID, ledgerID); err != nil {
		return cm.Errorf(cm.FailedPrecondition, err, "T???o phi???u chi th???t b???i (%v)", err.Error()).WithMetap("shop_id", event.ShopID).WithMetap("money_transaction_id", event.MoneyTxShippingID)
	}

	return nil
}

func (m *ProcessManager) createPayment(
	ctx context.Context,
	totalShippingFee int,
	fulfillments []*shippingcore.FulfillmentExtended,
	shopID, ledgerID dot.ID,
) error {
	receiptLines := []*receipting.ReceiptLine{}
	refIDs := []dot.ID{}
	for _, fulfillment := range fulfillments {
		shippingFee := fulfillment.ShippingFeeShop
		if shippingFee == 0 {
			continue
		}
		receiptLines = append(receiptLines, &receipting.ReceiptLine{
			RefID:  fulfillment.ID,
			Amount: shippingFee,
		})
		refIDs = append(refIDs, fulfillment.ID)
	}
	if len(receiptLines) == 0 {
		return nil
	}

	query := &receipting.ListReceiptsByRefsAndStatusQuery{
		ShopID:     shopID,
		RefIDs:     refIDs,
		RefType:    receipt_ref.Fulfillment,
		Status:     int(status3.P),
		IsContains: true,
	}
	if err := m.receiptQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	if len(query.Result.Receipts) > 0 {
		// ???? t???o phi???u thanh to??n ph?? v???n chuy???n Topship
		// Kh??ng t???o th??m n???a
		return nil
	}

	// T???o m???i phi???u thanh to??n ph?? v???n chuy???n
	cmd := &receipting.CreateReceiptCommand{
		ShopID:      shopID,
		TraderID:    idutil.TopShipID,
		Title:       "Thanh to??n ph?? v???n chuy???n Topship",
		Description: "Phi???u ???????c t???o t??? ?????ng qua th??ng qua ?????i so??t Topship",
		Type:        receipt_type.Payment,
		Status:      int(status3.P),
		Amount:      totalShippingFee,
		LedgerID:    ledgerID,
		Lines:       receiptLines,
		PaidAt:      time.Now(),
		Mode:        receipt_mode.Auto,
		ConfirmedAt: time.Now(),
		RefType:     receipt_ref.Fulfillment,
	}
	if err := m.receiptAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	return nil
}

func (m *ProcessManager) createReceipts(
	ctx context.Context,
	mapOrderFulfillment map[dot.ID]*shippingcore.FulfillmentExtended,
	mapOrderAndReceivedAmount map[dot.ID]int,
	mapOrder map[dot.ID]ordermodelx.OrderWithFulfillments,
	shopID, ledgerID dot.ID) error {
	for key, order := range mapOrder {
		totalAmount := order.TotalAmount
		// s??? ti???n thu th???c t???
		receiptCOD := 0
		title := ""
		// remainAmount: S??? ti???n c??n l???i (c???n thu)
		remainAmount := totalAmount - mapOrderAndReceivedAmount[key]
		ffm := mapOrderFulfillment[key]
		switch ffm.ShippingState {
		case shipping.Delivered:
			receiptCOD = ffm.TotalCODAmount
			title = "Nh???n thu h??? ????n h??ng giao th??nh c??ng"
		case shipping.Undeliverable:
			receiptCOD = ffm.ActualCompensationAmount
			title = "Nh???n ho??n ti???n ????n m???t h??ng"
		default:
			continue
		}
		receiptLines := []*receipting.ReceiptLine{}
		// TH: Ti???n thu th???c t??? > s??? ti???n c???n thu -> c???ng l???i v??o t??i kho???n c???a kh??ch s??? ti???n = Ti???n thu th???c t??? - s??? ti???n c???n thu
		if remainAmount == 0 || receiptCOD == 0 {
			continue
		}
		if receiptCOD <= remainAmount {
			line := &receipting.ReceiptLine{
				RefID:  key,
				Amount: receiptCOD,
			}
			receiptLines = append(receiptLines, line)
		} else {
			line1 := &receipting.ReceiptLine{
				RefID:  key,
				Amount: remainAmount,
			}
			line2 := &receipting.ReceiptLine{
				Title:  "C???ng v??o t??i kho???n kh??ch h??ng",
				Amount: receiptCOD - remainAmount,
			}
			receiptLines = append(receiptLines, line1, line2)
		}

		traderID := order.CustomerID
		if traderID == 0 {
			traderID = customering.CustomerAnonymous
		}

		cmd := &receipting.CreateReceiptCommand{
			ShopID:      shopID,
			TraderID:    traderID,
			Title:       title,
			Description: "Phi???u ???????c t???o t??? ?????ng qua th??ng qua ?????i so??t Topship",
			Type:        receipt_type.Receipt,
			Status:      int(status3.P),
			Amount:      receiptCOD,
			LedgerID:    ledgerID,
			RefIDs:      []dot.ID{key},
			RefType:     receipt_ref.Order,
			Lines:       receiptLines,
			PaidAt:      time.Now(),
			Mode:        receipt_mode.Auto,
			ConfirmedAt: time.Now(),
		}
		if err := m.receiptAggr.Dispatch(ctx, cmd); err != nil {
			return err
		}
	}
	return nil
}

func (m *ProcessManager) getOrCreateLedgerID(ctx context.Context, bankAccount *identitytypes.BankAccount, shopID dot.ID) (dot.ID, error) {
	if bankAccount == nil {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Thi???u th??ng tin t??i kho???n ng??n h??ng. Vui l??ng ki???m tra l???i")
	}
	// Check accountNumber exists into ledgers, if it isn't then create
	query := &ledgering.GetLedgerByAccountNumberQuery{
		AccountNumber: bankAccount.AccountNumber,
		ShopID:        shopID,
	}
	err := m.ledgerQuery.Dispatch(ctx, query)
	switch cm.ErrorCode(err) {
	case cm.NotFound:
		// Create ledger
		cmd := &ledgering.CreateLedgerCommand{
			ShopID:      shopID,
			Name:        fmt.Sprintf("[%v] %v", bankAccount.Branch, bankAccount.AccountName),
			BankAccount: bankAccount,
			Note:        "T??i kho???n thanh to??n t??? t???o",
			Type:        ledger_type.LedgerTypeBank,
		}
		if err := m.ledgerAggr.Dispatch(ctx, cmd); err != nil {
			return 0, err
		}
		return cmd.Result.ID, nil
	case cm.NoError:
		return query.Result.ID, nil
	default:
		return 0, err
	}
}

func (m *ProcessManager) MoneyTxShippingEtopConfirmed(ctx context.Context, event *moneytx.MoneyTxShippingEtopConfirmedEvent) error {
	query := &moneytx.ListMoneyTxShippingsQuery{
		MoneyTxShippingEtopID: event.MoneyTxShippingEtopID,
	}
	if err := m.moneyTxQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	moneyTxs := query.Result.MoneyTxShippings
	for _, moneyTx := range moneyTxs {
		cmd := &moneytx.MoneyTxShippingConfirmedEvent{
			ShopID:            moneyTx.ShopID,
			MoneyTxShippingID: moneyTx.ID,
		}
		if err := m.MoneyTransactionConfirmed(ctx, cmd); err != nil {
			return err
		}
	}
	return nil
}
