package pricelist

import (
	"context"

	"o.o/api/main/shipmentpricing/pricelist"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/shipmentpricing/pricelist/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/redis"
	"o.o/capi/dot"
)

var _ pricelist.QueryService = &QueryService{}

type QueryService struct {
	shipmentPriceListStore sqlstore.ShipmentPriceListStoreFactory
	redisStore             redis.Store
}

func NewQueryService(db com.MainDB, redisStore redis.Store) *QueryService {
	return &QueryService{
		shipmentPriceListStore: sqlstore.NewShipmentPriceListStore(db),
		redisStore:             redisStore,
	}
}

func QueryServiceMessageBus(q *QueryService) pricelist.QueryBus {
	b := bus.New()
	return pricelist.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *QueryService) GetShipmentPriceList(ctx context.Context, id dot.ID) (*pricelist.ShipmentPriceList, error) {
	return q.shipmentPriceListStore(ctx).ID(id).GetShipmentPriceList()
}

func (q *QueryService) GetActiveShipmentPriceList(ctx context.Context, connectionID dot.ID) (*pricelist.ShipmentPriceList, error) {
	return q.shipmentPriceListStore(ctx).ConnectionID(connectionID).IsActive(true).GetShipmentPriceList()
}

func (q *QueryService) ListShipmentPriceLists(ctx context.Context,
	args *pricelist.ListShipmentPriceListsArgs) ([]*pricelist.ShipmentPriceList, error) {
	query := q.shipmentPriceListStore(ctx).OptionalConnectionID(args.ConnectionID)
	if args.IsActive.Valid {
		query = query.IsActive(args.IsActive.Bool)
	}
	return query.ListShipmentPriceLists()
}
