package query

import (
	"context"

	"etop.vn/api/shopping"
	"etop.vn/api/shopping/tradering"
	"etop.vn/backend/com/shopping/tradering/sqlstore"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/common/bus"
)

var _ tradering.QueryService = &TraderQuery{}

type TraderQuery struct {
	store sqlstore.TraderStoreFactory
}

func NewTraderQuery(db cmsql.Database) *TraderQuery {
	return &TraderQuery{
		store: sqlstore.NewTraderStore(db),
	}
}

func (q *TraderQuery) MessageBus() tradering.QueryBus {
	b := bus.New()
	return tradering.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *TraderQuery) GetTraderByID(
	ctx context.Context, args *shopping.IDQueryShopArg,
) (*tradering.ShopTrader, error) {
	return q.store(ctx).ID(args.ID).OptionalShopID(args.ShopID).GetTrader()
}

func (q *TraderQuery) ListTradersByIDs(
	ctx context.Context, args *shopping.IDsQueryShopArgs,
) (*tradering.TradersResponse, error) {
	traders, err := q.store(ctx).ShopID(args.ShopID).IDs(args.IDs...).ListTraders()
	if err != nil {
		return nil, err
	}
	return &tradering.TradersResponse{Traders: traders}, nil
}
