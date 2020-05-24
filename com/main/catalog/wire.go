// +build wireinject

package catalog

import (
	"github.com/google/wire"
	"o.o/backend/com/main/catalog/aggregate"
	"o.o/backend/com/main/catalog/query"
)

var WireSet = wire.NewSet(
	aggregate.New, aggregate.AggregateMessageBus,
	query.New, query.QueryServiceMessageBus,
)
