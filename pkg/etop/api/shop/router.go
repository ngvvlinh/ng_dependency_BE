package shop

import (
	"o.o/backend/pkg/common/apifw/idemp"
	"o.o/backend/pkg/common/redis"
	"o.o/capi/httprpc"
	"o.o/common/l"
)

var ll = l.New()
var idempgroup *idemp.RedisGroup

type Servers []httprpc.Server

// workaround for imcsv
// TODO: remove this
var ProductServiceImpl *ProductService
var StocktakeServiceImpl *StocktakeService
var InventoryServiceImpl *InventoryService

func InitIdemp(rd redis.Store) {
	idempgroup = idemp.NewRedisGroup(rd, "idemp_shop", 30)
}
