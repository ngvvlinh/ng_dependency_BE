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

type ShopCustomers []*ShopCustomer

const __sqlShopCustomer_Table = "shop_customer"
const __sqlShopCustomer_ListCols = "\"external_id\",\"external_code\",\"partner_id\",\"id\",\"shop_id\",\"code\",\"full_name\",\"gender\",\"type\",\"birthday\",\"note\",\"phone\",\"email\",\"status\",\"created_at\",\"updated_at\",\"rid\""
const __sqlShopCustomer_ListColsOnConflict = "\"external_id\" = EXCLUDED.\"external_id\",\"external_code\" = EXCLUDED.\"external_code\",\"partner_id\" = EXCLUDED.\"partner_id\",\"id\" = EXCLUDED.\"id\",\"shop_id\" = EXCLUDED.\"shop_id\",\"code\" = EXCLUDED.\"code\",\"full_name\" = EXCLUDED.\"full_name\",\"gender\" = EXCLUDED.\"gender\",\"type\" = EXCLUDED.\"type\",\"birthday\" = EXCLUDED.\"birthday\",\"note\" = EXCLUDED.\"note\",\"phone\" = EXCLUDED.\"phone\",\"email\" = EXCLUDED.\"email\",\"status\" = EXCLUDED.\"status\",\"created_at\" = EXCLUDED.\"created_at\",\"updated_at\" = EXCLUDED.\"updated_at\",\"rid\" = EXCLUDED.\"rid\""
const __sqlShopCustomer_Insert = "INSERT INTO \"shop_customer\" (" + __sqlShopCustomer_ListCols + ") VALUES"
const __sqlShopCustomer_Select = "SELECT " + __sqlShopCustomer_ListCols + " FROM \"shop_customer\""
const __sqlShopCustomer_Select_history = "SELECT " + __sqlShopCustomer_ListCols + " FROM history.\"shop_customer\""
const __sqlShopCustomer_UpdateAll = "UPDATE \"shop_customer\" SET (" + __sqlShopCustomer_ListCols + ")"
const __sqlShopCustomer_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT shop_customer_pkey DO UPDATE SET"

func (m *ShopCustomer) SQLTableName() string  { return "shop_customer" }
func (m *ShopCustomers) SQLTableName() string { return "shop_customer" }
func (m *ShopCustomer) SQLListCols() string   { return __sqlShopCustomer_ListCols }

func (m *ShopCustomer) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlShopCustomer_ListCols + " FROM \"shop_customer\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func (m *ShopCustomer) Migration(db *cmsql.Database) {
	var mDBColumnNameAndType map[string]string
	if val, err := migration.GetColumnNamesAndTypes(db, "shop_customer"); err != nil {
		db.RecordError(err)
		return
	} else {
		mDBColumnNameAndType = val
	}
	mModelColumnNameAndType := map[string]migration.ColumnDef{
		"external_id": {
			ColumnName:       "external_id",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"external_code": {
			ColumnName:       "external_code",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"partner_id": {
			ColumnName:       "partner_id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
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
		"code": {
			ColumnName:       "code",
			ColumnType:       "string",
			ColumnDBType:     "string",
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
		"gender": {
			ColumnName:       "gender",
			ColumnType:       "gender.Gender",
			ColumnDBType:     "enum",
			ColumnTag:        "enum(gender_type)",
			ColumnEnumValues: []string{"unknown", "male", "female", "other"},
		},
		"type": {
			ColumnName:       "type",
			ColumnType:       "customer_type.CustomerType",
			ColumnDBType:     "enum",
			ColumnTag:        "enum(customer_type)",
			ColumnEnumValues: []string{"unknown", "individual", "organization", "anonymous", "independent"},
		},
		"birthday": {
			ColumnName:       "birthday",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "date",
			ColumnEnumValues: []string{},
		},
		"note": {
			ColumnName:       "note",
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
		"email": {
			ColumnName:       "email",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"status": {
			ColumnName:       "status",
			ColumnType:       "int",
			ColumnDBType:     "int",
			ColumnTag:        "int2",
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
		"rid": {
			ColumnName:       "rid",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
	}
	if err := migration.Compare(db, "shop_customer", mModelColumnNameAndType, mDBColumnNameAndType); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*ShopCustomer)(nil))
}

func (m *ShopCustomer) SQLArgs(opts core.Opts, create bool) []interface{} {
	return []interface{}{
		core.String(m.ExternalID),
		core.String(m.ExternalCode),
		m.PartnerID,
		m.ID,
		m.ShopID,
		core.String(m.Code),
		core.String(m.FullName),
		m.Gender,
		m.Type,
		core.String(m.Birthday),
		core.String(m.Note),
		core.String(m.Phone),
		core.String(m.Email),
		core.Int(m.Status),
		core.Time(m.CreatedAt),
		core.Time(m.UpdatedAt),
		m.Rid,
	}
}

func (m *ShopCustomer) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		(*core.String)(&m.ExternalID),
		(*core.String)(&m.ExternalCode),
		&m.PartnerID,
		&m.ID,
		&m.ShopID,
		(*core.String)(&m.Code),
		(*core.String)(&m.FullName),
		&m.Gender,
		&m.Type,
		(*core.String)(&m.Birthday),
		(*core.String)(&m.Note),
		(*core.String)(&m.Phone),
		(*core.String)(&m.Email),
		(*core.Int)(&m.Status),
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		&m.Rid,
	}
}

func (m *ShopCustomer) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *ShopCustomers) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(ShopCustomers, 0, 128)
	for rows.Next() {
		m := new(ShopCustomer)
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

func (_ *ShopCustomer) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCustomer_Select)
	return nil
}

func (_ *ShopCustomers) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCustomer_Select)
	return nil
}

func (m *ShopCustomer) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCustomer_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(17)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms ShopCustomers) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCustomer_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(17)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *ShopCustomer) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlShopCustomer_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlShopCustomer_ListColsOnConflict)
	return nil
}

func (ms ShopCustomers) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlShopCustomer_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlShopCustomer_ListColsOnConflict)
	return nil
}

func (m *ShopCustomer) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("shop_customer")
	w.WriteRawString(" SET ")
	if m.ExternalID != "" {
		flag = true
		w.WriteName("external_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ExternalID)
	}
	if m.ExternalCode != "" {
		flag = true
		w.WriteName("external_code")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ExternalCode)
	}
	if m.PartnerID != 0 {
		flag = true
		w.WriteName("partner_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.PartnerID)
	}
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
	if m.Code != "" {
		flag = true
		w.WriteName("code")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Code)
	}
	if m.FullName != "" {
		flag = true
		w.WriteName("full_name")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.FullName)
	}
	if m.Gender != 0 {
		flag = true
		w.WriteName("gender")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Gender)
	}
	if m.Type != 0 {
		flag = true
		w.WriteName("type")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Type)
	}
	if m.Birthday != "" {
		flag = true
		w.WriteName("birthday")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Birthday)
	}
	if m.Note != "" {
		flag = true
		w.WriteName("note")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Note)
	}
	if m.Phone != "" {
		flag = true
		w.WriteName("phone")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Phone)
	}
	if m.Email != "" {
		flag = true
		w.WriteName("email")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Email)
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
		w.WriteArg(m.UpdatedAt)
	}
	if m.Rid != 0 {
		flag = true
		w.WriteName("rid")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Rid)
	}
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *ShopCustomer) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCustomer_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(17)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type ShopCustomerHistory map[string]interface{}
type ShopCustomerHistories []map[string]interface{}

func (m *ShopCustomerHistory) SQLTableName() string  { return "history.\"shop_customer\"" }
func (m ShopCustomerHistories) SQLTableName() string { return "history.\"shop_customer\"" }

func (m *ShopCustomerHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCustomer_Select_history)
	return nil
}

func (m ShopCustomerHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCustomer_Select_history)
	return nil
}

func (m ShopCustomerHistory) ExternalID() core.Interface   { return core.Interface{m["external_id"]} }
func (m ShopCustomerHistory) ExternalCode() core.Interface { return core.Interface{m["external_code"]} }
func (m ShopCustomerHistory) PartnerID() core.Interface    { return core.Interface{m["partner_id"]} }
func (m ShopCustomerHistory) ID() core.Interface           { return core.Interface{m["id"]} }
func (m ShopCustomerHistory) ShopID() core.Interface       { return core.Interface{m["shop_id"]} }
func (m ShopCustomerHistory) Code() core.Interface         { return core.Interface{m["code"]} }
func (m ShopCustomerHistory) FullName() core.Interface     { return core.Interface{m["full_name"]} }
func (m ShopCustomerHistory) Gender() core.Interface       { return core.Interface{m["gender"]} }
func (m ShopCustomerHistory) Type() core.Interface         { return core.Interface{m["type"]} }
func (m ShopCustomerHistory) Birthday() core.Interface     { return core.Interface{m["birthday"]} }
func (m ShopCustomerHistory) Note() core.Interface         { return core.Interface{m["note"]} }
func (m ShopCustomerHistory) Phone() core.Interface        { return core.Interface{m["phone"]} }
func (m ShopCustomerHistory) Email() core.Interface        { return core.Interface{m["email"]} }
func (m ShopCustomerHistory) Status() core.Interface       { return core.Interface{m["status"]} }
func (m ShopCustomerHistory) CreatedAt() core.Interface    { return core.Interface{m["created_at"]} }
func (m ShopCustomerHistory) UpdatedAt() core.Interface    { return core.Interface{m["updated_at"]} }
func (m ShopCustomerHistory) Rid() core.Interface          { return core.Interface{m["rid"]} }

func (m *ShopCustomerHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 17)
	args := make([]interface{}, 17)
	for i := 0; i < 17; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(ShopCustomerHistory, 17)
	res["external_id"] = data[0]
	res["external_code"] = data[1]
	res["partner_id"] = data[2]
	res["id"] = data[3]
	res["shop_id"] = data[4]
	res["code"] = data[5]
	res["full_name"] = data[6]
	res["gender"] = data[7]
	res["type"] = data[8]
	res["birthday"] = data[9]
	res["note"] = data[10]
	res["phone"] = data[11]
	res["email"] = data[12]
	res["status"] = data[13]
	res["created_at"] = data[14]
	res["updated_at"] = data[15]
	res["rid"] = data[16]
	*m = res
	return nil
}

func (ms *ShopCustomerHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 17)
	args := make([]interface{}, 17)
	for i := 0; i < 17; i++ {
		args[i] = &data[i]
	}
	res := make(ShopCustomerHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(ShopCustomerHistory)
		m["external_id"] = data[0]
		m["external_code"] = data[1]
		m["partner_id"] = data[2]
		m["id"] = data[3]
		m["shop_id"] = data[4]
		m["code"] = data[5]
		m["full_name"] = data[6]
		m["gender"] = data[7]
		m["type"] = data[8]
		m["birthday"] = data[9]
		m["note"] = data[10]
		m["phone"] = data[11]
		m["email"] = data[12]
		m["status"] = data[13]
		m["created_at"] = data[14]
		m["updated_at"] = data[15]
		m["rid"] = data[16]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
