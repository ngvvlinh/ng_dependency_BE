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

type TicketProviderWebhooks []*TicketProviderWebhook

const __sqlTicketProviderWebhook_Table = "ticket_provider_webhook"
const __sqlTicketProviderWebhook_ListCols = "\"id\",\"ticket_provider\",\"data\",\"external_status\",\"external_type\",\"client_id\",\"created_at\",\"error\",\"connection_id\""
const __sqlTicketProviderWebhook_ListColsOnConflict = "\"id\" = EXCLUDED.\"id\",\"ticket_provider\" = EXCLUDED.\"ticket_provider\",\"data\" = EXCLUDED.\"data\",\"external_status\" = EXCLUDED.\"external_status\",\"external_type\" = EXCLUDED.\"external_type\",\"client_id\" = EXCLUDED.\"client_id\",\"created_at\" = EXCLUDED.\"created_at\",\"error\" = EXCLUDED.\"error\",\"connection_id\" = EXCLUDED.\"connection_id\""
const __sqlTicketProviderWebhook_Insert = "INSERT INTO \"ticket_provider_webhook\" (" + __sqlTicketProviderWebhook_ListCols + ") VALUES"
const __sqlTicketProviderWebhook_Select = "SELECT " + __sqlTicketProviderWebhook_ListCols + " FROM \"ticket_provider_webhook\""
const __sqlTicketProviderWebhook_Select_history = "SELECT " + __sqlTicketProviderWebhook_ListCols + " FROM history.\"ticket_provider_webhook\""
const __sqlTicketProviderWebhook_UpdateAll = "UPDATE \"ticket_provider_webhook\" SET (" + __sqlTicketProviderWebhook_ListCols + ")"
const __sqlTicketProviderWebhook_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT ticket_provider_webhook_pkey DO UPDATE SET"

func (m *TicketProviderWebhook) SQLTableName() string  { return "ticket_provider_webhook" }
func (m *TicketProviderWebhooks) SQLTableName() string { return "ticket_provider_webhook" }
func (m *TicketProviderWebhook) SQLListCols() string   { return __sqlTicketProviderWebhook_ListCols }

func (m *TicketProviderWebhook) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlTicketProviderWebhook_ListCols + " FROM \"ticket_provider_webhook\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func (m *TicketProviderWebhook) Migration(db *cmsql.Database) {
	var mDBColumnNameAndType map[string]string
	if val, err := migration.GetColumnNamesAndTypes(db, "ticket_provider_webhook"); err != nil {
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
		"ticket_provider": {
			ColumnName:       "ticket_provider",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"data": {
			ColumnName:       "data",
			ColumnType:       "json.RawMessage",
			ColumnDBType:     "[]byte",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"external_status": {
			ColumnName:       "external_status",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"external_type": {
			ColumnName:       "external_type",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"client_id": {
			ColumnName:       "client_id",
			ColumnType:       "string",
			ColumnDBType:     "string",
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
		"error": {
			ColumnName:       "error",
			ColumnType:       "*etopmodel.Error",
			ColumnDBType:     "*struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"connection_id": {
			ColumnName:       "connection_id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
	}
	if err := migration.Compare(db, "ticket_provider_webhook", mModelColumnNameAndType, mDBColumnNameAndType); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*TicketProviderWebhook)(nil))
}

func (m *TicketProviderWebhook) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		m.ID,
		core.String(m.TicketProvider),
		core.JSON{m.Data},
		core.String(m.ExternalStatus),
		core.String(m.ExternalType),
		core.String(m.ClientID),
		core.Now(m.CreatedAt, now, create),
		core.JSON{m.Error},
		m.ConnectionID,
	}
}

func (m *TicketProviderWebhook) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		(*core.String)(&m.TicketProvider),
		core.JSON{&m.Data},
		(*core.String)(&m.ExternalStatus),
		(*core.String)(&m.ExternalType),
		(*core.String)(&m.ClientID),
		(*core.Time)(&m.CreatedAt),
		core.JSON{&m.Error},
		&m.ConnectionID,
	}
}

func (m *TicketProviderWebhook) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *TicketProviderWebhooks) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(TicketProviderWebhooks, 0, 128)
	for rows.Next() {
		m := new(TicketProviderWebhook)
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

func (_ *TicketProviderWebhook) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlTicketProviderWebhook_Select)
	return nil
}

func (_ *TicketProviderWebhooks) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlTicketProviderWebhook_Select)
	return nil
}

func (m *TicketProviderWebhook) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlTicketProviderWebhook_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(9)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms TicketProviderWebhooks) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlTicketProviderWebhook_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(9)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *TicketProviderWebhook) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlTicketProviderWebhook_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlTicketProviderWebhook_ListColsOnConflict)
	return nil
}

func (ms TicketProviderWebhooks) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlTicketProviderWebhook_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlTicketProviderWebhook_ListColsOnConflict)
	return nil
}

func (m *TicketProviderWebhook) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("ticket_provider_webhook")
	w.WriteRawString(" SET ")
	if m.ID != 0 {
		flag = true
		w.WriteName("id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ID)
	}
	if m.TicketProvider != "" {
		flag = true
		w.WriteName("ticket_provider")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.TicketProvider)
	}
	if m.Data != nil {
		flag = true
		w.WriteName("data")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.Data})
	}
	if m.ExternalStatus != "" {
		flag = true
		w.WriteName("external_status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ExternalStatus)
	}
	if m.ExternalType != "" {
		flag = true
		w.WriteName("external_type")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ExternalType)
	}
	if m.ClientID != "" {
		flag = true
		w.WriteName("client_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ClientID)
	}
	if !m.CreatedAt.IsZero() {
		flag = true
		w.WriteName("created_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CreatedAt)
	}
	if m.Error != nil {
		flag = true
		w.WriteName("error")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.Error})
	}
	if m.ConnectionID != 0 {
		flag = true
		w.WriteName("connection_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ConnectionID)
	}
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *TicketProviderWebhook) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlTicketProviderWebhook_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(9)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type TicketProviderWebhookHistory map[string]interface{}
type TicketProviderWebhookHistories []map[string]interface{}

func (m *TicketProviderWebhookHistory) SQLTableName() string {
	return "history.\"ticket_provider_webhook\""
}
func (m TicketProviderWebhookHistories) SQLTableName() string {
	return "history.\"ticket_provider_webhook\""
}

func (m *TicketProviderWebhookHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlTicketProviderWebhook_Select_history)
	return nil
}

func (m TicketProviderWebhookHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlTicketProviderWebhook_Select_history)
	return nil
}

func (m TicketProviderWebhookHistory) ID() core.Interface { return core.Interface{m["id"]} }
func (m TicketProviderWebhookHistory) TicketProvider() core.Interface {
	return core.Interface{m["ticket_provider"]}
}
func (m TicketProviderWebhookHistory) Data() core.Interface { return core.Interface{m["data"]} }
func (m TicketProviderWebhookHistory) ExternalStatus() core.Interface {
	return core.Interface{m["external_status"]}
}
func (m TicketProviderWebhookHistory) ExternalType() core.Interface {
	return core.Interface{m["external_type"]}
}
func (m TicketProviderWebhookHistory) ClientID() core.Interface {
	return core.Interface{m["client_id"]}
}
func (m TicketProviderWebhookHistory) CreatedAt() core.Interface {
	return core.Interface{m["created_at"]}
}
func (m TicketProviderWebhookHistory) Error() core.Interface { return core.Interface{m["error"]} }
func (m TicketProviderWebhookHistory) ConnectionID() core.Interface {
	return core.Interface{m["connection_id"]}
}

func (m *TicketProviderWebhookHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 9)
	args := make([]interface{}, 9)
	for i := 0; i < 9; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(TicketProviderWebhookHistory, 9)
	res["id"] = data[0]
	res["ticket_provider"] = data[1]
	res["data"] = data[2]
	res["external_status"] = data[3]
	res["external_type"] = data[4]
	res["client_id"] = data[5]
	res["created_at"] = data[6]
	res["error"] = data[7]
	res["connection_id"] = data[8]
	*m = res
	return nil
}

func (ms *TicketProviderWebhookHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 9)
	args := make([]interface{}, 9)
	for i := 0; i < 9; i++ {
		args[i] = &data[i]
	}
	res := make(TicketProviderWebhookHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(TicketProviderWebhookHistory)
		m["id"] = data[0]
		m["ticket_provider"] = data[1]
		m["data"] = data[2]
		m["external_status"] = data[3]
		m["external_type"] = data[4]
		m["client_id"] = data[5]
		m["created_at"] = data[6]
		m["error"] = data[7]
		m["connection_id"] = data[8]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
