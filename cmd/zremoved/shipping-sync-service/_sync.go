package shipping_sync_service

import (
	"context"
	"math/rand"
	"time"

	"o.o/api/top/types/etc/shipping_provider"
	shipmodel "o.o/backend/com/main/shipping/model"
	shipmodelx "o.o/backend/com/main/shipping/modelx"
	"o.o/backend/pkg/common/apifw/scheduler"
	"o.o/backend/pkg/common/bus"
	ghnupdatev1 "o.o/backend/pkg/integration/shipping/ghn/update/v1"
	"o.o/backend/pkg/integration/shipping/ghtk"
	"o.o/common/l"
)

var (
	gScheduler       *scheduler.Scheduler
	defaultRecurr    = 2 * time.Hour
	defaultErrRecurr = 3 * time.Hour
)

const (
	defaultNumWorkers = 32
)

func SyncIncompleteFfms() {
	gScheduler = scheduler.New(defaultNumWorkers)

	// now := time.Now()
	// endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 0, 0, 0, time.Local)
	// if now.Before(endDate) {
	// 	gScheduler.AddAfter(0, endDate.Sub(now), syncUnCompleteFfms)
	// } else {
	// 	t := rand.Intn(int(time.Second))
	// 	gScheduler.AddAfter(0, time.Duration(t), syncUnCompleteFfms)
	// }

	t := rand.Intn(int(time.Second))
	gScheduler.AddAfter(0, time.Duration(t), syncUnCompleteFfms)
	gScheduler.Start()
}

func StopSync() {
	gScheduler.Stop()
}

func syncUnCompleteFfms(id interface{}, p scheduler.Planner) (_err error) {
	ll.S.Info("run syncUnCompleteFfms", time.Now())
	defer func() {
		err := recover()
		if err != nil {
			ll.S.Info("Add after err :: ", defaultErrRecurr)
			p.AddAfter(id, defaultErrRecurr, syncUnCompleteFfms)
		} else {
			// now := time.Now()
			// nextRunAt := now.Add(time.Hour * 24)
			// p.AddAfter(id, nextRunAt.Sub(now), syncUnCompleteFfms).Recurrent(nextRunAt.Sub(now))
			ll.S.Info("Add after success :: ", defaultRecurr)
			p.AddAfter(id, defaultRecurr, syncUnCompleteFfms)
		}
	}()

	ctx := context.Background()
	cmd := &shipmodelx.GetUnCompleteFulfillmentsQuery{
		ShippingProviders: []shipping_provider.ShippingProvider{
			shipping_provider.GHN, shipping_provider.GHTK,
		},
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	var updateFfms, updateFfmsGHN, updateFfmsGHTK []*shipmodel.Fulfillment
	ll.S.Info("uncomplete order :: ", len(cmd.Result))
	for _, ffm := range cmd.Result {
		switch ffm.ShippingProvider {
		case shipping_provider.GHN:
			updateFfmsGHN = append(updateFfmsGHN, ffm)
		case shipping_provider.GHTK:
			updateFfmsGHTK = append(updateFfmsGHTK, ffm)
			// TODO
		default:
			// Nothing
		}
	}
	if len(updateFfmsGHN) > 0 {
		ffms, _ := ghnupdatev1.SyncTrackingOrders(updateFfmsGHN)
		for _, ffm := range ffms {
			updateFfms = append(updateFfms, &shipmodel.Fulfillment{
				ID:                   ffm.ID,
				ExternalShippingLogs: ffm.ExternalShippingLogs,
			})
		}
	}
	if len(updateFfmsGHTK) > 0 {
		ffms, err := ghtk.SyncOrders(updateFfmsGHTK)
		if err == nil {
			updateFfms = append(updateFfms, ffms...)
		}
	}
	if len(updateFfms) > 0 {
		cmdUpdate := &shipmodelx.UpdateFulfillmentsWithoutTransactionCommand{
			Fulfillments: updateFfms,
		}
		if err := bus.Dispatch(ctx, cmdUpdate); err != nil {
			ll.Debug("Không thể cập nhật ffm", l.Error(err))
		}
	}
	return nil
}
