package pm

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/api/main/inventory"
	"etop.vn/api/main/ordering"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/capi"
)

type ProcessManager struct {
	eventBus           capi.EventBus
	shopVariantQuery   catalog.QueryBus
	inventoryAggregate inventory.CommandBus
	orderQuery         ordering.QueryBus
}

func New(
	eventBusArgs capi.EventBus,
	shopVariantQueryArgs catalog.QueryBus,
	orderQuery ordering.QueryBus,
	inventoryAggregate inventory.CommandBus,
) *ProcessManager {
	return &ProcessManager{
		eventBus:           eventBusArgs,
		shopVariantQuery:   shopVariantQueryArgs,
		orderQuery:         orderQuery,
		inventoryAggregate: inventoryAggregate,
	}
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.ValidateVariant)
}

func (m *ProcessManager) ValidateVariant(ctx context.Context, event *inventory.InventoryVoucherCreatedEvent) error {
	var variantId []int64
	for _, value := range event.Line {
		variantId = append(variantId, value.VariantID)
	}
	query := catalog.ValidateVariantIDsQuery{
		ShopId:         event.ShopID,
		ShopVariantIds: variantId,
	}
	return m.shopVariantQuery.Dispatch(ctx, &query)
}

// TODO: handle OrderCreatedEvent later
func (m *ProcessManager) ListenOrderCreatedEvent(ctx context.Context, event *ordering.OrderCreatedEvent) error {
	query := ordering.GetOrderByIDQuery{ID: event.OrderID}
	err := m.orderQuery.Dispatch(ctx, &query)
	if err != nil {
		return err
	}
	result := query.Result.Lines
	for _, value := range result {
		cmd := &inventory.CreateInventoryVariantCommand{
			ShopID:    event.ShopID,
			VariantID: value.VariantId,
		}
		err = m.inventoryAggregate.Dispatch(ctx, cmd)
		if err != nil {
			return err
		}
	}
	return nil
}
