package reporting

import (
	"context"
	"time"

	"o.o/api/top/types/etc/report_time_filter"
	"o.o/capi/dot"
)

// +gen:api

type QueryService interface {
	ReportOrders(context.Context, *ReportOrdersArgs) ([]*ReportOrder, error)
	ReportIncomeStatement(context.Context, *ReportIncomeStatementArgs) (map[int]*ReportIncomeStatement, error)
}

type ReportOrdersArgs struct {
	ShopID        dot.ID
	CreatedAtFrom time.Time
	CreatedAtTo   time.Time
	CreatedBy     dot.ID
}

type ReportIncomeStatementArgs struct {
	ShopID     dot.ID
	Year       int
	TimeFilter report_time_filter.TimeFilter
}
