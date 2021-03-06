package main

import (
	"context"
	"flag"

	"o.o/api/shopping/customering"
	"o.o/api/top/types/etc/customer_type"
	"o.o/backend/cmd/etop-server/config"
	identitymodel "o.o/backend/com/main/identity/model"
	customeraggregate "o.o/backend/com/shopping/customering/aggregate"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	ll       = l.New()
	cfg      config.Config
	db       *cmsql.Database
	eventBus capi.EventBus
)

func main() {
	cc.InitFlags()
	flag.Parse()

	var err error
	if cfg, err = config.Load(false); err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}
	cmenv.SetEnvironment(cfg.Env)

	postgres := cfg.Postgres
	postgres.Database = "etop"
	postgres.Port = 5432
	postgres.Username = "postgres"
	postgres.Password = "postgres"

	if db, err = cmsql.Connect(postgres); err != nil {
		ll.Fatal("Error while connecting database", l.Error(err))
	}

	customerAggr := customeraggregate.NewCustomerAggregate(eventBus, db).MessageBus()
	{
		var arrayShops []identitymodel.Shops
		var fromID dot.ID
		count, errCount, createdCount, maxGoroutines := 0, 0, 0, 8
		ch := make(chan dot.ID, maxGoroutines)
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
				go func(shopID dot.ID) (_err error) {
					ctx, ctxCancel := context.WithCancel(context.Background())
					defer func() {
						<-ch
						chInsert <- _err
						ctxCancel()
					}()

					createCustomerCmd := &customering.CreateCustomerCommand{
						ShopID:   shopID,
						FullName: "Khách lẻ",
						Type:     customer_type.Independent,
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

func scanShops(fromID dot.ID) (shops identitymodel.Shops, err error) {
	err = db.
		Where("id > ?", fromID).
		OrderBy("id").
		Limit(1000).
		Find(&shops)
	return
}
