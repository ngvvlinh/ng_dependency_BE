// +build wireinject

package summary

import (
	"github.com/google/wire"

	"o.o/backend/com/summary/query"
)

var WireSet = wire.NewSet(
	query.NewDashboardQuery, query.DashboardQueryMessageBus,
)
