// +build wireinject

package fabo

import (
	"github.com/google/wire"

	"o.o/backend/com/fabo/main/fbmessaging"
	"o.o/backend/com/fabo/main/fbpage"
	"o.o/backend/com/fabo/main/fbuser"
	"o.o/backend/com/fabo/pkg/fbclient"
	"o.o/backend/com/fabo/pkg/redis"
	"o.o/backend/com/fabo/pkg/webhook"
)

var WireSet = wire.NewSet(
	webhook.New,
	fbclient.New,
	redis.NewFaboRedis,
	fbmessaging.WireSet,
	fbpage.WireSet,
	fbuser.WireSet,
)
