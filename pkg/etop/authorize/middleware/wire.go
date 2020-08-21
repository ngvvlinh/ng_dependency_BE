// +build wireinject

package middleware

import "github.com/google/wire"

var WireSet = wire.NewSet(
	New,
)
