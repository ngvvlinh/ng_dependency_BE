package main

import (
	"context"
	"flag"
	"strconv"

	"o.o/backend/cmd/etop-server/config"
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
		var mapshopCode = make(map[string]int)

		for {
			variants, err := scanVariantDuplicate(fromID)
			if err != nil {
				ll.Fatal("Error", l.Error(err))
			}
			if len(variants) == 0 {
				ll.S.Infof("Done: updated %v/%v", updated, count)
				break
			}
			fromID = variants[len(variants)-1].VariantID

			variant_ids := []dot.ID{}
			for _, variant := range variants {
				defaultCode := variant.Code
				key := variant.ShopID.String() + variant.Code
				mapshopCode[key]++
				for mapshopCode[key] > 1 {
					variant.Code = variant.Code + strconv.Itoa(mapshopCode[key])
					key = variant.ShopID.String() + variant.Code
					mapshopCode[key]++
				}
				if defaultCode == variant.Code {
					continue
				}
				count++
				variant_ids = append(variant_ids, variant.VariantID)

				ch <- variant.VariantID
				go func(variantID dot.ID, code string) (_err error) {
					_, ctxCancel := context.WithCancel(context.Background())
					defer func() {
						<-ch
						chInsert <- _err
						ctxCancel()
					}()
					update := make(map[string]interface{})
					update["code"] = code
					if len(update) > 0 {
						_err = db.
							Table("shop_variant").
							Where("variant_id = ?", variantID).
							ShouldUpdateMap(update)
					}
					return _err
				}(variant.VariantID, variant.Code)
			}
			for i := 0; i < len(variant_ids); i++ {
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

func scanVariantDuplicate(fromID dot.ID) (variants catalogmodel.ShopVariants, err error) {
	err = db.
		Where(`variant_id > ? and shop_id in (
			select
				shop_id
			from
				shop_variant
			where
				code is not null 
			group by
				code,
				shop_id
			having
				count(*) > 1
			order by
				count(*) desc
		) `, fromID.String()).
		OrderBy("variant_id").
		Limit(1000).
		Find(&variants)
	return
}
