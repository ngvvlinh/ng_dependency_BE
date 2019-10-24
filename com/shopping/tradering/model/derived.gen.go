// Code generated by goderive DO NOT EDIT.

package model

import (
	"database/sql"
	"sync"
	"time"

	"etop.vn/backend/pkg/common/cmsql"
	core "etop.vn/backend/pkg/common/sq/core"
)

var __sqlModels []interface{ SQLVerifySchema(db *cmsql.Database) }
var __sqlonce sync.Once

func SQLVerifySchema(db *cmsql.Database) {
	__sqlonce.Do(func() {
		for _, m := range __sqlModels {
			m.SQLVerifySchema(db)
		}
	})
}

type SQLWriter = core.SQLWriter

// Type ShopTrader represents table shop_trader
func sqlgenShopTrader(_ *ShopTrader) bool { return true }

type ShopTraders []*ShopTrader

const __sqlShopTrader_Table = "shop_trader"
const __sqlShopTrader_ListCols = "\"id\",\"shop_id\",\"type\""
const __sqlShopTrader_Insert = "INSERT INTO \"shop_trader\" (" + __sqlShopTrader_ListCols + ") VALUES"
const __sqlShopTrader_Select = "SELECT " + __sqlShopTrader_ListCols + " FROM \"shop_trader\""
const __sqlShopTrader_Select_history = "SELECT " + __sqlShopTrader_ListCols + " FROM history.\"shop_trader\""
const __sqlShopTrader_UpdateAll = "UPDATE \"shop_trader\" SET (" + __sqlShopTrader_ListCols + ")"

func (m *ShopTrader) SQLTableName() string  { return "shop_trader" }
func (m *ShopTraders) SQLTableName() string { return "shop_trader" }
func (m *ShopTrader) SQLListCols() string   { return __sqlShopTrader_ListCols }

func (m *ShopTrader) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlShopTrader_ListCols + " FROM \"shop_trader\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*ShopTrader)(nil))
}

func (m *ShopTrader) SQLArgs(opts core.Opts, create bool) []interface{} {
	return []interface{}{
		core.Int64(m.ID),
		core.Int64(m.ShopID),
		core.String(m.Type),
	}
}

func (m *ShopTrader) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		(*core.Int64)(&m.ID),
		(*core.Int64)(&m.ShopID),
		(*core.String)(&m.Type),
	}
}

func (m *ShopTrader) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *ShopTraders) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(ShopTraders, 0, 128)
	for rows.Next() {
		m := new(ShopTrader)
		args := m.SQLScanArgs(opts)
		if err := rows.Scan(args...); err != nil {
			return err
		}
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}

func (_ *ShopTrader) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopTrader_Select)
	return nil
}

func (_ *ShopTraders) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopTrader_Select)
	return nil
}

func (m *ShopTrader) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShopTrader_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(3)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms ShopTraders) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShopTrader_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(3)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *ShopTrader) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("shop_trader")
	w.WriteRawString(" SET ")
	if m.ID != 0 {
		flag = true
		w.WriteName("id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ID)
	}
	if m.ShopID != 0 {
		flag = true
		w.WriteName("shop_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShopID)
	}
	if m.Type != "" {
		flag = true
		w.WriteName("type")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Type)
	}
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *ShopTrader) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlShopTrader_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(3)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type ShopTraderHistory map[string]interface{}
type ShopTraderHistories []map[string]interface{}

func (m *ShopTraderHistory) SQLTableName() string  { return "history.\"shop_trader\"" }
func (m ShopTraderHistories) SQLTableName() string { return "history.\"shop_trader\"" }

func (m *ShopTraderHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopTrader_Select_history)
	return nil
}

func (m ShopTraderHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopTrader_Select_history)
	return nil
}

func (m ShopTraderHistory) ID() core.Interface     { return core.Interface{m["id"]} }
func (m ShopTraderHistory) ShopID() core.Interface { return core.Interface{m["shop_id"]} }
func (m ShopTraderHistory) Type() core.Interface   { return core.Interface{m["type"]} }

func (m *ShopTraderHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 3)
	args := make([]interface{}, 3)
	for i := 0; i < 3; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(ShopTraderHistory, 3)
	res["id"] = data[0]
	res["shop_id"] = data[1]
	res["type"] = data[2]
	*m = res
	return nil
}

func (ms *ShopTraderHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 3)
	args := make([]interface{}, 3)
	for i := 0; i < 3; i++ {
		args[i] = &data[i]
	}
	res := make(ShopTraderHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(ShopTraderHistory)
		m["id"] = data[0]
		m["shop_id"] = data[1]
		m["type"] = data[2]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
