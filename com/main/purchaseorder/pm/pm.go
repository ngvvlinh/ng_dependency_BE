package pm

import (
	"context"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/purchaseorder"
	"etop.vn/api/main/receipting"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/capi/dot"
)

type ProcessManager struct {
	purchaseOrderQuery *purchaseorder.QueryBus
	receiptQuery       *receipting.QueryBus
}

func New(
	purchaseOrderQ *purchaseorder.QueryBus,
	receiptQ *receipting.QueryBus,
) *ProcessManager {
	return &ProcessManager{
		purchaseOrderQuery: purchaseOrderQ,
		receiptQuery:       receiptQ,
	}
}

func (p *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(p.ReceiptCreating)
}

func (p *ProcessManager) ReceiptCreating(
	ctx context.Context, event *receipting.ReceiptCreatingEvent,
) error {
	var purchaseOrders []*purchaseorder.PurchaseOrder
	mPurchaseOrder := make(map[dot.ID]*purchaseorder.PurchaseOrder)
	refIDs := event.RefIDs
	receipt := event.Receipt
	mapRefIDAmount := event.MapRefIDAmount

	if receipt.RefType != receipting.ReceiptRefTypePurchaseOrder {
		return nil
	}

	// List purchases depend on refIDs
	query := &purchaseorder.GetPurchaseOrdersByIDsQuery{
		IDs:    refIDs,
		ShopID: receipt.ShopID,
	}
	if err := p.purchaseOrderQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	purchaseOrders = query.Result.PurchaseOrders
	for _, purchaseOrder := range purchaseOrders {
		if purchaseOrder.Status == etop.S3Negative {
			continue
		}
		mPurchaseOrder[purchaseOrder.ID] = purchaseOrder
		if purchaseOrder.SupplierID != receipt.TraderID {
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn nhập hàng %v không thuộc đối tác đã chọn", purchaseOrder.Code)
		}
	}

	// Check refIDs and orderIDs (result of query above)
	if len(refIDs) != len(purchaseOrders) {
		for _, refID := range refIDs {
			if _, ok := mPurchaseOrder[refID]; !ok {
				return cm.Errorf(cm.FailedPrecondition, nil, "ref_id %d không tìm thấy", refID)
			}
		}
	}

	// List receipts by refIDs
	listReceiptsQuery := &receipting.ListReceiptsByRefsAndStatusQuery{
		ShopID:  receipt.ShopID,
		RefIDs:  refIDs,
		RefType: receipting.ReceiptRefTypePurchaseOrder,
		Status:  int32(etop.S3Positive),
	}
	if err := p.receiptQuery.Dispatch(ctx, listReceiptsQuery); err != nil {
		return err
	}
	receipts := listReceiptsQuery.Result.Receipts

	// Get total amount each purchaseOrderID
	// Map of [ purchaseOrderID ] amount of receiptLines (current receipts into DB)
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
				case receipting.ReceiptTypePayment:
					mapRefIDAmountOld[receiptLine.RefID] += receiptLine.Amount
				}
			}
		}
	}

	// Check each amount of receiptLine (param) with (total amount of old receiptLines + total amount of order)
	for key, value := range mapRefIDAmount {
		if value > mPurchaseOrder[key].TotalAmount-mapRefIDAmountOld[key] {
			return cm.Errorf(cm.InvalidArgument, nil, "Giá trị của đơn hàng không hợp lệ, Vui lòng tải lại trang và thử lại")
		}
	}
	return nil
}
