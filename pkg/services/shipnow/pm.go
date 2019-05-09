package shipnow

import (
	"context"

	"etop.vn/api/main/order"
	"etop.vn/api/main/shipnow"
)

type ProcessManager struct {
	order order.Aggregate
}

func NewProcessManager(order order.Aggregate) ProcessManager {
	return ProcessManager{
		order: order,
	}
}

func (pm *ProcessManager) HandleCreation(ctx context.Context, ffm *shipnow.ShipnowFulfillment) error {
	// TODO: validation with order

	return nil
}
