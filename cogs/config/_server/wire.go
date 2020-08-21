// +build wireinject

package config_server

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	NewSession,
	WireSAdminToken,
)
