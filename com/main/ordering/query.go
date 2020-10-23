package ordering

import (
	"context"

	"o.o/api/main/ordering"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status5"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/ordering/convert"
	"o.o/backend/com/main/ordering/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/capi/dot"
)

var _ ordering.QueryService = &QueryService{}

type QueryService struct {
	store sqlstore.OrderStoreFactory
}

func NewQueryService(db com.MainDB) *QueryService {
	return &QueryService{
		store: sqlstore.NewOrderStore(db),
	}
}

func QueryServiceMessageBus(q *QueryService) ordering.QueryBus {
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
	ctx context.Context, shopID dot.ID, IDs []dot.ID, customerID dot.ID,
) (*ordering.OrdersResponse, error) {
	statuses := []status5.Status{status5.Z, status5.P, status5.S}
	orders, err := q.store(ctx).IDs(IDs...).CustomerID(customerID).Statuses(statuses).ListOrders()
	if err != nil {
		return nil, err
	}
	return &ordering.OrdersResponse{Orders: convert.Orders(orders)}, nil
}

func (q *QueryService) GetOrderByCode(ctx context.Context, args *ordering.GetOrderByCodeArgs) (*ordering.Order, error) {
	return q.store(ctx).Code(args.Code).GetOrder()
}

func (q *QueryService) ListOrdersByCustomerID(ctx context.Context, shopID, customerID dot.ID) (*ordering.OrdersResponse, error) {
	statuses := []status5.Status{status5.Z, status5.P, status5.S}
	orders, err := q.store(ctx).CustomerID(customerID).Statuses(statuses).ListOrders()
	if err != nil {
		return nil, err
	}
	return &ordering.OrdersResponse{Orders: convert.Orders(orders)}, nil
}

func (q *QueryService) ListOrdersByCustomerIDs(ctx context.Context, shopID dot.ID, customerIDs []dot.ID) (*ordering.OrdersResponse, error) {
	statuses := []status5.Status{status5.Z, status5.P, status5.S}
	orders, err := q.store(ctx).CustomerIDs(customerIDs...).Statuses(statuses).ListOrders()
	if err != nil {
		return nil, err
	}
	return &ordering.OrdersResponse{Orders: convert.Orders(orders)}, nil
}

func (q *QueryService) ListOrdersConfirmed(
	ctx context.Context, args *ordering.ListOrdersConfirmedArgs,
) ([]*ordering.Order, error) {
	query := q.store(ctx)
	query = query.ConfirmStatus(status3.P)
	query = query.Statuses([]status5.Status{status5.NS, status5.Z, status5.P, status5.S})
	query = query.CreatedAtFromAndTo(args.CreatedAtFrom, args.CreatedAtTo)

	if args.CreatedBy != 0 {
		query = query.CreatedBy(args.CreatedBy)
	}

	orders, err := query.ListOrders()
	if err != nil {
		return nil, err
	}
	return convert.Orders(orders), nil
}
