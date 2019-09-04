package main

import (
	"context"
	"math/rand"
	"time"

	"etop.vn/api/supporting/crm/vtiger"
	"etop.vn/backend/pkg/common/scheduler"
)

var (
	vtigerAggr                vtiger.CommandBus
	vtigerQS                  vtiger.QueryBus
	vtigerDefaultRecurr       = 30 * time.Second
	vtigerDefaultErrRecurr    = 5 * time.Minute
	LastTimeSyncVtigerContact = time.Now().Add(time.Duration(-24*30*8) * time.Hour)
)

func SyncVtiger(vtigerAggregate vtiger.CommandBus, vtigerQuery vtiger.QueryBus) {
	vtigerAggr = vtigerAggregate
	vtigerQS = vtigerQuery
	gScheduler = scheduler.New(defaultNumWorkers)

	t := rand.Intn(int(time.Second))
	gScheduler.AddAfter(0, time.Duration(t), SyncVtigerData)
	gScheduler.Start()
}

func SyncVtigerData(id interface{}, p scheduler.Planner) (_err error) {
	ll.S.Info("run SyncVtigerData", time.Now())
	defer func() {
		err := recover()
		GetLastVtigerModifytimeSyncInDB()
		if err != nil {
			ll.S.Info("Add after err :: ", vtigerDefaultErrRecurr)
			p.AddAfter(id, vtigerDefaultErrRecurr, SyncVtigerData)
		} else {
			ll.S.Info("Add after success :: ", vtigerDefaultRecurr)
			p.AddAfter(id, vtigerDefaultRecurr, SyncVtigerData)
		}
	}()
	if err := SyncVtigerAccount(); err != nil {
		return err
	}
	if err := SynVtigetContact(); err != nil {
		return err
	}
	return nil
}

// TODO crm is curently using this function new feature
func SyncVtigerAccount() error {
	return nil
}

func SynVtigetContact() error {
	ctx = context.Background()
	cmd := &vtiger.SyncContactCommand{
		SyncTime: LastTimeSyncVtigerContact,
	}
	err := vtigerAggr.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
}

func GetLastVtigerModifytimeSyncInDB() {
	query := &vtiger.GetLastTimeModifyQuery{
		Offset: 0,
		Limit:  1,
	}
	err := vtigerQS.Dispatch(ctx, query)
	if err == nil && query.Result != nil {
		LastTimeSyncVtigerContact = query.Result.Modifiedtime
	}
}
