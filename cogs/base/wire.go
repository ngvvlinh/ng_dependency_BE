package _base

import (
	"github.com/google/wire"

	"o.o/backend/pkg/common/authorization/auth"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/authorize/tokens"
)

var WireSet = wire.NewSet(
	tokens.NewTokenStore,
	redis.Connect,
	auth.NewGenerator,
)
