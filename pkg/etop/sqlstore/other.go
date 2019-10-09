package sqlstore

import (
	"context"
	"encoding/json"
	"fmt"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/µjson"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/common/jsonx"
)

func init() {
	bus.AddHandlers("sql", GetHistory)
}

func validateColumn(name string) bool {
	if len(name) == 0 {
		return false
	}
	for i, l := 0, len(name); i < l; i++ {
		c := name[i]
		switch {
		case c >= 'a' && c <= 'z':
		case c == '_':
		default:
			return false
		}
	}
	return true
}

func GetHistory(ctx context.Context, query *model.GetHistoryQuery) error {
	if query.Table == "" {
		return cm.Error(cm.InvalidArgument, "Missing table", nil)
	}
	if !validateColumn(query.Table) {
		return cm.Errorf(cm.InvalidArgument, nil, "Invalid table name: %v", query.Table)
	}

	s := x.
		Select(`row_to_json("` + query.Table + `")`).
		From(`history."` + query.Table + `"`)
	s, err := LimitSort(s, query.Paging, Ms{"rid": ""})
	if err != nil {
		return err
	}

	for name, value := range query.Filters {
		if !validateColumn(name) {
			return cm.Errorf(cm.InvalidArgument, nil, "Tên không hợp lệ: %v", name)
		}
		s = s.Where(`"`+name+`" = ?`, value)
	}

	sql, args, err := s.Build()
	if err != nil {
		return err
	}
	rows, err := x.DB().Query(sql, args...)
	if err != nil {
		return err
	}

	if query.KeepRaw {
		var items []json.RawMessage
		if query.Paging != nil {
			items = make([]json.RawMessage, 0, query.Paging.Limit)
		}
		for rows.Next() {
			var item []byte
			if err = rows.Scan(&item); err != nil {
				return err
			}
			items = append(items, item)
		}
		query.Result.Raws = items
		query.Result.Len = len(items)

	} else {
		c := 0
		data := make([]byte, 0, 1024)
		data = append(data, '[')
		for rows.Next() {
			if c != 0 {
				data = append(data, ',')
			}
			c++
			var item []byte
			if err = rows.Scan(&item); err != nil {
				return err
			}
			data, err = µjson.FilterAndRename(data, item)
			if err != nil {
				return err
			}
		}
		data = append(data, ']')
		query.Result.Data = data
		query.Result.Len = c

		if cm.IsDev() {
			var v interface{}
			if err := jsonx.Unmarshal(data, &v); err != nil {
				ll.Error("Invalid json output")
				fmt.Printf("-- output\n%s\n\n", data)
			}
		}
	}
	return nil
}
