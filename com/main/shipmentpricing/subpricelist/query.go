package subpricelist

import (
	"context"

	"o.o/api/main/shipmentpricing/subpricelist"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/shipmentpricing/subpricelist/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/capi/dot"
)

var _ subpricelist.QueryService = &QueryService{}

type QueryService struct {
	shipmentSubPriceListStore sqlstore.ShipmentSubPriceListStoreFactory
}

func NewQueryService(db com.MainDB) *QueryService {
	return &QueryService{
		shipmentSubPriceListStore: sqlstore.NewShipmentSubPriceListStore(db),
	}
}

func QueryServiceMessageBus(q *QueryService) subpricelist.QueryBus {
	b := bus.New()
	return subpricelist.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *QueryService) ListShipmentSubPriceList(ctx context.Context, args *subpricelist.ListSubPriceListArgs) ([]*subpricelist.ShipmentSubPriceList, error) {
	query := q.shipmentSubPriceListStore(ctx).OptionalConnectionID(args.ConnectionID)
	if args.Status.Valid {
		query = query.Status(args.Status.Enum)
	}
	return query.ListShipmentSubPriceLists()
}

func (q *QueryService) ListShipmentSubPriceListByIDs(ctx context.Context, ids []dot.ID) ([]*subpricelist.ShipmentSubPriceList, error) {
	return q.shipmentSubPriceListStore(ctx).IDs(ids...).ListShipmentSubPriceLists()
}

func (q *QueryService) GetShipmentSubPriceList(ctx context.Context, id dot.ID) (*subpricelist.ShipmentSubPriceList, error) {
	return q.shipmentSubPriceListStore(ctx).ID(id).GetShipmentSubPriceList()
}
