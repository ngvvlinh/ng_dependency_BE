// +build wireinject

package server_max

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	BuildIntHandlers,
	BuildExtHandlers,
)
