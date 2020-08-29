package main

import (
	"flag"
	"fmt"

	"o.o/api/top/types/etc/status3"
	"o.o/backend/cmd/etop-server/config"
	"o.o/backend/com/main/address/model"
	shop "o.o/backend/com/main/identity/model"
	"o.o/backend/pkg/common/bus"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
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

	postgres := cfg.Databases.Postgres

	if db, err = cmsql.Connect(postgres); err != nil {
		ll.Fatal("Error while connecting database", l.Error(err))
	}

	{
		ll.S.Info("Migrate is_default from address")
		var fromShopID dot.ID
		count, updated, errors := 0, 0, 0
		for {
			shops, err := scanShopExistShipFromAddressID(fromShopID)
			if err != nil {
				ll.Fatal("Error", l.Error(err))
			}

			count += len(shops)

			if len(shops) == 0 {
				ll.S.Infof("Done: updated %v/%v", updated, count)
				break
			}

			var mapShopIDs = make(map[dot.ID]dot.ID)
			for _, item := range shops {
				if _, ok := mapShopIDs[item.ID]; !ok {
					mapShopIDs[item.ID] = item.ShipFromAddressID
				}
			}

			var addressIDs []dot.ID
			for _, addressID := range mapShopIDs {
				addressIDs = append(addressIDs, addressID)
			}

			fromShopID = shops[len(shops)-1].ID

			addresses, err := scanAddressWithoutDefault(addressIDs)
			if err != nil {
				ll.Fatal("Error", l.Error(err))
			}

			if len(addresses) == 0 {
				continue
			}

			var listAddresses []dot.ID
			for _, addr := range addresses {
				listAddresses = append(listAddresses, addr.ID)
			}

			numUpdated, numErrors := SetDefaultAddress(db, addresses)

			updated += numUpdated
			errors += numErrors
		}

		ll.S.Infof("Updated is_default %v/%v", updated, count)
		ll.S.Infof("Errors is_default %v", errors)
	}

	{
		ll.S.Info("Migrate ship_from_address_id is NULL from shop")
		var fromID dot.ID
		count, updated, errors := 0, 0, 0

		for {
			addresses, err := scanAddress(fromID)
			if err != nil {
				ll.Fatal("Error", l.Error(err))
			}

			if len(addresses) == 0 {
				ll.S.Infof("Done: updated %v/%v", updated, count)
				break
			}

			fromID = addresses[len(addresses)-1].ID

			var mapShopIDs = make(map[dot.ID]dot.ID)
			for _, item := range addresses {
				if _, ok := mapShopIDs[item.AccountID]; !ok {
					mapShopIDs[item.AccountID] = item.AccountID
				}
			}

			var shopIDs []dot.ID
			for _, shopID := range mapShopIDs {
				shopIDs = append(shopIDs, shopID)
			}

			shops, err := scanShop(shopIDs)
			if err != nil {
				ll.Fatal("Error", l.Error(err))
			}

			if len(shops) == 0 {
				continue
			}
			// get map shop info
			var mapShopInfo = make(map[dot.ID]dot.ID)
			for _, shop := range shops {
				if _, ok := mapShopInfo[shop.ID]; !ok {
					mapShopInfo[shop.ID] = shop.ID
				}
			}

			count += len(shopIDs)

			numUpdated, numErrors := UpdateDefaultAddress(db, addresses, mapShopInfo)

			updated += numUpdated
			errors += numErrors
		}
		ll.S.Infof("Updated %v/%v", updated, count)
		ll.S.Infof("Errors %v", errors)
	}
}

func scanShopExistShipFromAddressID(fromID dot.ID) (shops shop.Shops, err error) {
	err = db.
		Where("id > ? AND ship_from_address_id IS NOT NULL", fromID.Int64()).
		Where("deleted_at IS NULL").
		OrderBy("id").
		Limit(500).
		Find(&shops)

	return
}

func scanAddressWithoutDefault(addressIDs []dot.ID) (addresses model.Addresses, err error) {
	err = db.
		Where("type = ?", "shipfrom").
		Where("is_default = ?", false).
		Where(sq.In("id", addressIDs)).
		OrderBy("id").
		Find(&addresses)
	return
}

func scanShop(shopIDs []dot.ID) (shops shop.Shops, err error) {
	err = db.
		Where("ship_from_address_id IS NULL").
		Where("status = ?", status3.P).
		Where("deleted_at IS NULL").
		Where(sq.In("id", shopIDs)).
		OrderBy("id").
		Find(&shops)

	return
}

var sqlScanAddress = fmt.Sprintf(`
SELECT %v FROM (
	SELECT DISTINCT ON (account_id) * FROM address
	WHERE (id > ? AND type = ? AND is_default = ?)
	ORDER BY account_id ASC, created_at ASC LIMIT 500
) t ORDER BY created_at ASC
`, (*model.Address)(nil).SQLListCols())

func scanAddress(fromID dot.ID) (addresses model.Addresses, err error) {
	rows, err := db.SQL(sqlScanAddress, fromID.Int64(), "shipfrom", false).Query()
	if err != nil {
		return nil, err
	}
	err = addresses.SQLScan(db.Opts(), rows)
	return addresses, err
}

func SetDefaultAddress(db *cmsql.Database, addresses model.Addresses) (updated, errors int) {
	chErr := make(chan error, len(addresses))
	for _, address := range addresses {
		go func(_addr *model.Address) (_err error) {
			defer func() {
				chErr <- _err
			}()

			_err = db.
				Table("address").
				Where("id = ?", _addr.ID).
				ShouldUpdateMap(map[string]interface{}{
					"is_default": true,
				})
			if _err != nil {
				ll.Debug("err", l.Error(_err))
			}
			return
		}(address)
	}
	for i, n := 0, len(addresses); i < n; i++ {
		err := <-chErr
		if err == nil {
			updated++
		} else {
			errors++
		}
	}
	ll.S.Infof("update default addresses :: created %v/%v, errored %v/%v",
		updated, len(addresses),
		errors, len(addresses))
	return
}

func UpdateDefaultAddress(db *cmsql.Database, addresses model.Addresses, mapShopInfo map[dot.ID]dot.ID) (updated, errors int) {
	for _, address := range addresses {
		if _, ok := mapShopInfo[address.AccountID]; !ok {
			continue
		}

		if err := db.InTransaction(bus.Ctx(), func(tx cmsql.QueryInterface) error {
			if _, _err := tx.Table("shop").Where("ship_from_address_id IS NULL").Where("id = ?", address.AccountID).UpdateMap(map[string]interface{}{
				"ship_from_address_id": address.ID,
			}); _err != nil {
				ll.Debug("err update ship_from_address_id", l.Error(_err))
				return _err
			}

			if _, _err := tx.Table("address").Where("id = ?", address.ID).UpdateMap(map[string]interface{}{
				"is_default": true,
			}); _err != nil {
				ll.Debug("err set is_default", l.Error(_err))
				return _err
			}
			return nil
		}); err != nil {
			errors++
			ll.Debug("err", l.Error(err))
		} else {
			updated++
		}
	}

	ll.S.Infof("update default addresses :: created %v/%v, errored %v/%v",
		updated, len(addresses),
		errors, len(addresses))
	return
}
