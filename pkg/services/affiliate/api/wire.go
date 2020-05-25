// +build wireinject

package api

import (
	"github.com/google/wire"
	"o.o/api/main/catalog"
	"o.o/api/main/identity"
	"o.o/api/services/affiliate"
)

var WireSet = wire.NewSet(
	wire.Struct(new(UserService)),
	wire.Struct(new(TradingService)),
	wire.Struct(new(ShopService)),
	wire.Struct(new(AffiliateService)),
	NewServers,
)

func BuildServers(
	secret Secret,
	affAggr affiliate.CommandBus,
	affQuery affiliate.QueryBus,
	catQuery catalog.QueryBus,
	idenQuery identity.QueryBus,
) Servers {
	panic(wire.Build(WireSet))
}
