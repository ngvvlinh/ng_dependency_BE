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

// Type Invitation represents table invitation
func sqlgenInvitation(_ *Invitation) bool { return true }

type Invitations []*Invitation

const __sqlInvitation_Table = "invitation"
const __sqlInvitation_ListCols = "\"id\",\"account_id\",\"email\",\"roles\",\"token\",\"status\",\"invited_by\",\"accepted_at\",\"rejected_at\",\"expires_at\",\"created_at\",\"updated_at\",\"deleted_at\""
const __sqlInvitation_Insert = "INSERT INTO \"invitation\" (" + __sqlInvitation_ListCols + ") VALUES"
const __sqlInvitation_Select = "SELECT " + __sqlInvitation_ListCols + " FROM \"invitation\""
const __sqlInvitation_Select_history = "SELECT " + __sqlInvitation_ListCols + " FROM history.\"invitation\""
const __sqlInvitation_UpdateAll = "UPDATE \"invitation\" SET (" + __sqlInvitation_ListCols + ")"

func (m *Invitation) SQLTableName() string  { return "invitation" }
func (m *Invitations) SQLTableName() string { return "invitation" }
func (m *Invitation) SQLListCols() string   { return __sqlInvitation_ListCols }

func (m *Invitation) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlInvitation_ListCols + " FROM \"invitation\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*Invitation)(nil))
}

func (m *Invitation) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		core.Int64(m.ID),
		core.Int64(m.AccountID),
		core.String(m.Email),
		core.Array{m.Roles, opts},
		core.String(m.Token),
		core.Int32(m.Status),
		core.Int64(m.InvitedBy),
		core.Time(m.AcceptedAt),
		core.Time(m.RejectedAt),
		core.Time(m.ExpiresAt),
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
		core.Time(m.DeletedAt),
	}
}

func (m *Invitation) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		(*core.Int64)(&m.ID),
		(*core.Int64)(&m.AccountID),
		(*core.String)(&m.Email),
		core.Array{&m.Roles, opts},
		(*core.String)(&m.Token),
		(*core.Int32)(&m.Status),
		(*core.Int64)(&m.InvitedBy),
		(*core.Time)(&m.AcceptedAt),
		(*core.Time)(&m.RejectedAt),
		(*core.Time)(&m.ExpiresAt),
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		(*core.Time)(&m.DeletedAt),
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
	w.WriteMarkers(13)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms Invitations) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlInvitation_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(13)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
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
		w.WriteArg(int32(m.Status))
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

func (m *Invitation) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlInvitation_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(13)
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

func (m *InvitationHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 13)
	args := make([]interface{}, 13)
	for i := 0; i < 13; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(InvitationHistory, 13)
	res["id"] = data[0]
	res["account_id"] = data[1]
	res["email"] = data[2]
	res["roles"] = data[3]
	res["token"] = data[4]
	res["status"] = data[5]
	res["invited_by"] = data[6]
	res["accepted_at"] = data[7]
	res["rejected_at"] = data[8]
	res["expires_at"] = data[9]
	res["created_at"] = data[10]
	res["updated_at"] = data[11]
	res["deleted_at"] = data[12]
	*m = res
	return nil
}

func (ms *InvitationHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 13)
	args := make([]interface{}, 13)
	for i := 0; i < 13; i++ {
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
		m["roles"] = data[3]
		m["token"] = data[4]
		m["status"] = data[5]
		m["invited_by"] = data[6]
		m["accepted_at"] = data[7]
		m["rejected_at"] = data[8]
		m["expires_at"] = data[9]
		m["created_at"] = data[10]
		m["updated_at"] = data[11]
		m["deleted_at"] = data[12]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}