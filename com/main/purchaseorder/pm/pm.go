package pm

import (
	"context"

	"o.o/api/main/purchaseorder"
	"o.o/api/main/receipting"
	"o.o/api/top/types/etc/receipt_ref"
	"o.o/api/top/types/etc/receipt_type"
	"o.o/api/top/types/etc/status3"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/capi/dot"
)

type ProcessManager struct {
	purchaseOrderQuery purchaseorder.QueryBus
	receiptQuery       receipting.QueryBus
}

func New(
	purchaseOrderQ purchaseorder.QueryBus,
	receiptQ receipting.QueryBus,
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

	if receipt.RefType != receipt_ref.PurchaseOrder {
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
		if purchaseOrder.Status == status3.N {
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
		RefType: receipt_ref.PurchaseOrder,
		Status:  int(status3.P),
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
				case receipt_type.Payment:
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
