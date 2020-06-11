package pm

import (
	"context"

	"o.o/api/main/shipmentpricing/pricelist"
	"o.o/api/main/shipmentpricing/shopshipmentpricelist"
	"o.o/backend/com/main/shipmentpricing/shipmentprice"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/redis"
)

type ProcessManager struct {
	redisStore      redis.Store
	priceListQS     pricelist.QueryBus
	shopPriceListQS shopshipmentpricelist.QueryBus
}

func New(redisStore redis.Store, eventBus bus.EventRegistry, priceListQS pricelist.QueryBus, shopPriceListQS shopshipmentpricelist.QueryBus) *ProcessManager {
	p := &ProcessManager{
		redisStore:      redisStore,
		priceListQS:     priceListQS,
		shopPriceListQS: shopPriceListQS,
	}
	p.registerEventHandlers(eventBus)
	return p
}

func (m *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.DeleteCachePriceList)
}

func (m *ProcessManager) DeleteCachePriceList(ctx context.Context, event *pricelist.DeleteCachePriceListEvent) error {
	// xóa cache danh sach shipmentprices
	return shipmentprice.DeleteRedisCache(ctx, m.redisStore, event.ShipmentPriceListID)
}
