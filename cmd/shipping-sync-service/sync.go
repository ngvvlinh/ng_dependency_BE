package main

import (
	"context"
	"math/rand"
	"time"

	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/common/scheduler"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/integration/ghn/update"
	"etop.vn/backend/pkg/integration/ghtk"
)

var (
	gScheduler       *scheduler.Scheduler
	defaultRecurr    = 2 * time.Hour
	defaultErrRecurr = 3 * time.Hour
)

const (
	defaultNumWorkers = 32
)

func SyncUnCompleteFfms() {
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
	cmd := &model.GetUnCompleteFulfillmentsQuery{
		ShippingProviders: []model.ShippingProvider{
			model.TypeGHN, model.TypeGHTK,
		},
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	var updateFfms, updateFfmsGHN, updateFfmsGHTK []*model.Fulfillment
	ll.S.Info("uncomplete order :: ", len(cmd.Result))
	for _, ffm := range cmd.Result {
		switch ffm.ShippingProvider {
		case model.TypeGHN:
			updateFfmsGHN = append(updateFfmsGHN, ffm)
		case model.TypeGHTK:
			updateFfmsGHTK = append(updateFfmsGHTK, ffm)
			// TODO
		default:
			// Nothing
		}
	}
	if len(updateFfmsGHN) > 0 {
		ffms, _ := update.SyncTrackingOrders(updateFfmsGHN)
		for _, ffm := range ffms {
			updateFfms = append(updateFfms, &model.Fulfillment{
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
		cmdUpdate := &model.UpdateFulfillmentsWithoutTransactionCommand{
			Fulfillments: updateFfms,
		}
		if err := bus.Dispatch(ctx, cmdUpdate); err != nil {
			ll.Debug("Không thể cập nhật ffm", l.Error(err))
		}
	}
	return nil
}
