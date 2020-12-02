// +build wireinject

package provider

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewTelecomManager,
)
