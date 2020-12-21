package summary

import (
	"context"
	"time"

	"o.o/api/summary"
	"o.o/capi/dot"
)

// +gen:api

type QueryService interface {
	Summary(context.Context, *SummaryArgs) (*SummaryResponse, error)
}

type SummaryArgs struct {
	ShopID   dot.ID
	DateFrom time.Time
	DateTo   time.Time
}

type SummaryResponse struct {
	ListTable []*summary.SummaryTable `json:"list_table"`
}
