package pm

import (
	"context"

	"etop.vn/api/main/ordering"
	"etop.vn/api/main/shipnow"
)

type ProcessManager struct {
	order   ordering.Aggregate
	shipnow shipnow.Aggregate
}

func New(
	orderAggr ordering.Aggregate,
	shipnowAggr shipnow.Aggregate,
) *ProcessManager {
	return &ProcessManager{
		order:   orderAggr,
		shipnow: shipnowAggr,
	}
}

func (pm *ProcessManager) ValidateOrders(ctx context.Context, args *ordering.ValidateOrdersForShippingCommand) error {
	return pm.order.ValidateOrders(ctx, args)
}

func (pm *ProcessManager) GetOrders(ctx context.Context, args *ordering.GetOrdersArgs) ([]*ordering.Order, error) {
	return pm.order.GetOrders(ctx, args)
}
