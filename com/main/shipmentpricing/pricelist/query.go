package pricelist

import (
	"context"

	"etop.vn/api/main/shipmentpricing/pricelist"
	"etop.vn/api/meta"
	"etop.vn/backend/com/main/shipmentpricing/pricelist/sqlstore"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/redis"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi/dot"
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
