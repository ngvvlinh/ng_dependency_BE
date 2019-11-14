package pm

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/api/main/inventory"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/capi"
)

type ProcessManager struct {
	eventBus capi.EventBus
	catalogQ catalog.QueryBus
}

func New(
	eventBusArgs capi.EventBus,
	catalogQuery catalog.QueryBus,
) *ProcessManager {
	return &ProcessManager{
		eventBus: eventBusArgs,
		catalogQ: catalogQuery,
	}
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.InventoryVoucherUpdating)
	eventBus.AddEventListener(m.InventoryVoucherCreating)
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
