package shop

import (
	"o.o/backend/pkg/common/apifw/idemp"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/api/shop/inventory"
	"o.o/backend/pkg/etop/api/shop/product"
	"o.o/backend/pkg/etop/api/shop/stocktake"
	"o.o/capi/httprpc"
	"o.o/common/l"
)

var ll = l.New()
var Idempgroup *idemp.RedisGroup

type Servers []httprpc.Server

// workaround for imcsv
// TODO: remove this
var ProductServiceImpl *product.ProductService
var StocktakeServiceImpl *stocktake.StocktakeService
var InventoryServiceImpl *inventory.InventoryService

func InitIdemp(rd redis.Store) {
	Idempgroup = idemp.NewRedisGroup(rd, "idemp_shop", 30)
}
