package model

import (
	"fmt"
	"net/url"
	"time"

	"o.o/capi/dot"
)

const (
	Limit  = "limit"
	Before = "before"
	After  = "after"
	Offset = "offset"
	Until  = "until"
	Since  = "since"
)

type FacebookPagingResponse struct {
	Cursors  *CursorPaginationResponse `json:"cursors"`
	Previous string                    `json:"previous"`
	Next     string                    `json:"next"`
}

func (f *FacebookPagingResponse) ToPagingRequestAfter(limit int) *FacebookPagingRequest {
	if f == nil {
		return &FacebookPagingRequest{}
	}
	return &FacebookPagingRequest{
		Limit: dot.Int(limit),
		CursorPagination: &CursorPaginationRequest{
			After: f.Cursors.After,
		},
	}
}

func (f *FacebookPagingResponse) ToPagingRequestBefore(limit int) *FacebookPagingRequest {
	if f == nil {
		return &FacebookPagingRequest{}
	}
	return &FacebookPagingRequest{
		Limit: dot.Int(limit),
		CursorPagination: &CursorPaginationRequest{
			Before: f.Cursors.Before,
		},
	}
}

func (res *FacebookPagingResponse) CompareFacebookPagingRequest(req *FacebookPagingRequest) bool {
	if res == nil && req == nil {
		return true
	}
	if res != nil && req == nil {
		return false
	}
	if res == nil && req != nil {
		return false
	}

	if req.CursorPagination == nil {
		return false
	}

	if (req.CursorPagination.After != "" && req.CursorPagination.After != res.Cursors.After) ||
		(req.CursorPagination.Before != "" && req.CursorPagination.Before != res.Cursors.Before) {
		return false
	}
	return true
}

type CursorPaginationResponse struct {
	Before string `json:"before"`
	After  string `json:"after"`
}

type FacebookPagingRequest struct {
	Limit            dot.NullInt
	CursorPagination *CursorPaginationRequest
	OffsetPagination *OffsetPaginationRequest
	TimePagination   *TimePaginationRequest
}

type CursorPaginationRequest struct {
	Before string
	After  string
}

type OffsetPaginationRequest struct {
	Offset int
}

type TimePaginationRequest struct {
	Until time.Time
	Since time.Time
}

func (p *FacebookPagingRequest) AddQueryParams(currentURL string, includeLimit bool, defaultPaging int) string {
	if p == nil {
		return currentURL
	}

	URL, err := url.Parse(currentURL)
	if err != nil {
		panic(err)
	}

	query, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		panic(err)
	}

	if p.Limit.Valid && includeLimit {
		if p.Limit.Int < 0 || p.Limit.Int > defaultPaging {
			p.Limit.Int = defaultPaging
		}

		query.Add(Limit, fmt.Sprintf("%d", p.Limit.Int))
	}

	countPagination := 0
	if p.CursorPagination != nil {
		if p.CursorPagination.After != "" && p.CursorPagination.Before != "" {
			panic("After and Before of CursorPagination couldn't have values at the same time.")
		}
		if p.CursorPagination.Before != "" {
			query.Add(Before, p.CursorPagination.Before)
		}
		if p.CursorPagination.After != "" {
			query.Add(After, p.CursorPagination.After)
		}
		countPagination += 1
	}
	if p.OffsetPagination != nil {
		query.Add(Offset, fmt.Sprintf("%d", p.OffsetPagination))
		countPagination += 1
	}
	if p.TimePagination != nil {
		// TODO: check
		if !p.TimePagination.Until.IsZero() && !p.TimePagination.Since.IsZero() {
			panic("Since and Until of TimePagination couldn't have values at the same time.")
		}
		if !p.TimePagination.Since.IsZero() {
			query.Add(Since, fmt.Sprintf("%d", p.TimePagination.Since.Unix()))
		}
		if !p.TimePagination.Until.IsZero() {
			query.Add(Until, fmt.Sprintf("%d", p.TimePagination.Until.Unix()))
		}
		countPagination += 1
	}

	if countPagination > 1 {
		panic("More than 1 pagination")
	}

	URL.RawQuery = query.Encode()
	return URL.String()
}
