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

// Type ShopCarrier represents table shop_carrier
func sqlgenShopCarrier(_ *ShopCarrier) bool { return true }

type ShopCarriers []*ShopCarrier

const __sqlShopCarrier_Table = "shop_carrier"
const __sqlShopCarrier_ListCols = "\"id\",\"shop_id\",\"full_name\",\"note\",\"status\",\"created_at\",\"updated_at\",\"deleted_at\""
const __sqlShopCarrier_Insert = "INSERT INTO \"shop_carrier\" (" + __sqlShopCarrier_ListCols + ") VALUES"
const __sqlShopCarrier_Select = "SELECT " + __sqlShopCarrier_ListCols + " FROM \"shop_carrier\""
const __sqlShopCarrier_Select_history = "SELECT " + __sqlShopCarrier_ListCols + " FROM history.\"shop_carrier\""
const __sqlShopCarrier_UpdateAll = "UPDATE \"shop_carrier\" SET (" + __sqlShopCarrier_ListCols + ")"

func (m *ShopCarrier) SQLTableName() string  { return "shop_carrier" }
func (m *ShopCarriers) SQLTableName() string { return "shop_carrier" }
func (m *ShopCarrier) SQLListCols() string   { return __sqlShopCarrier_ListCols }

func (m *ShopCarrier) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlShopCarrier_ListCols + " FROM \"shop_carrier\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*ShopCarrier)(nil))
}

func (m *ShopCarrier) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		core.Int64(m.ID),
		core.Int64(m.ShopID),
		core.String(m.FullName),
		core.String(m.Note),
		core.Int32(m.Status),
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
		core.Time(m.DeletedAt),
	}
}

func (m *ShopCarrier) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		(*core.Int64)(&m.ID),
		(*core.Int64)(&m.ShopID),
		(*core.String)(&m.FullName),
		(*core.String)(&m.Note),
		(*core.Int32)(&m.Status),
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		(*core.Time)(&m.DeletedAt),
	}
}

func (m *ShopCarrier) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *ShopCarriers) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(ShopCarriers, 0, 128)
	for rows.Next() {
		m := new(ShopCarrier)
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

func (_ *ShopCarrier) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCarrier_Select)
	return nil
}

func (_ *ShopCarriers) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCarrier_Select)
	return nil
}

func (m *ShopCarrier) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCarrier_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(8)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms ShopCarriers) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCarrier_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(8)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *ShopCarrier) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("shop_carrier")
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
	if m.FullName != "" {
		flag = true
		w.WriteName("full_name")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.FullName)
	}
	if m.Note != "" {
		flag = true
		w.WriteName("note")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Note)
	}
	if m.Status != 0 {
		flag = true
		w.WriteName("status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Status)
	}
	if !m.CreatedAt.IsZero() {
		flag = true
		w.WriteName("created_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CreatedAt)
	}
	if !m.UpdatedAt.IsZero() {
		flag = true
		w.WriteName("updated_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Now(m.UpdatedAt, time.Now(), true))
	}
	if !m.DeletedAt.IsZero() {
		flag = true
		w.WriteName("deleted_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.DeletedAt)
	}
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *ShopCarrier) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCarrier_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(8)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type ShopCarrierHistory map[string]interface{}
type ShopCarrierHistories []map[string]interface{}

func (m *ShopCarrierHistory) SQLTableName() string  { return "history.\"shop_carrier\"" }
func (m ShopCarrierHistories) SQLTableName() string { return "history.\"shop_carrier\"" }

func (m *ShopCarrierHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCarrier_Select_history)
	return nil
}

func (m ShopCarrierHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCarrier_Select_history)
	return nil
}

func (m ShopCarrierHistory) ID() core.Interface        { return core.Interface{m["id"]} }
func (m ShopCarrierHistory) ShopID() core.Interface    { return core.Interface{m["shop_id"]} }
func (m ShopCarrierHistory) FullName() core.Interface  { return core.Interface{m["full_name"]} }
func (m ShopCarrierHistory) Note() core.Interface      { return core.Interface{m["note"]} }
func (m ShopCarrierHistory) Status() core.Interface    { return core.Interface{m["status"]} }
func (m ShopCarrierHistory) CreatedAt() core.Interface { return core.Interface{m["created_at"]} }
func (m ShopCarrierHistory) UpdatedAt() core.Interface { return core.Interface{m["updated_at"]} }
func (m ShopCarrierHistory) DeletedAt() core.Interface { return core.Interface{m["deleted_at"]} }

func (m *ShopCarrierHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 8)
	args := make([]interface{}, 8)
	for i := 0; i < 8; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(ShopCarrierHistory, 8)
	res["id"] = data[0]
	res["shop_id"] = data[1]
	res["full_name"] = data[2]
	res["note"] = data[3]
	res["status"] = data[4]
	res["created_at"] = data[5]
	res["updated_at"] = data[6]
	res["deleted_at"] = data[7]
	*m = res
	return nil
}

func (ms *ShopCarrierHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 8)
	args := make([]interface{}, 8)
	for i := 0; i < 8; i++ {
		args[i] = &data[i]
	}
	res := make(ShopCarrierHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(ShopCarrierHistory)
		m["id"] = data[0]
		m["shop_id"] = data[1]
		m["full_name"] = data[2]
		m["note"] = data[3]
		m["status"] = data[4]
		m["created_at"] = data[5]
		m["updated_at"] = data[6]
		m["deleted_at"] = data[7]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
