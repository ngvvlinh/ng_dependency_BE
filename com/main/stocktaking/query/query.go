package query

import (
	"context"

	"etop.vn/api/meta"

	cm "etop.vn/backend/pkg/common"

	"etop.vn/backend/pkg/common/cmsql"

	st "etop.vn/api/main/stocktaking"
	"etop.vn/backend/com/main/stocktaking/sqlstore"
	"etop.vn/backend/pkg/common/bus"
)

var _ st.QueryService = &StocktakeQuery{}

type StocktakeQuery struct {
	StocktakeStore sqlstore.ShopStocktakeFactory
	EventBus       bus.Bus
}

func NewQueryStocktake(db *cmsql.Database) *StocktakeQuery {
	return &StocktakeQuery{
		StocktakeStore: sqlstore.NewStocktakeStore(db),
	}
}

func (q *StocktakeQuery) MessageBus() st.QueryBus {
	b := bus.New()
	return st.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *StocktakeQuery) GetStocktakeByID(ctx context.Context, id int64, shopID int64) (*st.ShopStocktake, error) {
	if shopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	result, err := q.StocktakeStore(ctx).ShopID(shopID).ID(id).GetShopStocktake()
	if err != nil {
		return nil, err
	}
	return result, err
}

func (q *StocktakeQuery) GetStocktakesByIDs(ctx context.Context, ids []int64, shopID int64) ([]*st.ShopStocktake, error) {
	if shopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	result, err := q.StocktakeStore(ctx).ShopID(shopID).IDs(ids...).ListShopStocktake()
	return result, err
}

func (q *StocktakeQuery) ListStocktake(ctx context.Context, args *st.ListStocktakeRequest) (*st.ListStocktakeResponse, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	query := q.StocktakeStore(ctx).ShopID(args.ShopID).Paging(&args.Page)
	result, err := query.ListShopStocktake()
	if err != nil {
		return nil, err
	}
	total, err := query.Count()
	if err != nil {
		return nil, err
	}
	return &st.ListStocktakeResponse{
		Stocktakes: result,
		PageInfo: meta.PageInfo{
			Offset: args.Page.Offset,
			Limit:  args.Page.Limit,
			Sort:   args.Page.Sort,
		},
		Total: int32(total),
	}, err
}