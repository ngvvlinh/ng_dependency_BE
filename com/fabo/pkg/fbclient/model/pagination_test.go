package model

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"o.o/capi/dot"
)

type queryParams struct {
	Limit  string
	Before string
	After  string
	Offset string
	Since  string
	Until  string
}

func TestApplyQueryParams(t *testing.T) {
	Convey("test applyQueryParams", t, func() {
		fbPagingRequest := FacebookPagingRequest{
			Limit: dot.Int(10),
			CursorPagination: &CursorPaginationRequest{
				Before: "12345",
			},
		}

		var params queryParams
		fbPagingRequest.ApplyQueryParams(true, 100, &params)

		So(params.Limit, ShouldEqual, "10")
		So(params.Before, ShouldEqual, "12345")
	})
}
