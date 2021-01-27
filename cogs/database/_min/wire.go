// +build wireinject

package database_min

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	wire.FieldsOf(new(Databases), "main", "log", "notifier", "webhook"),
	BuildDatabases,
)
