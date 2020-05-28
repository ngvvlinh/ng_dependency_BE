package manager

import "github.com/google/wire"

var WireSet = wire.NewSet(
	ManagerMesssageBus,
	NewManager,
)
