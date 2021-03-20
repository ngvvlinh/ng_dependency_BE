package reportserver

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValidateSchema(t *testing.T) {
	Convey("convertColName", t, func() {
		So(convertColName(1), ShouldEqual, "B")
		So(convertColName(27), ShouldEqual, "AB")
		So(convertColName(52), ShouldEqual, "BA")
	})
}
