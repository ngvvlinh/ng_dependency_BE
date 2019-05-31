package pm

import (
	"etop.vn/api/main/ordering"
	"etop.vn/api/main/shipnow"
)

type ProcessManager struct {
	order   ordering.Aggregate
	shipnow shipnow.Aggregate
}

func New(
	orderAggr ordering.Aggregate,
) *ProcessManager {
	return &ProcessManager{
		order: orderAggr,
	}
}
