package pm

import (
	"context"

	"o.o/api/main/shipmentpricing/pricelist"
	"o.o/api/main/shipmentpricing/shopshipmentpricelist"
	"o.o/api/main/shipmentpricing/subpricelist"
	"o.o/backend/com/main/shipmentpricing/shipmentprice"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/redis"
	"o.o/capi/dot"
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
	eventBus.AddEventListener(m.ShipmentSubPriceListUpdated)
	eventBus.AddEventListener(m.ShipmentSubPriceListDeleting)
	eventBus.AddEventListener(m.ShipmentSubPriceListDeleted)
	eventBus.AddEventListener(m.DeleteCachePriceList)
}

func (m *ProcessManager) DeleteCachePriceList(ctx context.Context, event *pricelist.DeleteCachePriceListEvent) error {
	// xóa cache danh sach shipmentprices
	return shipmentprice.DeleteRedisCache(ctx, m.redisStore, event.ShipmentPriceListID)
}

func (m *ProcessManager) ShipmentSubPriceListUpdated(ctx context.Context, event *subpricelist.ShipmentSubPriceListUpdatedEvent) error {
	// xóa cache danh sach shipmentprices
	query := &pricelist.ListShipmentPriceListsQuery{
		SubShipmentPriceListIDs: []dot.ID{event.ID},
	}
	if err := m.priceListQS.Dispatch(ctx, query); err != nil {
		return err
	}
	for _, priceList := range query.Result {
		if err := shipmentprice.DeleteRedisCache(ctx, m.redisStore, priceList.ID); err != nil {
			return err
		}
	}
	return nil
}

func (m *ProcessManager) ShipmentSubPriceListDeleting(ctx context.Context, event *subpricelist.ShipmentSubPriceListDeletingEvent) error {
	queryPriceList := &pricelist.ListShipmentPriceListsQuery{
		SubShipmentPriceListIDs: []dot.ID{event.ID},
	}
	if err := m.priceListQS.Dispatch(ctx, queryPriceList); err != nil {
		return err
	}
	priceListIDs := make([]dot.ID, len(queryPriceList.Result))
	for i, pl := range queryPriceList.Result {
		priceListIDs[i] = pl.ID
	}
	queryShopPriceList := &shopshipmentpricelist.ListShopShipmentPriceListsByPriceListIDsQuery{
		PriceListIDs: priceListIDs,
	}
	if err := m.shopPriceListQS.Dispatch(ctx, queryShopPriceList); err != nil {
		return err
	}
	if len(queryShopPriceList.Result) != 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "Sub price list is used. Can not delete id.")
	}
	return nil
}

func (m *ProcessManager) ShipmentSubPriceListDeleted(ctx context.Context, event *subpricelist.ShipmentSubPriceListDeletedEvent) error {
	// xóa cache danh sach shipmentprices
	query := &pricelist.ListShipmentPriceListsQuery{
		SubShipmentPriceListIDs: []dot.ID{event.ID},
	}
	if err := m.priceListQS.Dispatch(ctx, query); err != nil {
		return err
	}
	for _, priceList := range query.Result {
		if err := shipmentprice.DeleteRedisCache(ctx, m.redisStore, priceList.ID); err != nil {
			return err
		}
	}
	return nil
}
