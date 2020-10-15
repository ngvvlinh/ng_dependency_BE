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

type Invitations []*Invitation

const __sqlInvitation_Table = "invitation"
const __sqlInvitation_ListCols = "\"id\",\"account_id\",\"email\",\"phone\",\"full_name\",\"short_name\",\"roles\",\"token\",\"status\",\"invited_by\",\"accepted_at\",\"rejected_at\",\"expires_at\",\"created_at\",\"updated_at\",\"deleted_at\",\"rid\""
const __sqlInvitation_ListColsOnConflict = "\"id\" = EXCLUDED.\"id\",\"account_id\" = EXCLUDED.\"account_id\",\"email\" = EXCLUDED.\"email\",\"phone\" = EXCLUDED.\"phone\",\"full_name\" = EXCLUDED.\"full_name\",\"short_name\" = EXCLUDED.\"short_name\",\"roles\" = EXCLUDED.\"roles\",\"token\" = EXCLUDED.\"token\",\"status\" = EXCLUDED.\"status\",\"invited_by\" = EXCLUDED.\"invited_by\",\"accepted_at\" = EXCLUDED.\"accepted_at\",\"rejected_at\" = EXCLUDED.\"rejected_at\",\"expires_at\" = EXCLUDED.\"expires_at\",\"created_at\" = EXCLUDED.\"created_at\",\"updated_at\" = EXCLUDED.\"updated_at\",\"deleted_at\" = EXCLUDED.\"deleted_at\",\"rid\" = EXCLUDED.\"rid\""
const __sqlInvitation_Insert = "INSERT INTO \"invitation\" (" + __sqlInvitation_ListCols + ") VALUES"
const __sqlInvitation_Select = "SELECT " + __sqlInvitation_ListCols + " FROM \"invitation\""
const __sqlInvitation_Select_history = "SELECT " + __sqlInvitation_ListCols + " FROM history.\"invitation\""
const __sqlInvitation_UpdateAll = "UPDATE \"invitation\" SET (" + __sqlInvitation_ListCols + ")"
const __sqlInvitation_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT invitation_pkey DO UPDATE SET"

func (m *Invitation) SQLTableName() string  { return "invitation" }
func (m *Invitations) SQLTableName() string { return "invitation" }
func (m *Invitation) SQLListCols() string   { return __sqlInvitation_ListCols }

func (m *Invitation) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlInvitation_ListCols + " FROM \"invitation\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func (m *Invitation) Migration(db *cmsql.Database) {
	var mDBColumnNameAndType map[string]string
	if val, err := migration.GetColumnNamesAndTypes(db, "invitation"); err != nil {
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
		"email": {
			ColumnName:       "email",
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
		"roles": {
			ColumnName:       "roles",
			ColumnType:       "[]string",
			ColumnDBType:     "[]string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"token": {
			ColumnName:       "token",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"status": {
			ColumnName:       "status",
			ColumnType:       "status3.Status",
			ColumnDBType:     "enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{"Z", "P", "N"},
		},
		"invited_by": {
			ColumnName:       "invited_by",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"accepted_at": {
			ColumnName:       "accepted_at",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"rejected_at": {
			ColumnName:       "rejected_at",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"expires_at": {
			ColumnName:       "expires_at",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
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
		"rid": {
			ColumnName:       "rid",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
	}
	if err := migration.Compare(db, "invitation", mModelColumnNameAndType, mDBColumnNameAndType); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*Invitation)(nil))
}

func (m *Invitation) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		m.ID,
		m.AccountID,
		core.String(m.Email),
		core.String(m.Phone),
		core.String(m.FullName),
		core.String(m.ShortName),
		core.Array{m.Roles, opts},
		core.String(m.Token),
		m.Status,
		m.InvitedBy,
		core.Time(m.AcceptedAt),
		core.Time(m.RejectedAt),
		core.Time(m.ExpiresAt),
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
		core.Time(m.DeletedAt),
		m.Rid,
	}
}

func (m *Invitation) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		&m.AccountID,
		(*core.String)(&m.Email),
		(*core.String)(&m.Phone),
		(*core.String)(&m.FullName),
		(*core.String)(&m.ShortName),
		core.Array{&m.Roles, opts},
		(*core.String)(&m.Token),
		&m.Status,
		&m.InvitedBy,
		(*core.Time)(&m.AcceptedAt),
		(*core.Time)(&m.RejectedAt),
		(*core.Time)(&m.ExpiresAt),
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		(*core.Time)(&m.DeletedAt),
		&m.Rid,
	}
}

func (m *Invitation) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *Invitations) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(Invitations, 0, 128)
	for rows.Next() {
		m := new(Invitation)
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

func (_ *Invitation) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlInvitation_Select)
	return nil
}

func (_ *Invitations) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlInvitation_Select)
	return nil
}

func (m *Invitation) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlInvitation_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(17)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms Invitations) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlInvitation_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(17)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *Invitation) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlInvitation_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlInvitation_ListColsOnConflict)
	return nil
}

func (ms Invitations) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlInvitation_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlInvitation_ListColsOnConflict)
	return nil
}

func (m *Invitation) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("invitation")
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
	if m.Email != "" {
		flag = true
		w.WriteName("email")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Email)
	}
	if m.Phone != "" {
		flag = true
		w.WriteName("phone")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Phone)
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
	if m.Roles != nil {
		flag = true
		w.WriteName("roles")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.Roles, opts})
	}
	if m.Token != "" {
		flag = true
		w.WriteName("token")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Token)
	}
	if m.Status != 0 {
		flag = true
		w.WriteName("status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Status)
	}
	if m.InvitedBy != 0 {
		flag = true
		w.WriteName("invited_by")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.InvitedBy)
	}
	if !m.AcceptedAt.IsZero() {
		flag = true
		w.WriteName("accepted_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.AcceptedAt)
	}
	if !m.RejectedAt.IsZero() {
		flag = true
		w.WriteName("rejected_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.RejectedAt)
	}
	if !m.ExpiresAt.IsZero() {
		flag = true
		w.WriteName("expires_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ExpiresAt)
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

func (m *Invitation) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlInvitation_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(17)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type InvitationHistory map[string]interface{}
type InvitationHistories []map[string]interface{}

func (m *InvitationHistory) SQLTableName() string  { return "history.\"invitation\"" }
func (m InvitationHistories) SQLTableName() string { return "history.\"invitation\"" }

func (m *InvitationHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlInvitation_Select_history)
	return nil
}

func (m InvitationHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlInvitation_Select_history)
	return nil
}

func (m InvitationHistory) ID() core.Interface         { return core.Interface{m["id"]} }
func (m InvitationHistory) AccountID() core.Interface  { return core.Interface{m["account_id"]} }
func (m InvitationHistory) Email() core.Interface      { return core.Interface{m["email"]} }
func (m InvitationHistory) Phone() core.Interface      { return core.Interface{m["phone"]} }
func (m InvitationHistory) FullName() core.Interface   { return core.Interface{m["full_name"]} }
func (m InvitationHistory) ShortName() core.Interface  { return core.Interface{m["short_name"]} }
func (m InvitationHistory) Roles() core.Interface      { return core.Interface{m["roles"]} }
func (m InvitationHistory) Token() core.Interface      { return core.Interface{m["token"]} }
func (m InvitationHistory) Status() core.Interface     { return core.Interface{m["status"]} }
func (m InvitationHistory) InvitedBy() core.Interface  { return core.Interface{m["invited_by"]} }
func (m InvitationHistory) AcceptedAt() core.Interface { return core.Interface{m["accepted_at"]} }
func (m InvitationHistory) RejectedAt() core.Interface { return core.Interface{m["rejected_at"]} }
func (m InvitationHistory) ExpiresAt() core.Interface  { return core.Interface{m["expires_at"]} }
func (m InvitationHistory) CreatedAt() core.Interface  { return core.Interface{m["created_at"]} }
func (m InvitationHistory) UpdatedAt() core.Interface  { return core.Interface{m["updated_at"]} }
func (m InvitationHistory) DeletedAt() core.Interface  { return core.Interface{m["deleted_at"]} }
func (m InvitationHistory) Rid() core.Interface        { return core.Interface{m["rid"]} }

func (m *InvitationHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 17)
	args := make([]interface{}, 17)
	for i := 0; i < 17; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(InvitationHistory, 17)
	res["id"] = data[0]
	res["account_id"] = data[1]
	res["email"] = data[2]
	res["phone"] = data[3]
	res["full_name"] = data[4]
	res["short_name"] = data[5]
	res["roles"] = data[6]
	res["token"] = data[7]
	res["status"] = data[8]
	res["invited_by"] = data[9]
	res["accepted_at"] = data[10]
	res["rejected_at"] = data[11]
	res["expires_at"] = data[12]
	res["created_at"] = data[13]
	res["updated_at"] = data[14]
	res["deleted_at"] = data[15]
	res["rid"] = data[16]
	*m = res
	return nil
}

func (ms *InvitationHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 17)
	args := make([]interface{}, 17)
	for i := 0; i < 17; i++ {
		args[i] = &data[i]
	}
	res := make(InvitationHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(InvitationHistory)
		m["id"] = data[0]
		m["account_id"] = data[1]
		m["email"] = data[2]
		m["phone"] = data[3]
		m["full_name"] = data[4]
		m["short_name"] = data[5]
		m["roles"] = data[6]
		m["token"] = data[7]
		m["status"] = data[8]
		m["invited_by"] = data[9]
		m["accepted_at"] = data[10]
		m["rejected_at"] = data[11]
		m["expires_at"] = data[12]
		m["created_at"] = data[13]
		m["updated_at"] = data[14]
		m["deleted_at"] = data[15]
		m["rid"] = data[16]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
