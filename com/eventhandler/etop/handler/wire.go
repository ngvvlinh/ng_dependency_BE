// +build wireinject

package handler

import "github.com/google/wire"

var WireSet = wire.NewSet(
	New,
)
