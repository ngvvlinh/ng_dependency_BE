package main

import (
	"context"
	"flag"

	"o.o/backend/cmd/etop-server/config"
	"o.o/backend/com/main/catalog/convert"
	catalogmodel "o.o/backend/com/main/catalog/model"
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
			variants, err := scanVariant(fromID)
			if err != nil {
				ll.Fatal("Error", l.Error(err))
			}
			if len(variants) == 0 {
				ll.S.Infof("Done: updated %v/%v", updated, count)
				break
			}
			fromID = variants[len(variants)-1].VariantID
			count += len(variants)
			for _, variant := range variants {

				ch <- variant.VariantID
				go func(p *catalogmodel.ShopVariant) (_err error) {
					_, ctxCancel := context.WithCancel(context.Background())
					defer func() {
						<-ch
						chInsert <- _err
						ctxCancel()
					}()

					attrNormKv := ""
					if p.Attributes != nil {
						_, attrNormKv = catalogmodel.NormalizeAttributes(convert.Convert_catalogmodel_ProductAttributes_catalogtypes_Attributes(p.Attributes))
					}
					update := make(map[string]interface{})
					if p.Attributes != nil && attrNormKv != "_" {
						update["attr_norm_kv"] = attrNormKv
					}

					if len(update) > 0 {
						_err = db.
							Table("shop_variant").
							Where("variant_id = ?", p.VariantID).
							ShouldUpdateMap(update)
					}
					return _err
				}(variant)
			}
			for i := 0; i < len(variants); i++ {
				err = <-chInsert
				if err != nil {
					errCount++
				} else {
					updated++
				}
			}
		}
		ll.S.Infof("Updated shop variant: success %v/%v, error %v/%v", updated, count, errCount, count)
	}
}

func scanVariant(fromID dot.ID) (variants catalogmodel.ShopVariants, err error) {
	err = db.
		Where("variant_id > ?", fromID.String()).
		OrderBy("variant_id").
		Limit(1000).
		Find(&variants)
	return
}
