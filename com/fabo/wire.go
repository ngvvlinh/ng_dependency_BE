// +build wireinject

package fabo

import (
	"github.com/google/wire"

	"o.o/backend/com/fabo/pkg/fbclient"
	"o.o/backend/com/fabo/pkg/redis"
	"o.o/backend/com/fabo/pkg/sync"
	"o.o/backend/com/fabo/pkg/webhook"
)

var WireSet = wire.NewSet(
	webhook.New,
	fbclient.New,
	redis.NewFaboRedis,
	sync.New,
)
