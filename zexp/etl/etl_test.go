package etl

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/tools/pkg/gen"
	"o.o/backend/zexp/etl/tests/testdst"
	"o.o/backend/zexp/etl/tests/testsrc"
	"o.o/backend/zexp/etl/types"
	"o.o/capi/dot"
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
		_, err := db.Insert(account)
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
		err := ng.scanModels(src.DB, plistSrc.Interface().(types.Model), ETLQuery{
			OrderBy: "id",
			Where:   []interface{}{sq.NewExpr("id >= ?", 10020)},
			Limit:   15,
		})
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
		err := ng.scanModels(src.DB, plistSrc.Interface().(types.Model), ETLQuery{
			OrderBy: "id",
			Where:   []interface{}{sq.NewExpr("id >= ?", 10020)},
			Limit:   15,
		})
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

func TestIsEqual(t *testing.T) {
	tests := []struct {
		obj         interface{}
		otherObject interface{}
		result      bool
	}{
		{
			obj: struct {
				FirstName string
				LastName  string
				Age       int
			}{
				FirstName: "A",
				LastName:  "Nguyen Van",
				Age:       18,
			},
			otherObject: struct {
				FirstName string
				LastName  string
				FullName  string
			}{
				FirstName: "A",
				LastName:  "Nguyen Van",
				FullName:  "Nguyen Van A",
			},
			result: true,
		},
		{
			obj: struct {
				FirstName string
				LastName  string
			}{
				FirstName: "A",
				LastName:  "Nguyen Van",
			},
			otherObject: struct {
				FirstName string
				LastName  string
			}{
				FirstName: "B",
				LastName:  "Nguyen Van",
			},
			result: false,
		},
	}
	Convey("isEqual", t, func() {
		for _, test := range tests {
			So(isEqual(reflect.ValueOf(test.obj), reflect.ValueOf(test.otherObject)), ShouldEqual, test.result)
		}
	})
}

func TestDeleteModelsByRids(t *testing.T) {
	accountModel := types.NewDataSource(srcDB, (*testsrc.Accounts)(nil)).Model
	Convey("test deleteModelsByRids", t, func() {
		var insertedRIDs []dot.ID
		for i := 0; i < 10; i++ {
			account := &testsrc.Account{
				ID:        101,
				FirstName: "a",
				LastName:  "b",
				Rid:       dot.ID(100 + i),
			}
			_, _err := srcDB.Insert(account)
			must(_err)

			insertedRIDs = append(insertedRIDs, dot.ID(100+i))
		}

		err := deleteModelsByRids(srcDB, accountModel, insertedRIDs)
		So(err, ShouldBeNil)

		rows, err := srcDB.Select("rid").From(accountModel.SQLTableName()).In("rid", insertedRIDs).Query()
		must(err)

		var rids []dot.ID
		var rid dot.ID
		for rows.Next() {
			err := rows.Scan(&rid)
			must(err)
			rids = append(rids, rid)
		}

		So(len(rids), ShouldEqual, 0)
	})
}
