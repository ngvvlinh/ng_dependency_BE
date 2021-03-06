// +build !generator

// Code generated by generator sqlgen. DO NOT EDIT.

package model

import (
	"database/sql"
	"sync"
	"time"

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

type Equations []*Equation

const __sqlEquation_Table = "equation"
const __sqlEquation_ListCols = "\"id\",\"equation\",\"result\",\"created_at\",\"updated_at\""
const __sqlEquation_ListColsOnConflict = "\"id\" = EXCLUDED.\"id\",\"equation\" = EXCLUDED.\"equation\",\"result\" = EXCLUDED.\"result\",\"created_at\" = EXCLUDED.\"created_at\",\"updated_at\" = EXCLUDED.\"updated_at\""
const __sqlEquation_Insert = "INSERT INTO \"equation\" (" + __sqlEquation_ListCols + ") VALUES"
const __sqlEquation_Select = "SELECT " + __sqlEquation_ListCols + " FROM \"equation\""
const __sqlEquation_Select_history = "SELECT " + __sqlEquation_ListCols + " FROM history.\"equation\""
const __sqlEquation_UpdateAll = "UPDATE \"equation\" SET (" + __sqlEquation_ListCols + ")"
const __sqlEquation_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT equation_pkey DO UPDATE SET"

func (m *Equation) SQLTableName() string  { return "equation" }
func (m *Equations) SQLTableName() string { return "equation" }
func (m *Equation) SQLListCols() string   { return __sqlEquation_ListCols }

func (m *Equation) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlEquation_ListCols + " FROM \"equation\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func (m *Equation) Migration(db *cmsql.Database) {
	var mDBColumnNameAndType map[string]string
	if val, err := migration.GetColumnNamesAndTypes(db, "equation"); err != nil {
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
		"equation": {
			ColumnName:       "equation",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"result": {
			ColumnName:       "result",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"created_at": {
			ColumnName:       "created_at",
			ColumnType:       "dot.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"updated_at": {
			ColumnName:       "updated_at",
			ColumnType:       "dot.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
	}
	if err := migration.Compare(db, "equation", mModelColumnNameAndType, mDBColumnNameAndType); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*Equation)(nil))
}

func (m *Equation) SQLArgs(opts core.Opts, create bool) []interface{} {
	return []interface{}{
		m.ID,
		core.String(m.Equation),
		core.String(m.Result),
		m.CreatedAt,
		m.UpdatedAt,
	}
}

func (m *Equation) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		(*core.String)(&m.Equation),
		(*core.String)(&m.Result),
		&m.CreatedAt,
		&m.UpdatedAt,
	}
}

func (m *Equation) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *Equations) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(Equations, 0, 128)
	for rows.Next() {
		m := new(Equation)
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

func (_ *Equation) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlEquation_Select)
	return nil
}

func (_ *Equations) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlEquation_Select)
	return nil
}

func (m *Equation) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlEquation_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(5)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms Equations) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlEquation_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(5)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *Equation) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlEquation_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlEquation_ListColsOnConflict)
	return nil
}

func (ms Equations) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlEquation_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlEquation_ListColsOnConflict)
	return nil
}

func (m *Equation) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("equation")
	w.WriteRawString(" SET ")
	if m.ID != 0 {
		flag = true
		w.WriteName("id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ID)
	}
	if m.Equation != "" {
		flag = true
		w.WriteName("equation")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Equation)
	}
	if m.Result != "" {
		flag = true
		w.WriteName("result")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Result)
	}
	if true {
		flag = true
		w.WriteName("created_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CreatedAt)
	}
	if true {
		flag = true
		w.WriteName("updated_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.UpdatedAt)
	}
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *Equation) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlEquation_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(5)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type EquationHistory map[string]interface{}
type EquationHistories []map[string]interface{}

func (m *EquationHistory) SQLTableName() string  { return "history.\"equation\"" }
func (m EquationHistories) SQLTableName() string { return "history.\"equation\"" }

func (m *EquationHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlEquation_Select_history)
	return nil
}

func (m EquationHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlEquation_Select_history)
	return nil
}

func (m EquationHistory) ID() core.Interface        { return core.Interface{m["id"]} }
func (m EquationHistory) Equation() core.Interface  { return core.Interface{m["equation"]} }
func (m EquationHistory) Result() core.Interface    { return core.Interface{m["result"]} }
func (m EquationHistory) CreatedAt() core.Interface { return core.Interface{m["created_at"]} }
func (m EquationHistory) UpdatedAt() core.Interface { return core.Interface{m["updated_at"]} }

func (m *EquationHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 5)
	args := make([]interface{}, 5)
	for i := 0; i < 5; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(EquationHistory, 5)
	res["id"] = data[0]
	res["equation"] = data[1]
	res["result"] = data[2]
	res["created_at"] = data[3]
	res["updated_at"] = data[4]
	*m = res
	return nil
}

func (ms *EquationHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 5)
	args := make([]interface{}, 5)
	for i := 0; i < 5; i++ {
		args[i] = &data[i]
	}
	res := make(EquationHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(EquationHistory)
		m["id"] = data[0]
		m["equation"] = data[1]
		m["result"] = data[2]
		m["created_at"] = data[3]
		m["updated_at"] = data[4]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
