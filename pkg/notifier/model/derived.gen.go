// Code generated by goderive DO NOT EDIT.

package model

import (
	"database/sql"
	"time"

	core "etop.vn/backend/pkg/common/sq/core"
)

type SQLWriter = core.SQLWriter

// Type Notification represents table notification
func sqlgenNotification(_ *Notification) bool { return true }

type Notifications []*Notification

const __sqlNotification_Table = "notification"
const __sqlNotification_ListCols = "\"id\",\"title\",\"message\",\"is_read\",\"entity_id\",\"entity\",\"account_id\",\"sync_status\",\"external_service_id\",\"external_noti_id\",\"send_notification\",\"synced_at\",\"seen_at\",\"created_at\",\"updated_at\""
const __sqlNotification_Insert = "INSERT INTO \"notification\" (" + __sqlNotification_ListCols + ") VALUES"
const __sqlNotification_Select = "SELECT " + __sqlNotification_ListCols + " FROM \"notification\""
const __sqlNotification_Select_history = "SELECT " + __sqlNotification_ListCols + " FROM history.\"notification\""
const __sqlNotification_UpdateAll = "UPDATE \"notification\" SET (" + __sqlNotification_ListCols + ")"

func (m *Notification) SQLTableName() string  { return "notification" }
func (m *Notifications) SQLTableName() string { return "notification" }
func (m *Notification) SQLListCols() string   { return __sqlNotification_ListCols }

func (m *Notification) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		core.Int64(m.ID),
		core.String(m.Title),
		core.String(m.Message),
		core.Bool(m.IsRead),
		core.Int64(m.EntityID),
		core.String(m.Entity),
		core.Int64(m.AccountID),
		core.Int(m.SyncStatus),
		core.Int(m.ExternalServiceID),
		core.String(m.ExternalNotiID),
		core.Bool(m.SendNotification),
		core.Time(m.SyncedAt),
		core.Time(m.SeenAt),
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
	}
}

func (m *Notification) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		(*core.Int64)(&m.ID),
		(*core.String)(&m.Title),
		(*core.String)(&m.Message),
		(*core.Bool)(&m.IsRead),
		(*core.Int64)(&m.EntityID),
		(*core.String)(&m.Entity),
		(*core.Int64)(&m.AccountID),
		(*core.Int)(&m.SyncStatus),
		(*core.Int)(&m.ExternalServiceID),
		(*core.String)(&m.ExternalNotiID),
		(*core.Bool)(&m.SendNotification),
		(*core.Time)(&m.SyncedAt),
		(*core.Time)(&m.SeenAt),
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
	}
}

func (m *Notification) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *Notifications) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(Notifications, 0, 128)
	for rows.Next() {
		m := new(Notification)
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

func (_ *Notification) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlNotification_Select)
	return nil
}

func (_ *Notifications) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlNotification_Select)
	return nil
}

func (m *Notification) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlNotification_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(15)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms Notifications) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlNotification_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(15)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *Notification) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("notification")
	w.WriteRawString(" SET ")
	if m.ID != 0 {
		flag = true
		w.WriteName("id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ID)
	}
	if m.Title != "" {
		flag = true
		w.WriteName("title")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Title)
	}
	if m.Message != "" {
		flag = true
		w.WriteName("message")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Message)
	}
	if m.IsRead {
		flag = true
		w.WriteName("is_read")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.IsRead)
	}
	if m.EntityID != 0 {
		flag = true
		w.WriteName("entity_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.EntityID)
	}
	if m.Entity != "" {
		flag = true
		w.WriteName("entity")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(string(m.Entity))
	}
	if m.AccountID != 0 {
		flag = true
		w.WriteName("account_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.AccountID)
	}
	if m.SyncStatus != 0 {
		flag = true
		w.WriteName("sync_status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(int(m.SyncStatus))
	}
	if m.ExternalServiceID != 0 {
		flag = true
		w.WriteName("external_service_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ExternalServiceID)
	}
	if m.ExternalNotiID != "" {
		flag = true
		w.WriteName("external_noti_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ExternalNotiID)
	}
	if m.SendNotification {
		flag = true
		w.WriteName("send_notification")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.SendNotification)
	}
	if !m.SyncedAt.IsZero() {
		flag = true
		w.WriteName("synced_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.SyncedAt)
	}
	if !m.SeenAt.IsZero() {
		flag = true
		w.WriteName("seen_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.SeenAt)
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

func (m *Notification) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlNotification_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(15)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type NotificationHistory map[string]interface{}
type NotificationHistories []map[string]interface{}

func (m *NotificationHistory) SQLTableName() string  { return "history.\"notification\"" }
func (m NotificationHistories) SQLTableName() string { return "history.\"notification\"" }

func (m *NotificationHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlNotification_Select_history)
	return nil
}

func (m NotificationHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlNotification_Select_history)
	return nil
}

func (m NotificationHistory) ID() core.Interface         { return core.Interface{m["id"]} }
func (m NotificationHistory) Title() core.Interface      { return core.Interface{m["title"]} }
func (m NotificationHistory) Message() core.Interface    { return core.Interface{m["message"]} }
func (m NotificationHistory) IsRead() core.Interface     { return core.Interface{m["is_read"]} }
func (m NotificationHistory) EntityID() core.Interface   { return core.Interface{m["entity_id"]} }
func (m NotificationHistory) Entity() core.Interface     { return core.Interface{m["entity"]} }
func (m NotificationHistory) AccountID() core.Interface  { return core.Interface{m["account_id"]} }
func (m NotificationHistory) SyncStatus() core.Interface { return core.Interface{m["sync_status"]} }
func (m NotificationHistory) ExternalServiceID() core.Interface {
	return core.Interface{m["external_service_id"]}
}
func (m NotificationHistory) ExternalNotiID() core.Interface {
	return core.Interface{m["external_noti_id"]}
}
func (m NotificationHistory) SendNotification() core.Interface {
	return core.Interface{m["send_notification"]}
}
func (m NotificationHistory) SyncedAt() core.Interface  { return core.Interface{m["synced_at"]} }
func (m NotificationHistory) SeenAt() core.Interface    { return core.Interface{m["seen_at"]} }
func (m NotificationHistory) CreatedAt() core.Interface { return core.Interface{m["created_at"]} }
func (m NotificationHistory) UpdatedAt() core.Interface { return core.Interface{m["updated_at"]} }

func (m *NotificationHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 15)
	args := make([]interface{}, 15)
	for i := 0; i < 15; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(NotificationHistory, 15)
	res["id"] = data[0]
	res["title"] = data[1]
	res["message"] = data[2]
	res["is_read"] = data[3]
	res["entity_id"] = data[4]
	res["entity"] = data[5]
	res["account_id"] = data[6]
	res["sync_status"] = data[7]
	res["external_service_id"] = data[8]
	res["external_noti_id"] = data[9]
	res["send_notification"] = data[10]
	res["synced_at"] = data[11]
	res["seen_at"] = data[12]
	res["created_at"] = data[13]
	res["updated_at"] = data[14]
	*m = res
	return nil
}

func (ms *NotificationHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 15)
	args := make([]interface{}, 15)
	for i := 0; i < 15; i++ {
		args[i] = &data[i]
	}
	res := make(NotificationHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(NotificationHistory)
		m["id"] = data[0]
		m["title"] = data[1]
		m["message"] = data[2]
		m["is_read"] = data[3]
		m["entity_id"] = data[4]
		m["entity"] = data[5]
		m["account_id"] = data[6]
		m["sync_status"] = data[7]
		m["external_service_id"] = data[8]
		m["external_noti_id"] = data[9]
		m["send_notification"] = data[10]
		m["synced_at"] = data[11]
		m["seen_at"] = data[12]
		m["created_at"] = data[13]
		m["updated_at"] = data[14]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}

// Type Device represents table device
func sqlgenDevice(_ *Device) bool { return true }

type Devices []*Device

const __sqlDevice_Table = "device"
const __sqlDevice_ListCols = "\"id\",\"device_id\",\"device_name\",\"external_device_id\",\"external_service_id\",\"account_id\",\"user_id\",\"created_at\",\"updated_at\",\"deactivated_at\",\"config\""
const __sqlDevice_Insert = "INSERT INTO \"device\" (" + __sqlDevice_ListCols + ") VALUES"
const __sqlDevice_Select = "SELECT " + __sqlDevice_ListCols + " FROM \"device\""
const __sqlDevice_Select_history = "SELECT " + __sqlDevice_ListCols + " FROM history.\"device\""
const __sqlDevice_UpdateAll = "UPDATE \"device\" SET (" + __sqlDevice_ListCols + ")"

func (m *Device) SQLTableName() string  { return "device" }
func (m *Devices) SQLTableName() string { return "device" }
func (m *Device) SQLListCols() string   { return __sqlDevice_ListCols }

func (m *Device) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		core.Int64(m.ID),
		core.String(m.DeviceID),
		core.String(m.DeviceName),
		core.String(m.ExternalDeviceID),
		core.Int(m.ExternalServiceID),
		core.Int64(m.AccountID),
		core.Int64(m.UserID),
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
		core.Time(m.DeactivatedAt),
		core.JSON{m.Config},
	}
}

func (m *Device) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		(*core.Int64)(&m.ID),
		(*core.String)(&m.DeviceID),
		(*core.String)(&m.DeviceName),
		(*core.String)(&m.ExternalDeviceID),
		(*core.Int)(&m.ExternalServiceID),
		(*core.Int64)(&m.AccountID),
		(*core.Int64)(&m.UserID),
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		(*core.Time)(&m.DeactivatedAt),
		core.JSON{&m.Config},
	}
}

func (m *Device) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *Devices) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(Devices, 0, 128)
	for rows.Next() {
		m := new(Device)
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

func (_ *Device) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlDevice_Select)
	return nil
}

func (_ *Devices) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlDevice_Select)
	return nil
}

func (m *Device) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlDevice_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(11)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms Devices) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlDevice_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(11)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *Device) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("device")
	w.WriteRawString(" SET ")
	if m.ID != 0 {
		flag = true
		w.WriteName("id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ID)
	}
	if m.DeviceID != "" {
		flag = true
		w.WriteName("device_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.DeviceID)
	}
	if m.DeviceName != "" {
		flag = true
		w.WriteName("device_name")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.DeviceName)
	}
	if m.ExternalDeviceID != "" {
		flag = true
		w.WriteName("external_device_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ExternalDeviceID)
	}
	if m.ExternalServiceID != 0 {
		flag = true
		w.WriteName("external_service_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ExternalServiceID)
	}
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
	if !m.DeactivatedAt.IsZero() {
		flag = true
		w.WriteName("deactivated_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.DeactivatedAt)
	}
	if m.Config != nil {
		flag = true
		w.WriteName("config")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.Config})
	}
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *Device) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlDevice_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(11)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type DeviceHistory map[string]interface{}
type DeviceHistories []map[string]interface{}

func (m *DeviceHistory) SQLTableName() string  { return "history.\"device\"" }
func (m DeviceHistories) SQLTableName() string { return "history.\"device\"" }

func (m *DeviceHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlDevice_Select_history)
	return nil
}

func (m DeviceHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlDevice_Select_history)
	return nil
}

func (m DeviceHistory) ID() core.Interface         { return core.Interface{m["id"]} }
func (m DeviceHistory) DeviceID() core.Interface   { return core.Interface{m["device_id"]} }
func (m DeviceHistory) DeviceName() core.Interface { return core.Interface{m["device_name"]} }
func (m DeviceHistory) ExternalDeviceID() core.Interface {
	return core.Interface{m["external_device_id"]}
}
func (m DeviceHistory) ExternalServiceID() core.Interface {
	return core.Interface{m["external_service_id"]}
}
func (m DeviceHistory) AccountID() core.Interface     { return core.Interface{m["account_id"]} }
func (m DeviceHistory) UserID() core.Interface        { return core.Interface{m["user_id"]} }
func (m DeviceHistory) CreatedAt() core.Interface     { return core.Interface{m["created_at"]} }
func (m DeviceHistory) UpdatedAt() core.Interface     { return core.Interface{m["updated_at"]} }
func (m DeviceHistory) DeactivatedAt() core.Interface { return core.Interface{m["deactivated_at"]} }
func (m DeviceHistory) Config() core.Interface        { return core.Interface{m["config"]} }

func (m *DeviceHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 11)
	args := make([]interface{}, 11)
	for i := 0; i < 11; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(DeviceHistory, 11)
	res["id"] = data[0]
	res["device_id"] = data[1]
	res["device_name"] = data[2]
	res["external_device_id"] = data[3]
	res["external_service_id"] = data[4]
	res["account_id"] = data[5]
	res["user_id"] = data[6]
	res["created_at"] = data[7]
	res["updated_at"] = data[8]
	res["deactivated_at"] = data[9]
	res["config"] = data[10]
	*m = res
	return nil
}

func (ms *DeviceHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 11)
	args := make([]interface{}, 11)
	for i := 0; i < 11; i++ {
		args[i] = &data[i]
	}
	res := make(DeviceHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(DeviceHistory)
		m["id"] = data[0]
		m["device_id"] = data[1]
		m["device_name"] = data[2]
		m["external_device_id"] = data[3]
		m["external_service_id"] = data[4]
		m["account_id"] = data[5]
		m["user_id"] = data[6]
		m["created_at"] = data[7]
		m["updated_at"] = data[8]
		m["deactivated_at"] = data[9]
		m["config"] = data[10]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
