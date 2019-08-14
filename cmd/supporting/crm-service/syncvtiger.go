package main

import (
	"math/rand"
	"time"

	"etop.vn/backend/pkg/common/cmsql"

	"etop.vn/backend/pkg/common/scheduler"
	mapVtiger "etop.vn/backend/pkg/services/crm-service/mapping"
	vtigerService "etop.vn/backend/pkg/services/crm-service/vtiger/service"
)

var (
	Vs *vtigerService.VtigerService
)

func SyncVtiger(db cmsql.Database, vConfig vtigerService.Config, fieldMap mapVtiger.ConfigMap) {
	Vs = vtigerService.NewSVtigerService(db, vConfig, fieldMap)
	_, err := Vs.Client.GetSessionKey(Vs.Cfg.ServiceURL, Vs.Cfg.Username, Vs.Cfg.APIKey)
	if err != nil {
		ll.Error("Can't connnect to vtiger")
		return
	}

	gScheduler = scheduler.New(defaultNumWorkers)

	t := rand.Intn(int(time.Second))
	gScheduler.AddAfter(0, time.Duration(t), SyncVtigerData)
	gScheduler.Start()
}

func SyncVtigerData(id interface{}, p scheduler.Planner) (_err error) {
	ll.S.Info("run SyncVtigerData", time.Now())
	defer func() {
		err := recover()
		if err != nil {
			ll.S.Info("Add after err :: ", defaultErrRecurr)
			p.AddAfter(id, defaultErrRecurr, SyncVtigerData)
		} else {
			ll.S.Info("Add after success :: ", defaultRecurr)
			p.AddAfter(id, defaultRecurr, SyncVtigerData)
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
	err := Vs.SyncContac()
	if err != nil {
		return err
	}
	return nil
}
