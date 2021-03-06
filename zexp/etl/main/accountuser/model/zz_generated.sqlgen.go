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

type AccountUsers []*AccountUser

const __sqlAccountUser_Table = "account_user"
const __sqlAccountUser_ListCols = "\"account_id\",\"user_id\",\"status\",\"response_status\",\"created_at\",\"updated_at\",\"roles\",\"permissions\",\"full_name\",\"short_name\",\"position\",\"invitation_sent_at\",\"invitation_sent_by\",\"invitation_accepted_at\",\"invitation_rejected_at\",\"disabled_at\",\"disabled_by\",\"disable_reason\",\"rid\""
const __sqlAccountUser_ListColsOnConflict = "\"account_id\" = EXCLUDED.\"account_id\",\"user_id\" = EXCLUDED.\"user_id\",\"status\" = EXCLUDED.\"status\",\"response_status\" = EXCLUDED.\"response_status\",\"created_at\" = EXCLUDED.\"created_at\",\"updated_at\" = EXCLUDED.\"updated_at\",\"roles\" = EXCLUDED.\"roles\",\"permissions\" = EXCLUDED.\"permissions\",\"full_name\" = EXCLUDED.\"full_name\",\"short_name\" = EXCLUDED.\"short_name\",\"position\" = EXCLUDED.\"position\",\"invitation_sent_at\" = EXCLUDED.\"invitation_sent_at\",\"invitation_sent_by\" = EXCLUDED.\"invitation_sent_by\",\"invitation_accepted_at\" = EXCLUDED.\"invitation_accepted_at\",\"invitation_rejected_at\" = EXCLUDED.\"invitation_rejected_at\",\"disabled_at\" = EXCLUDED.\"disabled_at\",\"disabled_by\" = EXCLUDED.\"disabled_by\",\"disable_reason\" = EXCLUDED.\"disable_reason\",\"rid\" = EXCLUDED.\"rid\""
const __sqlAccountUser_Insert = "INSERT INTO \"account_user\" (" + __sqlAccountUser_ListCols + ") VALUES"
const __sqlAccountUser_Select = "SELECT " + __sqlAccountUser_ListCols + " FROM \"account_user\""
const __sqlAccountUser_Select_history = "SELECT " + __sqlAccountUser_ListCols + " FROM history.\"account_user\""
const __sqlAccountUser_UpdateAll = "UPDATE \"account_user\" SET (" + __sqlAccountUser_ListCols + ")"
const __sqlAccountUser_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT account_user_pkey DO UPDATE SET"

func (m *AccountUser) SQLTableName() string  { return "account_user" }
func (m *AccountUsers) SQLTableName() string { return "account_user" }
func (m *AccountUser) SQLListCols() string   { return __sqlAccountUser_ListCols }

func (m *AccountUser) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlAccountUser_ListCols + " FROM \"account_user\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func (m *AccountUser) Migration(db *cmsql.Database) {
	var mDBColumnNameAndType map[string]string
	if val, err := migration.GetColumnNamesAndTypes(db, "account_user"); err != nil {
		db.RecordError(err)
		return
	} else {
		mDBColumnNameAndType = val
	}
	mModelColumnNameAndType := map[string]migration.ColumnDef{
		"account_id": {
			ColumnName:       "account_id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"user_id": {
			ColumnName:       "user_id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"status": {
			ColumnName:       "status",
			ColumnType:       "status3.Status",
			ColumnDBType:     "enum",
			ColumnTag:        "int2",
			ColumnEnumValues: []string{"Z", "P", "N"},
		},
		"response_status": {
			ColumnName:       "response_status",
			ColumnType:       "status3.Status",
			ColumnDBType:     "enum",
			ColumnTag:        "int2",
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
		"roles": {
			ColumnName:       "roles",
			ColumnType:       "[]string",
			ColumnDBType:     "[]string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"permissions": {
			ColumnName:       "permissions",
			ColumnType:       "[]string",
			ColumnDBType:     "[]string",
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
		"short_name": {
			ColumnName:       "short_name",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"position": {
			ColumnName:       "position",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"invitation_sent_at": {
			ColumnName:       "invitation_sent_at",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"invitation_sent_by": {
			ColumnName:       "invitation_sent_by",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"invitation_accepted_at": {
			ColumnName:       "invitation_accepted_at",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"invitation_rejected_at": {
			ColumnName:       "invitation_rejected_at",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"disabled_at": {
			ColumnName:       "disabled_at",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"disabled_by": {
			ColumnName:       "disabled_by",
			ColumnType:       "int8",
			ColumnDBType:     "int8",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"disable_reason": {
			ColumnName:       "disable_reason",
			ColumnType:       "string",
			ColumnDBType:     "string",
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
	if err := migration.Compare(db, "account_user", mModelColumnNameAndType, mDBColumnNameAndType); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*AccountUser)(nil))
}

func (m *AccountUser) SQLArgs(opts core.Opts, create bool) []interface{} {
	return []interface{}{
		m.AccountID,
		m.UserID,
		m.Status,
		m.ResponseStatus,
		core.Time(m.CreatedAt),
		core.Time(m.UpdatedAt),
		core.Array{m.Permission.Roles, opts},
		core.Array{m.Permission.Permissions, opts},
		core.String(m.FullName),
		core.String(m.ShortName),
		core.String(m.Position),
		core.Time(m.InvitationSentAt),
		m.InvitationSentBy,
		core.Time(m.InvitationAcceptedAt),
		core.Time(m.InvitationRejectedAt),
		core.Time(m.DisabledAt),
		core.Int8(m.DisabledBy),
		core.String(m.DisableReason),
		m.Rid,
	}
}

func (m *AccountUser) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.AccountID,
		&m.UserID,
		&m.Status,
		&m.ResponseStatus,
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		core.Array{&m.Permission.Roles, opts},
		core.Array{&m.Permission.Permissions, opts},
		(*core.String)(&m.FullName),
		(*core.String)(&m.ShortName),
		(*core.String)(&m.Position),
		(*core.Time)(&m.InvitationSentAt),
		&m.InvitationSentBy,
		(*core.Time)(&m.InvitationAcceptedAt),
		(*core.Time)(&m.InvitationRejectedAt),
		(*core.Time)(&m.DisabledAt),
		(*core.Int8)(&m.DisabledBy),
		(*core.String)(&m.DisableReason),
		&m.Rid,
	}
}

func (m *AccountUser) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *AccountUsers) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(AccountUsers, 0, 128)
	for rows.Next() {
		m := new(AccountUser)
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

func (_ *AccountUser) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlAccountUser_Select)
	return nil
}

func (_ *AccountUsers) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlAccountUser_Select)
	return nil
}

func (m *AccountUser) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlAccountUser_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(19)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms AccountUsers) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlAccountUser_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(19)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *AccountUser) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlAccountUser_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlAccountUser_ListColsOnConflict)
	return nil
}

func (ms AccountUsers) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlAccountUser_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlAccountUser_ListColsOnConflict)
	return nil
}

func (m *AccountUser) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("account_user")
	w.WriteRawString(" SET ")
	if m.AccountID != 0 {
		flag = true
		w.WriteName("account_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.AccountID)
	}
	if m.UserID != 0 {
		flag = true
		w.WriteName("user_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.UserID)
	}
	if m.Status != 0 {
		flag = true
		w.WriteName("status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Status)
	}
	if m.ResponseStatus != 0 {
		flag = true
		w.WriteName("response_status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ResponseStatus)
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
	if m.Permission.Roles != nil {
		flag = true
		w.WriteName("roles")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.Permission.Roles, opts})
	}
	if m.Permission.Permissions != nil {
		flag = true
		w.WriteName("permissions")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.Permission.Permissions, opts})
	}
	if m.FullName != "" {
		flag = true
		w.WriteName("full_name")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.FullName)
	}
	if m.ShortName != "" {
		flag = true
		w.WriteName("short_name")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShortName)
	}
	if m.Position != "" {
		flag = true
		w.WriteName("position")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Position)
	}
	if !m.InvitationSentAt.IsZero() {
		flag = true
		w.WriteName("invitation_sent_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.InvitationSentAt)
	}
	if m.InvitationSentBy != 0 {
		flag = true
		w.WriteName("invitation_sent_by")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.InvitationSentBy)
	}
	if !m.InvitationAcceptedAt.IsZero() {
		flag = true
		w.WriteName("invitation_accepted_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.InvitationAcceptedAt)
	}
	if !m.InvitationRejectedAt.IsZero() {
		flag = true
		w.WriteName("invitation_rejected_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.InvitationRejectedAt)
	}
	if !m.DisabledAt.IsZero() {
		flag = true
		w.WriteName("disabled_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.DisabledAt)
	}
	if m.DisabledBy != 0 {
		flag = true
		w.WriteName("disabled_by")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.DisabledBy)
	}
	if m.DisableReason != "" {
		flag = true
		w.WriteName("disable_reason")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.DisableReason)
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

func (m *AccountUser) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlAccountUser_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(19)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type AccountUserHistory map[string]interface{}
type AccountUserHistories []map[string]interface{}

func (m *AccountUserHistory) SQLTableName() string  { return "history.\"account_user\"" }
func (m AccountUserHistories) SQLTableName() string { return "history.\"account_user\"" }

func (m *AccountUserHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlAccountUser_Select_history)
	return nil
}

func (m AccountUserHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlAccountUser_Select_history)
	return nil
}

func (m AccountUserHistory) AccountID() core.Interface { return core.Interface{m["account_id"]} }
func (m AccountUserHistory) UserID() core.Interface    { return core.Interface{m["user_id"]} }
func (m AccountUserHistory) Status() core.Interface    { return core.Interface{m["status"]} }
func (m AccountUserHistory) ResponseStatus() core.Interface {
	return core.Interface{m["response_status"]}
}
func (m AccountUserHistory) CreatedAt() core.Interface   { return core.Interface{m["created_at"]} }
func (m AccountUserHistory) UpdatedAt() core.Interface   { return core.Interface{m["updated_at"]} }
func (m AccountUserHistory) Roles() core.Interface       { return core.Interface{m["roles"]} }
func (m AccountUserHistory) Permissions() core.Interface { return core.Interface{m["permissions"]} }
func (m AccountUserHistory) FullName() core.Interface    { return core.Interface{m["full_name"]} }
func (m AccountUserHistory) ShortName() core.Interface   { return core.Interface{m["short_name"]} }
func (m AccountUserHistory) Position() core.Interface    { return core.Interface{m["position"]} }
func (m AccountUserHistory) InvitationSentAt() core.Interface {
	return core.Interface{m["invitation_sent_at"]}
}
func (m AccountUserHistory) InvitationSentBy() core.Interface {
	return core.Interface{m["invitation_sent_by"]}
}
func (m AccountUserHistory) InvitationAcceptedAt() core.Interface {
	return core.Interface{m["invitation_accepted_at"]}
}
func (m AccountUserHistory) InvitationRejectedAt() core.Interface {
	return core.Interface{m["invitation_rejected_at"]}
}
func (m AccountUserHistory) DisabledAt() core.Interface { return core.Interface{m["disabled_at"]} }
func (m AccountUserHistory) DisabledBy() core.Interface { return core.Interface{m["disabled_by"]} }
func (m AccountUserHistory) DisableReason() core.Interface {
	return core.Interface{m["disable_reason"]}
}
func (m AccountUserHistory) Rid() core.Interface { return core.Interface{m["rid"]} }

func (m *AccountUserHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 19)
	args := make([]interface{}, 19)
	for i := 0; i < 19; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(AccountUserHistory, 19)
	res["account_id"] = data[0]
	res["user_id"] = data[1]
	res["status"] = data[2]
	res["response_status"] = data[3]
	res["created_at"] = data[4]
	res["updated_at"] = data[5]
	res["roles"] = data[6]
	res["permissions"] = data[7]
	res["full_name"] = data[8]
	res["short_name"] = data[9]
	res["position"] = data[10]
	res["invitation_sent_at"] = data[11]
	res["invitation_sent_by"] = data[12]
	res["invitation_accepted_at"] = data[13]
	res["invitation_rejected_at"] = data[14]
	res["disabled_at"] = data[15]
	res["disabled_by"] = data[16]
	res["disable_reason"] = data[17]
	res["rid"] = data[18]
	*m = res
	return nil
}

func (ms *AccountUserHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 19)
	args := make([]interface{}, 19)
	for i := 0; i < 19; i++ {
		args[i] = &data[i]
	}
	res := make(AccountUserHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(AccountUserHistory)
		m["account_id"] = data[0]
		m["user_id"] = data[1]
		m["status"] = data[2]
		m["response_status"] = data[3]
		m["created_at"] = data[4]
		m["updated_at"] = data[5]
		m["roles"] = data[6]
		m["permissions"] = data[7]
		m["full_name"] = data[8]
		m["short_name"] = data[9]
		m["position"] = data[10]
		m["invitation_sent_at"] = data[11]
		m["invitation_sent_by"] = data[12]
		m["invitation_accepted_at"] = data[13]
		m["invitation_rejected_at"] = data[14]
		m["disabled_at"] = data[15]
		m["disabled_by"] = data[16]
		m["disable_reason"] = data[17]
		m["rid"] = data[18]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
