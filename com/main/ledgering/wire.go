// +build wireinject

package ledgering

import (
	"github.com/google/wire"

	"o.o/backend/com/main/ledgering/aggregate"
	"o.o/backend/com/main/ledgering/query"
)

var WireSet = wire.NewSet(
	aggregate.NewLedgerAggregate, aggregate.LedgerAggregateMessageBus,
	query.NewLedgerQuery, query.LedgerQueryMessageBus,
)
