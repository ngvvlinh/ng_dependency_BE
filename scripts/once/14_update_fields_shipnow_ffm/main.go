package main

import (
	"strings"
	"sync"

	"o.o/backend/cmd/etop-server/config"
	model "o.o/backend/com/main/shipnow/model"
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
		ffms, err := scanShipnowFulfillments(fromID)
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
			go func(_ffm *model.ShipnowFulfillment) {
				wg.Add(1)
				phoneNorm, fullNameNorm := toAddressToPhoneAndAddressToFullNameNorm(_ffm.DeliveryPoints)
				update := map[string]interface{}{
					"address_to_phone":          phoneNorm,
					"address_to_full_name_norm": fullNameNorm,
				}
				err := db.Table("shipnow_fulfillment").In("id", _ffm.ID).ShouldUpdateMap(update)

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

func scanShipnowFulfillments(fromID dot.ID) (ffms model.ShipnowFulfillments, err error) {
	err = db.Where("id > ?", fromID.String()).
		OrderBy("id").
		Limit(1000).
		Find(&ffms)
	return
}

func toAddressToPhoneAndAddressToFullNameNorm(deliveryPoints []*model.DeliveryPoint) (addressToPhoneNorm, addressToFullNameNorm string) {
	var phones, fullNames []string

	for _, deliveryPoint := range deliveryPoints {
		phones = append(phones, deliveryPoint.ShippingAddress.Phone)
		fullNames = append(fullNames, deliveryPoint.ShippingAddress.FullName)
	}

	{
		tempNorm := validate.NormalizeSearch(strings.Join(phones, " "))
		hash := make(map[string]bool)

		var words []string
		for _, word := range strings.Split(tempNorm, " ") {
			if _, ok := hash[word]; ok {
				continue
			}
			words = append(words, word)
			hash[word] = true
		}

		addressToPhoneNorm = strings.Join(words, " ")
	}

	{
		tempNorm := validate.NormalizeSearch(strings.Join(fullNames, " "))
		hash := make(map[string]bool)

		var words []string
		for _, word := range strings.Split(tempNorm, " ") {
			if _, ok := hash[word]; ok {
				continue
			}
			words = append(words, word)
			hash[word] = true
		}

		addressToFullNameNorm = strings.Join(words, " ")
	}

	return
}
