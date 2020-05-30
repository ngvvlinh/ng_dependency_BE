package webserver

import (
	"github.com/google/wire"

	"o.o/backend/com/web/webserver/aggregate"
	"o.o/backend/com/web/webserver/query"
	"o.o/backend/pkg/common/sql/cmsql"
)

type WebServerDB *cmsql.Database

var WireSet = wire.NewSet(
	aggregate.New, aggregate.WebserverAggregateMessageBus,
	query.New, query.WebserverQueryServiceMessageBus,
)
