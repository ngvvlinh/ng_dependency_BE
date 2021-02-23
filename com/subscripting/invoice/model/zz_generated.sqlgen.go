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

type Invoices []*Invoice

const __sqlInvoice_Table = "invoice"
const __sqlInvoice_ListCols = "\"id\",\"account_id\",\"total_amount\",\"description\",\"payment_id\",\"payment_status\",\"status\",\"customer\",\"created_at\",\"updated_at\",\"deleted_at\",\"wl_partner_id\",\"referral_type\",\"referral_ids\""
const __sqlInvoice_ListColsOnConflict = "\"id\" = EXCLUDED.\"id\",\"account_id\" = EXCLUDED.\"account_id\",\"total_amount\" = EXCLUDED.\"total_amount\",\"description\" = EXCLUDED.\"description\",\"payment_id\" = EXCLUDED.\"payment_id\",\"payment_status\" = EXCLUDED.\"payment_status\",\"status\" = EXCLUDED.\"status\",\"customer\" = EXCLUDED.\"customer\",\"created_at\" = EXCLUDED.\"created_at\",\"updated_at\" = EXCLUDED.\"updated_at\",\"deleted_at\" = EXCLUDED.\"deleted_at\",\"wl_partner_id\" = EXCLUDED.\"wl_partner_id\",\"referral_type\" = EXCLUDED.\"referral_type\",\"referral_ids\" = EXCLUDED.\"referral_ids\""
const __sqlInvoice_Insert = "INSERT INTO \"invoice\" (" + __sqlInvoice_ListCols + ") VALUES"
const __sqlInvoice_Select = "SELECT " + __sqlInvoice_ListCols + " FROM \"invoice\""
const __sqlInvoice_Select_history = "SELECT " + __sqlInvoice_ListCols + " FROM history.\"invoice\""
const __sqlInvoice_UpdateAll = "UPDATE \"invoice\" SET (" + __sqlInvoice_ListCols + ")"
const __sqlInvoice_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT invoice_pkey DO UPDATE SET"

func (m *Invoice) SQLTableName() string  { return "invoice" }
func (m *Invoices) SQLTableName() string { return "invoice" }
func (m *Invoice) SQLListCols() string   { return __sqlInvoice_ListCols }

func (m *Invoice) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlInvoice_ListCols + " FROM \"invoice\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func (m *Invoice) Migration(db *cmsql.Database) {
	var mDBColumnNameAndType map[string]string
	if val, err := migration.GetColumnNamesAndTypes(db, "invoice"); err != nil {
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
		"account_id": {
			ColumnName:       "account_id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"total_amount": {
			ColumnName:       "total_amount",
			ColumnType:       "int",
			ColumnDBType:     "int",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"description": {
			ColumnName:       "description",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"payment_id": {
			ColumnName:       "payment_id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"payment_status": {
			ColumnName:       "payment_status",
			ColumnType:       "status4.Status",
			ColumnDBType:     "enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{"Z", "P", "S", "N"},
		},
		"status": {
			ColumnName:       "status",
			ColumnType:       "status4.Status",
			ColumnDBType:     "enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{"Z", "P", "S", "N"},
		},
		"customer": {
			ColumnName:       "customer",
			ColumnType:       "*sharemodel.CustomerInfo",
			ColumnDBType:     "*struct",
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
		"wl_partner_id": {
			ColumnName:       "wl_partner_id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"referral_type": {
			ColumnName:       "referral_type",
			ColumnType:       "subject_referral.SubjectReferral",
			ColumnDBType:     "enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{"credit", "invoice", "subscription"},
		},
		"referral_ids": {
			ColumnName:       "referral_ids",
			ColumnType:       "[]dot.ID",
			ColumnDBType:     "[]int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
	}
	if err := migration.Compare(db, "invoice", mModelColumnNameAndType, mDBColumnNameAndType); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*Invoice)(nil))
}

func (m *Invoice) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		m.ID,
		m.AccountID,
		core.Int(m.TotalAmount),
		core.String(m.Description),
		m.PaymentID,
		m.PaymentStatus,
		m.Status,
		core.JSON{m.Customer},
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, create),
		core.Time(m.DeletedAt),
		m.WLPartnerID,
		m.ReferralType,
		core.Array{m.ReferralIDs, opts},
	}
}

func (m *Invoice) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		&m.AccountID,
		(*core.Int)(&m.TotalAmount),
		(*core.String)(&m.Description),
		&m.PaymentID,
		&m.PaymentStatus,
		&m.Status,
		core.JSON{&m.Customer},
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		(*core.Time)(&m.DeletedAt),
		&m.WLPartnerID,
		&m.ReferralType,
		core.Array{&m.ReferralIDs, opts},
	}
}

func (m *Invoice) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *Invoices) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(Invoices, 0, 128)
	for rows.Next() {
		m := new(Invoice)
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

func (_ *Invoice) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlInvoice_Select)
	return nil
}

func (_ *Invoices) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlInvoice_Select)
	return nil
}

func (m *Invoice) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlInvoice_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(14)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms Invoices) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlInvoice_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(14)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *Invoice) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlInvoice_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlInvoice_ListColsOnConflict)
	return nil
}

func (ms Invoices) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlInvoice_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlInvoice_ListColsOnConflict)
	return nil
}

func (m *Invoice) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("invoice")
	w.WriteRawString(" SET ")
	if m.ID != 0 {
		flag = true
		w.WriteName("id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ID)
	}
	if m.AccountID != 0 {
		flag = true
		w.WriteName("account_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.AccountID)
	}
	if m.TotalAmount != 0 {
		flag = true
		w.WriteName("total_amount")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.TotalAmount)
	}
	if m.Description != "" {
		flag = true
		w.WriteName("description")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Description)
	}
	if m.PaymentID != 0 {
		flag = true
		w.WriteName("payment_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.PaymentID)
	}
	if m.PaymentStatus != 0 {
		flag = true
		w.WriteName("payment_status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.PaymentStatus)
	}
	if m.Status != 0 {
		flag = true
		w.WriteName("status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Status)
	}
	if m.Customer != nil {
		flag = true
		w.WriteName("customer")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.Customer})
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
	if !m.DeletedAt.IsZero() {
		flag = true
		w.WriteName("deleted_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.DeletedAt)
	}
	if m.WLPartnerID != 0 {
		flag = true
		w.WriteName("wl_partner_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.WLPartnerID)
	}
	if m.ReferralType != 0 {
		flag = true
		w.WriteName("referral_type")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ReferralType)
	}
	if m.ReferralIDs != nil {
		flag = true
		w.WriteName("referral_ids")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.ReferralIDs, opts})
	}
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *Invoice) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlInvoice_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(14)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type InvoiceHistory map[string]interface{}
type InvoiceHistories []map[string]interface{}

func (m *InvoiceHistory) SQLTableName() string  { return "history.\"invoice\"" }
func (m InvoiceHistories) SQLTableName() string { return "history.\"invoice\"" }

func (m *InvoiceHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlInvoice_Select_history)
	return nil
}

func (m InvoiceHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlInvoice_Select_history)
	return nil
}

func (m InvoiceHistory) ID() core.Interface            { return core.Interface{m["id"]} }
func (m InvoiceHistory) AccountID() core.Interface     { return core.Interface{m["account_id"]} }
func (m InvoiceHistory) TotalAmount() core.Interface   { return core.Interface{m["total_amount"]} }
func (m InvoiceHistory) Description() core.Interface   { return core.Interface{m["description"]} }
func (m InvoiceHistory) PaymentID() core.Interface     { return core.Interface{m["payment_id"]} }
func (m InvoiceHistory) PaymentStatus() core.Interface { return core.Interface{m["payment_status"]} }
func (m InvoiceHistory) Status() core.Interface        { return core.Interface{m["status"]} }
func (m InvoiceHistory) Customer() core.Interface      { return core.Interface{m["customer"]} }
func (m InvoiceHistory) CreatedAt() core.Interface     { return core.Interface{m["created_at"]} }
func (m InvoiceHistory) UpdatedAt() core.Interface     { return core.Interface{m["updated_at"]} }
func (m InvoiceHistory) DeletedAt() core.Interface     { return core.Interface{m["deleted_at"]} }
func (m InvoiceHistory) WLPartnerID() core.Interface   { return core.Interface{m["wl_partner_id"]} }
func (m InvoiceHistory) ReferralType() core.Interface  { return core.Interface{m["referral_type"]} }
func (m InvoiceHistory) ReferralIDs() core.Interface   { return core.Interface{m["referral_ids"]} }

func (m *InvoiceHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 14)
	args := make([]interface{}, 14)
	for i := 0; i < 14; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(InvoiceHistory, 14)
	res["id"] = data[0]
	res["account_id"] = data[1]
	res["total_amount"] = data[2]
	res["description"] = data[3]
	res["payment_id"] = data[4]
	res["payment_status"] = data[5]
	res["status"] = data[6]
	res["customer"] = data[7]
	res["created_at"] = data[8]
	res["updated_at"] = data[9]
	res["deleted_at"] = data[10]
	res["wl_partner_id"] = data[11]
	res["referral_type"] = data[12]
	res["referral_ids"] = data[13]
	*m = res
	return nil
}

func (ms *InvoiceHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 14)
	args := make([]interface{}, 14)
	for i := 0; i < 14; i++ {
		args[i] = &data[i]
	}
	res := make(InvoiceHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(InvoiceHistory)
		m["id"] = data[0]
		m["account_id"] = data[1]
		m["total_amount"] = data[2]
		m["description"] = data[3]
		m["payment_id"] = data[4]
		m["payment_status"] = data[5]
		m["status"] = data[6]
		m["customer"] = data[7]
		m["created_at"] = data[8]
		m["updated_at"] = data[9]
		m["deleted_at"] = data[10]
		m["wl_partner_id"] = data[11]
		m["referral_type"] = data[12]
		m["referral_ids"] = data[13]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}

type InvoiceLines []*InvoiceLine

const __sqlInvoiceLine_Table = "invoice_line"
const __sqlInvoiceLine_ListCols = "\"id\",\"line_amount\",\"price\",\"quantity\",\"description\",\"invoice_id\",\"referral_type\",\"referral_id\",\"created_at\",\"updated_at\""
const __sqlInvoiceLine_ListColsOnConflict = "\"id\" = EXCLUDED.\"id\",\"line_amount\" = EXCLUDED.\"line_amount\",\"price\" = EXCLUDED.\"price\",\"quantity\" = EXCLUDED.\"quantity\",\"description\" = EXCLUDED.\"description\",\"invoice_id\" = EXCLUDED.\"invoice_id\",\"referral_type\" = EXCLUDED.\"referral_type\",\"referral_id\" = EXCLUDED.\"referral_id\",\"created_at\" = EXCLUDED.\"created_at\",\"updated_at\" = EXCLUDED.\"updated_at\""
const __sqlInvoiceLine_Insert = "INSERT INTO \"invoice_line\" (" + __sqlInvoiceLine_ListCols + ") VALUES"
const __sqlInvoiceLine_Select = "SELECT " + __sqlInvoiceLine_ListCols + " FROM \"invoice_line\""
const __sqlInvoiceLine_Select_history = "SELECT " + __sqlInvoiceLine_ListCols + " FROM history.\"invoice_line\""
const __sqlInvoiceLine_UpdateAll = "UPDATE \"invoice_line\" SET (" + __sqlInvoiceLine_ListCols + ")"
const __sqlInvoiceLine_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT invoice_line_pkey DO UPDATE SET"

func (m *InvoiceLine) SQLTableName() string  { return "invoice_line" }
func (m *InvoiceLines) SQLTableName() string { return "invoice_line" }
func (m *InvoiceLine) SQLListCols() string   { return __sqlInvoiceLine_ListCols }

func (m *InvoiceLine) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlInvoiceLine_ListCols + " FROM \"invoice_line\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func (m *InvoiceLine) Migration(db *cmsql.Database) {
	var mDBColumnNameAndType map[string]string
	if val, err := migration.GetColumnNamesAndTypes(db, "invoice_line"); err != nil {
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
		"line_amount": {
			ColumnName:       "line_amount",
			ColumnType:       "int",
			ColumnDBType:     "int",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"price": {
			ColumnName:       "price",
			ColumnType:       "int",
			ColumnDBType:     "int",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"quantity": {
			ColumnName:       "quantity",
			ColumnType:       "int",
			ColumnDBType:     "int",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"description": {
			ColumnName:       "description",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"invoice_id": {
			ColumnName:       "invoice_id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"referral_type": {
			ColumnName:       "referral_type",
			ColumnType:       "subject_referral.SubjectReferral",
			ColumnDBType:     "enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{"credit", "invoice", "subscription"},
		},
		"referral_id": {
			ColumnName:       "referral_id",
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
	}
	if err := migration.Compare(db, "invoice_line", mModelColumnNameAndType, mDBColumnNameAndType); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*InvoiceLine)(nil))
}

func (m *InvoiceLine) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		m.ID,
		core.Int(m.LineAmount),
		core.Int(m.Price),
		core.Int(m.Quantity),
		core.String(m.Description),
		m.InvoiceID,
		m.ReferralType,
		m.ReferralID,
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
	}
}

func (m *InvoiceLine) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		(*core.Int)(&m.LineAmount),
		(*core.Int)(&m.Price),
		(*core.Int)(&m.Quantity),
		(*core.String)(&m.Description),
		&m.InvoiceID,
		&m.ReferralType,
		&m.ReferralID,
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
	}
}

func (m *InvoiceLine) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *InvoiceLines) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(InvoiceLines, 0, 128)
	for rows.Next() {
		m := new(InvoiceLine)
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

func (_ *InvoiceLine) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlInvoiceLine_Select)
	return nil
}

func (_ *InvoiceLines) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlInvoiceLine_Select)
	return nil
}

func (m *InvoiceLine) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlInvoiceLine_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(10)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms InvoiceLines) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlInvoiceLine_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(10)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *InvoiceLine) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlInvoiceLine_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlInvoiceLine_ListColsOnConflict)
	return nil
}

func (ms InvoiceLines) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlInvoiceLine_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlInvoiceLine_ListColsOnConflict)
	return nil
}

func (m *InvoiceLine) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("invoice_line")
	w.WriteRawString(" SET ")
	if m.ID != 0 {
		flag = true
		w.WriteName("id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ID)
	}
	if m.LineAmount != 0 {
		flag = true
		w.WriteName("line_amount")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.LineAmount)
	}
	if m.Price != 0 {
		flag = true
		w.WriteName("price")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Price)
	}
	if m.Quantity != 0 {
		flag = true
		w.WriteName("quantity")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Quantity)
	}
	if m.Description != "" {
		flag = true
		w.WriteName("description")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Description)
	}
	if m.InvoiceID != 0 {
		flag = true
		w.WriteName("invoice_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.InvoiceID)
	}
	if m.ReferralType != 0 {
		flag = true
		w.WriteName("referral_type")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ReferralType)
	}
	if m.ReferralID != 0 {
		flag = true
		w.WriteName("referral_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ReferralID)
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
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *InvoiceLine) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlInvoiceLine_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(10)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type InvoiceLineHistory map[string]interface{}
type InvoiceLineHistories []map[string]interface{}

func (m *InvoiceLineHistory) SQLTableName() string  { return "history.\"invoice_line\"" }
func (m InvoiceLineHistories) SQLTableName() string { return "history.\"invoice_line\"" }

func (m *InvoiceLineHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlInvoiceLine_Select_history)
	return nil
}

func (m InvoiceLineHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlInvoiceLine_Select_history)
	return nil
}

func (m InvoiceLineHistory) ID() core.Interface           { return core.Interface{m["id"]} }
func (m InvoiceLineHistory) LineAmount() core.Interface   { return core.Interface{m["line_amount"]} }
func (m InvoiceLineHistory) Price() core.Interface        { return core.Interface{m["price"]} }
func (m InvoiceLineHistory) Quantity() core.Interface     { return core.Interface{m["quantity"]} }
func (m InvoiceLineHistory) Description() core.Interface  { return core.Interface{m["description"]} }
func (m InvoiceLineHistory) InvoiceID() core.Interface    { return core.Interface{m["invoice_id"]} }
func (m InvoiceLineHistory) ReferralType() core.Interface { return core.Interface{m["referral_type"]} }
func (m InvoiceLineHistory) ReferralID() core.Interface   { return core.Interface{m["referral_id"]} }
func (m InvoiceLineHistory) CreatedAt() core.Interface    { return core.Interface{m["created_at"]} }
func (m InvoiceLineHistory) UpdatedAt() core.Interface    { return core.Interface{m["updated_at"]} }

func (m *InvoiceLineHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 10)
	args := make([]interface{}, 10)
	for i := 0; i < 10; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(InvoiceLineHistory, 10)
	res["id"] = data[0]
	res["line_amount"] = data[1]
	res["price"] = data[2]
	res["quantity"] = data[3]
	res["description"] = data[4]
	res["invoice_id"] = data[5]
	res["referral_type"] = data[6]
	res["referral_id"] = data[7]
	res["created_at"] = data[8]
	res["updated_at"] = data[9]
	*m = res
	return nil
}

func (ms *InvoiceLineHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 10)
	args := make([]interface{}, 10)
	for i := 0; i < 10; i++ {
		args[i] = &data[i]
	}
	res := make(InvoiceLineHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(InvoiceLineHistory)
		m["id"] = data[0]
		m["line_amount"] = data[1]
		m["price"] = data[2]
		m["quantity"] = data[3]
		m["description"] = data[4]
		m["invoice_id"] = data[5]
		m["referral_type"] = data[6]
		m["referral_id"] = data[7]
		m["created_at"] = data[8]
		m["updated_at"] = data[9]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
