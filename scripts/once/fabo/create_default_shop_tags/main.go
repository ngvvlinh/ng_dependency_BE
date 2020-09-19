package main

import (
	"flag"
	"fmt"

	"o.o/backend/cmd/fabo-server/config"
	fabouseringmodel "o.o/backend/com/fabo/main/fbuser/model"
	identifymodel "o.o/backend/com/main/identity/model"
	cm "o.o/backend/pkg/common"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	ll         = l.New()
	cfg        config.Config
	db         *cmsql.Database
	mapShopTag = map[dot.ID]struct{}{}
)

func main() {
	cc.InitFlags()
	flag.Parse()

	var err error
	if cfg, err = config.Load(); err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}

	postgres := cfg.Databases.Postgres
	if db, err = cmsql.Connect(postgres); err != nil {
		ll.Fatal("Error while connecting database", l.Error(err))
	}

	totalShop, _ := db.Table("shop").Count(&identifymodel.Shop{})
	if err != nil {
		ll.Fatal(err.Error())
	}

	initMap()
	count := 0
	var fromID dot.ID
	for {
		shops, err := getShops(fromID)
		if err != nil {
			ll.Fatal(err.Error())
		}
		if len(shops) == 0 {
			break
		}
		ll.Info("fetch 1000")
		for _, shop := range shops {
			err = createTags(shop)
			if err != nil {
				ll.Error(err.Error())
			}
			fromID = shop.ID
			ll.Info(fmt.Sprintf("create tags for shop %d/%d", count, totalShop))
			count++
		}
	}

}

func createTags(shop *identifymodel.Shop) error {
	shopID := shop.ID
	if _, ok := mapShopTag[shopID]; ok {
		return nil
	}
	tags := fabouseringmodel.FbShopUserTags{
		{
			ID:     cm.NewID(),
			Name:   "Chốt Đơn",
			Color:  "#3498db",
			ShopID: shopID,
		},
		{
			ID:     cm.NewID(),
			Name:   "Đã Ship",
			Color:  "#2ecc71",
			ShopID: shopID,
		},
		{
			ID:     cm.NewID(),
			Name:   "Hỏi Giá",
			Color:  "#95a5a6",
			ShopID: shopID,
		},
		{
			ID:     cm.NewID(),
			Name:   "Tư Vấn",
			Color:  "#e74c3c",
			ShopID: shopID,
		},
		{
			ID:     cm.NewID(),
			Name:   "Bank",
			Color:  "#9b59b6",
			ShopID: shopID,
		},
		{
			ID:     cm.NewID(),
			Name:   "COD",
			Color:  "#f39c12",
			ShopID: shopID,
		},
	}

	return db.Table("fb_shop_user_tag").ShouldInsert(&tags)
}

func getShops(fromID dot.ID) (shops identifymodel.Shops, err error) {
	err = db.
		Where("id > ?", fromID.Int64()).
		OrderBy("id").
		Limit(1000).
		Find(&shops)
	return
}

func initMap() {
	var tags fabouseringmodel.FbShopUserTags
	if err := db.Table("fb_shop_user_tag").Find(&tags); err != nil {
		ll.Fatal(err.Error())
	}

	for _, tag := range tags {
		mapShopTag[tag.ShopID] = struct{}{}
	}
}
