// +build wireinject

package service

import "github.com/google/wire"

var WireSet = wire.NewSet(
	wire.Struct(new(MCShipnowService), "*"),
)
