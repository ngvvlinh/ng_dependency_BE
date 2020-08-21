// +build wireinject

package shop_min

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	NewServers,
)
