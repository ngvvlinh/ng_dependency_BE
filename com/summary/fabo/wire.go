// +build wireinject

package fabo

import (
	"github.com/google/wire"
	"o.o/backend/com/summary/fabo/query"
)

var WireSet = wire.NewSet(
	query.NewDashboardQuery,
	query.DashboardQueryMessageBus,
)