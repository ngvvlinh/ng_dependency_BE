package main

import (
	"sync"

	"github.com/lib/pq"

	orderingtypes "o.o/api/main/ordering/types"
	"o.o/backend/cmd/etop-server/config"
	model "o.o/backend/com/main/ordering/model"
	shipnowfulfillmentmodel "o.o/backend/com/main/shipnow/model"
	"o.o/backend/pkg/common/sql/cmsql"
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
		orders, err := scanOrders(fromID)
		if err != nil {
			ll.Fatal("Error", l.Error(err))
		}
		ll.S.Infof("=========== Updated =========: %v/%v", updated, count)
		ll.S.Infof("=========== Error ========: %v/%v", errCount, count)
		if len(orders) == 0 {
			ll.S.Infof("Done")
			break
		}

		fromID = orders[len(orders)-1].ID
		count += len(orders)

		for _, order := range orders {
			go func(_order *model.Order) {
				wg.Add(1)

				var err error
				defer func() {
					mu.Lock()
					if err != nil {
						errCount += 1
					} else {
						updated += 1
					}
					mu.Unlock()

					wg.Done()
				}()

				var ffm shipnowfulfillmentmodel.ShipnowFulfillment
				_, err = db.Where("id = ?", _order.FulfillmentIDs[0]).Get(&ffm)
				if err != nil {
					return
				}

				update := map[string]interface{}{
					"fulfillment_shipping_codes": pq.StringArray([]string{ffm.ShippingCode}),
				}
				err = db.Table("order").In("id", _order.ID).ShouldUpdateMap(update)
				return
			}(order)
		}

		wg.Wait()
	}
}

func scanOrders(fromID dot.ID) (orders model.Orders, err error) {
	err = db.
		Where("id > ?", fromID.String()).
		Where("fulfillment_type = ?", orderingtypes.ShippingTypeShipnow.Enum()).
		OrderBy("id").
		Limit(1000).
		Find(&orders)
	return
}
