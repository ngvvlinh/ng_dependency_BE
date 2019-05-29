package ordering

import (
	"context"

	"github.com/k0kubun/pp"

	etoptypes "etop.vn/api/main/etop"
	"etop.vn/api/main/ordering"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/services/ordering/convert"
	"etop.vn/backend/pkg/services/ordering/pm"
	"etop.vn/backend/pkg/services/ordering/sqlstore"
)

var _ ordering.Aggregate = &Aggregate{}

type Aggregate struct {
	pm    *pm.ProcessManager
	store sqlstore.OrderStoreFactory
}

func NewAggregate(db cmsql.Database) *Aggregate {
	return &Aggregate{
		store: sqlstore.NewOrderStore(db),
	}
}

func (a *Aggregate) WithPM(pm *pm.ProcessManager) *Aggregate {
	a.pm = pm
	return a
}

func (a *Aggregate) MessageBus() ordering.AggregateBus {
	b := bus.New()
	ordering.NewAggregateHandler(a).RegisterHandlers(b)
	return ordering.AggregateBus{b}
}

func (a *Aggregate) GetOrderByID(ctx context.Context, args *ordering.GetOrderByIDArgs) (*ordering.Order, error) {
	ord, err := a.store(ctx).ID(args.ID).Get()
	if err != nil {
		return nil, err
	}
	return convert.Order(ord), nil
}

func (a *Aggregate) GetOrders(ctx context.Context, args *ordering.GetOrdersArgs) (*ordering.OrdersResponse, error) {
	orders, err := a.store(ctx).GetOrders(args)
	if err != nil {
		return nil, err
	}
	return &ordering.OrdersResponse{Orders: orders}, nil
}

func (a *Aggregate) ValidateOrders(ctx context.Context, args *ordering.ValidateOrdersForShippingArgs) (*ordering.ValidateOrdersResponse, error) {
	args1 := &ordering.GetOrdersArgs{
		IDs: args.OrderIDs,
	}
	orders, err := a.GetOrders(ctx, args1)
	if err != nil {
		return nil, err
	}
	for _, order := range orders.Orders {
		if err := ValidateOrderStatus(order); err != nil {
			return nil, err
		}
	}
	return &ordering.ValidateOrdersResponse{}, nil
}

func ValidateOrderStatus(order *ordering.Order) error {
	switch order.Status {
	case etoptypes.S5Negative:
		return cm.Error(cm.FailedPrecondition, "Đơn hàng đã huỷ.", nil).WithMetaID("orderID", order.ID)
	case etoptypes.S5Positive:
		return cm.Error(cm.FailedPrecondition, "Đơn hàng đã hoàn thành.", nil).WithMetaID("orderID", order.ID)
	case etoptypes.S5NegSuper:
		return cm.Error(cm.FailedPrecondition, "Đơn hàng đã trả hàng.", nil).WithMetaID("orderID", order.ID)
	}

	if order.ConfirmStatus == etoptypes.S3Negative {
		return cm.Error(cm.FailedPrecondition, "Đơn hàng đã huỷ.", nil).WithMetaID("orderID", order.ID)
	}
	return nil
}

func (a *Aggregate) ReserveOrdersForFfm(ctx context.Context, args *ordering.ReserveOrdersForFfmArgs) (*ordering.ReserveOrdersForFfmResponse, error) {
	orderIDs := args.OrderIDs
	orders, err := a.store(ctx).GetOrders(&ordering.GetOrdersArgs{
		IDs: orderIDs,
	})
	if err != nil {
		return nil, err
	}
	for _, order := range orders {
		if err := ValidateOrderForReserveFfm(order); err != nil {
			return nil, err
		}
	}

	update := sqlstore.UpdateOrdersForReserveOrdersArgs{
		OrderIDs:   orderIDs,
		Fulfill:    args.Fulfill,
		FulfillIDs: args.FulfillIDs,
	}
	pp.Println("update :: ", update)
	orders, err = a.store(ctx).UpdateOrdersForReverseOrders(update)
	if err != nil {
		return nil, err
	}
	return &ordering.ReserveOrdersForFfmResponse{
		Orders: orders,
	}, nil
}

func ValidateOrderForReserveFfm(order *ordering.Order) error {
	if err := ValidateOrderStatus(order); err != nil {
		return err
	}
	if len(order.FulfillIDs) != 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "Order has been reserved").WithMetaID("order_id", order.ID)
	}
	return nil
}
