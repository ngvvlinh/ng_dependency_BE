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
	cmenv.SetEnvironment(cfg.SharedConfig.Env)

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
			addresses, err := scanShopTraderAddress(fromID)
			if err != nil {
				ll.Fatal("Error", l.Error(err))
			}
			if len(addresses) == 0 {
				ll.S.Infof("Done: updated %v/%v", updated, count)
				break
			}
			shopSearchs, err := scanShopTraderAddressSearch(addresses[0].ID, addresses[len(addresses)-1].ID)
			if err != nil {
				ll.Fatal("Error", l.Error(err))
			}
			var AddressesSearchMap = make(map[dot.ID]*model.ShopTraderAddressSearch)
			for _, v := range shopSearchs {
				AddressesSearchMap[v.ID] = v
			}
			fromID = addresses[len(addresses)-1].ID
			count += len(addresses)
			for _, shop := range addresses {
				ch <- shop.ID
				go func(p *model.ShopTraderAddress, m map[dot.ID]*model.ShopTraderAddressSearch) (_err error) {
					_, ctxCancel := context.WithCancel(context.Background())
					defer func() {
						<-ch
						chInsert <- _err
						ctxCancel()
					}()

					phone, ok := validate.NormalizePhone(p.Phone)
					if !ok {
						return
					}
					phoneNorm := validate.NormalizeSearchPhone(phone.String())
					if m[p.ID] != nil {
						update := make(map[string]interface{})
						if phone != "" && m[p.ID].PhoneNorm == "" {
							update["phone_norm"] = phoneNorm
						}
						if len(update) > 0 {
							_err = db.
								Table("shop_trader_address_search").
								Where("id = ?", p.ID).
								ShouldUpdateMap(update)
						}
					} else {
						var shopSearch *model.ShopTraderAddressSearch
						shopSearch = &model.ShopTraderAddressSearch{
							ID:        p.ID,
							PhoneNorm: phoneNorm,
						}
						_, _err = db.
							Table("shop_trader_address_search").
							Insert(shopSearch)
					}
					return _err
				}(shop, AddressesSearchMap)
			}
			for i := 0; i < len(addresses); i++ {
				err = <-chInsert
				if err != nil {
					errCount++
				} else {
					updated++
				}
			}
		}
		ll.S.Infof("Updated shop shop: success %v/%v, error %v/%v", updated, count, errCount, count)
	}
}

func scanShopTraderAddress(fromID dot.ID) (shops model.ShopTraderAddresses, err error) {
	err = db.
		Where(`id > ? `, fromID.String()).
		OrderBy("id").
		Limit(1000).
		Find(&shops)
	return
}

func scanShopTraderAddressSearch(fromID dot.ID, toID dot.ID) (shopSearchs model.ShopTraderAddressSearchs, err error) {
	err = db.
		Where(`id >= ? and id <=  ?`, fromID.String(), toID.String()).
		OrderBy("id").
		Find(&shopSearchs)
	return
}
