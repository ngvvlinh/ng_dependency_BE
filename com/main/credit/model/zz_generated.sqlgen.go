// +build !generator

// Code generated by generator sqlgen. DO NOT EDIT.

package model

import (
	"database/sql"
	"sync"
	time "time"

	identitymodel "o.o/backend/com/main/identity/model"
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

type Credits []*Credit

const __sqlCredit_Table = "credit"
const __sqlCredit_ListCols = "\"id\",\"amount\",\"shop_id\",\"type\",\"status\",\"created_at\",\"updated_at\",\"deleted_at\",\"paid_at\",\"classify\",\"bank_statement_id\""
const __sqlCredit_ListColsOnConflict = "\"id\" = EXCLUDED.\"id\",\"amount\" = EXCLUDED.\"amount\",\"shop_id\" = EXCLUDED.\"shop_id\",\"type\" = EXCLUDED.\"type\",\"status\" = EXCLUDED.\"status\",\"created_at\" = EXCLUDED.\"created_at\",\"updated_at\" = EXCLUDED.\"updated_at\",\"deleted_at\" = EXCLUDED.\"deleted_at\",\"paid_at\" = EXCLUDED.\"paid_at\",\"classify\" = EXCLUDED.\"classify\",\"bank_statement_id\" = EXCLUDED.\"bank_statement_id\""
const __sqlCredit_Insert = "INSERT INTO \"credit\" (" + __sqlCredit_ListCols + ") VALUES"
const __sqlCredit_Select = "SELECT " + __sqlCredit_ListCols + " FROM \"credit\""
const __sqlCredit_Select_history = "SELECT " + __sqlCredit_ListCols + " FROM history.\"credit\""
const __sqlCredit_UpdateAll = "UPDATE \"credit\" SET (" + __sqlCredit_ListCols + ")"
const __sqlCredit_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT credit_pkey DO UPDATE SET"

func (m *Credit) SQLTableName() string  { return "credit" }
func (m *Credits) SQLTableName() string { return "credit" }
func (m *Credit) SQLListCols() string   { return __sqlCredit_ListCols }

func (m *Credit) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlCredit_ListCols + " FROM \"credit\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func (m *Credit) Migration(db *cmsql.Database) {
	var mDBColumnNameAndType map[string]string
	if val, err := migration.GetColumnNamesAndTypes(db, "credit"); err != nil {
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
		"amount": {
			ColumnName:       "amount",
			ColumnType:       "int",
			ColumnDBType:     "int",
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
		"type": {
			ColumnName:       "type",
			ColumnType:       "credit_type.CreditType",
			ColumnDBType:     "enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{"shop"},
		},
		"status": {
			ColumnName:       "status",
			ColumnType:       "status3.Status",
			ColumnDBType:     "enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{"Z", "P", "N"},
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
		"paid_at": {
			ColumnName:       "paid_at",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"classify": {
			ColumnName:       "classify",
			ColumnType:       "credit_type.CreditClassify",
			ColumnDBType:     "enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{"shipping", "telecom"},
		},
		"bank_statement_id": {
			ColumnName:       "bank_statement_id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
	}
	if err := migration.Compare(db, "credit", mModelColumnNameAndType, mDBColumnNameAndType); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*Credit)(nil))
}

func (m *Credit) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		m.ID,
		core.Int(m.Amount),
		m.ShopID,
		m.Type,
		m.Status,
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
		core.Time(m.DeletedAt),
		core.Time(m.PaidAt),
		m.Classify,
		m.BankStatementID,
	}
}

func (m *Credit) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		(*core.Int)(&m.Amount),
		&m.ShopID,
		&m.Type,
		&m.Status,
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		(*core.Time)(&m.DeletedAt),
		(*core.Time)(&m.PaidAt),
		&m.Classify,
		&m.BankStatementID,
	}
}

func (m *Credit) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *Credits) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(Credits, 0, 128)
	for rows.Next() {
		m := new(Credit)
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

func (_ *Credit) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlCredit_Select)
	return nil
}

func (_ *Credits) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlCredit_Select)
	return nil
}

func (m *Credit) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlCredit_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(11)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms Credits) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlCredit_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(11)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *Credit) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlCredit_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlCredit_ListColsOnConflict)
	return nil
}

func (ms Credits) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlCredit_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlCredit_ListColsOnConflict)
	return nil
}

func (m *Credit) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("credit")
	w.WriteRawString(" SET ")
	if m.ID != 0 {
		flag = true
		w.WriteName("id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ID)
	}
	if m.Amount != 0 {
		flag = true
		w.WriteName("amount")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Amount)
	}
	if m.ShopID != 0 {
		flag = true
		w.WriteName("shop_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShopID)
	}
	if m.Type != 0 {
		flag = true
		w.WriteName("type")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Type)
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
	if !m.PaidAt.IsZero() {
		flag = true
		w.WriteName("paid_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.PaidAt)
	}
	if m.Classify != 0 {
		flag = true
		w.WriteName("classify")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Classify)
	}
	if m.BankStatementID != 0 {
		flag = true
		w.WriteName("bank_statement_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.BankStatementID)
	}
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *Credit) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlCredit_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(11)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type CreditHistory map[string]interface{}
type CreditHistories []map[string]interface{}

func (m *CreditHistory) SQLTableName() string  { return "history.\"credit\"" }
func (m CreditHistories) SQLTableName() string { return "history.\"credit\"" }

func (m *CreditHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlCredit_Select_history)
	return nil
}

func (m CreditHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlCredit_Select_history)
	return nil
}

func (m CreditHistory) ID() core.Interface        { return core.Interface{m["id"]} }
func (m CreditHistory) Amount() core.Interface    { return core.Interface{m["amount"]} }
func (m CreditHistory) ShopID() core.Interface    { return core.Interface{m["shop_id"]} }
func (m CreditHistory) Type() core.Interface      { return core.Interface{m["type"]} }
func (m CreditHistory) Status() core.Interface    { return core.Interface{m["status"]} }
func (m CreditHistory) CreatedAt() core.Interface { return core.Interface{m["created_at"]} }
func (m CreditHistory) UpdatedAt() core.Interface { return core.Interface{m["updated_at"]} }
func (m CreditHistory) DeletedAt() core.Interface { return core.Interface{m["deleted_at"]} }
func (m CreditHistory) PaidAt() core.Interface    { return core.Interface{m["paid_at"]} }
func (m CreditHistory) Classify() core.Interface  { return core.Interface{m["classify"]} }
func (m CreditHistory) BankStatementID() core.Interface {
	return core.Interface{m["bank_statement_id"]}
}

func (m *CreditHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 11)
	args := make([]interface{}, 11)
	for i := 0; i < 11; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(CreditHistory, 11)
	res["id"] = data[0]
	res["amount"] = data[1]
	res["shop_id"] = data[2]
	res["type"] = data[3]
	res["status"] = data[4]
	res["created_at"] = data[5]
	res["updated_at"] = data[6]
	res["deleted_at"] = data[7]
	res["paid_at"] = data[8]
	res["classify"] = data[9]
	res["bank_statement_id"] = data[10]
	*m = res
	return nil
}

func (ms *CreditHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 11)
	args := make([]interface{}, 11)
	for i := 0; i < 11; i++ {
		args[i] = &data[i]
	}
	res := make(CreditHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(CreditHistory)
		m["id"] = data[0]
		m["amount"] = data[1]
		m["shop_id"] = data[2]
		m["type"] = data[3]
		m["status"] = data[4]
		m["created_at"] = data[5]
		m["updated_at"] = data[6]
		m["deleted_at"] = data[7]
		m["paid_at"] = data[8]
		m["classify"] = data[9]
		m["bank_statement_id"] = data[10]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}

type CreditExtendeds []*CreditExtended

func (m *CreditExtended) SQLTableName() string  { return "credit" }
func (m *CreditExtendeds) SQLTableName() string { return "credit" }

func (m *CreditExtended) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *CreditExtendeds) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(CreditExtendeds, 0, 128)
	for rows.Next() {
		m := new(CreditExtended)
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

func (m *CreditExtended) SQLSelect(w SQLWriter) error {
	(*CreditExtended)(nil).__sqlSelect(w)
	w.WriteByte(' ')
	(*CreditExtended)(nil).__sqlJoin(w)
	return nil
}

func (m *CreditExtendeds) SQLSelect(w SQLWriter) error {
	return (*CreditExtended)(nil).SQLSelect(w)
}

func (m *CreditExtended) SQLJoin(w SQLWriter) error {
	m.__sqlJoin(w)
	return nil
}

func (m *CreditExtendeds) SQLJoin(w SQLWriter) error {
	return (*CreditExtended)(nil).SQLJoin(w)
}

func (m *CreditExtended) __sqlSelect(w SQLWriter) {
	w.WriteRawString("SELECT ")
	core.WriteCols(w, "c", (*Credit)(nil).SQLListCols())
	w.WriteByte(',')
	core.WriteCols(w, "s", (*identitymodel.Shop)(nil).SQLListCols())
}

func (m *CreditExtended) __sqlJoin(w SQLWriter) {
	w.WriteRawString("FROM ")
	w.WriteName("credit")
	w.WriteRawString(" AS ")
	w.WriteName("c")
	w.WriteRawString(" LEFT JOIN ")
	w.WriteName((*identitymodel.Shop)(nil).SQLTableName())
	w.WriteRawString(" AS s ON")
	w.WriteQueryString(" s.id = c.shop_id")
}

func (m *CreditExtended) SQLScanArgs(opts core.Opts) []interface{} {
	args := make([]interface{}, 0, 64) // TODO: pre-calculate length
	m.Credit = new(Credit)
	args = append(args, m.Credit.SQLScanArgs(opts)...)
	m.Shop = new(identitymodel.Shop)
	args = append(args, m.Shop.SQLScanArgs(opts)...)
	return args
}
