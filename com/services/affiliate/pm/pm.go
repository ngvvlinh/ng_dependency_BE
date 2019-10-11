package pm

import (
	"context"

	"etop.vn/api/main/ordering"

	"etop.vn/api/services/affiliate"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/common/l"
)

type ProcessManager struct {
	affiliate affiliate.CommandBus
}

var (
	ll = l.New()
)

func New(affiliateAggr affiliate.CommandBus) *ProcessManager {
	return &ProcessManager{affiliate: affiliateAggr}
}

func (p *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(p.OrderPaymentSuccess)
}

func (p *ProcessManager) OrderPaymentSuccess(ctx context.Context, event *ordering.OrderPaymentSuccessEvent) error {
	ll.Info("OrderPaymentSuccess CALLED", l.Object("event", event))

	orderPaymentSuccessCmd := &affiliate.OrderPaymentSuccessCommand{
		OrderID: event.OrderID,
	}
	if err := p.affiliate.Dispatch(ctx, orderPaymentSuccessCmd); err != nil {
		return err
	}
	return nil
}
