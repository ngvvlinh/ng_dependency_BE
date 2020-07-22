package main

import (
	"context"
	"flag"

	"o.o/backend/cmd/etop-server/config"
	catalogmodel "o.o/backend/com/main/catalog/model"
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
			products, err := scanProduct(fromID)
			if err != nil {
				ll.Fatal("Error", l.Error(err))
			}
			if len(products) == 0 {
				ll.S.Infof("Done: updated %v/%v", updated, count)
				break
			}
			fromID = products[len(products)-1].ProductID
			count += len(products)
			for _, product := range products {

				ch <- product.ProductID
				go func(p *catalogmodel.ShopProduct) (_err error) {
					_, ctxCancel := context.WithCancel(context.Background())
					defer func() {
						<-ch
						chInsert <- _err
						ctxCancel()
					}()

					nameNorm := validate.NormalizeSearchCharacter(p.Name)
					if p.Code != "" {
						nameNorm = validate.NormalizeSearchCharacter(p.Name+" "+p.Code) + " " + validate.NormalizeSearchCode(p.Code)
					}
					update := make(map[string]interface{})
					if p.NameNorm != "" && p.Name != "" {
						update["name_norm"] = nameNorm
					}

					if len(update) > 0 {
						_err = db.
							Table("shop_product").
							Where("product_id = ?", p.ProductID).
							ShouldUpdateMap(update)
					}
					return _err
				}(product)
			}
			for i := 0; i < len(products); i++ {
				err = <-chInsert
				if err != nil {
					errCount++
				} else {
					updated++
				}
			}
		}
		ll.S.Infof("Updated shop product: success %v/%v, error %v/%v", updated, count, errCount, count)
	}
}

func scanProduct(fromID dot.ID) (products catalogmodel.ShopProducts, err error) {
	err = db.
		Where("product_id > ?", fromID.String()).
		OrderBy("product_id").
		Limit(1000).
		Find(&products)
	return
}
