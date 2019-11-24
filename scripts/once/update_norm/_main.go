package main

import (
	"flag"

	"etop.vn/backend/cmd/etop-server/config"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
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
	cm.SetEnvironment(cfg.Env)

	if db, err = cmsql.Connect(cfg.Postgres); err != nil {
		ll.Fatal("Error while connecting database", l.Error(err))
	}

	{
		fromID := int64(0)
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
			fromID = products[len(products)-1].ID
			count += len(products)

			for _, p := range products {
				nameNorm := validate.NormalizeSearch(p.Name)
				nameNormUa := validate.NormalizeUnaccent(p.Name)
				if split(nameNorm) != p.NameNorm || nameNormUa != p.NameNormUa {
					err = db.
						Table("product").
						Where("id = ?", p.ID).
						ShouldUpdateMap(M{
							"name_norm":    nameNorm,
							"name_norm_ua": nameNormUa,
						})
					if err != nil {
						ll.S.Fatalf("UpdateInfo product id=%v", p.ID)
					}
					updated++
				}
			}
			ll.S.Infof("Updated %v/%v", updated, count)
		}
	}
	{
		fromShopID, fromProductID := int64(0), int64(0)
		count, updated := 0, 0
		for {
			products, err := scanShopProduct(fromShopID, fromProductID)
			if err != nil {
				ll.Fatal("Error", l.Error(err))
			}
			if len(products) == 0 {
				ll.S.Infof("Done: updated %v/%v", updated, count)
				break
			}
			fromShopID = products[len(products)-1].ShopID
			fromProductID = products[len(products)-1].ProductID
			count += len(products)

			for _, p := range products {
				nameNorm := validate.NormalizeSearch(p.Name)
				if split(nameNorm) != p.NameNorm {
					err = db.
						Table("shop_product").
						Where("shop_id = ?", p.ShopID).
						Where("product_id = ?", p.ProductID).
						ShouldUpdateMap(M{
							"name_norm": nameNorm,
						})
					if err != nil {
						ll.S.Fatalf("UpdateInfo shop product %v-%v", p.ShopID, p.ProductID)
					}
					updated++
				}
			}
			ll.S.Infof("Updated %v/%v", updated, count)
		}
	}
	{
		fromID := int64(0)
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
			fromID = variants[len(variants)-1].ID
			count += len(variants)

			for _, v := range variants {
				attrsNorm, normKv := model.NormalizeAttributes(v.Attributes)
				if split(normKv) != v.AttrNormKv {
					update := M{
						"attr_norm_kv": normKv,
					}
					if len(attrsNorm) == 0 {
						update["attributes"] = nil
					} else {
						update["attributes"] = string(cm.ToJSON(attrsNorm))
					}
					err = db.
						Table("variant").
						Where("id = ?", v.ID).
						ShouldUpdateMap(update)
					if err != nil {
						ll.S.Fatalf("UpdateInfo variant id=%v", v.ID)
					}
					updated++
				}
			}
			ll.S.Infof("Updated %v/%v", updated, count)
		}
	}
}

func scanProduct(fromID dot.ID) (products model.Products, err error) {
	err = db.
		Where("id > ?", fromID).
		OrderBy("id").
		Limit(1000).
		Find(&products)
	return
}

func scanVariant(fromID dot.ID) (variants model.Variants, err error) {
	err = db.
		Where("id > ?", fromID).
		OrderBy("id").
		Limit(1000).
		Find(&variants)
	return
}

func scanShopProduct(shopID, productID dot.ID) (products model.ShopProducts, err error) {
	err = db.
		Where("(shop_id, product_id) > (?,?)", shopID, productID).
		OrderBy("shop_id, product_id").
		Limit(1000).
		Find(&products)
	return
}

func split(s string) string {
	return validate.NormalizedSearchToTsVector(s)
}
