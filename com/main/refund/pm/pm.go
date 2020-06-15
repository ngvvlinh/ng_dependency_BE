package pm

import (
	"context"
	"fmt"

	"o.o/api/main/ordering"
	"o.o/api/main/receipting"
	"o.o/api/main/refund"
	"o.o/api/top/types/etc/receipt_ref"
	"o.o/api/top/types/etc/receipt_type"
	"o.o/api/top/types/etc/status3"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/capi/dot"
)

type ProcessManager struct {
	refundQuery  refund.QueryBus
	receiptQuery receipting.QueryBus
	refundAggr   refund.CommandBus
}

func New(
	eventBus bus.EventRegistry,
	refundQ refund.QueryBus,
	receiptQ receipting.QueryBus,
	refundA refund.CommandBus,
) *ProcessManager {
	p := &ProcessManager{
		refundQuery:  refundQ,
		receiptQuery: receiptQ,
		refundAggr:   refundA,
	}
	p.registerEventHandlers(eventBus)
	return p
}

func (p *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(p.ReceiptCreating)
	eventBus.AddEventListener(p.OrderCancelledEvent)
}

func (p *ProcessManager) OrderCancelledEvent(ctx context.Context, event *ordering.OrderCancelledEvent) error {
	query := &refund.GetRefundsByOrderIDQuery{
		OrderID: event.OrderID,
		ShopID:  event.ShopID,
	}
	err := p.refundQuery.Dispatch(ctx, query)
	if err != nil {
		return err
	}
	if len(query.Result) == 0 {
		return nil
	}
	for _, value := range query.Result {
		cmd := &refund.CancelRefundCommand{
			ShopID:               event.ShopID,
			ID:                   value.ID,
			UpdatedBy:            event.UpdatedBy,
			CancelReason:         fmt.Sprintf("Hủy đơn hàng %v", event.OrderCode),
			AutoInventoryVoucher: event.AutoInventoryVoucher,
		}
		err = p.refundAggr.Dispatch(ctx, cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *ProcessManager) ReceiptCreating(
	ctx context.Context, event *receipting.ReceiptCreatingEvent,
) error {
	var refunds []*refund.Refund
	mRefund := make(map[dot.ID]*refund.Refund)
	refIDs := event.RefIDs
	receipt := event.Receipt
	mapRefIDAmount := event.MapRefIDAmount

	if receipt.RefType != receipt_ref.Refund {
		return nil
	}

	// List Refund on refIDs
	query := &refund.GetRefundsByIDsQuery{
		IDs:    refIDs,
		ShopID: receipt.ShopID,
	}
	if err := p.refundQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	refunds = query.Result
	for _, v := range refunds {
		if v.Status == status3.N {
			continue
		}
		mRefund[v.ID] = v
		if v.CustomerID != receipt.TraderID {
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn trả hàng %v không thuộc đối tác đã chọn", v.Code)
		}
	}

	// Check refIDs and orderIDs (result of query above)
	if len(refIDs) != len(refunds) {
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

	// Get total amount each refund
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
