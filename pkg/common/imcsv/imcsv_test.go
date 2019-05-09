package imcsv

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var schema1 = Schema{
	{
		Name:    "col0",
		Display: "Column 0",
		Norm:    "column 0",
	},
	{
		Name:    "col1",
		Display: "Column 1",
		Norm:    "column 1",
		Line:    true,
	},
	{
		Name:     "col2",
		Display:  "Column 2",
		Norm:     "column 2",
		Optional: true,
	},
	{
		Name:    "col3",
		Display: "Column 3",
		Norm:    "column 3",
	},
}

func TestValidateSchema(t *testing.T) {
	Convey("Optional", t, func() {
		Convey("All columns", func() {
			header := strings.Split("Column 0,Column 1,Column 2,Column 3", ",")
			idx, errs, err := schema1.ValidateSchema(&header)
			So(err, ShouldBeNil)
			So(errs, ShouldBeEmpty)
			So(idx.mapIndex, ShouldResemble, []int{0, 1, 2, 3})
		})
		Convey("Missing optional columns", func() {
			header := strings.Split("Column 0,Column 1,Column 3", ",")
			idx, errs, err := schema1.ValidateSchema(&header)
			So(err, ShouldBeNil)
			So(errs, ShouldBeEmpty)
			So(idx.mapIndex, ShouldResemble, []int{0, 1, -1, 2})
		})
	})
}
