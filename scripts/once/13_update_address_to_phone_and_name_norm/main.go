package main

import (
	"sync"

	"o.o/backend/cmd/etop-server/config"
	"o.o/backend/com/main/shipping/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/validate"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	cfg config.Config
	ll  = l.New()
	db  *cmsql.Database
)

func main() {
	var err error
	if cfg, err = config.Load(false); err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}
	postgres := cfg.Databases.Postgres

	if db, err = cmsql.Connect(postgres); err != nil {
		ll.Fatal("Error while connection database", l.Error(err))
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	var fromID dot.ID = 0
	count, updated, errCount := 0, 0, 0
	for {
		ffms, err := scanFulfillments(fromID)
		if err != nil {
			ll.Fatal("Error", l.Error(err))
		}
		ll.S.Infof("=========== Updated =========: %v/%v", updated, count)
		ll.S.Infof("=========== Error ========: %v/%v", errCount, count)
		if len(ffms) == 0 {
			ll.S.Infof("Done")
			break
		}

		fromID = ffms[len(ffms)-1].ID
		count += len(ffms)

		for _, ffm := range ffms {
			go func(_ffm *model.Fulfillment) {
				wg.Add(1)
				update := map[string]interface{}{
					"address_to_phone":          _ffm.AddressTo.Phone,
					"address_to_full_name_norm": validate.NormalizeSearch(_ffm.AddressTo.FullName),
				}
				err := db.Table("fulfillment").In("id", _ffm.ID).ShouldUpdateMap(update)

				mu.Lock()
				if err != nil {
					errCount += 1
				} else {
					updated += 1
				}
				mu.Unlock()

				wg.Done()
			}(ffm)
		}

		wg.Wait()
	}
}

func scanFulfillments(fromID dot.ID) (ffms model.Fulfillments, err error) {
	err = db.Where("id > ?", fromID.String()).
		Where("address_to IS NOT NULL").
		OrderBy("id").
		Limit(1000).
		Find(&ffms)
	return
}
