// +build wireinject

package reporting

import (
	"github.com/google/wire"
	"o.o/backend/com/main/reporting/query"
)

var WireSet = wire.NewSet(
	query.NewReportQuery, query.ReportQueryMessageBus,
)
