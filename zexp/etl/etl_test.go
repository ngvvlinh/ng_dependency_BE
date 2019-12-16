package etl

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/tools/pkg/gen"
	"etop.vn/backend/zexp/etl/tests/testdst"
	"etop.vn/backend/zexp/etl/tests/testsrc"
	"etop.vn/backend/zexp/etl/types"
	"etop.vn/capi/dot"
)

var srcDB, dstDB *cmsql.Database

func init() {
	path := filepath.Join(gen.ProjectPath(), "zexp/etl")

	cfg := cc.DefaultPostgres()
	db := cmsql.MustConnect(cfg)
	_, _ = db.Exec(`CREATE DATABASE test_src`)
	_, _ = db.Exec(`CREATE DATABASE test_dst`)

	cfg.Database = "test_src"
	srcDB = cmsql.MustConnect(cfg)
	srcDB.MustExec(`DROP SCHEMA IF EXISTS public CASCADE`)
	srcDB.MustExec(`CREATE SCHEMA public`)

	cfg.Database = "test_dst"
	dstDB = cmsql.MustConnect(cfg)
	dstDB.MustExec(`DROP SCHEMA IF EXISTS public CASCADE`)
	dstDB.MustExec(`CREATE SCHEMA public`)

	srcSQL, err := ioutil.ReadFile(filepath.Join(path, "tests/test_src.sql"))
	must(err)
	srcDB.MustExec(string(srcSQL))

	dstSQL, err := ioutil.ReadFile(filepath.Join(path, "tests/test_dst.sql"))
	must(err)
	dstDB.MustExec(string(dstSQL))

	populateSampleData(srcDB)
}

func populateSampleData(db *cmsql.Database) {
	for i := 0; i < 200; i++ {
		account := &testsrc.Account{
			ID:        dot.ID(10000 + i),
			FirstName: fmt.Sprintf("Al%02dce", i),
			LastName:  fmt.Sprintf("St%02drk", i),
		}
		_, err := srcDB.Insert(account)
		must(err)
	}
}

func TestEtl(t *testing.T) {
	reset := func() { dstDB.MustExec(`TRUNCATE account`) }
	reset()
}

func TestScan(t *testing.T) {
	src := types.NewDataSource(srcDB, (*testsrc.Accounts)(nil))

	ng := NewETLEngine(nil)
	ng.RegisterConversion(testdst.RegisterConversions)
	ng.Bootstrap()

	Convey("scan account", t, func() {
		plistSrc := reflect.New(reflect.TypeOf(src.Model).Elem())
		err := ng.scanModels(src.DB, plistSrc.Interface().(types.Model), 10020, 15)
		So(err, ShouldBeNil)

		accounts := plistSrc.Elem().Interface().(testsrc.Accounts)
		So(accounts, ShouldHaveLength, 15)
		So(accounts[0].ID, ShouldEqual, 10020)
		So(accounts[14].ID, ShouldEqual, 10034)
	})
}

func TestTransform(t *testing.T) {
	src := types.NewDataSource(srcDB, (*testsrc.Accounts)(nil))
	dst := types.NewDataSource(dstDB, (*testdst.Accounts)(nil))

	ng := NewETLEngine(nil)
	ng.RegisterConversion(testdst.RegisterConversions)
	ng.Bootstrap()

	Convey("scan & transform", t, func() {
		plistSrc := reflect.New(reflect.TypeOf(src.Model).Elem())
		err := ng.scanModels(src.DB, plistSrc.Interface().(types.Model), 10020, 15)
		So(err, ShouldBeNil)

		Convey("transform", func() {
			plistDst := reflect.New(reflect.TypeOf(dst.Model).Elem())
			err := ng.transform(plistSrc.Elem().Interface(), plistDst.Interface())
			So(err, ShouldBeNil)

			accounts := plistDst.Elem().Interface().(testdst.Accounts)
			So(accounts, ShouldHaveLength, 15)
			So(accounts[0].ID, ShouldEqual, 10020)
			So(accounts[0].FullName, ShouldEqual, "Al20ce St20rk")
		})
	})
}
