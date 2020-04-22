package pm

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/main/inventory"
	"o.o/api/main/ordering"
	"o.o/api/main/purchaseorder"
	"o.o/api/main/purchaserefund"
	"o.o/api/main/refund"
	stocktake "o.o/api/main/stocktaking"
	"o.o/api/top/types/etc/inventory_auto"
	"o.o/api/top/types/etc/inventory_type"
	"o.o/api/top/types/etc/inventory_voucher_ref"
	"o.o/backend/pkg/common/bus"
	"o.o/capi"
	"o.o/capi/dot"
)

type ProcessManager struct {
	eventBus     capi.EventBus
	catalogQ     catalog.QueryBus
	inventoryAgg inventory.CommandBus
	orderQuery   ordering.QueryBus
}

func New(
	eventBus bus.EventRegistry,
	shopVariantQueryArgs catalog.QueryBus,
	orderQuery ordering.QueryBus,
	inventoryAggregate inventory.CommandBus,
) *ProcessManager {
	p := &ProcessManager{
		eventBus:     eventBus,
		catalogQ:     shopVariantQueryArgs,
		orderQuery:   orderQuery,
		inventoryAgg: inventoryAggregate,
	}
	p.registerEventHandlers(eventBus)
	return p
}

func (m *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.PurchaseOrderConfirmed)
	eventBus.AddEventListener(m.OrderConfirmingEvent)
	eventBus.AddEventListener(m.OrderConfirmedEvent)
	eventBus.AddEventListener(m.StocktakeConfirmed)
	eventBus.AddEventListener(m.OrderCancelledEvent)
	eventBus.AddEventListener(m.RefundConfirmedEvent)
	eventBus.AddEventListener(m.PurchaseRefundConfirmedEvent)
	eventBus.AddEventListener(m.PurchaseOrderCancelledEvent)
	eventBus.AddEventListener(m.PurchaseRefundCancelledEvent)
	eventBus.AddEventListener(m.RefundCancelledEvent)
}

func (m *ProcessManager) OrderCancelledEvent(ctx context.Context, event *ordering.OrderCancelledEvent) error {
	if event.AutoInventoryVoucher == inventory_auto.Unknown {
		return nil
	}
	cmd := &inventory.CancelInventoryByRefIDCommand{
		RefID:                event.OrderID,
		ShopID:               event.ShopID,
		InventoryOverStock:   true,
		AutoInventoryVoucher: event.AutoInventoryVoucher,
		UpdateBy:             event.UpdatedBy,
	}
	err := m.inventoryAgg.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (m *ProcessManager) RefundCancelledEvent(ctx context.Context, event *refund.RefundCancelledEvent) error {
	cmd := &inventory.CancelInventoryByRefIDCommand{
		RefID:                event.RefundID,
		ShopID:               event.ShopID,
		InventoryOverStock:   true,
		AutoInventoryVoucher: event.AutoInventoryVoucher,
		UpdateBy:             event.UpdatedBy,
	}
	err := m.inventoryAgg.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (m *ProcessManager) PurchaseOrderCancelledEvent(ctx context.Context, event *purchaseorder.PurchaseOrderCancelledEvent) error {
	cmd := &inventory.CancelInventoryByRefIDCommand{
		RefType:              inventory_voucher_ref.PurchaseOrder,
		RefID:                event.PurchaseOrderID,
		ShopID:               event.ShopID,
		InventoryOverStock:   true,
		AutoInventoryVoucher: event.AutoInventoryVoucher,
		UpdateBy:             event.UpdatedBy,
	}
	err := m.inventoryAgg.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (m *ProcessManager) PurchaseRefundCancelledEvent(ctx context.Context, event *purchaserefund.PurchaseRefundCancelledEvent) error {
	cmd := &inventory.CancelInventoryByRefIDCommand{
		RefID:                event.PurchaseRefundID,
		ShopID:               event.ShopID,
		InventoryOverStock:   true,
		AutoInventoryVoucher: event.AutoInventoryVoucher,
		UpdateBy:             event.UpdatedBy,
	}
	err := m.inventoryAgg.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (m *ProcessManager) PurchaseOrderConfirmed(ctx context.Context, event *purchaseorder.PurchaseOrderConfirmedEvent) error {
	if event.AutoInventoryVoucher == inventory_auto.Unknown {
		return nil
	}
	var isCreate, isConfirm bool
	if event.AutoInventoryVoucher == inventory_auto.Create {
		isCreate = true
	}
	if event.AutoInventoryVoucher == inventory_auto.Confirm {
		isCreate = true
		isConfirm = true
	}
	var inventoryVoucherID dot.ID
	if isCreate {
		cmd := &inventory.CreateInventoryVoucherByReferenceCommand{
			RefType:   inventory_voucher_ref.PurchaseOrder,
			RefID:     event.PurchaseOrderID,
			ShopID:    event.ShopID,
			UserID:    event.UserID,
			OverStock: false,
		}
		if err := m.inventoryAgg.Dispatch(ctx, cmd); err != nil {
			return err
		}
		inventoryVoucherID = cmd.Result[0].ID
	}
	if isConfirm {
		cmd := &inventory.ConfirmInventoryVoucherCommand{
			ShopID:    event.ShopID,
			ID:        inventoryVoucherID,
			UpdatedBy: event.UserID,
		}
		if err := m.inventoryAgg.Dispatch(ctx, cmd); err != nil {
			return err
		}
	}

	return nil
}

func (m *ProcessManager) OrderConfirmedEvent(ctx context.Context, event *ordering.OrderConfirmedEvent) error {
	if event.AutoInventoryVoucher == inventory_auto.Unknown {
		return nil
	}
	var isCreate, isConfirm bool
	if event.AutoInventoryVoucher == inventory_auto.Create {
		isCreate = true
	}
	if event.AutoInventoryVoucher == inventory_auto.Confirm {
		isCreate = true
		isConfirm = true
	}
	var err error
	var inventoryVoucherID dot.ID
	if isCreate {
		cmdCreate := &inventory.CreateInventoryVoucherByReferenceCommand{
			RefType: inventory_voucher_ref.Order,
			RefID:   event.OrderID,
			ShopID:  event.ShopID,
			UserID:  event.UpdatedBy,
			Type:    inventory_type.Out,

			OverStock: event.InventoryOverStock,
		}
		err = m.inventoryAgg.Dispatch(ctx, cmdCreate)
		if err != nil {
			return err
		}
		if len(cmdCreate.Result) == 0 {
			return nil
		}
		inventoryVoucherID = cmdCreate.Result[0].ID
	}

	if isConfirm {
		cmdConfirm := &inventory.ConfirmInventoryVoucherCommand{
			ShopID:    event.ShopID,
			ID:        inventoryVoucherID,
			UpdatedBy: event.UpdatedBy,
		}
		err = m.inventoryAgg.Dispatch(ctx, cmdConfirm)
	}
	return err
}

// OrderConfirmingEvent
// Create InventoryVariant if not exist
// Validate quantity in case of InventoryVoucherTypeOut
func (m *ProcessManager) OrderConfirmingEvent(ctx context.Context, event *ordering.OrderConfirmingEvent) error {
	if event.AutoInventoryVoucher == inventory_auto.Unknown {
		return nil
	}
	inventoryVoucherLines := []*inventory.InventoryVoucherItem{}
	var variantIDs []dot.ID
	for _, line := range event.Lines {
		inventoryVoucherLines = append(inventoryVoucherLines, &inventory.InventoryVoucherItem{
			VariantID: line.VariantID,
			Quantity:  line.Quantity,
		})
		variantIDs = append(variantIDs, line.VariantID)
	}
	query := catalog.ValidateVariantIDsQuery{
		ShopId:         event.ShopID,
		ShopVariantIds: variantIDs,
	}
	err := m.catalogQ.Dispatch(ctx, &query)
	if err != nil {
		return err
	}
	cmd := &inventory.CheckInventoryVariantsQuantityCommand{
		InventoryOverStock: event.InventoryOverStock,
		ShopID:             event.ShopID,
		Type:               inventory_type.Out,
		Lines:              inventoryVoucherLines,
	}
	err = m.inventoryAgg.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (m *ProcessManager) StocktakeConfirmed(ctx context.Context, event *stocktake.StocktakeConfirmedEvent) error {
	if event.AutoInventoryVoucher == inventory_auto.Unknown {
		return nil
	}
	var isCreate, isConfirm bool
	if event.AutoInventoryVoucher == inventory_auto.Create {
		isCreate = true
	}
	if event.AutoInventoryVoucher == inventory_auto.Confirm {
		isCreate = true
		isConfirm = true
	}

	var inventoryVoucherIDs []dot.ID
	if isCreate {
		cmdCreate := &inventory.CreateInventoryVoucherByReferenceCommand{
			RefType:   inventory_voucher_ref.StockTake,
			RefID:     event.StocktakeID,
			ShopID:    event.ShopID,
			UserID:    event.ConfirmedBy,
			Type:      inventory_type.Out,
			OverStock: event.Overstock,
		}
		err := m.inventoryAgg.Dispatch(ctx, cmdCreate)
		if err != nil {
			return err
		}
		for _, value := range cmdCreate.Result {
			inventoryVoucherIDs = append(inventoryVoucherIDs, value.ID)
		}
	}
	if isConfirm {
		for _, value := range inventoryVoucherIDs {
			cmdConfirm := &inventory.ConfirmInventoryVoucherCommand{
				ShopID:    event.ShopID,
				ID:        value,
				UpdatedBy: event.ConfirmedBy,
			}
			err := m.inventoryAgg.Dispatch(ctx, cmdConfirm)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *ProcessManager) RefundConfirmedEvent(ctx context.Context, event *refund.RefundConfirmedEvent) error {
	if event.AutoInventoryVoucher == inventory_auto.Unknown {
		return nil
	}
	var isCreate, isConfirm bool
	if event.AutoInventoryVoucher == inventory_auto.Create {
		isCreate = true
	}
	if event.AutoInventoryVoucher == inventory_auto.Confirm {
		isCreate = true
		isConfirm = true
	}

	var inventoryVoucherID dot.ID
	if isCreate {
		cmd := &inventory.CreateInventoryVoucherByReferenceCommand{
			RefType:   inventory_voucher_ref.Refund,
			RefID:     event.RefundID,
			ShopID:    event.ShopID,
			UserID:    event.UpdatedBy,
			OverStock: false,
		}
		if err := m.inventoryAgg.Dispatch(ctx, cmd); err != nil {
			return err
		}
		inventoryVoucherID = cmd.Result[0].ID
	}

	if isConfirm {
		cmd := &inventory.ConfirmInventoryVoucherCommand{
			ShopID:    event.ShopID,
			ID:        inventoryVoucherID,
			UpdatedBy: event.UpdatedBy,
		}
		if err := m.inventoryAgg.Dispatch(ctx, cmd); err != nil {
			return err
		}
	}
	return nil
}

func (m *ProcessManager) PurchaseRefundConfirmedEvent(ctx context.Context, event *purchaserefund.ConfirmedPurchaseRefundEvent) error {
	if event.AutoInventoryVoucher == inventory_auto.Unknown {
		return nil
	}
	var isCreate, isConfirm bool
	if event.AutoInventoryVoucher == inventory_auto.Create {
		isCreate = true
	}
	if event.AutoInventoryVoucher == inventory_auto.Confirm {
		isCreate = true
		isConfirm = true
	}

	var inventoryVoucherID dot.ID
	if isCreate {
		cmd := &inventory.CreateInventoryVoucherByReferenceCommand{
			RefType:   inventory_voucher_ref.PurchaseRefund,
			RefID:     event.PurchaseRefundID,
			ShopID:    event.ShopID,
			UserID:    event.UpdatedBy,
			OverStock: event.InventoryOverStock,
		}
		if err := m.inventoryAgg.Dispatch(ctx, cmd); err != nil {
			return err
		}
		inventoryVoucherID = cmd.Result[0].ID
	}

	if isConfirm {
		cmd := &inventory.ConfirmInventoryVoucherCommand{
			ShopID:    event.ShopID,
			ID:        inventoryVoucherID,
			UpdatedBy: event.UpdatedBy,
		}
		if err := m.inventoryAgg.Dispatch(ctx, cmd); err != nil {
			return err
		}
	}
	return nil
}
