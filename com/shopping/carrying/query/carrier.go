package query

import (
	"context"

	"etop.vn/api/shopping"
	"etop.vn/api/shopping/carrying"
	"etop.vn/backend/com/shopping/carrying/sqlstore"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
)

var _ carrying.QueryService = &CarrierQuery{}

type CarrierQuery struct {
	store sqlstore.CarrierStoreFactory
}

func NewCarrierQuery(db cmsql.Database) *CarrierQuery {
	return &CarrierQuery{
		store: sqlstore.NewCarrierStore(db),
	}
}

func (q *CarrierQuery) MessageBus() carrying.QueryBus {
	b := bus.New()
	return carrying.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *CarrierQuery) GetCarrierByID(
	ctx context.Context, args *shopping.IDQueryShopArg,
) (*carrying.ShopCarrier, error) {
	return q.store(ctx).ID(args.ID).OptionalShopID(args.ShopID).GetCarrier()
}

func (q *CarrierQuery) ListCarriers(
	ctx context.Context, args *shopping.ListQueryShopArgs,
) (*carrying.CarriersResponse, error) {
	query := q.store(ctx).ShopID(args.ShopID).Paging(args.Paging).Filters(args.Filters)
	carriers, err := query.ListCarriers()
	if err != nil {
		return nil, err
	}
	count, err := query.Count()
	if err != nil {
		return nil, err
	}
	return &carrying.CarriersResponse{
		Carriers: carriers,
		Count:    int32(count),
	}, nil
}

func (q *CarrierQuery) ListCarriersByIDs(
	ctx context.Context, args *shopping.IDsQueryShopArgs,
) (*carrying.CarriersResponse, error) {
	carries, err := q.store(ctx).ShopID(args.ShopID).IDs(args.IDs...).ListCarriers()
	if err != nil {
		return nil, err
	}
	return &carrying.CarriersResponse{Carriers: carries}, nil
}
