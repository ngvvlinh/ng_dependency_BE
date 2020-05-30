// +build wireinject

package suppliering

import (
	"github.com/google/wire"

	"o.o/backend/com/shopping/suppliering/aggregate"
	"o.o/backend/com/shopping/suppliering/query"
)

var WireSet = wire.NewSet(
	aggregate.NewSupplierAggregate, aggregate.SupplierAggregateMessageBus,
	query.NewSupplierQuery, query.SupplierQueryMessageBus,
)
