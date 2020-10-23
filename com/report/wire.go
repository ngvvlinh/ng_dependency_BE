package report

import (
	"github.com/google/wire"
	"o.o/backend/com/report/reportserver"
)

var WireSet = wire.NewSet(
	reportserver.BuildReportServer,
	wire.Struct(new(reportserver.ReportService), "*"),
)
