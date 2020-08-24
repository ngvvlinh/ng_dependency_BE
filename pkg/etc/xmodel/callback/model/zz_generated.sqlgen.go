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

type Callbacks []*Callback

const __sqlCallback_Table = "callback"
const __sqlCallback_ListCols = "\"id\",\"webhook_id\",\"account_id\",\"created_at\",\"changes\",\"result\""
const __sqlCallback_ListColsOnConflict = "\"id\" = EXCLUDED.\"id\",\"webhook_id\" = EXCLUDED.\"webhook_id\",\"account_id\" = EXCLUDED.\"account_id\",\"created_at\" = EXCLUDED.\"created_at\",\"changes\" = EXCLUDED.\"changes\",\"result\" = EXCLUDED.\"result\""
const __sqlCallback_Insert = "INSERT INTO \"callback\" (" + __sqlCallback_ListCols + ") VALUES"
const __sqlCallback_Select = "SELECT " + __sqlCallback_ListCols + " FROM \"callback\""
const __sqlCallback_Select_history = "SELECT " + __sqlCallback_ListCols + " FROM history.\"callback\""
const __sqlCallback_UpdateAll = "UPDATE \"callback\" SET (" + __sqlCallback_ListCols + ")"
const __sqlCallback_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT callback_pkey DO UPDATE SET"

func (m *Callback) SQLTableName() string  { return "callback" }
func (m *Callbacks) SQLTableName() string { return "callback" }
func (m *Callback) SQLListCols() string   { return __sqlCallback_ListCols }

func (m *Callback) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlCallback_ListCols + " FROM \"callback\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func (m *Callback) Migration(db *cmsql.Database) {
	var mDBColumnNameAndType map[string]string
	if val, err := migration.GetColumnNamesAndTypes(db, "callback"); err != nil {
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
		"webhook_id": {
			ColumnName:       "webhook_id",
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
		"created_at": {
			ColumnName:       "created_at",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"changes": {
			ColumnName:       "changes",
			ColumnType:       "json.RawMessage",
			ColumnDBType:     "[]byte",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"result": {
			ColumnName:       "result",
			ColumnType:       "json.RawMessage",
			ColumnDBType:     "[]byte",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
	}
	if err := migration.Compare(db, "callback", mModelColumnNameAndType, mDBColumnNameAndType); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*Callback)(nil))
}

func (m *Callback) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		m.ID,
		m.WebhookID,
		m.AccountID,
		core.Now(m.CreatedAt, now, create),
		core.JSON{m.Changes},
		core.JSON{m.Result},
	}
}

func (m *Callback) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		&m.WebhookID,
		&m.AccountID,
		(*core.Time)(&m.CreatedAt),
		core.JSON{&m.Changes},
		core.JSON{&m.Result},
	}
}

func (m *Callback) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *Callbacks) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(Callbacks, 0, 128)
	for rows.Next() {
		m := new(Callback)
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

func (_ *Callback) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlCallback_Select)
	return nil
}

func (_ *Callbacks) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlCallback_Select)
	return nil
}

func (m *Callback) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlCallback_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(6)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms Callbacks) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlCallback_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(6)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *Callback) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlCallback_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlCallback_ListColsOnConflict)
	return nil
}

func (ms Callbacks) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlCallback_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlCallback_ListColsOnConflict)
	return nil
}

func (m *Callback) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("callback")
	w.WriteRawString(" SET ")
	if m.ID != 0 {
		flag = true
		w.WriteName("id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ID)
	}
	if m.WebhookID != 0 {
		flag = true
		w.WriteName("webhook_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.WebhookID)
	}
	if m.AccountID != 0 {
		flag = true
		w.WriteName("account_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.AccountID)
	}
	if !m.CreatedAt.IsZero() {
		flag = true
		w.WriteName("created_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CreatedAt)
	}
	if m.Changes != nil {
		flag = true
		w.WriteName("changes")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.Changes})
	}
	if m.Result != nil {
		flag = true
		w.WriteName("result")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.Result})
	}
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *Callback) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlCallback_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(6)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type CallbackHistory map[string]interface{}
type CallbackHistories []map[string]interface{}

func (m *CallbackHistory) SQLTableName() string  { return "history.\"callback\"" }
func (m CallbackHistories) SQLTableName() string { return "history.\"callback\"" }

func (m *CallbackHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlCallback_Select_history)
	return nil
}

func (m CallbackHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlCallback_Select_history)
	return nil
}

func (m CallbackHistory) ID() core.Interface        { return core.Interface{m["id"]} }
func (m CallbackHistory) WebhookID() core.Interface { return core.Interface{m["webhook_id"]} }
func (m CallbackHistory) AccountID() core.Interface { return core.Interface{m["account_id"]} }
func (m CallbackHistory) CreatedAt() core.Interface { return core.Interface{m["created_at"]} }
func (m CallbackHistory) Changes() core.Interface   { return core.Interface{m["changes"]} }
func (m CallbackHistory) Result() core.Interface    { return core.Interface{m["result"]} }

func (m *CallbackHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 6)
	args := make([]interface{}, 6)
	for i := 0; i < 6; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(CallbackHistory, 6)
	res["id"] = data[0]
	res["webhook_id"] = data[1]
	res["account_id"] = data[2]
	res["created_at"] = data[3]
	res["changes"] = data[4]
	res["result"] = data[5]
	*m = res
	return nil
}

func (ms *CallbackHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 6)
	args := make([]interface{}, 6)
	for i := 0; i < 6; i++ {
		args[i] = &data[i]
	}
	res := make(CallbackHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(CallbackHistory)
		m["id"] = data[0]
		m["webhook_id"] = data[1]
		m["account_id"] = data[2]
		m["created_at"] = data[3]
		m["changes"] = data[4]
		m["result"] = data[5]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}

type Webhooks []*Webhook

const __sqlWebhook_Table = "webhook"
const __sqlWebhook_ListCols = "\"id\",\"account_id\",\"entities\",\"fields\",\"url\",\"metadata\",\"created_at\",\"updated_at\",\"deleted_at\""
const __sqlWebhook_ListColsOnConflict = "\"id\" = EXCLUDED.\"id\",\"account_id\" = EXCLUDED.\"account_id\",\"entities\" = EXCLUDED.\"entities\",\"fields\" = EXCLUDED.\"fields\",\"url\" = EXCLUDED.\"url\",\"metadata\" = EXCLUDED.\"metadata\",\"created_at\" = EXCLUDED.\"created_at\",\"updated_at\" = EXCLUDED.\"updated_at\",\"deleted_at\" = EXCLUDED.\"deleted_at\""
const __sqlWebhook_Insert = "INSERT INTO \"webhook\" (" + __sqlWebhook_ListCols + ") VALUES"
const __sqlWebhook_Select = "SELECT " + __sqlWebhook_ListCols + " FROM \"webhook\""
const __sqlWebhook_Select_history = "SELECT " + __sqlWebhook_ListCols + " FROM history.\"webhook\""
const __sqlWebhook_UpdateAll = "UPDATE \"webhook\" SET (" + __sqlWebhook_ListCols + ")"
const __sqlWebhook_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT webhook_pkey DO UPDATE SET"

func (m *Webhook) SQLTableName() string  { return "webhook" }
func (m *Webhooks) SQLTableName() string { return "webhook" }
func (m *Webhook) SQLListCols() string   { return __sqlWebhook_ListCols }

func (m *Webhook) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlWebhook_ListCols + " FROM \"webhook\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func (m *Webhook) Migration(db *cmsql.Database) {
	var mDBColumnNameAndType map[string]string
	if val, err := migration.GetColumnNamesAndTypes(db, "webhook"); err != nil {
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
		"entities": {
			ColumnName:       "entities",
			ColumnType:       "[]entity_type.EntityType",
			ColumnDBType:     "[]enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"fields": {
			ColumnName:       "fields",
			ColumnType:       "[]string",
			ColumnDBType:     "[]string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"url": {
			ColumnName:       "url",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"metadata": {
			ColumnName:       "metadata",
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
	if err := migration.Compare(db, "webhook", mModelColumnNameAndType, mDBColumnNameAndType); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*Webhook)(nil))
}

func (m *Webhook) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		m.ID,
		m.AccountID,
		core.Array{m.Entities, opts},
		core.Array{m.Fields, opts},
		core.String(m.URL),
		core.String(m.Metadata),
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
		core.Time(m.DeletedAt),
	}
}

func (m *Webhook) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		&m.AccountID,
		core.Array{&m.Entities, opts},
		core.Array{&m.Fields, opts},
		(*core.String)(&m.URL),
		(*core.String)(&m.Metadata),
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		(*core.Time)(&m.DeletedAt),
	}
}

func (m *Webhook) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *Webhooks) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(Webhooks, 0, 128)
	for rows.Next() {
		m := new(Webhook)
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

func (_ *Webhook) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlWebhook_Select)
	return nil
}

func (_ *Webhooks) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlWebhook_Select)
	return nil
}

func (m *Webhook) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlWebhook_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(9)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms Webhooks) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlWebhook_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(9)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *Webhook) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlWebhook_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlWebhook_ListColsOnConflict)
	return nil
}

func (ms Webhooks) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlWebhook_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlWebhook_ListColsOnConflict)
	return nil
}

func (m *Webhook) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("webhook")
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
	if m.Entities != nil {
		flag = true
		w.WriteName("entities")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.Entities, opts})
	}
	if m.Fields != nil {
		flag = true
		w.WriteName("fields")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.Fields, opts})
	}
	if m.URL != "" {
		flag = true
		w.WriteName("url")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.URL)
	}
	if m.Metadata != "" {
		flag = true
		w.WriteName("metadata")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Metadata)
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

func (m *Webhook) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlWebhook_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(9)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type WebhookHistory map[string]interface{}
type WebhookHistories []map[string]interface{}

func (m *WebhookHistory) SQLTableName() string  { return "history.\"webhook\"" }
func (m WebhookHistories) SQLTableName() string { return "history.\"webhook\"" }

func (m *WebhookHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlWebhook_Select_history)
	return nil
}

func (m WebhookHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlWebhook_Select_history)
	return nil
}

func (m WebhookHistory) ID() core.Interface        { return core.Interface{m["id"]} }
func (m WebhookHistory) AccountID() core.Interface { return core.Interface{m["account_id"]} }
func (m WebhookHistory) Entities() core.Interface  { return core.Interface{m["entities"]} }
func (m WebhookHistory) Fields() core.Interface    { return core.Interface{m["fields"]} }
func (m WebhookHistory) URL() core.Interface       { return core.Interface{m["url"]} }
func (m WebhookHistory) Metadata() core.Interface  { return core.Interface{m["metadata"]} }
func (m WebhookHistory) CreatedAt() core.Interface { return core.Interface{m["created_at"]} }
func (m WebhookHistory) UpdatedAt() core.Interface { return core.Interface{m["updated_at"]} }
func (m WebhookHistory) DeletedAt() core.Interface { return core.Interface{m["deleted_at"]} }

func (m *WebhookHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 9)
	args := make([]interface{}, 9)
	for i := 0; i < 9; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(WebhookHistory, 9)
	res["id"] = data[0]
	res["account_id"] = data[1]
	res["entities"] = data[2]
	res["fields"] = data[3]
	res["url"] = data[4]
	res["metadata"] = data[5]
	res["created_at"] = data[6]
	res["updated_at"] = data[7]
	res["deleted_at"] = data[8]
	*m = res
	return nil
}

func (ms *WebhookHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 9)
	args := make([]interface{}, 9)
	for i := 0; i < 9; i++ {
		args[i] = &data[i]
	}
	res := make(WebhookHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(WebhookHistory)
		m["id"] = data[0]
		m["account_id"] = data[1]
		m["entities"] = data[2]
		m["fields"] = data[3]
		m["url"] = data[4]
		m["metadata"] = data[5]
		m["created_at"] = data[6]
		m["updated_at"] = data[7]
		m["deleted_at"] = data[8]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
