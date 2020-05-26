// +build wireinject

package shipnow_carrier

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewManager,
)
