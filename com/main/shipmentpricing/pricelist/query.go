package pricelist

import (
	"context"

	"o.o/api/main/shipmentpricing/pricelist"
	"o.o/api/meta"
	"o.o/backend/com/main/shipmentpricing/pricelist/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

var _ pricelist.QueryService = &QueryService{}

type QueryService struct {
	shipmentPriceListStore sqlstore.ShipmentPriceListStoreFactory
	redisStore             redis.Store
}

func NewQueryService(db *cmsql.Database, redisStore redis.Store) *QueryService {
	return &QueryService{
		shipmentPriceListStore: sqlstore.NewShipmentPriceListStore(db),
		redisStore:             redisStore,
	}
}

func (q *QueryService) MessageBus() pricelist.QueryBus {
	b := bus.New()
	return pricelist.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *QueryService) GetShipmentPriceList(ctx context.Context, id dot.ID) (*pricelist.ShipmentPriceList, error) {
	return q.shipmentPriceListStore(ctx).ID(id).GetShipmentPriceList()
}

func (q *QueryService) GetActiveShipmentPriceList(ctx context.Context, _ *meta.Empty) (*pricelist.ShipmentPriceList, error) {
	return q.shipmentPriceListStore(ctx).IsActive(true).GetShipmentPriceList()
}

func (q *QueryService) ListShipmentPriceList(ctx context.Context,
	_ *meta.Empty) ([]*pricelist.ShipmentPriceList, error) {
	return q.shipmentPriceListStore(ctx).ListShipmentPriceLists()
}
