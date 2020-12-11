package summary

import (
	"context"
	"time"

	"o.o/api/summary"
	"o.o/capi/dot"
)

// +gen:api

type QueryService interface {
	SummaryShop(context.Context, *SummaryShopArgs) (*SummaryShopResponse, error)
}

type SummaryShopArgs struct {
	ShopID   dot.ID
	DateFrom time.Time
	DateTo   time.Time
}

type SummaryShopResponse struct {
	ListTable []*summary.SummaryTable `json:"list_table"`
}
