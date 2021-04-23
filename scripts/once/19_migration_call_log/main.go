package main

import (
	"flag"
	"sync"

	"o.o/backend/com/etelecom/model"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/scripts/once/19_migration_call_log/config"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	ll         = l.New()
	cfg        config.Config
	DBMain     *cmsql.Database
	DBEtelecom *cmsql.Database
)

func main() {
	cc.InitFlags()
	flag.Parse()

	var err error
	if cfg, err = config.Load(); err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}

	if DBMain, err = cmsql.Connect(cfg.Postgres); err != nil {
		ll.Fatal("Error while loading new database", l.Error(err))
	}
	if DBEtelecom, err = cmsql.Connect(cfg.PostgresEtelecom); err != nil {
		ll.Fatal("Error while loading new database etelecom", l.Error(err))
	}

	mapExtension := make(map[dot.ID]*model.Extension)
	mapCallLog := make(map[dot.ID]*model.CallLog)
	mapHotline := make(map[dot.ID]*model.Hotline)
	{
		var fromID dot.ID
		for {
			callLogs, err := scanCallLogs(fromID)
			if err != nil {
				ll.Fatal("Error", l.Error(err))
			}
			if len(callLogs) == 0 {
				break
			}
			hotlineIDs := make([]string, 0, len(callLogs))
			extensionIDs := make([]string, 0, len(callLogs))
			for _, callLog := range callLogs {
				mapCallLog[callLog.ID] = callLog
				if callLog.ExtensionID != 0 {
					extensionIDs = append(extensionIDs, callLog.ExtensionID.String())
				}

				if callLog.HotlineID != 0 {
					hotlineIDs = append(hotlineIDs, callLog.HotlineID.String())
				}
			}

			hotlines, err := getHotlinesByIDs(hotlineIDs)
			if err != nil {
				ll.Fatal("Error", l.Error(err))
			}
			for _, hotline := range hotlines {
				if _, ok := mapHotline[hotline.ID]; !ok {
					mapHotline[hotline.ID] = hotline
				}
			}

			extensions, err := getExtensionsByIDs(extensionIDs)
			if err != nil {
				ll.Fatal("Error", l.Error(err))
			}
			for _, extension := range extensions {
				if _, ok := mapExtension[extension.ID]; !ok {
					mapExtension[extension.ID] = extension
				}
			}

			fromID = callLogs[len(callLogs)-1].ID
		}
	}

	{
		count, errCount, updatedCount := len(mapCallLog), 0, 0
		var (
			mu sync.Mutex
			wg sync.WaitGroup
		)

		maxGoroutines := 8
		ch := make(chan dot.ID, maxGoroutines)
		wg.Add(len(mapCallLog))
		for _, callLog := range mapCallLog {
			ch <- callLog.ID
			var userID, ownerID dot.ID

			if _, ok := mapHotline[callLog.HotlineID]; ok {
				ownerID = mapHotline[callLog.HotlineID].OwnerID
			}
			if _, ok := mapExtension[callLog.ExtensionID]; ok {
				userID = mapExtension[callLog.ExtensionID].UserID
			}

			go func(id, userID dot.ID, ownerID dot.ID) {
				defer func() {
					<-ch
					wg.Done()
				}()
				err = updateCallLog(id, userID, ownerID)
				if err != nil {
					mu.Lock()
					errCount++
					mu.Unlock()
				}
			}(callLog.ID, userID, ownerID)
		}

		wg.Wait()
		updatedCount = count - errCount
		ll.S.Infof("Update owner_id, user_id for callLogs: updated %v/%v, error %v/%v", updatedCount, count, errCount, count)
	}
}

func scanCallLogs(fromID dot.ID) (callLogs model.CallLogs, err error) {
	err = DBEtelecom.Where("id > ?", fromID.String()).
		Where("user_id is NULL OR owner_id is NULL").
		OrderBy("id").
		Limit(1000).
		Find(&callLogs)
	return
}

func updateCallLog(ID dot.ID, userID dot.ID, ownerID dot.ID) (err error) {
	update := map[string]interface{}{
		"user_id":  userID,
		"owner_id": ownerID,
	}

	err = DBEtelecom.
		Table("call_log").
		Where("id=?", ID).
		ShouldUpdateMap(update)
	return err
}

func getExtensionsByIDs(ids []string) (extensions model.Extensions, err error) {
	err = DBEtelecom.
		From("extension").
		In("id", ids).
		Find(&extensions)
	return
}

func getHotlinesByIDs(ids []string) (hotlines model.Hotlines, err error) {
	err = DBEtelecom.
		From("hotline").
		In("id", ids).
		Find(&hotlines)
	return
}
