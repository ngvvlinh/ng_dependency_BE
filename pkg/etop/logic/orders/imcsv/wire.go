package imcsv

import "github.com/google/wire"

var WireSet = wire.NewSet(
	New,
)
