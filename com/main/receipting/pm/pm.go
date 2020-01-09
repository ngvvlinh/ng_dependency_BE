package pm

import (
	"context"
	"fmt"
	"time"

	"etop.vn/api/main/identity"
	"etop.vn/api/main/ledgering"
	"etop.vn/api/main/moneytx"
	"etop.vn/api/main/receipting"
	"etop.vn/api/shopping/customering"
	"etop.vn/api/top/types/etc/ledger_type"
	"etop.vn/api/top/types/etc/receipt_mode"
	"etop.vn/api/top/types/etc/receipt_ref"
	"etop.vn/api/top/types/etc/receipt_type"
	"etop.vn/api/top/types/etc/shipping"
	"etop.vn/api/top/types/etc/status3"
	identityconvert "etop.vn/backend/com/main/identity/convert"
	identitysharemodel "etop.vn/backend/com/main/identity/sharemodel"
	"etop.vn/backend/com/main/moneytx/modelx"
	txmodelx "etop.vn/backend/com/main/moneytx/modelx"
	ordermodelx "etop.vn/backend/com/main/ordering/modelx"
	"etop.vn/backend/com/main/shipping/modely"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/capi"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

type ProcessManager struct {
	eventBus capi.EventBus

	receiptQuery  receipting.QueryBus
	receiptAggr   receipting.CommandBus
	ledgerQuery   ledgering.QueryBus
	ledgerAggr    ledgering.CommandBus
	identityQuery identity.QueryBus
}

func New(
	eventBus capi.EventBus,
	receiptQuery receipting.QueryBus,
	receiptAggregate receipting.CommandBus,
	ledgerQuery ledgering.QueryBus,
	ledgerAggregate ledgering.CommandBus,
	identityQuery identity.QueryBus,
) *ProcessManager {
	return &ProcessManager{
		eventBus:      eventBus,
		receiptQuery:  receiptQuery,
		receiptAggr:   receiptAggregate,
		ledgerQuery:   ledgerQuery,
		ledgerAggr:    ledgerAggregate,
		identityQuery: identityQuery,
	}
}

var (
	ll = l.New()
)

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.MoneyTransactionConfirmed)
	eventBus.AddEventListener(m.MoneyTxShippingEtopConfirmed)
}

func (m *ProcessManager) MoneyTransactionConfirmed(ctx context.Context, event *moneytx.MoneyTxShippingConfirmedEvent) error {
	var (
		ledgerID         dot.ID
		totalShippingFee int
		fulfillments     []*modely.FulfillmentExtended
		orderIDs         []dot.ID
	)
	mapOrderAndReceivedAmount := make(map[dot.ID]int)
	mapOrder := make(map[dot.ID]ordermodelx.OrderWithFulfillments)
	mapOrderFulfillment := make(map[dot.ID]*modely.FulfillmentExtended)

	getMoneyTransaction := &modelx.GetMoneyTransaction{
		ID:     event.MoneyTxShippingID,
		ShopID: event.ShopID,
	}
	if err := sqlstore.GetMoneyTransaction(ctx, getMoneyTransaction); err != nil {
		return err
	}
	for _, fulfillment := range getMoneyTransaction.Result.Fulfillments {
		fulfillments = append(fulfillments, fulfillment)
		orderIDs = append(orderIDs, fulfillment.OrderID)
		totalShippingFee += fulfillment.ShippingFeeShop
		mapOrderFulfillment[fulfillment.OrderID] = fulfillment
	}

	if len(orderIDs) == 0 {
		return nil
	}
	// số tiền thực tế của mỗi đơn hàng
	getOrdersQuery := &ordermodelx.GetOrdersQuery{
		IDs: orderIDs,
	}
	if err := bus.Dispatch(ctx, getOrdersQuery); err != nil {
		return err
	}
	for _, order := range getOrdersQuery.Result.Orders {
		mapOrderAndReceivedAmount[order.ID] = 0
		mapOrder[order.ID] = order
	}
	// Tính ReceivedAmount cho mỗi Order
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
		shopBankAccount := query.Result.BankAccount
		bankAccount = identityconvert.Convert_identitytypes_BankAccount_sharemodel_BankAccount(shopBankAccount, bankAccount)
	}
	if bankAccount == nil {
		// Bỏ qua trường hợp không tìm thấy sổ quỹ
		// Một số shop nạp tiền trước (credit) để xài nên không cập nhật thông tin tài khoản ngân hàng
		// -> Giải pháp tạm thời: bỏ qua, ko tạo receipt
		ll.Error("MoneyTxShippingConfirmedEvent failed: không tìm thấy tài khoản ngân hàng", l.ID("shop_id", event.ShopID), l.ID("money_transaction_id", event.MoneyTxShippingID))
		return nil
	}
	ledgerID, err := m.getOrCreateLedgerID(ctx, bankAccount, event.ShopID)
	if err != nil {
		return cm.Errorf(cm.NotFound, err, "Không tìm thấy sổ quỹ").WithMetap("shop_id", event.ShopID).WithMetap("money_transaction_id", event.MoneyTxShippingID)
	}

	if err := m.createReceipts(ctx, mapOrderFulfillment, mapOrderAndReceivedAmount, mapOrder, event.ShopID, ledgerID); err != nil {
		return cm.Errorf(cm.FailedPrecondition, err, "Tạo phiếu thu thất bại (%v)", err.Error()).WithMetap("shop_id", event.ShopID).WithMetap("money_transaction_id", event.MoneyTxShippingID)
	}

	// Create receipt type payment
	if err := m.createPayment(ctx, totalShippingFee, fulfillments, event.ShopID, ledgerID); err != nil {
		return cm.Errorf(cm.FailedPrecondition, err, "Tạo phiếu chi thất bại (%v)", err.Error()).WithMetap("shop_id", event.ShopID).WithMetap("money_transaction_id", event.MoneyTxShippingID)
	}

	return nil
}

func (m *ProcessManager) createPayment(
	ctx context.Context,
	totalShippingFee int,
	fulfillments []*modely.FulfillmentExtended,
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
		// Đã tạo phiếu thanh toán phí vận chuyển Topship
		// Không tạo thêm nữa
		return nil
	}

	// Tạo mới phiếu thanh toán phí vận chuyển
	cmd := &receipting.CreateReceiptCommand{
		ShopID:      shopID,
		TraderID:    model.TopShipID,
		Title:       "Thanh toán phí vận chuyển Topship",
		Description: "Phiếu được tạo tự động qua thông qua đối soát Topship",
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
	mapOrderFulfillment map[dot.ID]*modely.FulfillmentExtended,
	mapOrderAndReceivedAmount map[dot.ID]int,
	mapOrder map[dot.ID]ordermodelx.OrderWithFulfillments,
	shopID, ledgerID dot.ID) error {
	for key, order := range mapOrder {
		totalAmount := order.TotalAmount
		// số tiền thu thực tế
		receiptCOD := 0
		title := ""
		// remainAmount: Số tiền còn lại (cần thu)
		remainAmount := totalAmount - mapOrderAndReceivedAmount[key]
		ffm := mapOrderFulfillment[key]
		switch ffm.ShippingState {
		case shipping.Delivered:
			receiptCOD = ffm.TotalCODAmount
			title = "Nhận thu hộ đơn hàng giao thành công"
		case shipping.Undeliverable:
			receiptCOD = ffm.ActualCompensationAmount
			title = "Nhận hoàn tiền đơn mất hàng"
		default:
			continue
		}
		receiptLines := []*receipting.ReceiptLine{}
		// TH: Tiền thu thực tế > số tiền cần thu -> cộng lại vào tài khoản của khách số tiền = Tiền thu thực tế - số tiền cần thu
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
				Title:  "Cộng vào tài khoản khách hàng",
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
			Description: "Phiếu được tạo tự động qua thông qua đối soát Topship",
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

func (m *ProcessManager) getOrCreateLedgerID(ctx context.Context, bankAccount *identitysharemodel.BankAccount, shopID dot.ID) (dot.ID, error) {
	if bankAccount == nil {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Thiếu thông tin tài khoản ngân hàng. Vui lòng kiểm tra lại")
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
			BankAccount: identityconvert.BankAccount(bankAccount),
			Note:        "Tài khoản thanh toán tự tạo",
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
	query := &txmodelx.GetMoneyTxsByMoneyTxShippingEtopID{
		MoneyTxShippingEtopID: event.MoneyTxShippingEtopID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	moneyTxs := query.Result.MoneyTransactions
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
