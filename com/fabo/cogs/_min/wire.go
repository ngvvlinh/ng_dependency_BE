// +build wireinject

package _all

import (
	"github.com/google/wire"

	"o.o/backend/com/fabo/pkg/fbclient"
	"o.o/backend/com/fabo/pkg/redis"
	"o.o/backend/com/fabo/pkg/sync"
	"o.o/backend/com/fabo/pkg/webhook"
)

var WireSet = wire.NewSet(
	fbclient.New,
	redis.NewFaboRedis,
	sync.New,
)
