// +build wireinject

package ticket

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	NewAggregate,
)
