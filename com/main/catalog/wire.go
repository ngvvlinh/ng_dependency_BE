// +build wireinject

package catalog

import (
	"github.com/google/wire"

	"o.o/backend/com/main/catalog/aggregate"
	"o.o/backend/com/main/catalog/pm"
	"o.o/backend/com/main/catalog/query"
)

var WireSet = wire.NewSet(
	pm.New,
	aggregate.New, aggregate.AggregateMessageBus,
	query.New, query.QueryServiceMessageBus,
)
