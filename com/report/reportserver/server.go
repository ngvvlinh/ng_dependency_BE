package reportserver

import (
	"net/http"

	"o.o/backend/com/report/middlewares"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/etop/authorize/session"
)

type ReportServer httpx.Server

func BuildReportServer(
	reportService ReportService,
	ss session.Session,
) ReportServer {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/report/income-statement", mware(ss, reportService.ReportIncomeStatement))
	mux.HandleFunc("/api/report/export/income-statement", mware(ss, reportService.ExportReportIncomeStatement))
	mux.HandleFunc("/api/report/orders/end-of-day", mware(ss, reportService.ReportOrdersEndOfDay))
	mux.HandleFunc("/api/report/export/orders/end-of-day", mware(ss, reportService.ExportReportOrdersEndOfDay))
	return httpx.MakeServer("/api/report/", mux)
}

func mware(
	ss session.Session,
	fn func(w http.ResponseWriter, r *http.Request),
) func(w http.ResponseWriter, r *http.Request) {
	return middlewares.Authorization(ss)(fn)
}
