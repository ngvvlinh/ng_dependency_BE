// +build wireinject

package etelecom

import (
	"github.com/google/wire"

	"o.o/backend/com/summary/etelecom/query"
)

var WireSet = wire.NewSet(
	query.NewSummaryQuery, query.SummaryQueryMessageBus,
)
