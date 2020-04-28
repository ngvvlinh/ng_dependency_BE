package sqlstore

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	. "o.o/backend/pkg/common/testing"
	"o.o/backend/zexp/sample/model"
)

var db *cmsql.Database

func init() {
	db = cmsql.MustConnect(cc.DefaultPostgres())
	db.MustExec(`
DROP TABLE IF EXISTS foo;
CREATE TABLE foo (
  id int8,
  account_id int8,
  abc_2 text,
  def2 text,
  created_at timestamptz,
  updated_at timestamptz
)
`)
}

func TestModel(t *testing.T) {
	Convey("Foo", t, func() {
		Reset(func() {
			db.MustExec("truncate foo")
		})
		var ft FooFilters

		foo := &model.Foo{
			ID:        123,
			AccountID: 0,
			ABC:       "ABC",
			Def2:      "xyz",
		}
		_, err := db.Insert(foo)
		So(err, ShouldBeNil)

		Convey("Get", func() {
			var foo2 model.Foo
			ok, err := db.Where(ft.ByID(123)).Get(&foo2)
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)

			foo2.CreatedAt = foo.CreatedAt
			foo2.UpdatedAt = foo.UpdatedAt
			So(&foo2, ShouldDeepEqual, foo)
		})
		Convey("Update", func() {
			foo2 := &model.Foo{
				AccountID: 999,
			}
			err := db.Where("id = ?", foo.ID).ShouldUpdate(foo2)
			So(err, ShouldBeNil)

			Convey("Get again", func() {
				var foo3 model.Foo
				err := db.Where("id = ?", foo.ID).ShouldGet(&foo3)
				So(err, ShouldBeNil)

				foo.AccountID = foo2.AccountID
				foo3.UpdatedAt = foo2.UpdatedAt
				foo.UpdatedAt = foo2.UpdatedAt
				foo.CreatedAt = foo3.CreatedAt
				So(&foo3, ShouldDeepEqual, foo)
			})
		})
		Convey("Get 2", func() {
			var foo2 model.Foo
			ok, err := db.Where("id = ?", foo.ID).Get(&foo2)
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)
			So(foo2.AccountID, ShouldEqual, 0)
		})
		Convey("Update 2", func() {
			foo2 := &model.Foo{
				ABC: "000",
			}
			err := db.Where("id = ?", foo.ID).ShouldUpdate(foo2)
			So(err, ShouldBeNil)

			Convey("Get again", func() {
				var foo3 model.Foo
				err := db.Where("id = ?", foo.ID).ShouldGet(&foo3)
				So(err, ShouldBeNil)

				foo.ABC = foo2.ABC
				foo3.UpdatedAt = foo2.UpdatedAt
				foo.UpdatedAt = foo2.UpdatedAt
				foo.CreatedAt = foo3.CreatedAt
				So(&foo3, ShouldDeepEqual, foo)
				So(foo3.AccountID, ShouldEqual, 0)
			})
		})
	})
}
