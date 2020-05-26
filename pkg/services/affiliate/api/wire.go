// +build wireinject

package api

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	wire.Struct(new(UserService), "*"),
	wire.Struct(new(TradingService), "*"),
	wire.Struct(new(ShopService), "*"),
	wire.Struct(new(AffiliateService), "*"),
	NewServers,
)
