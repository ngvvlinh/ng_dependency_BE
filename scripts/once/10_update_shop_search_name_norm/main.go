package main

import (
	"context"
	"flag"

	"o.o/backend/cmd/etop-server/config"
	identitymodel "o.o/backend/com/main/identity/model"
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
			shops, err := scanShop(fromID)
			if err != nil {
				ll.Fatal("Error", l.Error(err))
			}
			if len(shops) == 0 {
				ll.S.Infof("Done: updated %v/%v", updated, count)
				break
			}
			shopSearchs, err := scanShopSearch(shops[0].ID, shops[len(shops)-1].ID)
			if err != nil {
				ll.Fatal("Error", l.Error(err))
			}
			var shopSearchMap = make(map[dot.ID]*identitymodel.ShopSearch)
			for _, v := range shopSearchs {
				shopSearchMap[v.ID] = v
			}
			fromID = shops[len(shops)-1].ID
			count += len(shops)
			for _, shop := range shops {
				ch <- shop.ID
				go func(p *identitymodel.Shop, m map[dot.ID]*identitymodel.ShopSearch) (_err error) {
					_, ctxCancel := context.WithCancel(context.Background())
					defer func() {
						<-ch
						chInsert <- _err
						ctxCancel()
					}()

					nameNorm := validate.NormalizeSearchCharacter(p.Name)
					if p.Code != "" {
						nameNorm += " " + validate.NormalizeSearchCode(p.Code)
					}
					var shopSearch *identitymodel.ShopSearch
					if p.Name != "" {
						shopSearch = &identitymodel.ShopSearch{
							ID:       p.ID,
							Name:     p.Name,
							NameNorm: nameNorm,
						}
					}
					if m[p.ID] != nil {
						// update
						update := make(map[string]interface{})
						if p.Name != "" {
							update["name_norm"] = nameNorm
						}

						if len(update) > 0 {
							_err = db.
								Table("shop_search").
								Where("id = ?", p.ID).
								ShouldUpdateMap(update)
						}
					} else {
						if shopSearch != nil {
							_, _err = db.
								Table("shop_search").
								Insert(shopSearch)
						}
					}
					return _err
				}(shop, shopSearchMap)
			}
			for i := 0; i < len(shops); i++ {
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

func scanShop(fromID dot.ID) (shops identitymodel.Shops, err error) {
	err = db.
		Where(`id > ? `, fromID.String()).
		OrderBy("id").
		Limit(1000).
		Find(&shops)
	return
}

func scanShopSearch(fromID dot.ID, toID dot.ID) (shopSearchs identitymodel.ShopSearchs, err error) {
	err = db.
		Where(`id >= ? and id <=  ?`, fromID.String(), toID.String()).
		OrderBy("id").
		Find(&shopSearchs)
	return
}
