package main

import (
	"flag"

	"etop.vn/backend/cmd/etop-server/config"
	customering "etop.vn/backend/com/shopping/customering/model"
	cm "etop.vn/backend/pkg/common"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
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
	cm.SetEnvironment(cfg.Env)

	postgres := cfg.Postgres

	if db, err = cmsql.Connect(postgres); err != nil {
		ll.Fatal("Error while connecting database", l.Error(err))
	}
	{
		var fromID dot.ID
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
				phoneNorm := validate.NormalizeSearchPhone(customer.Phone)
				fullNameNorm := validate.NormalizeSearch(customer.FullName)
				update := M{}
				if customer.PhoneNorm == "" || customer.PhoneNorm != phoneNorm {
					update["phone_norm"] = phoneNorm
				}
				if customer.FullNameNorm == "" || customer.FullNameNorm != fullNameNorm {
					update["full_name_norm"] = fullNameNorm
				}
				if len(update) > 0 {
					err = db.
						Table("shop_customer").
						Where("id = ?", customer.ID).
						ShouldUpdateMap(update)
					if err != nil {
						ll.S.Fatalf("can't update shop_customer id=%v", customer.ID)
					}
					updated++
				}
			}
			ll.S.Infof("Updated %v/%v", updated, count)
		}
	}
}

func scanCustomer(fromID dot.ID) (customers customering.ShopCustomers, err error) {
	err = db.
		Where("id > ?", fromID).
		OrderBy("id").
		Limit(1000).
		Find(&customers)
	return
}
