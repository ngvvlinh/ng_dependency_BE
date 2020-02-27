package pm

import (
	"context"

	"etop.vn/api/main/connectioning"
	"etop.vn/api/main/moneytx"
	"etop.vn/api/main/shipping"
	"etop.vn/backend/com/main/shipping/carrier"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/redis"
	"etop.vn/capi"
)

type ProcessManager struct {
	eventBus   capi.EventBus
	shippingA  shipping.CommandBus
	redisStore redis.Store
}

func New(eventBus capi.EventBus, shippingAggregate shipping.CommandBus, redisS redis.Store) *ProcessManager {
	return &ProcessManager{
		eventBus:   eventBus,
		shippingA:  shippingAggregate,
		redisStore: redisS,
	}
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.MoneyTxShippingExternalCreated)
	eventBus.AddEventListener(m.ConnectionUpdated)
}

func (m *ProcessManager) MoneyTxShippingExternalCreated(ctx context.Context, event *moneytx.MoneyTransactionShippingExternalCreatedEvent) error {
	if event.MoneyTxShippingExternalID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Event MoneyTransactionShippingExternalCreated missing ID")
	}
	if len(event.FulfillementIDs) == 0 {
		return nil
	}
	cmd := &shipping.UpdateFulfillmentsMoneyTxShippingExternalIDCommand{
		FulfillmentIDs:            event.FulfillementIDs,
		MoneyTxShippingExternalID: event.MoneyTxShippingExternalID,
	}
	if err := m.shippingA.Dispatch(ctx, cmd); err != nil {
		return err
	}
	return nil
}

func (m *ProcessManager) ConnectionUpdated(ctx context.Context, event *connectioning.ConnectionUpdatedEvent) error {
	if event.ConnectionID == 0 {
		return nil
	}
	// Delete cache connection in carrier manager
	key := carrier.GetRedisConnectionKeyByID(event.ConnectionID)
	return m.redisStore.Del(key)
}
