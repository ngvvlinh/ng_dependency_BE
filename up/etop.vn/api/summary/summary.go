package summary

import (
	"context"
	"time"

	"etop.vn/capi/dot"
)

// +gen:api

type QueryService interface {
	SummaryPOS(context.Context, *SummaryPOSRequest) (*SummaryPOSResponse, error)
	SummaryTopShip(context.Context, *SummaryTopShipRequest) (*SummaryTopShipResponse, error)
}

type SummaryPOSRequest struct {
	DateFrom time.Time
	DateTo   time.Time
	ShopID   dot.ID
}

type SummaryPOSResponse struct {
	ListTable []*SummaryTable
}

type SummaryTopShipRequest struct {
	DateFrom time.Time
	DateTo   time.Time
	ShopID   dot.ID
}

type SummaryTopShipResponse struct {
	ListTable []*SummaryTable
}

type SummaryTable struct {
	Label string
	Tags  []string
	Cols  []SummaryColRow
	Rows  []SummaryColRow
	Data  []SummaryItem
}

type SummaryColRow struct {
	Label  string
	Spec   string
	Unit   string
	Indent int
}

type SummaryItem struct {
	ImageUrls []string
	Label     string
	Spec      string
	Value     int
	Unit      string
}
