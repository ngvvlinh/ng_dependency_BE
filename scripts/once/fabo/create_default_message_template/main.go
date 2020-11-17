package main

import (
	"flag"
	"sync"
	"time"

	"o.o/backend/cmd/fabo-server/config"
	fbmessagetemplatemodel "o.o/backend/com/fabo/main/fbmessagetemplate/model"
	identifymodel "o.o/backend/com/main/identity/model"
	cm "o.o/backend/pkg/common"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	ll  = l.New()
	cfg config.Config
	db  *cmsql.Database
	mu  sync.Mutex
	wg  sync.WaitGroup
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

	createdCount := 0
	errCount := 0

	var fromID dot.ID
	for {
		shops, err := getShops(fromID)
		if err != nil {
			ll.Fatal(err.Error())
		}
		if len(shops) == 0 {
			break
		}
		for _, shop := range shops {
			wg.Add(1)
			messageTemplates := generateMessageTemplates(shop)
			go func(_messageTemplates fbmessagetemplatemodel.FbMessageTemplates) {
				defer func() {
					wg.Done()
				}()

				err := db.Table("fb_message_template").ShouldInsert(&_messageTemplates)
				mu.Lock()
				if err != nil {
					errCount += 1
				} else {
					createdCount += 1
				}
				mu.Unlock()

				return
			}(messageTemplates)
		}
		wg.Wait()
		fromID = shops[len(shops)-1].ID
		ll.S.Infof("create fb_message_template for shop: success %v/%v, error %v/%v", createdCount, totalShop, errCount, totalShop)
	}

}

func generateMessageTemplates(shop *identifymodel.Shop) fbmessagetemplatemodel.FbMessageTemplates {
	shopID := shop.ID
	now := time.Now()
	return fbmessagetemplatemodel.FbMessageTemplates{
		{
			ID:        cm.NewID(),
			ShopID:    shopID,
			Template:  "Chào bạn [Tên khách hàng], shop có thể giúp gì cho bạn?",
			ShortCode: "Xin chào",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        cm.NewID(),
			ShopID:    shopID,
			Template:  "Sản phẩm đã hết hàng, shop sẽ liên hệ bạn khi hàng về ạ!",
			ShortCode: "Hết hàng",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        cm.NewID(),
			ShopID:    shopID,
			Template:  "Dạ chào anh/chị! Shop sẽ kiểm tra lại hàng và báo anh/chị ngay!",
			ShortCode: "Kiểm tra hàng",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func getShops(fromID dot.ID) (shops identifymodel.Shops, err error) {
	err = db.
		Where("id > ?", fromID.Int64()).
		OrderBy("id").
		Limit(1000).
		Find(&shops)
	return
}
