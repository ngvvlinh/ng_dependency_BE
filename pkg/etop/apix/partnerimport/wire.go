package partnerimport

import "github.com/google/wire"

var WireSet = wire.NewSet(
	New,
	NewServers,
)
