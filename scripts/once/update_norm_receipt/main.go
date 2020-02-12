package main

import (
	"flag"

	"etop.vn/backend/cmd/etop-server/config"
	receipting "etop.vn/backend/com/main/receipting/model"
	"etop.vn/backend/pkg/common/cmenv"
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
	cmenv.SetEnvironment(cfg.Env)

	postgres := cfg.Postgres

	if db, err = cmsql.Connect(postgres); err != nil {
		ll.Fatal("Error while connecting database", l.Error(err))
	}
	var fromID dot.ID
	count, updated := 0, 0
	var arrayReceipt []*receipting.Receipt
	for {
		receipts, err := scanReceipt(fromID)
		if err != nil {
			ll.Fatal("Error", l.Error(err))
		}
		if len(receipts) == 0 {
			break
		}
		fromID = receipts[len(receipts)-1].ID
		count += len(receipts)
		arrayReceipt = append(arrayReceipt, receipts...)
	}
	for _, receipt := range arrayReceipt {
		traderFullNameNorm := validate.NormalizeSearch(receipt.Trader.FullName)
		update := M{}

		if receipt.TraderFullNameNorm == "" || receipt.TraderFullNameNorm != traderFullNameNorm {
			update["trader_full_name_norm"] = traderFullNameNorm
		}
		if len(update) > 0 {
			err = db.
				Table("receipt").
				Where("id = ?", receipt.ID).
				ShouldUpdateMap(update)
			if err != nil {
				ll.S.Fatalf("can't update receipt id=%v", receipt.ID)
			}
			updated++
		}
	}
	ll.S.Infof("Updated %v/%v", updated, count)
}

func scanReceipt(fromID dot.ID) (receipts receipting.Receipts, err error) {
	err = db.
		Where("id > ?", fromID).
		OrderBy("id").
		Limit(1000).
		Find(&receipts)
	return
}
