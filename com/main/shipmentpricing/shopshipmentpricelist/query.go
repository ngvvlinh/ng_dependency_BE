package shopshipmentpricelist

import (
	"context"

	"o.o/api/main/shipmentpricing/shopshipmentpricelist"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/shipmentpricing/shopshipmentpricelist/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/redis"
	"o.o/capi/dot"
)

var _ shopshipmentpricelist.QueryService = &QueryService{}

type QueryService struct {
	shopPriceListStore sqlstore.ShopPriceListStoreFactory
	redisStore         redis.Store
}

func NewQueryService(db com.MainDB, redisStore redis.Store) *QueryService {
	return &QueryService{
		shopPriceListStore: sqlstore.NewShopPriceListStore(db),
		redisStore:         redisStore,
	}
}

func QueryServiceMessageBus(q *QueryService) shopshipmentpricelist.QueryBus {
	b := bus.New()
	return shopshipmentpricelist.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *QueryService) ListShopShipmentPriceLists(ctx context.Context, args *shopshipmentpricelist.GetShopShipmentPriceListsArgs) (*shopshipmentpricelist.GetShopShipmentPriceListsResponse, error) {
	query := q.shopPriceListStore(ctx).WithPaging(args.Paging).OptionalShipmentPriceListID(args.ShipmentPriceListID)
	priceLists, err := query.ListShopPriceLists()
	if err != nil {
		return nil, err
	}
	return &shopshipmentpricelist.GetShopShipmentPriceListsResponse{
		ShopShipmentPriceLists: priceLists,
		Paging:                 query.GetPaging(),
	}, nil
}

func (q *QueryService) ListShopShipmentPriceListsByPriceListIDs(ctx context.Context, priceListIDs []dot.ID) ([]*shopshipmentpricelist.ShopShipmentPriceList, error) {
	return q.shopPriceListStore(ctx).ShipmentPriceListIDs(priceListIDs).ListShopPriceLists()
}

func (q *QueryService) GetShopShipmentPriceList(ctx context.Context, shopID dot.ID) (*shopshipmentpricelist.ShopShipmentPriceList, error) {
	return q.shopPriceListStore(ctx).ShopID(shopID).GetShopPriceList()
}
