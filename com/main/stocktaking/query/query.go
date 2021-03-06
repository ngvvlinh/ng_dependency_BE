package query

import (
	"context"

	st "o.o/api/main/stocktaking"
	"o.o/api/meta"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/stocktaking/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/capi/dot"
)

var _ st.QueryService = &StocktakeQuery{}

type StocktakeQuery struct {
	StocktakeStore sqlstore.ShopStocktakeFactory
	EventBus       bus.Bus
}

func NewQueryStocktake(db com.MainDB) *StocktakeQuery {
	return &StocktakeQuery{
		StocktakeStore: sqlstore.NewStocktakeStore(db),
	}
}

func StocktakeQueryMessageBus(q *StocktakeQuery) st.QueryBus {
	b := bus.New()
	return st.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *StocktakeQuery) GetStocktakeByID(ctx context.Context, id dot.ID, shopID dot.ID) (*st.ShopStocktake, error) {
	if shopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	result, err := q.StocktakeStore(ctx).ShopID(shopID).ID(id).GetShopStocktake()
	if err != nil {
		return nil, err
	}
	return result, err
}

func (q *StocktakeQuery) GetStocktakesByIDs(ctx context.Context, ids []dot.ID, shopID dot.ID) ([]*st.ShopStocktake, error) {
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
	query := q.StocktakeStore(ctx).ShopID(args.ShopID).WithPaging(&args.Page).Filters(args.Filter)
	if args.Type.Valid {
		query = query.Type(args.Type.Enum)
	}
	if !args.CreatedAtFrom.IsZero() && !args.CreatedAtTo.IsZero() {
		query = query.CreatedAtFromAndTo(args.CreatedAtFrom, args.CreatedAtTo)
	}
	result, err := query.ListShopStocktake()
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
	}, err
}
