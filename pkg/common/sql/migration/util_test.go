package migration

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
)

var db *cmsql.Database

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	cfg := cc.DefaultPostgres()
	db = cmsql.MustConnect(cfg)
	_, _ = db.Exec(`CREATE DATABASE test_migration`)

	cfg.Database = "test_migration"
	db = cmsql.MustConnect(cfg)
	db.MustExec(`DROP SCHEMA IF EXISTS public CASCADE`)
	db.MustExec(`CREATE SCHEMA public`)

	db.MustExec(`CREATE TYPE "mood" AS ENUM ('sad', 'ok', 'happy');`)
	db.MustExec(`
		CREATE TABLE account (
			id int8,
			ids int8[],
			names text[],
			first_name text,
			current_mood mood,
			created_at timestamptz,
			price int4,
			is_test boolean,
			info jsonb
		);`)
	db.MustExec(`
		INSERT INTO account
		VALUES(
			101,
			ARRAY[102, 103],
			ARRAY['a', 'bc', 'def'],
			'nguyen van a',
			'sad',
			now(),
			1024,
			true,
			'{ "a": "b",  "c": "d" }'
		);`)
}

func TestGetColumnNamesAndTypes(t *testing.T) {
	Convey("Test getColumnNamesAndTypes", t, func() {
		mapColumnNameAndType, err := GetColumnNamesAndTypes(db, "account")
		must(err)

		So(mapColumnNameAndType["id"], ShouldEqual, "int8")
		So(mapColumnNameAndType["ids"], ShouldEqual, "ARRAY")
		So(mapColumnNameAndType["names"], ShouldEqual, "ARRAY")
		So(mapColumnNameAndType["first_name"], ShouldEqual, "text")
		So(mapColumnNameAndType["current_mood"], ShouldEqual, "USER-DEFINED")
		So(mapColumnNameAndType["price"], ShouldEqual, "int4")
		So(mapColumnNameAndType["created_at"], ShouldEqual, "timestamptz")
		So(mapColumnNameAndType["is_test"], ShouldEqual, "boolean")
		So(mapColumnNameAndType["info"], ShouldEqual, "jsonb")
	})
}
