package pm

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/main/inventory"
	"o.o/api/main/purchaseorder"
	"o.o/api/shopping/tradering"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/capi"
	"o.o/capi/dot"
)

type ProcessManager struct {
	eventBus capi.EventBus
	catalogQ catalog.QueryBus
	catalogA catalog.CommandBus
}

func New(
	eventBus bus.EventRegistry,
	catalogQuery catalog.QueryBus,
	catalogAgg catalog.CommandBus,
) *ProcessManager {
	p := &ProcessManager{
		eventBus: eventBus,
		catalogQ: catalogQuery,
		catalogA: catalogAgg,
	}
	p.registerEventHandlers(eventBus)
	return p
}

func (m *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.InventoryVoucherUpdating)
	eventBus.AddEventListener(m.InventoryVoucherCreating)
	eventBus.AddEventListener(m.DeleteVariantSupplier)
	eventBus.AddEventListener(m.CreateVariantSupplier)
}

func (m *ProcessManager) InventoryVoucherCreating(ctx context.Context, event *inventory.InventoryVoucherCreatingEvent) error {
	var variantIDs []dot.ID
	for _, value := range event.Line {
		variantIDs = append(variantIDs, value.VariantID)
	}
	return m.VerifyVariantIDs(ctx, event.ShopID, variantIDs)
}

func (m *ProcessManager) InventoryVoucherUpdating(ctx context.Context, event *inventory.InventoryVoucherUpdatingEvent) error {
	var variantIDs []dot.ID
	for _, value := range event.Line {
		variantIDs = append(variantIDs, value.VariantID)
	}
	return m.VerifyVariantIDs(ctx, event.ShopID, variantIDs)
}

func (m *ProcessManager) VerifyVariantIDs(ctx context.Context, ShopID dot.ID, variantIDs []dot.ID) error {
	query := catalog.ValidateVariantIDsQuery{
		ShopId:         ShopID,
		ShopVariantIds: variantIDs,
	}
	return m.catalogQ.Dispatch(ctx, &query)
}

func (m *ProcessManager) DeleteVariantSupplier(ctx context.Context, event *tradering.TraderDeletedEvent) error {
	query := &catalog.GetVariantsBySupplierIDQuery{
		SupplierID: event.TraderID,
		ShopID:     event.ShopID,
	}
	if err := m.catalogQ.Dispatch(ctx, query); err != nil {
		if cm.ErrorCode(err) == cm.NotFound {
			return nil
		}
		return err
	}
	if len(query.Result.Variants) == 0 {
		return nil
	}
	cmd := &catalog.DeleteVariantSupplierCommand{
		VariantID:  0,
		SupplierID: event.TraderID,
		ShopID:     event.ShopID,
	}
	return m.catalogA.Dispatch(ctx, cmd)
}

func (m *ProcessManager) CreateVariantSupplier(ctx context.Context, event *purchaseorder.PurchaseOrderConfirmedEvent) error {
	mapVariantSupplier := make(map[dot.ID]dot.ID)
	query := catalog.GetVariantsBySupplierIDQuery{
		SupplierID: event.TraderID,
		ShopID:     event.ShopID,
	}
	if err := m.catalogQ.Dispatch(ctx, &query); err != nil {
		return err
	}
	for _, value := range query.Result.Variants {
		mapVariantSupplier[value.VariantID] = event.TraderID
	}
	var variantIDs []dot.ID
	for _, variant := range event.Lines {
		if mapVariantSupplier[variant.VariantID] == 0 {
			variantIDs = append(variantIDs, variant.VariantID)
		}
	}
	if len(variantIDs) == 0 {
		return nil
	}
	cmd := catalog.CreateVariantsSupplierCommand{
		ShopID:     event.ShopID,
		SupplierID: event.TraderID,
		VariantIDs: variantIDs,
	}
	return m.catalogA.Dispatch(ctx, &cmd)
}
