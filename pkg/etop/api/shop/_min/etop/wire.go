// +build wireinject

package etop

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	NewServers,
)
