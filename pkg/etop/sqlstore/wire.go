// +build wireinject

package sqlstore

import "github.com/google/wire"

var WireSet = wire.NewSet(
	New,
)
