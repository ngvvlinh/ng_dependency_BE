package pm

import (
	"context"

	"etop.vn/api/main/shipmentpricing/pricelist"
	"etop.vn/backend/com/main/shipmentpricing/shipmentprice"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/redis"
)

type ProcessManager struct {
	redisStore redis.Store
}

func New(redisStore redis.Store) *ProcessManager {
	return &ProcessManager{
		redisStore: redisStore,
	}
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.ShipmentPriceListActivated)
}

func (m *ProcessManager) ShipmentPriceListActivated(ctx context.Context, event *pricelist.ShipmentPriceListActivatedEvent) error {
	// xóa cache danh sach shipmentprices
	return shipmentprice.DeleteRedisCache(ctx, m.redisStore)
}
