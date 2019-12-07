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

// Type ExternalAccountAhamove represents table external_account_ahamove
func sqlgenExternalAccountAhamove(_ *ExternalAccountAhamove) bool { return true }

type ExternalAccountAhamoves []*ExternalAccountAhamove

const __sqlExternalAccountAhamove_Table = "external_account_ahamove"
const __sqlExternalAccountAhamove_ListCols = "\"id\",\"owner_id\",\"phone\",\"name\",\"external_id\",\"external_verified\",\"external_created_at\",\"external_token\",\"created_at\",\"updated_at\",\"last_send_verified_at\",\"external_ticket_id\",\"id_card_front_img\",\"id_card_back_img\",\"portrait_img\",\"website_url\",\"fanpage_url\",\"company_imgs\",\"business_license_imgs\",\"external_data_verified\",\"uploaded_at\""
const __sqlExternalAccountAhamove_Insert = "INSERT INTO \"external_account_ahamove\" (" + __sqlExternalAccountAhamove_ListCols + ") VALUES"
const __sqlExternalAccountAhamove_Select = "SELECT " + __sqlExternalAccountAhamove_ListCols + " FROM \"external_account_ahamove\""
const __sqlExternalAccountAhamove_Select_history = "SELECT " + __sqlExternalAccountAhamove_ListCols + " FROM history.\"external_account_ahamove\""
const __sqlExternalAccountAhamove_UpdateAll = "UPDATE \"external_account_ahamove\" SET (" + __sqlExternalAccountAhamove_ListCols + ")"

func (m *ExternalAccountAhamove) SQLTableName() string  { return "external_account_ahamove" }
func (m *ExternalAccountAhamoves) SQLTableName() string { return "external_account_ahamove" }
func (m *ExternalAccountAhamove) SQLListCols() string   { return __sqlExternalAccountAhamove_ListCols }

func (m *ExternalAccountAhamove) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlExternalAccountAhamove_ListCols + " FROM \"external_account_ahamove\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*ExternalAccountAhamove)(nil))
}

func (m *ExternalAccountAhamove) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		m.ID,
		m.OwnerID,
		core.String(m.Phone),
		core.String(m.Name),
		core.String(m.ExternalID),
		core.Bool(m.ExternalVerified),
		core.Time(m.ExternalCreatedAt),
		core.String(m.ExternalToken),
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
		core.Time(m.LastSendVerifiedAt),
		core.String(m.ExternalTicketID),
		core.String(m.IDCardFrontImg),
		core.String(m.IDCardBackImg),
		core.String(m.PortraitImg),
		core.String(m.WebsiteURL),
		core.String(m.FanpageURL),
		core.Array{m.CompanyImgs, opts},
		core.Array{m.BusinessLicenseImgs, opts},
		core.JSON{m.ExternalDataVerified},
		core.Time(m.UploadedAt),
	}
}

func (m *ExternalAccountAhamove) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		&m.OwnerID,
		(*core.String)(&m.Phone),
		(*core.String)(&m.Name),
		(*core.String)(&m.ExternalID),
		(*core.Bool)(&m.ExternalVerified),
		(*core.Time)(&m.ExternalCreatedAt),
		(*core.String)(&m.ExternalToken),
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		(*core.Time)(&m.LastSendVerifiedAt),
		(*core.String)(&m.ExternalTicketID),
		(*core.String)(&m.IDCardFrontImg),
		(*core.String)(&m.IDCardBackImg),
		(*core.String)(&m.PortraitImg),
		(*core.String)(&m.WebsiteURL),
		(*core.String)(&m.FanpageURL),
		core.Array{&m.CompanyImgs, opts},
		core.Array{&m.BusinessLicenseImgs, opts},
		core.JSON{&m.ExternalDataVerified},
		(*core.Time)(&m.UploadedAt),
	}
}

func (m *ExternalAccountAhamove) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *ExternalAccountAhamoves) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(ExternalAccountAhamoves, 0, 128)
	for rows.Next() {
		m := new(ExternalAccountAhamove)
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

func (_ *ExternalAccountAhamove) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlExternalAccountAhamove_Select)
	return nil
}

func (_ *ExternalAccountAhamoves) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlExternalAccountAhamove_Select)
	return nil
}

func (m *ExternalAccountAhamove) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlExternalAccountAhamove_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(21)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms ExternalAccountAhamoves) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlExternalAccountAhamove_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(21)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *ExternalAccountAhamove) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("external_account_ahamove")
	w.WriteRawString(" SET ")
	if m.ID != 0 {
		flag = true
		w.WriteName("id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ID)
	}
	if m.OwnerID != 0 {
		flag = true
		w.WriteName("owner_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.OwnerID)
	}
	if m.Phone != "" {
		flag = true
		w.WriteName("phone")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Phone)
	}
	if m.Name != "" {
		flag = true
		w.WriteName("name")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Name)
	}
	if m.ExternalID != "" {
		flag = true
		w.WriteName("external_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ExternalID)
	}
	if m.ExternalVerified {
		flag = true
		w.WriteName("external_verified")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ExternalVerified)
	}
	if !m.ExternalCreatedAt.IsZero() {
		flag = true
		w.WriteName("external_created_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ExternalCreatedAt)
	}
	if m.ExternalToken != "" {
		flag = true
		w.WriteName("external_token")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ExternalToken)
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
	if !m.LastSendVerifiedAt.IsZero() {
		flag = true
		w.WriteName("last_send_verified_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.LastSendVerifiedAt)
	}
	if m.ExternalTicketID != "" {
		flag = true
		w.WriteName("external_ticket_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ExternalTicketID)
	}
	if m.IDCardFrontImg != "" {
		flag = true
		w.WriteName("id_card_front_img")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.IDCardFrontImg)
	}
	if m.IDCardBackImg != "" {
		flag = true
		w.WriteName("id_card_back_img")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.IDCardBackImg)
	}
	if m.PortraitImg != "" {
		flag = true
		w.WriteName("portrait_img")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.PortraitImg)
	}
	if m.WebsiteURL != "" {
		flag = true
		w.WriteName("website_url")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.WebsiteURL)
	}
	if m.FanpageURL != "" {
		flag = true
		w.WriteName("fanpage_url")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.FanpageURL)
	}
	if m.CompanyImgs != nil {
		flag = true
		w.WriteName("company_imgs")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.CompanyImgs, opts})
	}
	if m.BusinessLicenseImgs != nil {
		flag = true
		w.WriteName("business_license_imgs")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.BusinessLicenseImgs, opts})
	}
	if m.ExternalDataVerified != nil {
		flag = true
		w.WriteName("external_data_verified")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.ExternalDataVerified})
	}
	if !m.UploadedAt.IsZero() {
		flag = true
		w.WriteName("uploaded_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.UploadedAt)
	}
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *ExternalAccountAhamove) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlExternalAccountAhamove_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(21)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type ExternalAccountAhamoveHistory map[string]interface{}
type ExternalAccountAhamoveHistories []map[string]interface{}

func (m *ExternalAccountAhamoveHistory) SQLTableName() string {
	return "history.\"external_account_ahamove\""
}
func (m ExternalAccountAhamoveHistories) SQLTableName() string {
	return "history.\"external_account_ahamove\""
}

func (m *ExternalAccountAhamoveHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlExternalAccountAhamove_Select_history)
	return nil
}

func (m ExternalAccountAhamoveHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlExternalAccountAhamove_Select_history)
	return nil
}

func (m ExternalAccountAhamoveHistory) ID() core.Interface      { return core.Interface{m["id"]} }
func (m ExternalAccountAhamoveHistory) OwnerID() core.Interface { return core.Interface{m["owner_id"]} }
func (m ExternalAccountAhamoveHistory) Phone() core.Interface   { return core.Interface{m["phone"]} }
func (m ExternalAccountAhamoveHistory) Name() core.Interface    { return core.Interface{m["name"]} }
func (m ExternalAccountAhamoveHistory) ExternalID() core.Interface {
	return core.Interface{m["external_id"]}
}
func (m ExternalAccountAhamoveHistory) ExternalVerified() core.Interface {
	return core.Interface{m["external_verified"]}
}
func (m ExternalAccountAhamoveHistory) ExternalCreatedAt() core.Interface {
	return core.Interface{m["external_created_at"]}
}
func (m ExternalAccountAhamoveHistory) ExternalToken() core.Interface {
	return core.Interface{m["external_token"]}
}
func (m ExternalAccountAhamoveHistory) CreatedAt() core.Interface {
	return core.Interface{m["created_at"]}
}
func (m ExternalAccountAhamoveHistory) UpdatedAt() core.Interface {
	return core.Interface{m["updated_at"]}
}
func (m ExternalAccountAhamoveHistory) LastSendVerifiedAt() core.Interface {
	return core.Interface{m["last_send_verified_at"]}
}
func (m ExternalAccountAhamoveHistory) ExternalTicketID() core.Interface {
	return core.Interface{m["external_ticket_id"]}
}
func (m ExternalAccountAhamoveHistory) IDCardFrontImg() core.Interface {
	return core.Interface{m["id_card_front_img"]}
}
func (m ExternalAccountAhamoveHistory) IDCardBackImg() core.Interface {
	return core.Interface{m["id_card_back_img"]}
}
func (m ExternalAccountAhamoveHistory) PortraitImg() core.Interface {
	return core.Interface{m["portrait_img"]}
}
func (m ExternalAccountAhamoveHistory) WebsiteURL() core.Interface {
	return core.Interface{m["website_url"]}
}
func (m ExternalAccountAhamoveHistory) FanpageURL() core.Interface {
	return core.Interface{m["fanpage_url"]}
}
func (m ExternalAccountAhamoveHistory) CompanyImgs() core.Interface {
	return core.Interface{m["company_imgs"]}
}
func (m ExternalAccountAhamoveHistory) BusinessLicenseImgs() core.Interface {
	return core.Interface{m["business_license_imgs"]}
}
func (m ExternalAccountAhamoveHistory) ExternalDataVerified() core.Interface {
	return core.Interface{m["external_data_verified"]}
}
func (m ExternalAccountAhamoveHistory) UploadedAt() core.Interface {
	return core.Interface{m["uploaded_at"]}
}

func (m *ExternalAccountAhamoveHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 21)
	args := make([]interface{}, 21)
	for i := 0; i < 21; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(ExternalAccountAhamoveHistory, 21)
	res["id"] = data[0]
	res["owner_id"] = data[1]
	res["phone"] = data[2]
	res["name"] = data[3]
	res["external_id"] = data[4]
	res["external_verified"] = data[5]
	res["external_created_at"] = data[6]
	res["external_token"] = data[7]
	res["created_at"] = data[8]
	res["updated_at"] = data[9]
	res["last_send_verified_at"] = data[10]
	res["external_ticket_id"] = data[11]
	res["id_card_front_img"] = data[12]
	res["id_card_back_img"] = data[13]
	res["portrait_img"] = data[14]
	res["website_url"] = data[15]
	res["fanpage_url"] = data[16]
	res["company_imgs"] = data[17]
	res["business_license_imgs"] = data[18]
	res["external_data_verified"] = data[19]
	res["uploaded_at"] = data[20]
	*m = res
	return nil
}

func (ms *ExternalAccountAhamoveHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 21)
	args := make([]interface{}, 21)
	for i := 0; i < 21; i++ {
		args[i] = &data[i]
	}
	res := make(ExternalAccountAhamoveHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(ExternalAccountAhamoveHistory)
		m["id"] = data[0]
		m["owner_id"] = data[1]
		m["phone"] = data[2]
		m["name"] = data[3]
		m["external_id"] = data[4]
		m["external_verified"] = data[5]
		m["external_created_at"] = data[6]
		m["external_token"] = data[7]
		m["created_at"] = data[8]
		m["updated_at"] = data[9]
		m["last_send_verified_at"] = data[10]
		m["external_ticket_id"] = data[11]
		m["id_card_front_img"] = data[12]
		m["id_card_back_img"] = data[13]
		m["portrait_img"] = data[14]
		m["website_url"] = data[15]
		m["fanpage_url"] = data[16]
		m["company_imgs"] = data[17]
		m["business_license_imgs"] = data[18]
		m["external_data_verified"] = data[19]
		m["uploaded_at"] = data[20]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}

// Type Affiliate represents table affiliate
func sqlgenSale(_ *Affiliate) bool { return true }

type Affiliates []*Affiliate

const __sqlAffiliate_Table = "affiliate"
const __sqlAffiliate_ListCols = "\"id\",\"owner_id\",\"name\",\"phone\",\"email\",\"is_test\",\"status\",\"created_at\",\"updated_at\",\"deleted_at\",\"bank_account\""
const __sqlAffiliate_Insert = "INSERT INTO \"affiliate\" (" + __sqlAffiliate_ListCols + ") VALUES"
const __sqlAffiliate_Select = "SELECT " + __sqlAffiliate_ListCols + " FROM \"affiliate\""
const __sqlAffiliate_Select_history = "SELECT " + __sqlAffiliate_ListCols + " FROM history.\"affiliate\""
const __sqlAffiliate_UpdateAll = "UPDATE \"affiliate\" SET (" + __sqlAffiliate_ListCols + ")"

func (m *Affiliate) SQLTableName() string  { return "affiliate" }
func (m *Affiliates) SQLTableName() string { return "affiliate" }
func (m *Affiliate) SQLListCols() string   { return __sqlAffiliate_ListCols }

func (m *Affiliate) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlAffiliate_ListCols + " FROM \"affiliate\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*Affiliate)(nil))
}

func (m *Affiliate) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		m.ID,
		m.OwnerID,
		core.String(m.Name),
		core.String(m.Phone),
		core.String(m.Email),
		core.Int(m.IsTest),
		m.Status,
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
		core.Time(m.DeletedAt),
		core.JSON{m.BankAccount},
	}
}

func (m *Affiliate) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		&m.OwnerID,
		(*core.String)(&m.Name),
		(*core.String)(&m.Phone),
		(*core.String)(&m.Email),
		(*core.Int)(&m.IsTest),
		&m.Status,
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		(*core.Time)(&m.DeletedAt),
		core.JSON{&m.BankAccount},
	}
}

func (m *Affiliate) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *Affiliates) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(Affiliates, 0, 128)
	for rows.Next() {
		m := new(Affiliate)
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

func (_ *Affiliate) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlAffiliate_Select)
	return nil
}

func (_ *Affiliates) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlAffiliate_Select)
	return nil
}

func (m *Affiliate) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlAffiliate_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(11)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms Affiliates) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlAffiliate_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(11)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *Affiliate) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("affiliate")
	w.WriteRawString(" SET ")
	if m.ID != 0 {
		flag = true
		w.WriteName("id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ID)
	}
	if m.OwnerID != 0 {
		flag = true
		w.WriteName("owner_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.OwnerID)
	}
	if m.Name != "" {
		flag = true
		w.WriteName("name")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Name)
	}
	if m.Phone != "" {
		flag = true
		w.WriteName("phone")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Phone)
	}
	if m.Email != "" {
		flag = true
		w.WriteName("email")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Email)
	}
	if m.IsTest != 0 {
		flag = true
		w.WriteName("is_test")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.IsTest)
	}
	if m.Status != 0 {
		flag = true
		w.WriteName("status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Status)
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
	if m.BankAccount != nil {
		flag = true
		w.WriteName("bank_account")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.BankAccount})
	}
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *Affiliate) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlAffiliate_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(11)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type AffiliateHistory map[string]interface{}
type AffiliateHistories []map[string]interface{}

func (m *AffiliateHistory) SQLTableName() string  { return "history.\"affiliate\"" }
func (m AffiliateHistories) SQLTableName() string { return "history.\"affiliate\"" }

func (m *AffiliateHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlAffiliate_Select_history)
	return nil
}

func (m AffiliateHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlAffiliate_Select_history)
	return nil
}

func (m AffiliateHistory) ID() core.Interface          { return core.Interface{m["id"]} }
func (m AffiliateHistory) OwnerID() core.Interface     { return core.Interface{m["owner_id"]} }
func (m AffiliateHistory) Name() core.Interface        { return core.Interface{m["name"]} }
func (m AffiliateHistory) Phone() core.Interface       { return core.Interface{m["phone"]} }
func (m AffiliateHistory) Email() core.Interface       { return core.Interface{m["email"]} }
func (m AffiliateHistory) IsTest() core.Interface      { return core.Interface{m["is_test"]} }
func (m AffiliateHistory) Status() core.Interface      { return core.Interface{m["status"]} }
func (m AffiliateHistory) CreatedAt() core.Interface   { return core.Interface{m["created_at"]} }
func (m AffiliateHistory) UpdatedAt() core.Interface   { return core.Interface{m["updated_at"]} }
func (m AffiliateHistory) DeletedAt() core.Interface   { return core.Interface{m["deleted_at"]} }
func (m AffiliateHistory) BankAccount() core.Interface { return core.Interface{m["bank_account"]} }

func (m *AffiliateHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 11)
	args := make([]interface{}, 11)
	for i := 0; i < 11; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(AffiliateHistory, 11)
	res["id"] = data[0]
	res["owner_id"] = data[1]
	res["name"] = data[2]
	res["phone"] = data[3]
	res["email"] = data[4]
	res["is_test"] = data[5]
	res["status"] = data[6]
	res["created_at"] = data[7]
	res["updated_at"] = data[8]
	res["deleted_at"] = data[9]
	res["bank_account"] = data[10]
	*m = res
	return nil
}

func (ms *AffiliateHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 11)
	args := make([]interface{}, 11)
	for i := 0; i < 11; i++ {
		args[i] = &data[i]
	}
	res := make(AffiliateHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(AffiliateHistory)
		m["id"] = data[0]
		m["owner_id"] = data[1]
		m["name"] = data[2]
		m["phone"] = data[3]
		m["email"] = data[4]
		m["is_test"] = data[5]
		m["status"] = data[6]
		m["created_at"] = data[7]
		m["updated_at"] = data[8]
		m["deleted_at"] = data[9]
		m["bank_account"] = data[10]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
