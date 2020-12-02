// +build wireinject

package database_all

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	wire.FieldsOf(new(Databases), "main", "log", "notifier", "affiliate", "webserver", "etelecom"),
	BuildDatabases,
)
