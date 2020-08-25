// +build wireinject

package server_admin

import "github.com/google/wire"

var WireSet = wire.NewSet(
	BuildImportHandlers,
)
