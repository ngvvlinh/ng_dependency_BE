// Code generated by goderive DO NOT EDIT.

package model

import (
	"database/sql"
	"sync"
	"time"

	"etop.vn/backend/pkg/common/sql/cmsql"
	core "etop.vn/backend/pkg/common/sql/sq/core"
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

// Type Connection represents table connection
func sqlgenConnection(_ *Connection) bool { return true }

type Connections []*Connection

const __sqlConnection_Table = "connection"
const __sqlConnection_ListCols = "\"id\",\"name\",\"status\",\"partner_id\",\"created_at\",\"updated_at\",\"deleted_at\",\"driver_config\",\"driver\",\"connection_type\",\"connection_subtype\",\"connection_method\",\"connection_provider\",\"etop_affiliate_account\",\"code\",\"image_url\""
const __sqlConnection_ListColsOnConflict = "\"id\" = EXCLUDED.\"id\",\"name\" = EXCLUDED.\"name\",\"status\" = EXCLUDED.\"status\",\"partner_id\" = EXCLUDED.\"partner_id\",\"created_at\" = EXCLUDED.\"created_at\",\"updated_at\" = EXCLUDED.\"updated_at\",\"deleted_at\" = EXCLUDED.\"deleted_at\",\"driver_config\" = EXCLUDED.\"driver_config\",\"driver\" = EXCLUDED.\"driver\",\"connection_type\" = EXCLUDED.\"connection_type\",\"connection_subtype\" = EXCLUDED.\"connection_subtype\",\"connection_method\" = EXCLUDED.\"connection_method\",\"connection_provider\" = EXCLUDED.\"connection_provider\",\"etop_affiliate_account\" = EXCLUDED.\"etop_affiliate_account\",\"code\" = EXCLUDED.\"code\",\"image_url\" = EXCLUDED.\"image_url\""
const __sqlConnection_Insert = "INSERT INTO \"connection\" (" + __sqlConnection_ListCols + ") VALUES"
const __sqlConnection_Select = "SELECT " + __sqlConnection_ListCols + " FROM \"connection\""
const __sqlConnection_Select_history = "SELECT " + __sqlConnection_ListCols + " FROM history.\"connection\""
const __sqlConnection_UpdateAll = "UPDATE \"connection\" SET (" + __sqlConnection_ListCols + ")"
const __sqlConnection_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT connection_pkey DO UPDATE SET"

func (m *Connection) SQLTableName() string  { return "connection" }
func (m *Connections) SQLTableName() string { return "connection" }
func (m *Connection) SQLListCols() string   { return __sqlConnection_ListCols }

func (m *Connection) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlConnection_ListCols + " FROM \"connection\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*Connection)(nil))
}

func (m *Connection) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		m.ID,
		core.String(m.Name),
		m.Status,
		m.PartnerID,
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
		core.Time(m.DeletedAt),
		core.JSON{m.DriverConfig},
		core.String(m.Driver),
		m.ConnectionType,
		m.ConnectionSubtype,
		m.ConnectionMethod,
		m.ConnectionProvider,
		core.JSON{m.EtopAffiliateAccount},
		core.String(m.Code),
		core.String(m.ImageURL),
	}
}

func (m *Connection) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		(*core.String)(&m.Name),
		&m.Status,
		&m.PartnerID,
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		(*core.Time)(&m.DeletedAt),
		core.JSON{&m.DriverConfig},
		(*core.String)(&m.Driver),
		&m.ConnectionType,
		&m.ConnectionSubtype,
		&m.ConnectionMethod,
		&m.ConnectionProvider,
		core.JSON{&m.EtopAffiliateAccount},
		(*core.String)(&m.Code),
		(*core.String)(&m.ImageURL),
	}
}

func (m *Connection) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *Connections) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(Connections, 0, 128)
	for rows.Next() {
		m := new(Connection)
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

func (_ *Connection) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlConnection_Select)
	return nil
}

func (_ *Connections) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlConnection_Select)
	return nil
}

func (m *Connection) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlConnection_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(16)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms Connections) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlConnection_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(16)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *Connection) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlConnection_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlConnection_ListColsOnConflict)
	return nil
}

func (ms Connections) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlConnection_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlConnection_ListColsOnConflict)
	return nil
}

func (m *Connection) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("connection")
	w.WriteRawString(" SET ")
	if m.ID != 0 {
		flag = true
		w.WriteName("id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ID)
	}
	if m.Name != "" {
		flag = true
		w.WriteName("name")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Name)
	}
	if m.Status != 0 {
		flag = true
		w.WriteName("status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Status)
	}
	if m.PartnerID != 0 {
		flag = true
		w.WriteName("partner_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.PartnerID)
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
	if m.DriverConfig != nil {
		flag = true
		w.WriteName("driver_config")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.DriverConfig})
	}
	if m.Driver != "" {
		flag = true
		w.WriteName("driver")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Driver)
	}
	if m.ConnectionType != 0 {
		flag = true
		w.WriteName("connection_type")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ConnectionType)
	}
	if m.ConnectionSubtype != 0 {
		flag = true
		w.WriteName("connection_subtype")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ConnectionSubtype)
	}
	if m.ConnectionMethod != 0 {
		flag = true
		w.WriteName("connection_method")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ConnectionMethod)
	}
	if m.ConnectionProvider != 0 {
		flag = true
		w.WriteName("connection_provider")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ConnectionProvider)
	}
	if m.EtopAffiliateAccount != nil {
		flag = true
		w.WriteName("etop_affiliate_account")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.EtopAffiliateAccount})
	}
	if m.Code != "" {
		flag = true
		w.WriteName("code")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Code)
	}
	if m.ImageURL != "" {
		flag = true
		w.WriteName("image_url")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ImageURL)
	}
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *Connection) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlConnection_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(16)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type ConnectionHistory map[string]interface{}
type ConnectionHistories []map[string]interface{}

func (m *ConnectionHistory) SQLTableName() string  { return "history.\"connection\"" }
func (m ConnectionHistories) SQLTableName() string { return "history.\"connection\"" }

func (m *ConnectionHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlConnection_Select_history)
	return nil
}

func (m ConnectionHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlConnection_Select_history)
	return nil
}

func (m ConnectionHistory) ID() core.Interface           { return core.Interface{m["id"]} }
func (m ConnectionHistory) Name() core.Interface         { return core.Interface{m["name"]} }
func (m ConnectionHistory) Status() core.Interface       { return core.Interface{m["status"]} }
func (m ConnectionHistory) PartnerID() core.Interface    { return core.Interface{m["partner_id"]} }
func (m ConnectionHistory) CreatedAt() core.Interface    { return core.Interface{m["created_at"]} }
func (m ConnectionHistory) UpdatedAt() core.Interface    { return core.Interface{m["updated_at"]} }
func (m ConnectionHistory) DeletedAt() core.Interface    { return core.Interface{m["deleted_at"]} }
func (m ConnectionHistory) DriverConfig() core.Interface { return core.Interface{m["driver_config"]} }
func (m ConnectionHistory) Driver() core.Interface       { return core.Interface{m["driver"]} }
func (m ConnectionHistory) ConnectionType() core.Interface {
	return core.Interface{m["connection_type"]}
}
func (m ConnectionHistory) ConnectionSubtype() core.Interface {
	return core.Interface{m["connection_subtype"]}
}
func (m ConnectionHistory) ConnectionMethod() core.Interface {
	return core.Interface{m["connection_method"]}
}
func (m ConnectionHistory) ConnectionProvider() core.Interface {
	return core.Interface{m["connection_provider"]}
}
func (m ConnectionHistory) EtopAffiliateAccount() core.Interface {
	return core.Interface{m["etop_affiliate_account"]}
}
func (m ConnectionHistory) Code() core.Interface     { return core.Interface{m["code"]} }
func (m ConnectionHistory) ImageURL() core.Interface { return core.Interface{m["image_url"]} }

func (m *ConnectionHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 16)
	args := make([]interface{}, 16)
	for i := 0; i < 16; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(ConnectionHistory, 16)
	res["id"] = data[0]
	res["name"] = data[1]
	res["status"] = data[2]
	res["partner_id"] = data[3]
	res["created_at"] = data[4]
	res["updated_at"] = data[5]
	res["deleted_at"] = data[6]
	res["driver_config"] = data[7]
	res["driver"] = data[8]
	res["connection_type"] = data[9]
	res["connection_subtype"] = data[10]
	res["connection_method"] = data[11]
	res["connection_provider"] = data[12]
	res["etop_affiliate_account"] = data[13]
	res["code"] = data[14]
	res["image_url"] = data[15]
	*m = res
	return nil
}

func (ms *ConnectionHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 16)
	args := make([]interface{}, 16)
	for i := 0; i < 16; i++ {
		args[i] = &data[i]
	}
	res := make(ConnectionHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(ConnectionHistory)
		m["id"] = data[0]
		m["name"] = data[1]
		m["status"] = data[2]
		m["partner_id"] = data[3]
		m["created_at"] = data[4]
		m["updated_at"] = data[5]
		m["deleted_at"] = data[6]
		m["driver_config"] = data[7]
		m["driver"] = data[8]
		m["connection_type"] = data[9]
		m["connection_subtype"] = data[10]
		m["connection_method"] = data[11]
		m["connection_provider"] = data[12]
		m["etop_affiliate_account"] = data[13]
		m["code"] = data[14]
		m["image_url"] = data[15]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}

// Type ShopConnection represents table shop_connection
func sqlgenShopConnection(_ *ShopConnection) bool { return true }

type ShopConnections []*ShopConnection

const __sqlShopConnection_Table = "shop_connection"
const __sqlShopConnection_ListCols = "\"shop_id\",\"connection_id\",\"token\",\"token_expires_at\",\"status\",\"connection_states\",\"created_at\",\"updated_at\",\"deleted_at\",\"is_global\",\"external_data\""
const __sqlShopConnection_ListColsOnConflict = "\"shop_id\" = EXCLUDED.\"shop_id\",\"connection_id\" = EXCLUDED.\"connection_id\",\"token\" = EXCLUDED.\"token\",\"token_expires_at\" = EXCLUDED.\"token_expires_at\",\"status\" = EXCLUDED.\"status\",\"connection_states\" = EXCLUDED.\"connection_states\",\"created_at\" = EXCLUDED.\"created_at\",\"updated_at\" = EXCLUDED.\"updated_at\",\"deleted_at\" = EXCLUDED.\"deleted_at\",\"is_global\" = EXCLUDED.\"is_global\",\"external_data\" = EXCLUDED.\"external_data\""
const __sqlShopConnection_Insert = "INSERT INTO \"shop_connection\" (" + __sqlShopConnection_ListCols + ") VALUES"
const __sqlShopConnection_Select = "SELECT " + __sqlShopConnection_ListCols + " FROM \"shop_connection\""
const __sqlShopConnection_Select_history = "SELECT " + __sqlShopConnection_ListCols + " FROM history.\"shop_connection\""
const __sqlShopConnection_UpdateAll = "UPDATE \"shop_connection\" SET (" + __sqlShopConnection_ListCols + ")"
const __sqlShopConnection_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT shop_connection_pkey DO UPDATE SET"

func (m *ShopConnection) SQLTableName() string  { return "shop_connection" }
func (m *ShopConnections) SQLTableName() string { return "shop_connection" }
func (m *ShopConnection) SQLListCols() string   { return __sqlShopConnection_ListCols }

func (m *ShopConnection) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlShopConnection_ListCols + " FROM \"shop_connection\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*ShopConnection)(nil))
}

func (m *ShopConnection) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		m.ShopID,
		m.ConnectionID,
		core.String(m.Token),
		core.Time(m.TokenExpiresAt),
		m.Status,
		core.JSON{m.ConnectionStates},
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
		core.Time(m.DeletedAt),
		core.Bool(m.IsGlobal),
		core.JSON{m.ExternalData},
	}
}

func (m *ShopConnection) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ShopID,
		&m.ConnectionID,
		(*core.String)(&m.Token),
		(*core.Time)(&m.TokenExpiresAt),
		&m.Status,
		core.JSON{&m.ConnectionStates},
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		(*core.Time)(&m.DeletedAt),
		(*core.Bool)(&m.IsGlobal),
		core.JSON{&m.ExternalData},
	}
}

func (m *ShopConnection) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *ShopConnections) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(ShopConnections, 0, 128)
	for rows.Next() {
		m := new(ShopConnection)
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

func (_ *ShopConnection) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopConnection_Select)
	return nil
}

func (_ *ShopConnections) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopConnection_Select)
	return nil
}

func (m *ShopConnection) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShopConnection_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(11)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms ShopConnections) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShopConnection_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(11)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *ShopConnection) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlShopConnection_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlShopConnection_ListColsOnConflict)
	return nil
}

func (ms ShopConnections) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlShopConnection_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlShopConnection_ListColsOnConflict)
	return nil
}

func (m *ShopConnection) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("shop_connection")
	w.WriteRawString(" SET ")
	if m.ShopID != 0 {
		flag = true
		w.WriteName("shop_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShopID)
	}
	if m.ConnectionID != 0 {
		flag = true
		w.WriteName("connection_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ConnectionID)
	}
	if m.Token != "" {
		flag = true
		w.WriteName("token")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Token)
	}
	if !m.TokenExpiresAt.IsZero() {
		flag = true
		w.WriteName("token_expires_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.TokenExpiresAt)
	}
	if m.Status != 0 {
		flag = true
		w.WriteName("status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Status)
	}
	if m.ConnectionStates != nil {
		flag = true
		w.WriteName("connection_states")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.ConnectionStates})
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
	if m.IsGlobal {
		flag = true
		w.WriteName("is_global")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.IsGlobal)
	}
	if m.ExternalData != nil {
		flag = true
		w.WriteName("external_data")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.ExternalData})
	}
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *ShopConnection) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlShopConnection_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(11)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type ShopConnectionHistory map[string]interface{}
type ShopConnectionHistories []map[string]interface{}

func (m *ShopConnectionHistory) SQLTableName() string  { return "history.\"shop_connection\"" }
func (m ShopConnectionHistories) SQLTableName() string { return "history.\"shop_connection\"" }

func (m *ShopConnectionHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopConnection_Select_history)
	return nil
}

func (m ShopConnectionHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopConnection_Select_history)
	return nil
}

func (m ShopConnectionHistory) ShopID() core.Interface { return core.Interface{m["shop_id"]} }
func (m ShopConnectionHistory) ConnectionID() core.Interface {
	return core.Interface{m["connection_id"]}
}
func (m ShopConnectionHistory) Token() core.Interface { return core.Interface{m["token"]} }
func (m ShopConnectionHistory) TokenExpiresAt() core.Interface {
	return core.Interface{m["token_expires_at"]}
}
func (m ShopConnectionHistory) Status() core.Interface { return core.Interface{m["status"]} }
func (m ShopConnectionHistory) ConnectionStates() core.Interface {
	return core.Interface{m["connection_states"]}
}
func (m ShopConnectionHistory) CreatedAt() core.Interface { return core.Interface{m["created_at"]} }
func (m ShopConnectionHistory) UpdatedAt() core.Interface { return core.Interface{m["updated_at"]} }
func (m ShopConnectionHistory) DeletedAt() core.Interface { return core.Interface{m["deleted_at"]} }
func (m ShopConnectionHistory) IsGlobal() core.Interface  { return core.Interface{m["is_global"]} }
func (m ShopConnectionHistory) ExternalData() core.Interface {
	return core.Interface{m["external_data"]}
}

func (m *ShopConnectionHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 11)
	args := make([]interface{}, 11)
	for i := 0; i < 11; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(ShopConnectionHistory, 11)
	res["shop_id"] = data[0]
	res["connection_id"] = data[1]
	res["token"] = data[2]
	res["token_expires_at"] = data[3]
	res["status"] = data[4]
	res["connection_states"] = data[5]
	res["created_at"] = data[6]
	res["updated_at"] = data[7]
	res["deleted_at"] = data[8]
	res["is_global"] = data[9]
	res["external_data"] = data[10]
	*m = res
	return nil
}

func (ms *ShopConnectionHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 11)
	args := make([]interface{}, 11)
	for i := 0; i < 11; i++ {
		args[i] = &data[i]
	}
	res := make(ShopConnectionHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(ShopConnectionHistory)
		m["shop_id"] = data[0]
		m["connection_id"] = data[1]
		m["token"] = data[2]
		m["token_expires_at"] = data[3]
		m["status"] = data[4]
		m["connection_states"] = data[5]
		m["created_at"] = data[6]
		m["updated_at"] = data[7]
		m["deleted_at"] = data[8]
		m["is_global"] = data[9]
		m["external_data"] = data[10]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
