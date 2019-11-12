package main

import (
	"context"
	"flag"

	"etop.vn/api/shopping/customering"
	"etop.vn/backend/cmd/etop-server/config"
	customeraggregate "etop.vn/backend/com/shopping/customering/aggregate"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/common/l"
)

var (
	ll  = l.New()
	cfg config.Config
	db  *cmsql.Database
)

func main() {
	cc.InitFlags()
	flag.Parse()

	var err error
	if cfg, err = config.Load(false); err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}
	cm.SetEnvironment(cfg.Env)

	postgres := cfg.Postgres

	if db, err = cmsql.Connect(postgres); err != nil {
		ll.Fatal("Error while connecting database", l.Error(err))
	}

	customerAggr := customeraggregate.NewCustomerAggregate(db).MessageBus()
	{
		var arrayShops []model.Shops
		fromID := int64(0)
		count, errCount, createdCount, maxGoroutines := 0, 0, 0, 8
		ch := make(chan int64, maxGoroutines)
		chInsert := make(chan error, maxGoroutines)
		for {
			shops, err := scanShops(fromID)
			if err != nil {
				ll.Fatal("Error", l.Error(err))
			}
			if len(shops) == 0 {
				break
			}
			fromID = shops[len(shops)-1].ID

			count += len(shops)
			arrayShops = append(arrayShops, shops)
		}
		for _, shops := range arrayShops {
			for _, shop := range shops {
				ch <- shop.ID
				go func(shopID int64) (_err error) {
					ctx, ctxCancel := context.WithCancel(context.Background())
					defer func() {
						<-ch
						chInsert <- _err
						ctxCancel()
					}()

					createCustomerCmd := &customering.CreateCustomerCommand{
						ShopID:   shopID,
						FullName: "Khách lẻ",
						Type:     customering.CustomerTypeIndependent,
					}
					_err = customerAggr.Dispatch(ctx, createCustomerCmd)

					return _err
				}(shop.ID)
			}
			for i := 0; i < len(shops); i++ {
				err := <-chInsert
				if err != nil {
					errCount++
				} else {
					createdCount++
				}
			}
		}
		ll.S.Infof("Created shop independent customer: success %v/%v, error %v/%v", createdCount, count, errCount, count)
	}
}

func scanShops(fromID int64) (shops model.Shops, err error) {
	err = db.
		Where("id > ?", fromID).
		OrderBy("id").
		Limit(1000).
		Find(&shops)
	return
}
