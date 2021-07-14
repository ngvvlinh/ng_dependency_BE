package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"o.o/backend/com/main/identity/model"
	cm "o.o/backend/pkg/common"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/scripts/once/23_add_id_for_account_user/config"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	ll     = l.New()
	cfg    config.Config
	DBMain *cmsql.Database
	mu     sync.Mutex
	wg     sync.WaitGroup
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
	start := time.Time{}
	count, updateCount, errCount := 0, 0, 0
	for {
		accUsers, err := scanAccountUsers(start)
		if err != nil {
			ll.Fatal(fmt.Sprintf("%v", err))
		}

		if len(accUsers) == 0 {
			break
		}
		count += len(accUsers)

		maxGoroutines := 8
		ch := make(chan dot.ID, maxGoroutines)
		wg.Add(len(accUsers))
		for _, accUser := range accUsers {
			id := dot.ID(cm.NewIDWithTime(accUser.CreatedAt))
			ch <- id
			go func(id, accountID, userID dot.ID) {
				defer func() {
					<-ch
					wg.Done()
				}()
				err = updateAccountUser(id, accountID, userID)
				if err != nil {
					mu.Lock()
					errCount++
					mu.Unlock()
				}
			}(id, accUser.AccountID, accUser.UserID)
			start = accUser.CreatedAt
		}
		wg.Wait()
	}
	updateCount = count - errCount
	ll.S.Infof("Update id for account_user: updated %v/%v, error %v/%v", updateCount, count, errCount, count)
}

func scanAccountUsers(from time.Time) (model.AccountUsers, error) {
	var accountUsers model.AccountUsers
	err := DBMain.
		Where("created_at > ? AND deleted_at IS NULL", from).
		OrderBy("created_at ASC").
		Limit(1000).
		Find(&accountUsers)
	if err != nil {
		return nil, err
	}
	return accountUsers, nil
}

func updateAccountUser(id, accountID, userID dot.ID) error {
	update := map[string]interface{}{
		"id": id,
	}
	err := DBMain.
		Table("account_user").
		Where("created_at IS NOT NULL AND account_id = ? AND user_id = ?", accountID, userID).
		ShouldUpdateMap(update)
	return err
}
