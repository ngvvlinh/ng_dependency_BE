package ordering

import (
	"context"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/ordering"
	"etop.vn/backend/com/main/ordering/convert"
	"etop.vn/backend/com/main/ordering/sqlstore"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
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

func (q *QueryService) MessageBus() ordering.QueryBus {
	b := bus.New()
	return ordering.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *QueryService) GetOrderByID(ctx context.Context, args *ordering.GetOrderByIDArgs) (*ordering.Order, error) {
	return q.store(ctx).ID(args.ID).GetOrder()
}

func (q *QueryService) GetOrders(ctx context.Context, args *ordering.GetOrdersArgs) (*ordering.OrdersResponse, error) {
	orders, err := q.store(ctx).GetOrders(args)
	if err != nil {
		return nil, err
	}
	return &ordering.OrdersResponse{Orders: orders}, nil
}

func (q *QueryService) GetOrdersByIDsAndCustomerID(
	ctx context.Context, shopID int64, IDs []int64, customerID int64,
) (*ordering.OrdersResponse, error) {
	statuses := []etop.Status5{etop.S5Zero, etop.S5Positive, etop.S5SuperPos}
	orders, err := q.store(ctx).IDs(IDs...).CustomerID(customerID).Statuses(statuses).ListOrders()
	if err != nil {
		return nil, err
	}
	return &ordering.OrdersResponse{Orders: convert.Orders(orders)}, nil
}

func (q *QueryService) GetOrderByCode(ctx context.Context, args *ordering.GetOrderByCodeArgs) (*ordering.Order, error) {
	return q.store(ctx).Code(args.Code).GetOrder()
}

func (q *QueryService) ListOrdersByCustomerID(ctx context.Context, shopID, customerID int64) (*ordering.OrdersResponse, error) {
	statuses := []etop.Status5{etop.S5Zero, etop.S5Positive, etop.S5SuperPos}
	orders, err := q.store(ctx).CustomerID(customerID).Statuses(statuses).ListOrders()
	if err != nil {
		return nil, err
	}
	return &ordering.OrdersResponse{Orders: convert.Orders(orders)}, nil
}

func (q *QueryService) ListOrdersByCustomerIDs(ctx context.Context, shopID int64, customerIDs []int64) (*ordering.OrdersResponse, error) {
	statuses := []etop.Status5{etop.S5Zero, etop.S5Positive, etop.S5SuperPos}
	orders, err := q.store(ctx).CustomerIDs(customerIDs...).Statuses(statuses).ListOrders()
	if err != nil {
		return nil, err
	}
	return &ordering.OrdersResponse{Orders: convert.Orders(orders)}, nil
}
