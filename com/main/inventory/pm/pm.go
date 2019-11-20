package pm

import (
	"context"

	"etop.vn/api/main/purchaseorder"

	"etop.vn/api/main/catalog"
	"etop.vn/api/main/inventory"
	"etop.vn/api/main/ordering"
	stocktake "etop.vn/api/main/stocktaking"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/capi"
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
	var isCreate, isConfirm bool
	if event.AutoInventoryVoucher == purchaseorder.AutoInventoryVoucherCreate {
		isCreate = true
	}
	if event.AutoInventoryVoucher == purchaseorder.AutoInventoryVoucherConfirm {
		isCreate = true
		isConfirm = true
	}

	var inventoryVoucherID int64
	if isCreate {
		cmd := &inventory.CreateInventoryVoucherCommand{
			Overstock:   false,
			ShopID:      event.ShopID,
			CreatedBy:   event.UserID,
			Title:       "Nhập kho khi kiểm hàng",
			RefID:       event.PurchaseOrderID,
			RefType:     inventory.RefTypePurchaseOrder,
			RefName:     inventory.RefNamePurchaseOrder,
			RefCode:     event.PurchaseOrderCode,
			TraderID:    event.TraderID,
			TotalAmount: int32(event.TotalAmount),
			Type:        inventory.InventoryVoucherTypeIn,
			Lines:       event.Lines,
		}
		if err := m.inventoryAgg.Dispatch(ctx, cmd); err != nil {
			return err
		}
		inventoryVoucherID = cmd.Result.ID
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
	// Create inventory voucher
	inventoryVoucherLines := []*inventory.InventoryVoucherItem{}
	for _, value := range event.Lines {
		inventoryVoucherLines = append(inventoryVoucherLines, &inventory.InventoryVoucherItem{
			VariantID: value.VariantId,
			Quantity:  value.Quantity,
		})
	}

	cmdCreate := &inventory.CreateInventoryVoucherCommand{
		Overstock: event.InventoryOverStock,
		ShopID:    event.ShopID,
		Title:     "Xuất kho khi bán hàng",
		RefCode:   event.OrderCode,
		RefID:     event.OrderID,
		RefType:   inventory.RefTypeOrder,
		TraderID:  event.CustomerID,
		Type:      inventory.InventoryVoucherTypeOut,
		Note:      "Tạo tự động khi xác nhận đơn hàng",
		Lines:     inventoryVoucherLines,
	}
	err := p.inventoryAgg.Dispatch(ctx, cmdCreate)
	if err != nil {
		return err
	}
	if event.AutoInventoryVoucher == inventory.AutoCreateAndConfirmInventory {
		cmdConfirm := &inventory.ConfirmInventoryVoucherCommand{
			ShopID: event.ShopID,
			ID:     cmdCreate.Result.ID,
			Result: nil,
		}
		err = p.inventoryAgg.Dispatch(ctx, cmdConfirm)
		if err != nil {
			return err
		}
	}
	return nil
}

// OrderConfirmingEvent
// Create InventoryVariant if not exist
// Validate quantity in case of InventoryVoucherTypeOut
func (p *ProcessManager) OrderConfirmingEvent(ctx context.Context, event *ordering.OrderConfirmingEvent) error {
	if !event.AutoInventoryVoucher.ValidateAutoInventoryVoucher() {
		return nil
	}
	inventoryVoucherLines := []*inventory.InventoryVoucherItem{}
	var variantIDs []int64
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
	if event.AutoInventoryVoucher == inventory.AutoCreateInventory || event.AutoInventoryVoucher == inventory.AutoCreateAndConfirmInventory {
		var inventoryVariantChange []*inventory.InventoryVariantQuantityChange
		for _, value := range event.Stocktake.Lines {
			inventoryVariantChange = append(inventoryVariantChange, &inventory.InventoryVariantQuantityChange{
				VariantID:      value.VariantID,
				QuantityChange: value.NewQuantity - value.OldQuantity,
			})
		}
		cmdCreate := &inventory.CreateInventoryVoucherByQuantityChangeCommand{
			ShopID:    event.Stocktake.ShopID,
			RefID:     event.Stocktake.ID,
			RefType:   inventory.RefTypeStockTake,
			RefName:   inventory.RefNameStockTake,
			RefCode:   event.Stocktake.Code,
			Note:      "Tạo tự động khi xác nhận phiếu kiểm kho",
			Title:     "Phiếu kiểm kho",
			Overstock: event.Overstock,
			CreatedBy: event.ConfirmedBy,
			Variants:  inventoryVariantChange,
		}
		err := p.inventoryAgg.Dispatch(ctx, cmdCreate)
		if err != nil {
			return err
		}
		if event.AutoInventoryVoucher == inventory.AutoCreateAndConfirmInventory {
			if cmdCreate.Result.TypeIn.ID != 0 {
				cmdConfirmInVoucher := &inventory.ConfirmInventoryVoucherCommand{
					ShopID:    event.Stocktake.ShopID,
					ID:        cmdCreate.Result.TypeIn.ID,
					UpdatedBy: event.ConfirmedBy,
				}
				err = p.inventoryAgg.Dispatch(ctx, cmdConfirmInVoucher)
				if err != nil {
					return err
				}
			}
			if cmdCreate.Result.TypeOut.ID != 0 {
				cmdConfirmOutVoucher := &inventory.ConfirmInventoryVoucherCommand{
					ShopID:    event.Stocktake.ShopID,
					ID:        cmdCreate.Result.TypeOut.ID,
					UpdatedBy: event.ConfirmedBy,
				}
				err = p.inventoryAgg.Dispatch(ctx, cmdConfirmOutVoucher)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (p *ProcessManager) OrderCancelledEvent(ctx context.Context, event *ordering.OrderCancelledEvent) error {
	if !event.AutoInventoryVoucher.ValidateAutoInventoryVoucher() {
		return nil
	}
	// Create inventory voucher
	inventoryVoucherLines := []*inventory.InventoryVoucherItem{}
	for _, value := range event.Lines {
		inventoryVoucherLines = append(inventoryVoucherLines, &inventory.InventoryVoucherItem{
			VariantID: value.VariantId,
			Quantity:  value.Quantity,
		})
	}

	cmdCreate := &inventory.CreateInventoryVoucherCommand{
		Overstock: false,
		ShopID:    event.ShopID,
		Title:     "Nhập kho khi hủy đơn hàng",
		RefID:     event.OrderID,
		RefType:   inventory.RefTypeOrder,
		TraderID:  event.CustomerID,
		Type:      inventory.InventoryVoucherTypeIn,
		Note:      "Tạo tự động khi hủy đơn hàng",
		Lines:     inventoryVoucherLines,
	}
	err := p.inventoryAgg.Dispatch(ctx, cmdCreate)
	if err != nil {
		return err
	}
	if event.AutoInventoryVoucher == inventory.AutoCreateAndConfirmInventory {
		cmdConfirm := &inventory.ConfirmInventoryVoucherCommand{
			ShopID: event.ShopID,
			ID:     cmdCreate.Result.ID,
			Result: nil,
		}
		err = p.inventoryAgg.Dispatch(ctx, cmdConfirm)
		if err != nil {
			return err
		}
	}
	return nil
}
