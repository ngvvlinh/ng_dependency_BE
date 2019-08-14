package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"etop.vn/api/meta"

	"etop.vn/backend/pkg/services/crm-service/model"

	cm "etop.vn/backend/pkg/common"

	"etop.vn/backend/pkg/common/cmsql"

	"etop.vn/backend/pkg/common/scheduler"
	vhtservice "etop.vn/backend/pkg/services/crm-service/vht/client"
)

var (
	gScheduler       *scheduler.Scheduler
	defaultRecurr    = 30 * time.Second
	defaultErrRecurr = 5 * time.Minute
	lasTimeSync      = time.Now().Add(time.Duration(-24*30*8) * time.Hour)
	VhtService       *vhtservice.Client
)

const (
	defaultNumWorkers = 32
)

func SyncCallHistoryVht(userName string, passWord string, db cmsql.Database) {
	VhtService = vhtservice.NewClient(userName, passWord, db)
	if err := VhtService.PingServerVht(); err != nil {
		ll.Error("Can't connect to vht server")
		return
	}
	var paging meta.Paging
	paging.Offset = 1
	paging.Limit = 1

	lastRecord, err := VhtService.VhtCallHistoryStore(ctx).Paging(paging).ByStatus("Done").SortBy("time_started desc").GetCallHistories()
	if err == nil {
		if len(lastRecord) > 0 {
			lasTimeSync = lastRecord[0].TimeStarted
		}
	}

	gScheduler = scheduler.New(defaultNumWorkers)

	t := rand.Intn(int(time.Second))
	gScheduler.AddAfter(0, time.Duration(t), SyncVhtCallHistoryData)
	gScheduler.Start()
}

func StopSync() {
	gScheduler.Stop()
}

func SyncVhtCallHistoryData(id interface{}, p scheduler.Planner) (_err error) {
	ll.S.Info("run syncUnCompleteFfms", time.Now())
	defer func() {
		err := recover()
		if err != nil {
			ll.S.Info("Add after err :: ", defaultErrRecurr)
			p.AddAfter(id, defaultErrRecurr, SyncVhtCallHistoryData)
		} else {
			ll.S.Info("Add after success :: ", defaultRecurr)
			p.AddAfter(id, defaultRecurr, SyncVhtCallHistoryData)
		}
	}()

	ctx = context.Background()
	err := SyncVhtCallHistory(ctx)
	if err != nil {
		return err
	}
	err = SyncVhtCallHistoryPending(ctx)
	if err != nil {
		return err
	}
	err = SyncVhtCallHistoryPending(ctx)
	if err != nil {
		return err
	}
	return nil
}

func SyncVhtCallHistory(ctx context.Context) error {
	fromDate := lasTimeSync.Unix()
	toDate := time.Now().Unix()
	lasTimeSync = time.Now()

	queryDTO := &vhtservice.VHTHistoryQueryDTO{
		Page:        1,
		Limit:       50,
		DateStarted: fromDate,
		DateEnded:   toDate,
		SortBy:      "time_started",
		SortType:    "ASC",
	}
	for true {
		result, err := VhtService.GetHistories(queryDTO)
		if err != nil {
			return err
		}
		for i := 0; i < len(result.Items); i++ {
			data := vhtservice.ConvertToModel(result.Items[i])
			err = data.BeforeInsertOrUpdate()
			if err != nil {
				return err
			}
			query := VhtService.VhtCallHistoryStore(ctx).ByCallID(data.CallID)
			var oldData *model.VhtCallHistory
			oldData, err = query.GetCallHistory()
			data.OData = ""
			if err != nil && cm.ErrorCode(err) == cm.NotFound {
				err = VhtService.VhtCallHistoryStore(ctx).CreateVhtCallHistory(data)
			} else if err == nil {
				if oldData.SyncStatus == "Done" {
					continue
				}
				var oldDataMarshal []byte
				oldDataMarshal, err = json.Marshal(oldData)
				data.OData = string(oldDataMarshal)
				err = VhtService.VhtCallHistoryStore(ctx).ByCallID(data.CallID).UpdateVhtCallHistory(data)
			} else {
				return err
			}
		}
		if len(result.Items) < 50 {
			break
		}
		queryDTO.Page = queryDTO.Page + 1
	}
	return nil
}

func SyncVhtCallHistoryPending(ctx context.Context) error {
	historiesPending, err := VhtService.VhtCallHistoryStore(ctx).ByStatus("Pending").GetCallHistories()
	if err != nil {
		return nil
	}
	for i := 0; i < len(historiesPending); i++ {
		var result *vhtservice.VhtCallHistory
		result, err = VhtService.GetHistoryBySDKCallID(historiesPending[i].SdkCallID)
		if err != nil {
			return err
		}
		data := vhtservice.ConvertToModel(result)
		err = VhtService.VhtCallHistoryStore(ctx).BySdkCallID(data.SdkCallID).UpdateVhtCallHistory(data)
		if err != nil {
			return err
		}
	}
	return nil
}
