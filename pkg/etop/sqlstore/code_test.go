package sqlstore

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/model"
)

func TestGenerateCodeWithinTransaction(t *testing.T) {
	ctx := context.Background()
	reset := func() { MustExec(db, `TRUNCATE "code"`) }
	reset()

	Convey("Init code", t, func() {
		Reset(reset)

		n, err := createCode(ctx, db, &model.CreateCodeCommand{
			Code: &model.Code{Code: "123", Type: model.CodeTypeShop},
		})
		So(err, ShouldBeNil)
		So(n, ShouldEqual, int64(1))

		Convey("Generate with retry (success)", func() {
			var c int
			fn := func() string {
				c++
				if c <= 2 {
					return "123" // duplicated code
				}
				return "ABC"
			}

			err := db.InTransaction(bus.Ctx(), func(x cmsql.QueryInterface) error {
				code, err := generateCode(ctx, x, model.CodeTypeShop, fn)
				So(err, ShouldBeNil)
				So(code, ShouldEqual, "ABC")
				return err
			})
			So(err, ShouldBeNil)
			So(c, ShouldEqual, 3)

			Convey("Get the code back", func() {
				var item model.Code
				err := db.Where("code = 'ABC'").ShouldGet(&item)
				So(err, ShouldBeNil)
			})
		})

		Convey("Generate with retry (error)", func() {
			fn := func() string {
				return "123" // always duplicate
			}
			err := db.InTransaction(bus.Ctx(), func(x cmsql.QueryInterface) error {
				_, err := generateCode(ctx, x, model.CodeTypeShop, fn)
				return err
			})
			So(err, ShouldBeError, "Can not generate code for type: shop")
		})
	})
}
