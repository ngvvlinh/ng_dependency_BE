package shipnow

import (
	"context"

	"etop.vn/api/main/ordering"
	"etop.vn/api/main/shipnow"
)

type ProcessManager struct {
	order ordering.Aggregate
}

func NewProcessManager(order ordering.Aggregate) ProcessManager {
	return ProcessManager{
		order: order,
	}
}

func (pm *ProcessManager) HandleCreation(ctx context.Context, ffm *shipnow.ShipnowFulfillment) error {
	// TODO: validation with order

	return nil
}
