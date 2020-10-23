package reporting

import (
	"context"
	"time"

	"o.o/capi/dot"
)

// +gen:api

type QueryService interface {
	ReportOrders(context.Context, *ReportOrdersArgs) ([]*ReportOrder, error)
}

type ReportOrdersArgs struct {
	ShopID        dot.ID
	CreatedAtFrom time.Time
	CreatedAtTo   time.Time
	CreatedBy     dot.ID
}
