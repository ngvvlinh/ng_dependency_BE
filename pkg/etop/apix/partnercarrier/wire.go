package partnercarrier

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewServers,
)
