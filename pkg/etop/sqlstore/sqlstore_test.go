package sqlstore

import (
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
)

var db *cmsql.Database

func init() {
	InitTest()
}

func InitTest() {
	db = cmsql.MustConnect(cc.DefaultPostgres())
	MustExec(db, "SELECT 1")
}

func MustExec(db *cmsql.Database, sql string) {
	_, err := db.Exec(sql)
	if err != nil {
		panic(err)
	}
}
