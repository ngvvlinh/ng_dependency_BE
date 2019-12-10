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

// Type Payment represents table payment
func sqlgenPayment(_ *Payment) bool { return true }

type Payments []*Payment

const __sqlPayment_Table = "payment"
const __sqlPayment_ListCols = "\"id\",\"data\",\"order_id\",\"payment_provider\",\"action\",\"created_at\",\"updated_at\""
const __sqlPayment_Insert = "INSERT INTO \"payment\" (" + __sqlPayment_ListCols + ") VALUES"
const __sqlPayment_Select = "SELECT " + __sqlPayment_ListCols + " FROM \"payment\""
const __sqlPayment_Select_history = "SELECT " + __sqlPayment_ListCols + " FROM history.\"payment\""
const __sqlPayment_UpdateAll = "UPDATE \"payment\" SET (" + __sqlPayment_ListCols + ")"

func (m *Payment) SQLTableName() string  { return "payment" }
func (m *Payments) SQLTableName() string { return "payment" }
func (m *Payment) SQLListCols() string   { return __sqlPayment_ListCols }

func (m *Payment) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlPayment_ListCols + " FROM \"payment\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*Payment)(nil))
}

func (m *Payment) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		m.ID,
		core.JSON{m.Data},
		core.String(m.OrderID),
		m.PaymentProvider,
		core.String(m.Action),
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
	}
}

func (m *Payment) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		core.JSON{&m.Data},
		(*core.String)(&m.OrderID),
		&m.PaymentProvider,
		(*core.String)(&m.Action),
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
	}
}

func (m *Payment) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *Payments) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(Payments, 0, 128)
	for rows.Next() {
		m := new(Payment)
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

func (_ *Payment) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlPayment_Select)
	return nil
}

func (_ *Payments) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlPayment_Select)
	return nil
}

func (m *Payment) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlPayment_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(7)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms Payments) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlPayment_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(7)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *Payment) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("payment")
	w.WriteRawString(" SET ")
	if m.ID != 0 {
		flag = true
		w.WriteName("id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ID)
	}
	if m.Data != nil {
		flag = true
		w.WriteName("data")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.Data})
	}
	if m.OrderID != "" {
		flag = true
		w.WriteName("order_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.OrderID)
	}
	if m.PaymentProvider != 0 {
		flag = true
		w.WriteName("payment_provider")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.PaymentProvider)
	}
	if m.Action != "" {
		flag = true
		w.WriteName("action")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(string(m.Action))
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
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *Payment) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlPayment_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(7)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type PaymentHistory map[string]interface{}
type PaymentHistories []map[string]interface{}

func (m *PaymentHistory) SQLTableName() string  { return "history.\"payment\"" }
func (m PaymentHistories) SQLTableName() string { return "history.\"payment\"" }

func (m *PaymentHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlPayment_Select_history)
	return nil
}

func (m PaymentHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlPayment_Select_history)
	return nil
}

func (m PaymentHistory) ID() core.Interface              { return core.Interface{m["id"]} }
func (m PaymentHistory) Data() core.Interface            { return core.Interface{m["data"]} }
func (m PaymentHistory) OrderID() core.Interface         { return core.Interface{m["order_id"]} }
func (m PaymentHistory) PaymentProvider() core.Interface { return core.Interface{m["payment_provider"]} }
func (m PaymentHistory) Action() core.Interface          { return core.Interface{m["action"]} }
func (m PaymentHistory) CreatedAt() core.Interface       { return core.Interface{m["created_at"]} }
func (m PaymentHistory) UpdatedAt() core.Interface       { return core.Interface{m["updated_at"]} }

func (m *PaymentHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 7)
	args := make([]interface{}, 7)
	for i := 0; i < 7; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(PaymentHistory, 7)
	res["id"] = data[0]
	res["data"] = data[1]
	res["order_id"] = data[2]
	res["payment_provider"] = data[3]
	res["action"] = data[4]
	res["created_at"] = data[5]
	res["updated_at"] = data[6]
	*m = res
	return nil
}

func (ms *PaymentHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 7)
	args := make([]interface{}, 7)
	for i := 0; i < 7; i++ {
		args[i] = &data[i]
	}
	res := make(PaymentHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(PaymentHistory)
		m["id"] = data[0]
		m["data"] = data[1]
		m["order_id"] = data[2]
		m["payment_provider"] = data[3]
		m["action"] = data[4]
		m["created_at"] = data[5]
		m["updated_at"] = data[6]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
