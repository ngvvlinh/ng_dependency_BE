package summary

import (
	"context"
	"time"

	"o.o/capi/dot"
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
	ListTable []*SummaryTable `json:"list_table"`
}

type SummaryTopShipRequest struct {
	DateFrom time.Time
	DateTo   time.Time
	ShopID   dot.ID
}

type SummaryTopShipResponse struct {
	ListTable []*SummaryTable `json:"list_table"`
}

type SummaryTable struct {
	Label string          `json:"label"`
	Tags  []string        `json:"tags"`
	Cols  []SummaryColRow `json:"cols"`
	Rows  []SummaryColRow `json:"rows"`
	Data  []SummaryItem   `json:"data"`
}

type SummaryColRow struct {
	Label  string `json:"label"`
	Spec   string `json:"spec"`
	Unit   string `json:"unit"`
	Indent int    `json:"indent"`
}

type SummaryItem struct {
	ImageUrls []string `json:"image_urls"`
	Label     string   `json:"label"`
	Spec      string   `json:"spec"`
	Value     int      `json:"value"`
	Unit      string   `json:"unit"`
}
