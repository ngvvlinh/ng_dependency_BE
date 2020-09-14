// +build wireinject

package orderS

import "github.com/google/wire"

var WireSet = wire.NewSet(
	wire.Struct(new(OrderLogic), "*"),
)
