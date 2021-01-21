// +build wireinject

package _publisher

import (
	"github.com/google/wire"
	"o.o/backend/com/fabo/pkg/fbclient"
	"o.o/backend/com/fabo/pkg/redis"
	"o.o/backend/com/fabo/pkg/sync"
	"o.o/backend/com/fabo/pkg/webhook"
)

var WireSet = wire.NewSet(
	webhook.NewWebhookHandler,
	fbclient.New,
	redis.NewFaboRedis,
	sync.New,
)
