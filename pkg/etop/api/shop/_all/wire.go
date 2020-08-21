// +build wireinject

package shop_all

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	NewServers,
)
