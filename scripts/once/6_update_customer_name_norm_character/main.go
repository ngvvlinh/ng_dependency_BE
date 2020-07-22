package main

import (
	"context"
	"flag"

	"o.o/backend/cmd/etop-server/config"
	"o.o/backend/com/shopping/customering/model"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/validate"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	ll  = l.New()
	cfg config.Config
	db  *cmsql.Database
)

type M map[string]interface{}

func main() {
	cc.InitFlags()
	flag.Parse()

	var err error
	if cfg, err = config.Load(false); err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}
	cmenv.SetEnvironment("script", cfg.SharedConfig.Env)

	postgres := cfg.Databases.Postgres

	errCount, maxGoroutines := 0, 8
	ch := make(chan dot.ID, maxGoroutines)
	chInsert := make(chan error, maxGoroutines)

	if db, err = cmsql.Connect(postgres); err != nil {
		ll.Fatal("Error while connecting database", l.Error(err))
	}
	{
		var fromID dot.ID = 0
		count, updated := 0, 0
		for {
			customers, err := scanCustomer(fromID)
			if err != nil {
				ll.Fatal("Error", l.Error(err))
			}
			if len(customers) == 0 {
				ll.S.Infof("Done: updated %v/%v", updated, count)
				break
			}
			fromID = customers[len(customers)-1].ID
			count += len(customers)
			for _, customer := range customers {

				ch <- customer.ID
				go func(p *model.ShopCustomer) (_err error) {
					_, ctxCancel := context.WithCancel(context.Background())
					defer func() {
						<-ch
						chInsert <- _err
						ctxCancel()
					}()

					nameNorm := validate.NormalizeSearchCharacter(p.FullName)
					update := make(map[string]interface{})
					if p.FullNameNorm != "" && p.FullName != "" {
						update["full_name_norm"] = nameNorm
					}
					if len(update) > 0 {
						_err = db.
							Table("shop_customer").
							Where("id = ?", p.ID).
							ShouldUpdateMap(update)
					}
					return _err
				}(customer)
			}
			for i := 0; i < len(customers); i++ {
				err = <-chInsert
				if err != nil {
					errCount++
				} else {
					updated++
				}
			}
		}
		ll.S.Infof("Updated shop customer: success %v/%v, error %v/%v", updated, count, errCount, count)
	}
}

func scanCustomer(fromID dot.ID) (customers model.ShopCustomers, err error) {
	err = db.
		Where("id > ?", fromID.String()).
		OrderBy("id").
		Limit(1000).
		Find(&customers)
	return
}
