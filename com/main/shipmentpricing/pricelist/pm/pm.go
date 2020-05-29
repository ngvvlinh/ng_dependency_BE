package pm

import (
	"context"

	"o.o/api/main/shipmentpricing/pricelist"
	"o.o/backend/com/main/shipmentpricing/shipmentprice"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/redis"
)

type ProcessManager struct {
	redisStore redis.Store
}

func New(redisStore redis.Store, eventBus bus.EventRegistry) *ProcessManager {
	p := &ProcessManager{
		redisStore: redisStore,
	}
	p.registerEventHandlers(eventBus)
	return p
}

func (m *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.ShipmentPriceListActivated)
}

func (m *ProcessManager) ShipmentPriceListActivated(ctx context.Context, event *pricelist.ShipmentPriceListActivatedEvent) error {
	// xóa cache danh sach shipmentprices
	return shipmentprice.DeleteRedisCache(ctx, m.redisStore)
}
