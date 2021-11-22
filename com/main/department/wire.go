package department

import (
	"github.com/google/wire"
	"o.o/backend/com/main/department/aggregate"
	"o.o/backend/com/main/department/query"
)

var WireSet = wire.NewSet(
	aggregate.NewDepartmentAggregate, aggregate.DepartmentAggregateMessageBus,
	query.NewDepartmentQuery, query.DepartmentQueryMessageBus,
)
