package ordering

import (
	"context"

	etoptypes "etop.vn/api/main/etop"
	"etop.vn/api/main/ordering"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/services/ordering/convert"
	"etop.vn/backend/pkg/services/ordering/pm"
	"etop.vn/backend/pkg/services/ordering/sqlstore"
)

var _ ordering.Aggregate = &Aggregate{}

type Aggregate struct {
	s  *sqlstore.OrderStore
	pm *pm.ProcessManager
}

func NewAggregate(db cmsql.Database) *Aggregate {
	return &Aggregate{
		s: sqlstore.NewOrderStore(db),
	}
}

func (a *Aggregate) WithPM(pm *pm.ProcessManager) *Aggregate {
	a.pm = pm
	return a
}

func (a *Aggregate) GetOrderByID(ctx context.Context, args ordering.GetOrderByIDArgs) (*ordering.Order, error) {
	ord, err := a.s.WithContext(ctx).ID(args.ID).Get()
	if err != nil {
		return nil, err
	}
	return convert.Order(ord), nil
}

func (a *Aggregate) GetOrders(ctx context.Context, args *ordering.GetOrdersArgs) ([]*ordering.Order, error) {
	orders, err := a.s.GetOrdes(args)
	if err != nil {
		return nil, err
	}
	return convert.Orders(orders), nil
}

func (a *Aggregate) ValidateOrders(ctx context.Context, args *ordering.ValidateOrdersForShippingCommand) error {
	args1 := &ordering.GetOrdersArgs{
		IDs: args.OrderIDs,
	}
	orders, err := a.GetOrders(ctx, args1)
	if err != nil {
		return err
	}
	for _, order := range orders {
		if err := ValidateOrder(order); err != nil {
			return err
		}
	}
	return nil
}

func ValidateOrder(order *ordering.Order) error {
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
