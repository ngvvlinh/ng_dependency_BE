package pm

import (
	"context"

	"o.o/api/main/purchaseorder"
	"o.o/api/main/purchaserefund"
	"o.o/api/main/receipting"
	"o.o/api/top/types/etc/receipt_ref"
	"o.o/api/top/types/etc/receipt_type"
	"o.o/api/top/types/etc/status3"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/capi/dot"
)

type ProcessManager struct {
	purchaseRefundQuery     *purchaserefund.QueryBus
	purchaseRefundAggregate *purchaserefund.CommandBus
	receiptQuery            *receipting.QueryBus
}

func New(
	purchaseRefundA *purchaserefund.CommandBus,
	purchaseRefundQ *purchaserefund.QueryBus,
	receiptQ *receipting.QueryBus,

) *ProcessManager {
	return &ProcessManager{
		purchaseRefundAggregate: purchaseRefundA,
		purchaseRefundQuery:     purchaseRefundQ,
		receiptQuery:            receiptQ,
	}
}
func (p *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(p.ReceiptCreating)
	eventBus.AddEventListener(p.PurchaseOrderCancelledEvent)
}

func (p *ProcessManager) PurchaseOrderCancelledEvent(ctx context.Context, event *purchaseorder.PurchaseOrderCancelledEvent) error {
	query := &purchaserefund.GetPurchaseRefundsByPurchaseOrderIDQuery{
		PurchaseOrderID: event.PurchaseOrderID,
		ShopID:          event.ShopID,
	}
	err := p.purchaseRefundQuery.Dispatch(ctx, query)
	if err != nil {
		return err
	}
	for _, value := range query.Result {
		cmd := &purchaserefund.CancelPurchaseRefundCommand{
			ShopID:               event.ShopID,
			ID:                   value.ID,
			UpdatedBy:            event.UpdatedBy,
			CancelReason:         "Cancel Purchase Order",
			AutoInventoryVoucher: event.AutoInventoryVoucher,
			InventoryOverStock:   event.InventoryOverStock,
		}
		err = p.purchaseRefundAggregate.Dispatch(ctx, cmd)
		if err != nil {
			return err
		}
	}
	return nil
}
func (p *ProcessManager) ReceiptCreating(
	ctx context.Context, event *receipting.ReceiptCreatingEvent,
) error {
	var purchaseRefunds []*purchaserefund.PurchaseRefund
	mRefund := make(map[dot.ID]*purchaserefund.PurchaseRefund)
	refIDs := event.RefIDs
	receipt := event.Receipt
	mapRefIDAmount := event.MapRefIDAmount

	if receipt.RefType != receipt_ref.PurchaseRefund {
		return nil
	}

	// List Refund on refIDs
	query := &purchaserefund.GetPurchaseRefundsByIDsQuery{
		IDs:    refIDs,
		ShopID: receipt.ShopID,
	}
	if err := p.purchaseRefundQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	purchaseRefunds = query.Result
	for _, v := range purchaseRefunds {
		if v.Status == status3.N {
			continue
		}
		mRefund[v.ID] = v
		if v.SupplierID != receipt.TraderID {
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn trả hàng %v không thuộc đối tác đã chọn", v.Code)
		}
	}

	// Check refIDs and orderIDs (result of query above)
	if len(refIDs) != len(purchaseRefunds) {
		for _, refID := range refIDs {
			if _, ok := mRefund[refID]; !ok {
				return cm.Errorf(cm.FailedPrecondition, nil, "ref_id %d không tìm thấy", refID)
			}
		}
	}

	// List receipts by refIDs
	listReceiptsQuery := &receipting.ListReceiptsByRefsAndStatusQuery{
		ShopID:  receipt.ShopID,
		RefIDs:  refIDs,
		RefType: receipt_ref.Refund,
		Status:  int(status3.P),
	}
	if err := p.receiptQuery.Dispatch(ctx, listReceiptsQuery); err != nil {
		return err
	}
	receipts := listReceiptsQuery.Result.Receipts

	// Get total amount each purchaserefund
	// Map of [ refundID ] amount of receiptLines (current receipts into DB)
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
				case receipt_type.Payment:
					mapRefIDAmountOld[receiptLine.RefID] += receiptLine.Amount
				}
			}
		}
	}

	// Check each amount of receiptLine (param) with (total amount of old receiptLines + total amount of order)
	for key, value := range mapRefIDAmount {
		receiptConfirmValue := mRefund[key].TotalAmount - mapRefIDAmountOld[key]
		if value > receiptConfirmValue {
			return cm.Errorf(cm.InvalidArgument, nil, "Giá trị của đơn trả hàng không hợp lệ, Vui lòng tải lại trang và thử lại")
		}
	}
	return nil
}
