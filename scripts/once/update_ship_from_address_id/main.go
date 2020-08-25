package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

	"o.o/backend/cmd/etop-server/config"
	"o.o/backend/com/main/address/model"
	shop "o.o/backend/com/main/identity/model"
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
		count, updated := 0, 0
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

			numUpdated, _ := SetDefaultAddress(db, addresses)

			updated += numUpdated
		}
	}

	{
		ll.S.Info("Migrate ship_from_address_id is NULL from shop")
		var fromID dot.ID
		count, updated := 0, 0

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

			fromID = addresses[len(addresses)-1].ID

			count += len(shopIDs)

			numUpdated, _ := UpdateDefaultAddress(db, addresses, mapShopInfo)

			updated += numUpdated
		}
		ll.S.Infof("Updated %v/%v", updated, count)
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
		Where("deleted_at IS NULL").
		Where(sq.In("id", shopIDs)).
		OrderBy("id").
		Find(&shops)

	return
}

func scanAddress(fromID dot.ID) (addresses model.Addresses, err error) {
	var res []*model.Address
	sql, args, err := db.SQL(`SELECT DISTINCT ON (account_id) * FROM "address"`).
		Where("id > ? AND type = ?", fromID.Int64(), "shipfrom").
		Limit(500).Build()

	if err != nil {
		return nil, err
	}

	sql2 := fmt.Sprintf(
		"SELECT * FROM (%v) AS s ORDER BY created_at ASC",
		sql,
	)
	rows, err := db.Query(sql2, args...)
	if err != nil {
		return nil, err
	}

	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtr := make([]interface{}, count)

	for rows.Next() {
		for i, _ := range columns {
			valuePtr[i] = &values[i]
		}
		var temp = make(map[string]interface{})
		rows.Scan(valuePtr...)

		// get only id and account_id
		for index, col := range columns {
			if col != "id" && col != "account_id" {
				continue
			}
			temp[col] = values[index]
		}

		address := model.Address{}
		dataStr, _ := json.Marshal(temp)
		err = json.Unmarshal(dataStr, &address)
		if err != nil {
			return nil, err
		}

		res = append(res, &address)
	}
	return res, err
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
		if _err := db.InTransaction(context.TODO(), func(tx cmsql.QueryInterface) error {
			if _, _err := db.Exec(`UPDATE shop SET ship_from_address_id = $1 WHERE ship_from_address_id IS NULL AND id = $2`, address.ID, address.AccountID); _err != nil {
				return _err
			}
			if _, _err := db.Exec(`UPDATE address SET is_default = true WHERE id = $1`, address.ID); _err != nil {
				return _err
			}
			return nil
		}); _err != nil {
			return
		}
	}

	ll.S.Infof("update default addresses :: created %v/%v, errored %v/%v",
		updated, len(addresses),
		errors, len(addresses))
	return
}
