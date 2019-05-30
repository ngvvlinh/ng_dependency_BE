package ordering

import (
	"context"

	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"

	"etop.vn/api/main/ordering"
	"etop.vn/backend/pkg/services/ordering/convert"
	"etop.vn/backend/pkg/services/ordering/sqlstore"
)

var _ ordering.QueryService = &QueryService{}

type QueryService struct {
	store sqlstore.OrderStoreFactory
}

func NewQueryService(db cmsql.Database) *QueryService {
	return &QueryService{
		store: sqlstore.NewOrderStore(db),
	}
}

func (a *QueryService) MessageBus() ordering.QueryBus {
	b := bus.New()
	return ordering.NewQueryServiceHandler(a).RegisterHandlers(b)
}

func (q *QueryService) GetOrderByID(ctx context.Context, args *ordering.GetOrderByIDArgs) (*ordering.Order, error) {
	ord, err := q.store(ctx).ID(args.ID).Get()
	if err != nil {
		return nil, err
	}
	return convert.Order(ord), nil
}

func (q *QueryService) GetOrders(ctx context.Context, args *ordering.GetOrdersArgs) (*ordering.OrdersResponse, error) {
	orders, err := q.store(ctx).GetOrders(args)
	if err != nil {
		return nil, err
	}
	return &ordering.OrdersResponse{Orders: orders}, nil
}
