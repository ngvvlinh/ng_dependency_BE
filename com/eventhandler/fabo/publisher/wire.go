// +build wireinject

package publisher

import "github.com/google/wire"

var WireSet = wire.NewSet(
	New,
)
