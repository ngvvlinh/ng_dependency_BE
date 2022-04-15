// +build wireinject

package _base

import (
	"github.com/google/wire"

	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/authorization/auth"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/elasticsearch"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/authorize/tokens"
)

var WireSet = wire.NewSet(
	tokens.NewTokenStore,
	redis.Connect,
	elasticsearch.Connect,
	auth.NewGenerator,
	health.New,
	bus.New, // event bus
)
