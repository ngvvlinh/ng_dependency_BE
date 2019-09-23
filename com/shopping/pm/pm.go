package pm

import (
	"context"

	"etop.vn/capi"

	"etop.vn/api/main/catalog"
	"etop.vn/api/shopping/vendoring"
	"etop.vn/backend/pkg/common/bus"
)

type ProcessManager struct {
	eventBus capi.EventBus

	vendorQuery vendoring.QueryBus
}

func New(
	eventBus capi.EventBus,
	vendorQuery vendoring.QueryBus,
) *ProcessManager {
	return &ProcessManager{
		eventBus:    eventBus,
		vendorQuery: vendorQuery,
	}
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.ShopProductCreating)
	eventBus.AddEventListener(m.ShopProductUpdating)
}

func (m *ProcessManager) ShopProductCreating(ctx context.Context, event *catalog.ShopProductCreatingEvent) error {
	if event.VendorID == 0 {
		return nil
	}
	// Call vendorQuery for getVendorByID
	cmd := &vendoring.GetVendorByIDQuery{
		ID:     event.VendorID,
		ShopID: event.ShopID,
	}
	if err := m.vendorQuery.Dispatch(ctx, cmd); err != nil {
		return err
	}
	return nil
}

func (m *ProcessManager) ShopProductUpdating(ctx context.Context, event *catalog.ShopProductUpdatingEvent) error {
	if event.VendorID == 0 {
		return nil
	}
	// Call vendorQuery for getVendorByID
	cmd := &vendoring.GetVendorByIDQuery{
		ID:     event.VendorID,
		ShopID: event.ShopID,
	}
	if err := m.vendorQuery.Dispatch(ctx, cmd); err != nil {
		return err
	}
	return nil
}
