// +build wireinject

package etop

import (
	"github.com/google/wire"

	"o.o/backend/com/summary/etop/query"
)

var WireSet = wire.NewSet(
	query.NewDashboardQuery, query.DashboardQueryMessageBus,
)
