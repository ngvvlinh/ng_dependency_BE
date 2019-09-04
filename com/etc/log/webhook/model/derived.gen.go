// Code generated by goderive DO NOT EDIT.

package model

import (
	"database/sql"
	"time"

	core "etop.vn/backend/pkg/common/sq/core"
)

type SQLWriter = core.SQLWriter

// Type ShippingProviderWebhook represents table shipping_provider_webhook
func sqlgenShippingProviderWebhook(_ *ShippingProviderWebhook) bool { return true }

type ShippingProviderWebhooks []*ShippingProviderWebhook

const __sqlShippingProviderWebhook_Table = "shipping_provider_webhook"
const __sqlShippingProviderWebhook_ListCols = "\"id\",\"shipping_provider\",\"data\",\"shipping_code\",\"shipping_state\",\"external_shipping_state\",\"external_shipping_sub_state\",\"created_at\",\"updated_at\""
const __sqlShippingProviderWebhook_Insert = "INSERT INTO \"shipping_provider_webhook\" (" + __sqlShippingProviderWebhook_ListCols + ") VALUES"
const __sqlShippingProviderWebhook_Select = "SELECT " + __sqlShippingProviderWebhook_ListCols + " FROM \"shipping_provider_webhook\""
const __sqlShippingProviderWebhook_Select_history = "SELECT " + __sqlShippingProviderWebhook_ListCols + " FROM history.\"shipping_provider_webhook\""
const __sqlShippingProviderWebhook_UpdateAll = "UPDATE \"shipping_provider_webhook\" SET (" + __sqlShippingProviderWebhook_ListCols + ")"

func (m *ShippingProviderWebhook) SQLTableName() string  { return "shipping_provider_webhook" }
func (m *ShippingProviderWebhooks) SQLTableName() string { return "shipping_provider_webhook" }
func (m *ShippingProviderWebhook) SQLListCols() string   { return __sqlShippingProviderWebhook_ListCols }

func (m *ShippingProviderWebhook) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		core.Int64(m.ID),
		core.String(m.ShippingProvider),
		core.JSON{m.Data},
		core.String(m.ShippingCode),
		core.String(m.ShippingState),
		core.String(m.ExternalShippingState),
		core.String(m.ExternalShippingSubState),
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
	}
}

func (m *ShippingProviderWebhook) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		(*core.Int64)(&m.ID),
		(*core.String)(&m.ShippingProvider),
		core.JSON{&m.Data},
		(*core.String)(&m.ShippingCode),
		(*core.String)(&m.ShippingState),
		(*core.String)(&m.ExternalShippingState),
		(*core.String)(&m.ExternalShippingSubState),
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
	}
}

func (m *ShippingProviderWebhook) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *ShippingProviderWebhooks) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(ShippingProviderWebhooks, 0, 128)
	for rows.Next() {
		m := new(ShippingProviderWebhook)
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

func (_ *ShippingProviderWebhook) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShippingProviderWebhook_Select)
	return nil
}

func (_ *ShippingProviderWebhooks) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShippingProviderWebhook_Select)
	return nil
}

func (m *ShippingProviderWebhook) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShippingProviderWebhook_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(9)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms ShippingProviderWebhooks) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShippingProviderWebhook_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(9)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *ShippingProviderWebhook) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("shipping_provider_webhook")
	w.WriteRawString(" SET ")
	if m.ID != 0 {
		flag = true
		w.WriteName("id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ID)
	}
	if m.ShippingProvider != "" {
		flag = true
		w.WriteName("shipping_provider")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShippingProvider)
	}
	if m.Data != nil {
		flag = true
		w.WriteName("data")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.Data})
	}
	if m.ShippingCode != "" {
		flag = true
		w.WriteName("shipping_code")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShippingCode)
	}
	if m.ShippingState != "" {
		flag = true
		w.WriteName("shipping_state")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShippingState)
	}
	if m.ExternalShippingState != "" {
		flag = true
		w.WriteName("external_shipping_state")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ExternalShippingState)
	}
	if m.ExternalShippingSubState != "" {
		flag = true
		w.WriteName("external_shipping_sub_state")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ExternalShippingSubState)
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

func (m *ShippingProviderWebhook) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlShippingProviderWebhook_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(9)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type ShippingProviderWebhookHistory map[string]interface{}
type ShippingProviderWebhookHistories []map[string]interface{}

func (m *ShippingProviderWebhookHistory) SQLTableName() string {
	return "history.\"shipping_provider_webhook\""
}
func (m ShippingProviderWebhookHistories) SQLTableName() string {
	return "history.\"shipping_provider_webhook\""
}

func (m *ShippingProviderWebhookHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShippingProviderWebhook_Select_history)
	return nil
}

func (m ShippingProviderWebhookHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShippingProviderWebhook_Select_history)
	return nil
}

func (m ShippingProviderWebhookHistory) ID() core.Interface { return core.Interface{m["id"]} }
func (m ShippingProviderWebhookHistory) ShippingProvider() core.Interface {
	return core.Interface{m["shipping_provider"]}
}
func (m ShippingProviderWebhookHistory) Data() core.Interface { return core.Interface{m["data"]} }
func (m ShippingProviderWebhookHistory) ShippingCode() core.Interface {
	return core.Interface{m["shipping_code"]}
}
func (m ShippingProviderWebhookHistory) ShippingState() core.Interface {
	return core.Interface{m["shipping_state"]}
}
func (m ShippingProviderWebhookHistory) ExternalShippingState() core.Interface {
	return core.Interface{m["external_shipping_state"]}
}
func (m ShippingProviderWebhookHistory) ExternalShippingSubState() core.Interface {
	return core.Interface{m["external_shipping_sub_state"]}
}
func (m ShippingProviderWebhookHistory) CreatedAt() core.Interface {
	return core.Interface{m["created_at"]}
}
func (m ShippingProviderWebhookHistory) UpdatedAt() core.Interface {
	return core.Interface{m["updated_at"]}
}

func (m *ShippingProviderWebhookHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 9)
	args := make([]interface{}, 9)
	for i := 0; i < 9; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(ShippingProviderWebhookHistory, 9)
	res["id"] = data[0]
	res["shipping_provider"] = data[1]
	res["data"] = data[2]
	res["shipping_code"] = data[3]
	res["shipping_state"] = data[4]
	res["external_shipping_state"] = data[5]
	res["external_shipping_sub_state"] = data[6]
	res["created_at"] = data[7]
	res["updated_at"] = data[8]
	*m = res
	return nil
}

func (ms *ShippingProviderWebhookHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 9)
	args := make([]interface{}, 9)
	for i := 0; i < 9; i++ {
		args[i] = &data[i]
	}
	res := make(ShippingProviderWebhookHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(ShippingProviderWebhookHistory)
		m["id"] = data[0]
		m["shipping_provider"] = data[1]
		m["data"] = data[2]
		m["shipping_code"] = data[3]
		m["shipping_state"] = data[4]
		m["external_shipping_state"] = data[5]
		m["external_shipping_sub_state"] = data[6]
		m["created_at"] = data[7]
		m["updated_at"] = data[8]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
