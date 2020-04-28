package sqlstore

import (
	"encoding/json"
	"time"

	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
)

func init() {
	InitTest()
}

func InitTest() {
	engine, err := cmsql.Connect(cc.DefaultPostgres())
	if err != nil {
		panic(err)
	}
	Init(engine)
	x = engine

	MustExec("SELECT 1")
}

func MustExec(sql string) {
	_, err := x.Exec(sql)
	if err != nil {
		panic(err)
	}
}

func jSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func toTime(s string) time.Time {
	var t time.Time
	err := json.Unmarshal([]byte(`"`+s+`"`), &t)
	if err != nil {
		panic(err)
	}
	return t
}

func toPTime(s string) *time.Time {
	t := toTime(s)
	return &t
}
