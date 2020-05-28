package pm

import (
	"context"

	"o.o/api/main/ordering"
	"o.o/api/services/affiliate"
	"o.o/backend/pkg/common/bus"
	"o.o/common/l"
)

type ProcessManager struct {
	affiliate affiliate.CommandBus
}

var (
	ll = l.New()
)

func New(eventBus bus.EventRegistry, affiliateAggr affiliate.CommandBus) *ProcessManager {
	p := &ProcessManager{affiliate: affiliateAggr}
	p.registerEventHandlers(eventBus)
	return p
}

func (p *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
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
