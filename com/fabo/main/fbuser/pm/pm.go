package pm

import (
	"context"

	"o.o/api/fabo/fbusering"
	"o.o/api/main/identity"
	"o.o/api/shopping/tradering"
	"o.o/backend/pkg/common/bus"
	"o.o/capi"
)

var (
	defaultTagTemplate = []*fbusering.FbShopUserTag{
		{
			Name:  "Chốt Đơn",
			Color: "#3498db",
		},
		{
			Name:  "Đã Ship",
			Color: "#2ecc71",
		},
		{
			Name:  "Hỏi Giá",
			Color: "#95a5a6",
		},
		{
			Name:  "Tư Vấn",
			Color: "#e74c3c",
		},
		{
			Name:  "Bank",
			Color: "#9b59b6",
		},
		{
			Name:  "COD",
			Color: "#f39c12",
		},
	}
)

type ProcessManager struct {
	eventBus   capi.EventBus
	fbUserAggr fbusering.CommandBus
}

func New(
	eventBusArgs bus.EventRegistry,
	fbUserA fbusering.CommandBus,
) *ProcessManager {
	p := &ProcessManager{
		eventBus:   eventBusArgs,
		fbUserAggr: fbUserA,
	}
	p.RegisterEventHandlers(eventBusArgs)
	return p
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.ShopCustomerDeletedEvent)
	eventBus.AddEventListener(m.AccountCreated)
}

func (m *ProcessManager) ShopCustomerDeletedEvent(ctx context.Context, event *tradering.TraderDeletedEvent) error {
	if event.TradingType != tradering.CustomerType {
		return nil
	}
	cmd := &fbusering.DeleteFbExternalUserShopCustomerCommand{
		ShopID:     event.ShopID,
		CustomerID: event.TraderID.Wrap(),
	}
	err := m.fbUserAggr.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (m *ProcessManager) AccountCreated(ctx context.Context, event *identity.AccountCreatedEvent) error {
	// creates default tag for shop
	shopID := event.ShopID
	for _, _tag := range defaultTagTemplate {
		cmd := &fbusering.CreateShopUserTagCommand{
			Name:   _tag.Name,
			Color:  _tag.Color,
			ShopID: shopID,
		}
		if err := m.fbUserAggr.Dispatch(ctx, cmd); err != nil {
			return err
		}
	}
	return nil
}
