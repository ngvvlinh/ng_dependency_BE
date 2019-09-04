package main

import (
	"context"
	"math/rand"
	"time"

	"etop.vn/api/supporting/crm/vht"
	"etop.vn/backend/pkg/common/scheduler"
)

var (
	gScheduler     *scheduler.Scheduler
	lasTimeSyncVht = time.Now().Add(time.Duration(-24*30*8) * time.Hour)
	VhtAggr        vht.CommandBus
	VhtQS          vht.QueryBus
)

const (
	defaultNumWorkers   = 32
	vhtDefaultRecurr    = 60 * time.Second
	vhtDefaultErrRecurr = 5 * time.Minute
)

func SyncCallHistoryVht(vhtAggr vht.CommandBus, vhtQS vht.QueryBus) {
	VhtQS = vhtQS
	ctx = context.Background()
	VhtAggr = vhtAggr
	cmd := vht.PingServerVhtCommand{}
	err := vhtAggr.Dispatch(ctx, &cmd)
	if err != nil {
		ll.Error("Can't connect to vht server")
		return
	}

	GetLasttimeSyncInDB()

	gScheduler = scheduler.New(defaultNumWorkers)

	t := rand.Intn(int(time.Second))
	gScheduler.AddAfter(0, time.Duration(t), SyncVhtCallHistoryData)
	gScheduler.Start()
}

func SyncVhtCallHistoryData(id interface{}, p scheduler.Planner) (_err error) {
	ll.S.Info("run syncUnCompleteFfms", time.Now())
	defer func() {
		GetLasttimeSyncInDB()
		err := recover()
		if err != nil {
			ll.S.Info("Add after err :: ", vhtDefaultErrRecurr)
			p.AddAfter(id, vhtDefaultErrRecurr, SyncVhtCallHistoryData)
		} else {
			ll.S.Info("Add after success :: ", vhtDefaultRecurr)
			p.AddAfter(id, vhtDefaultRecurr, SyncVhtCallHistoryData)
		}
	}()

	ctx = context.Background()
	cmd := &vht.SyncVhtCallHistoriesCommand{
		SyncTime: lasTimeSyncVht,
	}
	err := VhtAggr.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
}

func GetLasttimeSyncInDB() {
	query := &vht.GetLastCallHistoryQuery{
		Offset: 0,
		Limit:  1,
	}
	err := VhtQS.Dispatch(ctx, query)
	if err == nil && query.Result != nil {
		lasTimeSyncVht = query.Result.TimeStarted
	}
}
