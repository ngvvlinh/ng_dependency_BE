package main

import (
	"flag"
	"time"

	"etop.vn/api/main/ledgering"
	"etop.vn/backend/cmd/etop-server/config"
	ledgeringmodel "etop.vn/backend/com/main/ledgering/model"
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

	if db, err = cmsql.Connect(cfg.Postgres); err != nil {
		ll.Fatal("Error while connecting database", l.Error(err))
	}

	mapShopIDsHaveLedgers := make(map[int64]bool)
	{
		fromID := int64(0)

		for {
			shopLedgers, err := scanShopLedgersTypeCash(fromID)
			if err != nil {
				ll.Fatal("Error", l.Error(err))
			}
			if len(shopLedgers) == 0 {
				break
			}
			for _, shopLedger := range shopLedgers {
				mapShopIDsHaveLedgers[shopLedger.ShopID] = true
			}
			fromID = shopLedgers[len(shopLedgers)-1].ID
		}
	}

	// (Map) shopID (userID belongs to shopID) that need to create default shop ledger
	mapshopIDUserID := make(map[int64]int64)
	{
		fromID := int64(0)

		for {
			shops, err := scanShops(fromID)
			if err != nil {
				ll.Fatal("Error", l.Error(err))
			}
			if len(shops) == 0 {
				break
			}
			fromID = shops[len(shops)-1].ID

			for _, shop := range shops {
				if _, ok := mapShopIDsHaveLedgers[shop.ID]; !ok {
					mapshopIDUserID[shop.ID] = shop.OwnerID
				}
			}
		}
	}

	{
		count, created := len(mapshopIDUserID), 0
		for key, value := range mapshopIDUserID {
			_, err = db.
				Table("shop_ledger").
				Insert(&ledgeringmodel.ShopLedger{
					ID:          cm.NewID(),
					ShopID:      key,
					Name:        "Tiền mặt",
					BankAccount: nil,
					Note:        "Số quỹ mặc định",
					Type:        string(ledgering.LedgerTypeCash),
					Status:      0,
					CreatedBy:   value,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})

			if err != nil {
				ll.S.Errorf("Insert ledger for shop %v", value)
			} else {
				created++
			}
		}
		ll.S.Infof("Created %v/%v", created, count)
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

func scanShopLedgersTypeCash(fromID int64) (shopLedgers ledgeringmodel.ShopLedgers, err error) {
	err = db.
		Where("id > ? and type = ?", fromID, string(ledgering.LedgerTypeCash)).
		OrderBy("id").
		Limit(1000).
		Find(&shopLedgers)
	return
}
