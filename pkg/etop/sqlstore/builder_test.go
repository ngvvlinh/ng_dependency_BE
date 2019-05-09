package sqlstore

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSQLUpdateBuilder(t *testing.T) {
	Convey("sqlUpdateBuilder", t, func() {
		s := sqlUpdate("inverted").
			ColIndex("foo", 1, -1)

		Convey("Single elem", func() {
			sql, args := s.Build()
			So(sql, ShouldEqual, `UPDATE "inverted" SET foo[1] = ?`)
			So(args, ShouldResemble, []interface{}{-1})
		})
		Convey("Multiple elem", func() {
			s.ColIndex("bar", 2, -2)
			s.ColIndex("quix", 10, -10)

			Convey("Should build correctly", func() {
				sql, args := s.Build()
				So(sql, ShouldEqual,
					`UPDATE "inverted" SET foo[1] = ?, bar[2] = ?, quix[10] = ?`)
				So(args, ShouldResemble, []interface{}{-1, -2, -10})
			})
			Convey("With where", func() {
				s.Where("id = ?", 1001)
				s.Where("status = ?", "ok")
				sql, args := s.Build()
				So(sql, ShouldEqual,
					`UPDATE "inverted" SET foo[1] = ?, bar[2] = ?, quix[10] = ? WHERE id = ? AND status = ?`)
				So(args, ShouldResemble, []interface{}{-1, -2, -10, 1001, "ok"})
			})
		})
	})
}
