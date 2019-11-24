package pm

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/api/main/inventory"
	"etop.vn/api/main/ordering"
	"etop.vn/api/main/purchaseorder"
	stocktake "etop.vn/api/main/stocktaking"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/capi"
	"etop.vn/capi/dot"
)

type ProcessManager struct {
	eventBus     capi.EventBus
	catalogQ     catalog.QueryBus
	inventoryAgg inventory.CommandBus
	orderQuery   ordering.QueryBus
}

func New(
	eventBusArgs capi.EventBus,
	shopVariantQueryArgs catalog.QueryBus,
	orderQuery ordering.QueryBus,
	inventoryAggregate inventory.CommandBus,
) *ProcessManager {
	return &ProcessManager{
		eventBus:     eventBusArgs,
		catalogQ:     shopVariantQueryArgs,
		orderQuery:   orderQuery,
		inventoryAgg: inventoryAggregate,
	}
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.PurchaseOrderConfirmed)
	eventBus.AddEventListener(m.OrderConfirmingEvent)
	eventBus.AddEventListener(m.OrderConfirmedEvent)
	eventBus.AddEventListener(m.StocktakeConfirmed)
	eventBus.AddEventListener(m.OrderCancelledEvent)
}

func (m *ProcessManager) PurchaseOrderConfirmed(ctx context.Context, event *purchaseorder.PurchaseOrderConfirmedEvent) error {
	if !event.AutoInventoryVoucher.ValidateAutoInventoryVoucher() {
		return nil
	}
	var isCreate, isConfirm bool
	if event.AutoInventoryVoucher == inventory.AutoCreateInventory {
		isCreate = true
	}
	if event.AutoInventoryVoucher == inventory.AutoCreateAndConfirmInventory {
		isCreate = true
		isConfirm = true
	}

	var inventoryVoucherID dot.ID
	if isCreate {
		cmd := &inventory.CreateInventoryVoucherByReferenceCommand{
			RefType:   inventory.RefTypePurchaseOrder,
			RefID:     event.PurchaseOrderID,
			ShopID:    event.ShopID,
			UserID:    0,
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

func (p *ProcessManager) OrderConfirmedEvent(ctx context.Context, event *ordering.OrderConfirmedEvent) error {
	if !event.AutoInventoryVoucher.ValidateAutoInventoryVoucher() {
		return nil
	}
	var isCreate, isConfirm bool
	if event.AutoInventoryVoucher == inventory.AutoCreateInventory {
		isCreate = true
	}
	if event.AutoInventoryVoucher == inventory.AutoCreateAndConfirmInventory {
		isCreate = true
		isConfirm = true
	}
	var err error
	var inventoryVoucherID dot.ID
	if isCreate {
		cmdCreate := &inventory.CreateInventoryVoucherByReferenceCommand{
			RefType:   inventory.RefTypeOrder,
			RefID:     event.OrderID,
			ShopID:    event.ShopID,
			UserID:    event.UpdatedBy,
			Type:      inventory.InventoryVoucherTypeOut,
			OverStock: event.InventoryOverStock,
		}
		err = p.inventoryAgg.Dispatch(ctx, cmdCreate)
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
		err = p.inventoryAgg.Dispatch(ctx, cmdConfirm)
	}
	return err
}

// OrderConfirmingEvent
// Create InventoryVariant if not exist
// Validate quantity in case of InventoryVoucherTypeOut
func (p *ProcessManager) OrderConfirmingEvent(ctx context.Context, event *ordering.OrderConfirmingEvent) error {
	if !event.AutoInventoryVoucher.ValidateAutoInventoryVoucher() {
		return nil
	}
	inventoryVoucherLines := []*inventory.InventoryVoucherItem{}
	var variantIDs []dot.ID
	for _, line := range event.Lines {
		inventoryVoucherLines = append(inventoryVoucherLines, &inventory.InventoryVoucherItem{
			VariantID: line.VariantId,
			Quantity:  line.Quantity,
		})
		variantIDs = append(variantIDs, line.VariantId)
	}
	query := catalog.ValidateVariantIDsQuery{
		ShopId:         event.ShopID,
		ShopVariantIds: variantIDs,
	}
	err := p.catalogQ.Dispatch(ctx, &query)
	if err != nil {
		return err
	}
	cmd := &inventory.CheckInventoryVariantsQuantityCommand{
		InventoryOverStock: event.InventoryOverStock,
		ShopID:             event.ShopID,
		Type:               inventory.InventoryVoucherTypeOut,
		Lines:              inventoryVoucherLines,
	}
	err = p.inventoryAgg.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProcessManager) StocktakeConfirmed(ctx context.Context, event *stocktake.StocktakeConfirmedEvent) error {
	if !event.AutoInventoryVoucher.ValidateAutoInventoryVoucher() {
		return nil
	}
	var isCreate, isConfirm bool
	if event.AutoInventoryVoucher == inventory.AutoCreateInventory {
		isCreate = true
	}
	if event.AutoInventoryVoucher == inventory.AutoCreateAndConfirmInventory {
		isCreate = true
		isConfirm = true
	}

	var inventoryVoucherIDs []dot.ID
	if isCreate {
		cmdCreate := &inventory.CreateInventoryVoucherByReferenceCommand{
			RefType:   inventory.RefTypeStockTake,
			RefID:     event.StocktakeID,
			ShopID:    event.ShopID,
			UserID:    event.ConfirmedBy,
			Type:      inventory.InventoryVoucherTypeOut,
			OverStock: event.Overstock,
		}
		err := p.inventoryAgg.Dispatch(ctx, cmdCreate)
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
			err := p.inventoryAgg.Dispatch(ctx, cmdConfirm)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *ProcessManager) OrderCancelledEvent(ctx context.Context, event *ordering.OrderCancelledEvent) error {
	if !event.AutoInventoryVoucher.ValidateAutoInventoryVoucher() {
		return nil
	}
	var isCreate, isConfirm bool
	if event.AutoInventoryVoucher == inventory.AutoCreateInventory {
		isCreate = true
	}
	if event.AutoInventoryVoucher == inventory.AutoCreateAndConfirmInventory {
		isCreate = true
		isConfirm = true
	}
	var err error
	var inventoryVoucherID dot.ID
	if isCreate {
		cmdCreate := &inventory.CreateInventoryVoucherByReferenceCommand{
			RefType:   inventory.RefTypeOrder,
			RefID:     event.OrderID,
			ShopID:    event.ShopID,
			UserID:    event.UpdatedBy,
			Type:      inventory.InventoryVoucherTypeIn,
			OverStock: false,
		}
		err = p.inventoryAgg.Dispatch(ctx, cmdCreate)
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
		err = p.inventoryAgg.Dispatch(ctx, cmdConfirm)
	}
	return err
}
