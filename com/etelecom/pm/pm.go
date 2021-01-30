package pm

import (
	"context"
	"math"

	"o.o/api/etelecom"
	"o.o/api/etelecom/call_direction"
	"o.o/api/top/types/etc/status5"
	"o.o/backend/com/etelecom/postage"
	"o.o/backend/pkg/common/bus"
)

const SecondsPerMinute = 60

type ProcessManager struct {
	etelecomAggr  etelecom.CommandBus
	etelecomQuery etelecom.QueryBus
}

func New(eventBus bus.EventRegistry,
	etelecomA etelecom.CommandBus,
	etelecomQ etelecom.QueryBus,
) *ProcessManager {
	p := &ProcessManager{
		etelecomAggr:  etelecomA,
		etelecomQuery: etelecomQ,
	}
	p.registerEventHandlers(eventBus)
	return p
}

func (m *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.CallLogCalcPostage)
}

// calc postage
func (m *ProcessManager) CallLogCalcPostage(ctx context.Context, event *etelecom.CallLogCalcPostageEvent) error {
	if event.CallStatus != status5.P {
		return nil
	}

	if event.HotlineID != 0 {
		queryHotline := &etelecom.GetHotlineQuery{
			ID: event.HotlineID,
		}
		if err := m.etelecomQuery.Dispatch(ctx, queryHotline); err != nil {
			return err
		}
		if queryHotline.Result.IsFreeCharge.Bool {
			return nil
		}
	}

	// convert to minute
	// update 08-01-2021
	// Telecom change pricelist,
	// Postage will calculated by duration (second)
	// So don't need to use DurationForPostage (may be delete later)
	durationMinute := int(math.Ceil(float64(event.Duration) / SecondsPerMinute))
	update := &etelecom.UpdateCallLogPostageCommand{
		ID:                 event.ID,
		DurationForPostage: durationMinute,
	}
	defer func() error {
		return m.etelecomAggr.Dispatch(ctx, update)
	}()
	if event.Direction != call_direction.Out {
		return nil
	}
	phone := event.Callee
	if phone == "" {
		return nil
	}
	calcPostageArgs := postage.CalcPostageArgs{
		Phone:          phone,
		Direction:      event.Direction,
		DurationSecond: event.Duration,
	}
	fee := postage.CalcPostage(calcPostageArgs)
	update.Postage = fee
	return nil
}
