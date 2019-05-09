package sqlstore

import (
	"encoding/json"
	"time"

	"etop.vn/backend/cmd/etop-server/config"
	"etop.vn/backend/pkg/common/cmsql"
)

func init() {
	InitTest()
}

func InitTest() {
	cfg := config.DefaultTest()
	engine, err := cmsql.Connect(cmsql.ConfigPostgres(cfg.Postgres))
	if err != nil {
		panic(err)
	}
	Init(engine)

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
