// +build wireinject

package webserver

import (
	"github.com/google/wire"
	"o.o/backend/com/web/webserver/aggregate"
	"o.o/backend/com/web/webserver/query"
)

var WireSet = wire.NewSet(
	aggregate.New, aggregate.WebserverAggregateMessageBus,
	query.New, query.WebserverQueryServiceMessageBus,
)
