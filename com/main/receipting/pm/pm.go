package pm

import (
	"context"
	"fmt"
	"time"

	"etop.vn/api/main/ledgering"
	"etop.vn/api/main/receipting"
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
		ledgerID         int64
		totalShippingFee int32
		fulfillments     []*modely.FulfillmentExtended
		orderIDs         []int64
	)

	mapOrderAndTotalAmount := make(map[int64]int)
	mapOrderAndReceivedAmount := make(map[int64]int32)
	mapOrder := make(map[int64]ordermodelx.OrderWithFulfillments)

	getMoneyTransaction := &modelx.GetMoneyTransaction{
		ID:     event.MoneyTransactionID,
		ShopID: event.ShopID,
	}
	if err := sqlstore.GetMoneyTransaction(ctx, getMoneyTransaction); err != nil {
		return err
	}
	for _, fulfillment := range getMoneyTransaction.Result.Fulfillments {
		if fulfillment.ShippingState == etopmodel.StateDelivered && fulfillment.ShippingFeeShopTransferedAt.IsZero() {
			fulfillments = append(fulfillments, fulfillment)
			orderIDs = append(orderIDs, fulfillment.OrderID)
			totalShippingFee += int32(fulfillment.ShippingFeeShop)
		}
	}

	getOrdersQuery := &ordermodelx.GetOrdersQuery{
		IDs: orderIDs,
	}
	if err := bus.Dispatch(ctx, getOrdersQuery); err != nil {
		return err
	}
	for _, order := range getOrdersQuery.Result.Orders {
		mapOrderAndTotalAmount[order.ID] = order.TotalAmount
		mapOrder[order.ID] = order
	}

	getReceiptsByOrderIDs := &receipting.ListReceiptsByRefIDsAndStatusQuery{
		ShopID: event.ShopID,
		RefIDs: orderIDs,
		Status: int32(etopmodel.S3Positive),
	}
	if err := m.receiptQuery.Dispatch(ctx, getReceiptsByOrderIDs); err != nil {
		return err
	}
	for _, receipt := range getReceiptsByOrderIDs.Result.Receipts {
		for _, receiptLine := range receipt.Lines {
			if receiptLine.RefID == 0 {
				continue
			}
			if _, ok := mapOrderAndReceivedAmount[receiptLine.RefID]; ok {
				switch receipt.Type {
				case receipting.ReceiptTypeReceipt:
					mapOrderAndReceivedAmount[receiptLine.RefID] += receiptLine.Amount
				case receipting.ReceiptTypePayment:
					mapOrderAndReceivedAmount[receiptLine.RefID] -= receiptLine.Amount
				}
			}
		}
	}

	// Get bank_account
	var bankAccount *model.BankAccount
	var haveBankAccount bool
	ledgerID, err := m.getOrCreateBankAccount(getMoneyTransaction, bankAccount, haveBankAccount, event, ctx, ledgerID)
	if err != nil {
		return err
	}

	if err := createReceipts(mapOrderAndTotalAmount, mapOrderAndReceivedAmount, mapOrder, event, ledgerID, m, ctx); err != nil {
		return err
	}

	// Create receipt type payment
	if err := m.createPayment(totalShippingFee, event, ledgerID, ctx); err != nil {
		return err
	}

	return nil
}

func (m *ProcessManager) createPayment(totalShippingFee int32, event *receipting.MoneyTransactionConfirmedEvent, ledgerID int64, ctx context.Context) error {
	{
		receiptLines := []*receipting.ReceiptLine{}
		receiptLines = append(receiptLines, &receipting.ReceiptLine{
			Title:  "Thanh toán phí vận chuyển",
			Amount: totalShippingFee,
		})

		cmd := &receipting.CreateReceiptCommand{
			ShopID:      event.ShopID,
			TraderID:    model.TopShipID,
			Title:       "Thanh toán phí vận chuyển Topship",
			Description: "Phiếu được tạo tự động qua thông qua đối soát Topship",
			Type:        receipting.ReceiptTypePayment,
			Status:      int32(etopmodel.S3Positive),
			Amount:      totalShippingFee,
			LedgerID:    ledgerID,
			Lines:       receiptLines,
			PaidAt:      time.Now(),
			CreatedType: receipting.ReceiptCreatedTypeAuto,
			ConfirmedAt: time.Now(),
		}
		if err := m.receiptAggr.Dispatch(ctx, cmd); err != nil {
			return err
		}
	}
	return nil
}

func createReceipts(mapOrderAndTotalAmount map[int64]int, mapOrderAndReceivedAmount map[int64]int32, mapOrder map[int64]ordermodelx.OrderWithFulfillments, event *receipting.MoneyTransactionConfirmedEvent, ledgerID int64, m *ProcessManager, ctx context.Context) error {
	for key, value := range mapOrderAndTotalAmount {
		if int32(value)-mapOrderAndReceivedAmount[key] == 0 {
			continue
		}

		receiptLines := []*receipting.ReceiptLine{}
		receiptLines = append(receiptLines, &receipting.ReceiptLine{
			RefID:  key,
			Amount: int32(value) - mapOrderAndReceivedAmount[key],
		})

		traderID := mapOrder[key].CustomerID
		if traderID == 0 {
			traderID = model.IndependentCustomerID
		}

		cmd := &receipting.CreateReceiptCommand{
			ShopID:      event.ShopID,
			TraderID:    traderID,
			Title:       "Thanh toán đơn hàng",
			Description: "Phiếu được tạo tự động qua thông qua đối soát Topship",
			Type:        receipting.ReceiptTypeReceipt,
			Status:      int32(etopmodel.S3Positive),
			Amount:      int32(value) - mapOrderAndReceivedAmount[key],
			LedgerID:    ledgerID,
			RefIDs:      []int64{key},
			Lines:       receiptLines,
			PaidAt:      time.Now(),
			CreatedType: receipting.ReceiptCreatedTypeAuto,
			ConfirmedAt: time.Now(),
		}
		if err := m.receiptAggr.Dispatch(ctx, cmd); err != nil {
			return err
		}
	}
	return nil
}

func (m *ProcessManager) getOrCreateBankAccount(getMoneyTransaction *modelx.GetMoneyTransaction, bankAccount *etopmodel.BankAccount, haveBankAccount bool, event *receipting.MoneyTransactionConfirmedEvent, ctx context.Context, ledgerID int64) (int64, error) {
	if getMoneyTransaction.Result.BankAccount != nil {
		bankAccount = getMoneyTransaction.Result.BankAccount
		haveBankAccount = true
	}
	// Get default ledger (cash) when haven't bankAccount
	if !haveBankAccount {
		query := &ledgering.ListLedgersByTypeQuery{
			LedgerType: ledgering.LedgerTypeCash,
			ShopID:     event.ShopID,
		}
		if err := m.ledgerQuery.Dispatch(ctx, query); err != nil {
			return 0, err
		}
		ledgerID = query.Result.Ledgers[0].ID
	} else { // Check accountNumber exists into ledgers, if it isn't then create
		query := &ledgering.GetLedgerByAccountNumberQuery{
			AccountNumber: bankAccount.AccountNumber,
			ShopID:        event.ShopID,
		}
		err := m.ledgerQuery.Dispatch(ctx, query)
		switch cm.ErrorCode(err) {
		case cm.NotFound:
			// Create ledger
			cmd := &ledgering.CreateLedgerCommand{
				ShopID:      event.ShopID,
				Name:        fmt.Sprintf("[%v] %v", bankAccount.Branch, bankAccount.AccountName),
				BankAccount: identityconvert.BankAccount(bankAccount),
				Note:        "Sổ quỹ tự tạo",
				Type:        ledgering.LedgerTypeBank,
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
