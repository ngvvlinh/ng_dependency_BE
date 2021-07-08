package main

import (
	"flag"
	"fmt"
	"sync"

	"o.o/backend/com/main/contact/model"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/scripts/once/21_migration_contact/config"
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

	{
		count, updateCount, errCount := 0, 0, 0
		var fromID dot.ID
		for {
			contacts, err := scanContacts(fromID)
			if err != nil {
				ll.Fatal(fmt.Sprintf("%v", err))
			}
			if len(contacts) == 0 {
				break
			}
			count += len(contacts)
			maxGoroutines := 8
			ch := make(chan dot.ID, maxGoroutines)
			wg.Add(len(contacts))
			for _, contact := range contacts {
				ch <- contact.ID
				fullNameNorm := validate.NormalizeSearchCharacter(contact.FullName)
				go func(id dot.ID, fullNameNorm string) {
					defer func() {
						<-ch
						wg.Done()
					}()
					err = updateContact(id, fullNameNorm)
					if err != nil {
						mu.Lock()
						errCount++
						mu.Unlock()
					}
				}(contact.ID, fullNameNorm)
				fromID = contact.ID
			}
			wg.Wait()
		}
		updateCount = count - errCount
		ll.S.Infof("Update full_name_norm  for contact: updated %v/%v, error %v/%v", updateCount, count, errCount, count)
	}
}

func scanContacts(fromID dot.ID) (contacts model.Contacts, err error) {
	err = DBMain.Where("id > ? AND full_name_norm IS NULL", fromID.String()).
		OrderBy("id").
		Limit(1000).
		Find(&contacts)
	return
}

func updateContact(id dot.ID, fullNameNorm string) (err error) {
	update := map[string]interface{}{
		"full_name_norm": fullNameNorm,
	}

	err = DBMain.
		Table("contact").
		Where("id=? ", id).
		ShouldUpdateMap(update)
	return err
}
