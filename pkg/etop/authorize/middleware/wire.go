// +build wireinject

package middleware

import "github.com/google/wire"

var WireSet = wire.NewSet(
	wire.Struct(new(SessionStarter), "*"),
)
