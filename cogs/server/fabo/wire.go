// +build wireinject

package fabo

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	BuildFaboImageHandler,
)
