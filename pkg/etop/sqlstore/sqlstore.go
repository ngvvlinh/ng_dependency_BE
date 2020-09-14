package sqlstore

import (
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/common/l"
)

var ll = l.New()

type (
	M  map[string]interface{}
	Ms map[string]string

	Query = cmsql.Query
	Qx    = cmsql.QueryInterface
)

func inTransaction(db *cmsql.Database, callback func(cmsql.QueryInterface) error) (err error) {
	return db.InTransaction(bus.Ctx(), callback)
}

func IDs(items []int64) []interface{} {
	res := make([]interface{}, len(items))
	for i, item := range items {
		res[i] = item
	}
	return res
}
