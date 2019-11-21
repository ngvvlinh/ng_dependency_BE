package pm

import (
	"context"

	"etop.vn/api/main/purchaseorder"

	"etop.vn/api/shopping/suppliering"

	"etop.vn/api/main/catalog"
	"etop.vn/api/main/inventory"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/capi"
)

type ProcessManager struct {
	eventBus capi.EventBus
	catalogQ catalog.QueryBus
	catalogA catalog.CommandBus
}

func New(
	eventBusArgs capi.EventBus,
	catalogQuery catalog.QueryBus,
	catalogAgg catalog.CommandBus,
) *ProcessManager {
	return &ProcessManager{
		eventBus: eventBusArgs,
		catalogQ: catalogQuery,
		catalogA: catalogAgg,
	}
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.InventoryVoucherUpdating)
	eventBus.AddEventListener(m.InventoryVoucherCreating)
	eventBus.AddEventListener(m.DeleteVariantSupplier)
	eventBus.AddEventListener(m.CreateVariantSupplier)
}

func (m *ProcessManager) InventoryVoucherCreating(ctx context.Context, event *inventory.InventoryVoucherCreatingEvent) error {
	var variantIDs []int64
	for _, value := range event.Line {
		variantIDs = append(variantIDs, value.VariantID)
	}
	return m.VerifyVariantIDs(ctx, event.ShopID, variantIDs)
}

func (m *ProcessManager) InventoryVoucherUpdating(ctx context.Context, event *inventory.InventoryVoucherUpdatingEvent) error {
	var variantIDs []int64
	for _, value := range event.Line {
		variantIDs = append(variantIDs, value.VariantID)
	}
	return m.VerifyVariantIDs(ctx, event.ShopID, variantIDs)
}

func (m *ProcessManager) VerifyVariantIDs(ctx context.Context, ShopID int64, variantIDs []int64) error {
	query := catalog.ValidateVariantIDsQuery{
		ShopId:         ShopID,
		ShopVariantIds: variantIDs,
	}
	return m.catalogQ.Dispatch(ctx, &query)
}

func (m *ProcessManager) DeleteVariantSupplier(ctx context.Context, event *suppliering.VariantSupplierDeletedEvent) error {
	cmd := catalog.DeleteVariantSupplierCommand{
		VariantID:  0,
		SupplierID: event.SupplierID,
		ShopID:     event.ShopID,
	}
	return m.catalogA.Dispatch(ctx, &cmd)
}

func (m *ProcessManager) CreateVariantSupplier(ctx context.Context, event *purchaseorder.PurchaseOrderConfirmedEvent) error {
	var variantIDs []int64
	for _, variant := range event.Lines {
		variantIDs = append(variantIDs, variant.VariantID)
	}
	cmd := catalog.CreateVariantsSupplierCommand{
		ShopID:     event.ShopID,
		SupplierID: event.TraderID,
		VariantIDs: variantIDs,
	}
	return m.catalogA.Dispatch(ctx, &cmd)
}
