package admin

import (
	"o.o/backend/pkg/common/apifw/idemp"
	"o.o/backend/pkg/common/redis"
	"o.o/capi/httprpc"
	"o.o/common/l"
)

var ll = l.New()
var idempgroup *idemp.RedisGroup

type Servers []httprpc.Server

func InitIdemp(rd redis.Store) {
	idempgroup = idemp.NewRedisGroup(rd, "idemp_admin", 30)
}
