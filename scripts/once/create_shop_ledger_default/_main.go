package main

import (
	"flag"
	"time"

	"o.o/api/top/types/etc/ledger_type"
	"o.o/backend/cmd/etop-server/config"
	identitymodel "o.o/backend/com/main/identity/model"
	ledgeringmodel "o.o/backend/com/main/ledgering/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
	"o.o/common/l"
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
	cmenv.SetEnvironment(cfg.Env)

	postgres := cfg.Postgres

	if db, err = cmsql.Connect(postgres); err != nil {
		ll.Fatal("Error while connecting database", l.Error(err))
	}

	mapShopIDsHaveLedgers := make(map[dot.ID]bool)
	{
		var fromID dot.ID
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
	mapshopIDUserID := make(map[dot.ID]dot.ID)
	{
		var fromID dot.ID
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
		count, errCount, createdCount := len(mapshopIDUserID), 0, 0
		maxGoroutines := 8
		ch := make(chan dot.ID, maxGoroutines)
		chInsert := make(chan error, maxGoroutines)
		for key, value := range mapshopIDUserID {
			ch <- key
			shopID := key
			ownerID := value
			go func() (_err error) {
				defer func() {
					<-ch
					chInsert <- _err
				}()
				_, _err = db.
					Table("shop_ledger").
					Insert(&ledgeringmodel.ShopLedger{
						ID:          cm.NewID(),
						ShopID:      shopID,
						Name:        "Ti???n m???t",
						BankAccount: nil,
						Note:        "T??i kho???n thanh to??n m???c ?????nh",
						Type:        ledger_type.LedgerTypeCash,
						Status:      0,
						CreatedBy:   ownerID,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					})
				return _err
			}()
		}
		for i, n := 0, len(mapshopIDUserID); i < n; i++ {
			err := <-chInsert
			if err != nil {
				errCount++
			} else {
				createdCount++
			}
		}
		ll.S.Infof("Created shop ledger tien mat: success %v/%v, error %v/%v", createdCount, count, errCount, count)
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

func scanShopLedgersTypeCash(fromID dot.ID) (shopLedgers ledgeringmodel.ShopLedgers, err error) {
	err = db.
		Where("id > ? and type = ?", fromID, string(ledger_type.LedgerTypeCash)).
		OrderBy("id").
		Limit(1000).
		Find(&shopLedgers)
	return
}
