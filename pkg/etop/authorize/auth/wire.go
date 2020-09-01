// +build wireinject

package auth

import "github.com/google/wire"

var WireSet = wire.NewSet(
	New,
)
