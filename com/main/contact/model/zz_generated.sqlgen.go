// +build !generator

// Code generated by generator sqlgen. DO NOT EDIT.

package model

import (
	"database/sql"
	"sync"
	time "time"

	cmsql "o.o/backend/pkg/common/sql/cmsql"
	migration "o.o/backend/pkg/common/sql/migration"
	core "o.o/backend/pkg/common/sql/sq/core"
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

type Contacts []*Contact

const __sqlContact_Table = "contact"
const __sqlContact_ListCols = "\"id\",\"shop_id\",\"full_name\",\"phone\",\"wl_partner_id\",\"created_at\",\"updated_at\",\"deleted_at\""
const __sqlContact_ListColsOnConflict = "\"id\" = EXCLUDED.\"id\",\"shop_id\" = EXCLUDED.\"shop_id\",\"full_name\" = EXCLUDED.\"full_name\",\"phone\" = EXCLUDED.\"phone\",\"wl_partner_id\" = EXCLUDED.\"wl_partner_id\",\"created_at\" = EXCLUDED.\"created_at\",\"updated_at\" = EXCLUDED.\"updated_at\",\"deleted_at\" = EXCLUDED.\"deleted_at\""
const __sqlContact_Insert = "INSERT INTO \"contact\" (" + __sqlContact_ListCols + ") VALUES"
const __sqlContact_Select = "SELECT " + __sqlContact_ListCols + " FROM \"contact\""
const __sqlContact_Select_history = "SELECT " + __sqlContact_ListCols + " FROM history.\"contact\""
const __sqlContact_UpdateAll = "UPDATE \"contact\" SET (" + __sqlContact_ListCols + ")"
const __sqlContact_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT contact_pkey DO UPDATE SET"

func (m *Contact) SQLTableName() string  { return "contact" }
func (m *Contacts) SQLTableName() string { return "contact" }
func (m *Contact) SQLListCols() string   { return __sqlContact_ListCols }

func (m *Contact) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlContact_ListCols + " FROM \"contact\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func (m *Contact) Migration(db *cmsql.Database) {
	var mDBColumnNameAndType map[string]string
	if val, err := migration.GetColumnNamesAndTypes(db, "contact"); err != nil {
		db.RecordError(err)
		return
	} else {
		mDBColumnNameAndType = val
	}
	mModelColumnNameAndType := map[string]migration.ColumnDef{
		"id": {
			ColumnName:       "id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"shop_id": {
			ColumnName:       "shop_id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"full_name": {
			ColumnName:       "full_name",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"phone": {
			ColumnName:       "phone",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"wl_partner_id": {
			ColumnName:       "wl_partner_id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"created_at": {
			ColumnName:       "created_at",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"updated_at": {
			ColumnName:       "updated_at",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"deleted_at": {
			ColumnName:       "deleted_at",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
	}
	if err := migration.Compare(db, "contact", mModelColumnNameAndType, mDBColumnNameAndType); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*Contact)(nil))
}

func (m *Contact) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		m.ID,
		m.ShopID,
		core.String(m.FullName),
		core.String(m.Phone),
		m.WLPartnerID,
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
		core.Time(m.DeletedAt),
	}
}

func (m *Contact) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		&m.ShopID,
		(*core.String)(&m.FullName),
		(*core.String)(&m.Phone),
		&m.WLPartnerID,
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		(*core.Time)(&m.DeletedAt),
	}
}

func (m *Contact) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *Contacts) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(Contacts, 0, 128)
	for rows.Next() {
		m := new(Contact)
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

func (_ *Contact) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlContact_Select)
	return nil
}

func (_ *Contacts) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlContact_Select)
	return nil
}

func (m *Contact) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlContact_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(8)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms Contacts) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlContact_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(8)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *Contact) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlContact_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlContact_ListColsOnConflict)
	return nil
}

func (ms Contacts) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlContact_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlContact_ListColsOnConflict)
	return nil
}

func (m *Contact) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("contact")
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
	if m.Phone != "" {
		flag = true
		w.WriteName("phone")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Phone)
	}
	if m.WLPartnerID != 0 {
		flag = true
		w.WriteName("wl_partner_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.WLPartnerID)
	}
	if !m.CreatedAt.IsZero() {
		flag = true
		w.WriteName("created_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CreatedAt)
	}
	if true { // always update time
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

func (m *Contact) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlContact_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(8)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type ContactHistory map[string]interface{}
type ContactHistories []map[string]interface{}

func (m *ContactHistory) SQLTableName() string  { return "history.\"contact\"" }
func (m ContactHistories) SQLTableName() string { return "history.\"contact\"" }

func (m *ContactHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlContact_Select_history)
	return nil
}

func (m ContactHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlContact_Select_history)
	return nil
}

func (m ContactHistory) ID() core.Interface          { return core.Interface{m["id"]} }
func (m ContactHistory) ShopID() core.Interface      { return core.Interface{m["shop_id"]} }
func (m ContactHistory) FullName() core.Interface    { return core.Interface{m["full_name"]} }
func (m ContactHistory) Phone() core.Interface       { return core.Interface{m["phone"]} }
func (m ContactHistory) WLPartnerID() core.Interface { return core.Interface{m["wl_partner_id"]} }
func (m ContactHistory) CreatedAt() core.Interface   { return core.Interface{m["created_at"]} }
func (m ContactHistory) UpdatedAt() core.Interface   { return core.Interface{m["updated_at"]} }
func (m ContactHistory) DeletedAt() core.Interface   { return core.Interface{m["deleted_at"]} }

func (m *ContactHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 8)
	args := make([]interface{}, 8)
	for i := 0; i < 8; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(ContactHistory, 8)
	res["id"] = data[0]
	res["shop_id"] = data[1]
	res["full_name"] = data[2]
	res["phone"] = data[3]
	res["wl_partner_id"] = data[4]
	res["created_at"] = data[5]
	res["updated_at"] = data[6]
	res["deleted_at"] = data[7]
	*m = res
	return nil
}

func (ms *ContactHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 8)
	args := make([]interface{}, 8)
	for i := 0; i < 8; i++ {
		args[i] = &data[i]
	}
	res := make(ContactHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(ContactHistory)
		m["id"] = data[0]
		m["shop_id"] = data[1]
		m["full_name"] = data[2]
		m["phone"] = data[3]
		m["wl_partner_id"] = data[4]
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
