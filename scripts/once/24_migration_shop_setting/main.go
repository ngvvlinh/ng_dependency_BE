package main

import (
	"flag"
	"fmt"
	"sync"

	"go.uber.org/atomic"

	"o.o/backend/com/main/connectioning/model"
	settingmodel "o.o/backend/com/shopping/setting/model"
	cm "o.o/backend/pkg/common"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/scripts/once/21_migration_contact/config"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	ll     = l.New()
	cfg    config.Config
	DBMain *cmsql.Database
	wg     sync.WaitGroup
)

func main() {
	cc.InitFlags()
	flag.Parse()

	var err error
	if cfg, err = config.Load(); err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}

	if DBMain, err = cmsql.Connect(cfg.Postgres); err != nil {
		ll.Fatal("Error while loading new database", l.Error(err))
	}

	var fromID dot.ID
	var errorCount atomic.Int64
	totalCount := 0
	var wg sync.WaitGroup

	for {
		shopConnections, err := listShopConnections(fromID)
		ll.S.Info("shopConnections len :: ", len(shopConnections))
		if err != nil {
			ll.Fatal(fmt.Sprintf("%v", err))
		}
		if len(shopConnections) == 0 {
			break
		}
		for _, shopConnection := range shopConnections {
			wg.Add(1)
			go func(shopID dot.ID) {
				defer func() { wg.Done() }()
				_err := updateShopSetting(shopID)
				if _err != nil {
					errorCount.Inc()
					ll.S.Errorf("error %v", _err)
				}
			}(shopConnection.ShopID)
		}
		wg.Wait()
		fromID = shopConnections[len(shopConnections)-1].ShopID
		totalCount += len(shopConnections)
	}
	errCount := int(errorCount.Load())
	ll.S.Infof("Update shop setting result: ")
	ll.S.Infof("Error = %v / %v", errCount, totalCount)
	ll.S.Infof("Success = %v / %v", totalCount-errCount, totalCount)
}

func listShopConnections(fromID dot.ID) (sc model.ShopConnections, err error) {
	err = DBMain.Where("shop_id > ? AND shop_id IS NOT NULL AND deleted_at IS NULL", fromID.String()).
		OrderBy("shop_id").
		Limit(1000).
		Find(&sc)
	return
}

func updateShopSetting(shopID dot.ID) (err error) {
	update := &settingmodel.ShopSetting{
		ShopID:                     shopID,
		AllowConnectDirectShipment: true,
	}
	err = DBMain.
		Table("shop_setting").
		Where("shop_id = ?", shopID).
		ShouldUpdate(update)
	if err != nil && cm.ErrorCode(err) == cm.NotFound {
		return DBMain.Table("shop_setting").ShouldInsert(update)
	}

	return err
}
