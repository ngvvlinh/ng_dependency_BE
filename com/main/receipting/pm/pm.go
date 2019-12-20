package pm

import (
	"context"
	"fmt"
	"time"

	"etop.vn/api/main/ledgering"
	"etop.vn/api/main/receipting"
	"etop.vn/api/top/types/etc/ledger_type"
	"etop.vn/api/top/types/etc/receipt_mode"
	"etop.vn/api/top/types/etc/receipt_ref"
	"etop.vn/api/top/types/etc/receipt_type"
	"etop.vn/api/top/types/etc/status3"
	identityconvert "etop.vn/backend/com/main/identity/convert"
	"etop.vn/backend/com/main/moneytx/modelx"
	ordermodelx "etop.vn/backend/com/main/ordering/modelx"
	"etop.vn/backend/com/main/shipping/modely"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
	etopmodel "etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/capi"
	"etop.vn/capi/dot"
)

type ProcessManager struct {
	eventBus capi.EventBus

	receiptQuery receipting.QueryBus
	receiptAggr  receipting.CommandBus
	ledgerQuery  ledgering.QueryBus
	ledgerAggr   ledgering.CommandBus
}

func New(
	eventBus capi.EventBus,
	receiptQuery receipting.QueryBus,
	receiptAggregate receipting.CommandBus,
	ledgerQuery ledgering.QueryBus,
	ledgerAggregate ledgering.CommandBus,
) *ProcessManager {
	return &ProcessManager{
		eventBus:     eventBus,
		receiptQuery: receiptQuery,
		receiptAggr:  receiptAggregate,
		ledgerQuery:  ledgerQuery,
		ledgerAggr:   ledgerAggregate,
	}
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.MoneyTransactionConfirmed)
}

func (m *ProcessManager) MoneyTransactionConfirmed(ctx context.Context, event *receipting.MoneyTransactionConfirmedEvent) error {
	var (
		ledgerID         dot.ID
		totalShippingFee int
		fulfillments     []*modely.FulfillmentExtended
		orderIDs         []dot.ID
	)

	mapOrderAndTotalAmount := make(map[dot.ID]int)
	mapOrderAndReceivedAmount := make(map[dot.ID]int)
	mapOrder := make(map[dot.ID]ordermodelx.OrderWithFulfillments)

	getMoneyTransaction := &modelx.GetMoneyTransaction{
		ID:     event.MoneyTransactionID,
		ShopID: event.ShopID,
	}
	if err := sqlstore.GetMoneyTransaction(ctx, getMoneyTransaction); err != nil {
		return err
	}
	for _, fulfillment := range getMoneyTransaction.Result.Fulfillments {
		fulfillments = append(fulfillments, fulfillment)
		orderIDs = append(orderIDs, fulfillment.OrderID)
		totalShippingFee += fulfillment.ShippingFeeShop
	}

	if len(orderIDs) == 0 {
		return nil
	}

	getOrdersQuery := &ordermodelx.GetOrdersQuery{
		IDs: orderIDs,
	}
	if err := bus.Dispatch(ctx, getOrdersQuery); err != nil {
		return err
	}
	for _, order := range getOrdersQuery.Result.Orders {
		mapOrderAndTotalAmount[order.ID] = order.TotalAmount
		mapOrderAndReceivedAmount[order.ID] = 0
		mapOrder[order.ID] = order
	}

	getReceiptsByOrderIDs := &receipting.ListReceiptsByRefsAndStatusQuery{
		ShopID:  event.ShopID,
		RefIDs:  orderIDs,
		RefType: receipt_ref.ReceiptRefTypeOrder,
		Status:  int(status3.P),
	}
	if err := m.receiptQuery.Dispatch(ctx, getReceiptsByOrderIDs); err != nil {
		return err
	}
	for _, receipt := range getReceiptsByOrderIDs.Result.Receipts {
		if receipt.RefType != receipt_ref.ReceiptRefTypeOrder {
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
	var bankAccount *model.BankAccount
	var haveBankAccount bool
	ledgerID, err := m.getOrCreateBankAccount(getMoneyTransaction, bankAccount, haveBankAccount, event.ShopID, ctx, ledgerID)
	if err != nil {
		return err
	}

	if err := createReceipts(mapOrderAndTotalAmount, mapOrderAndReceivedAmount, mapOrder, event.ShopID, ledgerID, m, ctx); err != nil {
		return err
	}

	// Create receipt type payment
	if err := m.createPayment(totalShippingFee, fulfillments, event.ShopID, ledgerID, ctx); err != nil {
		return err
	}

	return nil
}

func (m *ProcessManager) createPayment(
	totalShippingFee int, fulfillments []*modely.FulfillmentExtended, shopID, ledgerID dot.ID, ctx context.Context,
) error {
	{
		receiptLines := []*receipting.ReceiptLine{}
		for _, fulfillment := range fulfillments {
			receiptLines = append(receiptLines, &receipting.ReceiptLine{
				RefID:  fulfillment.ID,
				Amount: fulfillment.ShippingFeeShop,
			})
		}

		cmd := &receipting.CreateReceiptCommand{
			ShopID:      shopID,
			TraderID:    model.TopShipID,
			Title:       "Thanh toán phí vận chuyển Topship",
			Description: "Phiếu được tạo tự động qua thông qua đối soát Topship",
			Type:        receipt_type.Payment,
			Status:      int(status3.P),
			Amount:      totalShippingFee,
			LedgerID:    ledgerID,
			RefType:     receipt_ref.ReceiptRefTypeFulfillment,
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

func createReceipts(
	mapOrderAndTotalAmount map[dot.ID]int, mapOrderAndReceivedAmount map[dot.ID]int,
	mapOrder map[dot.ID]ordermodelx.OrderWithFulfillments, shopID, ledgerID dot.ID,
	m *ProcessManager, ctx context.Context) error {
	for key, value := range mapOrderAndTotalAmount {
		if value-mapOrderAndReceivedAmount[key] == 0 {
			continue
		}

		receiptLines := []*receipting.ReceiptLine{}
		receiptLines = append(receiptLines, &receipting.ReceiptLine{
			RefID:  key,
			Amount: value - mapOrderAndReceivedAmount[key],
		})

		traderID := mapOrder[key].CustomerID
		if traderID == 0 {
			traderID = model.IndependentCustomerID
		}

		cmd := &receipting.CreateReceiptCommand{
			ShopID:      shopID,
			TraderID:    traderID,
			Title:       "Thanh toán đơn hàng",
			Description: "Phiếu được tạo tự động qua thông qua đối soát Topship",
			Type:        receipt_type.Receipt,
			Status:      int(status3.P),
			Amount:      value - mapOrderAndReceivedAmount[key],
			LedgerID:    ledgerID,
			RefIDs:      []dot.ID{key},
			RefType:     receipt_ref.ReceiptRefTypeOrder,
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

func (m *ProcessManager) getOrCreateBankAccount(
	getMoneyTransaction *modelx.GetMoneyTransaction, bankAccount *etopmodel.BankAccount,
	haveBankAccount bool, shopID dot.ID, ctx context.Context, ledgerID dot.ID) (dot.ID, error) {
	if getMoneyTransaction.Result.BankAccount != nil {
		bankAccount = getMoneyTransaction.Result.BankAccount
		haveBankAccount = true
	}
	// Get default ledger (cash) when haven't bankAccount
	if !haveBankAccount {
		query := &ledgering.ListLedgersByTypeQuery{
			LedgerType: ledger_type.LedgerTypeCash,
			ShopID:     shopID,
		}
		if err := m.ledgerQuery.Dispatch(ctx, query); err != nil {
			return 0, cm.MapError(err).
				Wrap(cm.NotFound, "tài khoản thanh toán tiền mặt không tìm thấy").
				Throw()
		}
		ledgerID = query.Result.Ledgers[0].ID
	} else { // Check accountNumber exists into ledgers, if it isn't then create
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
				Note:        "Sổ quỹ tự tạo",
				Type:        ledger_type.LedgerTypeBank,
			}
			if err := m.ledgerAggr.Dispatch(ctx, cmd); err != nil {
				return 0, err
			}
			ledgerID = cmd.Result.ID
		case cm.NoError:
			ledgerID = query.Result.ID
		default:
			return 0, err
		}
	}
	return ledgerID, nil
}
