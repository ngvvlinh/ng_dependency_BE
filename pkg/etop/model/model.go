package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"etop.vn/api/top/types/etc/account_type"
	"etop.vn/api/top/types/etc/ghn_note_code"
	"etop.vn/api/top/types/etc/order_source"
	"etop.vn/api/top/types/etc/payment_method"
	"etop.vn/api/top/types/etc/shipping"
	"etop.vn/api/top/types/etc/shipping_fee_type"
	"etop.vn/api/top/types/etc/shipping_provider"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/status4"
	"etop.vn/api/top/types/etc/status5"
	"etop.vn/api/top/types/etc/try_on"
	"etop.vn/api/top/types/etc/user_source"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/code/gencode"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

type (
	SubjectType         string
	FulfillmentEndpoint string
	ShippingPriceSource string

	RoleType  string
	FfmAction int
	CodeType  string

	ShippingRouteType string

	DBName string
)

// Type constants
const (
	DBMain     DBName = "main"
	DBNotifier DBName = "notifier"

	SubjectTypeAccount SubjectType = "account"
	SubjectTypeUser    SubjectType = "user"

	TypeShippingSourceEtop    ShippingPriceSource = "etop"
	TypeShippingSourceCarrier ShippingPriceSource = "carrier"

	TypeGHNCode  = "GHN"
	TypeGHTKCode = "GHTK"

	ShippingServiceNameStandard = "Chuẩn"
	ShippingServiceNameFaster   = "Nhanh"

	FFShop     FulfillmentEndpoint = "shop"
	FFCustomer FulfillmentEndpoint = "customer"

	// Don't change these values
	TagUser      = 17
	TagPartner   = 21
	TagShop      = 33
	TagAffiliate = 35
	TagEtop      = 101
	TagImport    = 111

	EtopAccountID         = TagEtop
	EtopTradingAccountID  = 1000015764575267699
	TopShipID             = 1000030662086749358
	IndependentCustomerID = 1000080135776788835

	CurrencyVND = "vnd"

	AddressTypeGeneral   = "general"
	AddressTypeWarehouse = "warehouse"
	AddressTypeShipTo    = "shipto"
	AddressTypeShipFrom  = "shipfrom"

	CodeTypeShop  = "shop"
	CodeTypeOrder = "order"

	CodeTypeMoneyTransaction         = "money_transaction_shipping"
	CodeTypeMoneyTransactionEtop     = "money_transaction_shipping_etop"
	CodeTypeMoneyTransactionExternal = "money_transaction_external"
	CodeTypeConnection               = "connection"

	// For import

	FfmActionNothing FfmAction = iota
	FfmActionCancel
	FfmActionReturn

	RouteSameProvince ShippingRouteType = "noi_tinh"
	RouteSameRegion   ShippingRouteType = "noi_mien"
	RouteNationWide   ShippingRouteType = "toan_quoc"
)

type NormalizedPhone = validate.NormalizedPhone
type NormalizedEmail = validate.NormalizedEmail

var ShippingFeeShopTypes = []shipping_fee_type.ShippingFeeType{
	shipping_fee_type.Main, shipping_fee_type.Return, shipping_fee_type.Adjustment,
	shipping_fee_type.AddressChange, shipping_fee_type.Cods, shipping_fee_type.Insurance, shipping_fee_type.Other, shipping_fee_type.Discount,
}

//        SyncAt: updated when sending data to external service successfully
//     TrySyncAt: updated when sending data to external service (may unsuccessfully)
//    LastSyncAt: updated when querying external data source successfully
// LastTrySypkg/etop/model/model.go:18:ncAt: updated when querying external data source (may unsuccessfully)

var ShippingStateMap = map[shipping.State]string{
	shipping.Default:       "Mặc định",
	shipping.Created:       "Mới",
	shipping.Picking:       "Đang lấy hàng",
	shipping.Holding:       "Đã lấy hàng",
	shipping.Delivering:    "Đang giao hàng",
	shipping.Returning:     "Đang trả hàng",
	shipping.Delivered:     "Đã giao hàng",
	shipping.Returned:      "Đã trả hàng",
	shipping.Cancelled:     "Hủy",
	shipping.Undeliverable: "Bồi hoàn",
	shipping.Unknown:       "Không xác định",
}

func EtopPaymentStatusLabel(s status4.Status) string {
	switch s {
	case status4.N:
		return "Không cần thanh toán"
	case status4.Z:
		return "Chưa đối soát"
	case status4.P:
		return "Đã thanh toán"
	case status4.S:
		return "Cần thanh toán"
	default:
		return "Không xác định"
	}
}

func OrderStatusLabel(s status5.Status) string {
	switch s {
	case status5.NS:
		return "Trả hàng"
	case status5.N:
		return "Huỷ"
	case status5.Z:
		return "Mới"
	case status5.P:
		return "Thành công"
	case status5.S:
		return "Đang xử lý"
	default:
		return "Không xác định"
	}
}

func PaymentStatusLabel(s status4.Status) string {
	switch s {
	case 0:
		return ""
	case 1:
		return "Đã đối soát"
	case 2:
		return "Chưa đối soát"
	default:
		return "Không xác định"
	}
}

func AccountTypeLabel(t account_type.AccountType) string {
	switch t {
	case account_type.Shop:
		return "cửa hàng"
	case account_type.Partner:
		return "tài khoản"
	case account_type.Etop:
		return "tài khoản"
	default:
		return "tài khoản"
	}
}

func (t CodeType) Validate() error {
	switch t {
	case "":
		return cm.Error(cm.InvalidArgument, "Missing Type", nil)
	case CodeTypeShop,
		CodeTypeOrder,
		CodeTypeMoneyTransaction,
		CodeTypeMoneyTransactionExternal,
		CodeTypeMoneyTransactionEtop,
		CodeTypeConnection:
		// no-op
	default:
		return cm.Error(cm.InvalidArgument, "Type is not valid", nil)
	}
	return nil
}

func VerifyPaymentMethod(s payment_method.PaymentMethod) bool {
	switch s {
	case payment_method.COD, payment_method.Bank,
		payment_method.Other, payment_method.VTPay:
		return true
	}
	return false
}

func VerifyOrderSource(s order_source.Source) bool {
	switch s {
	case order_source.EtopPOS, order_source.EtopPXS, order_source.EtopCMX, order_source.TSApp, order_source.API, order_source.Import, order_source.EtopApp:
		return true
	}
	return false
}

func VerifyShippingProvider(s shipping_provider.ShippingProvider) bool {
	if s == 0 {
		return true
	}
	switch s {
	case shipping_provider.GHN, shipping_provider.GHTK, shipping_provider.VTPost, shipping_provider.Manual:
		return true
	}
	return false
}

func TryOnFromGHNNoteCode(c ghn_note_code.GHNNoteCode) try_on.TryOnCode {
	switch c {
	case ghn_note_code.KHONGCHOXEMHANG:
		return try_on.None
	case ghn_note_code.CHOXEMHANGKHONGTHU:
		return try_on.Open
	case ghn_note_code.CHOTHUHANG:
		return try_on.Try
	default:
		return 0
	}
}

func GHNNoteCodeFromTryOn(to try_on.TryOnCode) ghn_note_code.GHNNoteCode {
	switch to {
	case try_on.None:
		return ghn_note_code.KHONGCHOXEMHANG
	case try_on.Open:
		return ghn_note_code.CHOXEMHANGKHONGTHU
	case try_on.Try:
		return ghn_note_code.CHOTHUHANG
	default:
		return 0
	}
}

var _ = sqlgenAccount(&Account{})

type Account struct {
	ID       dot.ID
	OwnerID  dot.ID
	Name     string
	Type     account_type.AccountType
	ImageURL string
	URLSlug  string
}

type AccountInterface interface {
	GetAccount() *Account
}

func (s *Shop) GetAccount() *Account {
	return &Account{
		ID:       s.ID,
		OwnerID:  s.OwnerID,
		Name:     s.Name,
		Type:     account_type.Shop,
		ImageURL: s.ImageURL,
		URLSlug:  "",
	}
}

func (s *Partner) GetAccount() *Account {
	return &Account{
		ID:       s.ID,
		OwnerID:  s.OwnerID,
		Name:     s.Name,
		Type:     account_type.Partner,
		ImageURL: s.ImageURL,
		URLSlug:  "",
	}
}

var _ = sqlgenShop(&Shop{})

type Shop struct {
	ID      dot.ID
	Name    string
	OwnerID dot.ID
	IsTest  int

	AddressID         dot.ID
	ShipToAddressID   dot.ID
	ShipFromAddressID dot.ID
	Phone             string
	BankAccount       *BankAccount
	WebsiteURL        string
	ImageURL          string
	Email             string
	Code              string
	AutoCreateFFM     bool

	OrderSourceID dot.ID

	Status    status3.Status
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time

	Address *Address `sq:"-"`

	RecognizedHosts []string

	// @deprecated use try_on instead
	GhnNoteCode ghn_note_code.GHNNoteCode
	TryOn       try_on.TryOnCode
	CompanyInfo *CompanyInfo
	// MoneyTransactionRRule format:
	// FREQ=DAILY
	// FREQ=WEEKLY;BYDAY=TU,TH,SA
	// FREQ=MONTHLY;BYDAY=-1FR
	MoneyTransactionRRule         string `sq:"'money_transaction_rrule'"`
	SurveyInfo                    []*SurveyInfo
	ShippingServiceSelectStrategy []*ShippingServiceSelectStrategyItem

	InventoryOverstock dot.NullBool
}

func (s *Shop) CheckInfo() error {
	if s.MoneyTransactionRRule != "" {
		acceptStrings := []string{"FREQ", "BYMONTHDAY", "BYDAY", "BYSETPOS"}
		ss := strings.Split(s.MoneyTransactionRRule, ";")
		for _, s := range ss {
			prefix := strings.Split(s, "=")[0]
			if !cm.StringsContain(acceptStrings, prefix) {
				return cm.Errorf(cm.InvalidArgument, nil, "Định dạng MoneyTransactionRRule không hợp lệ. Định dạng đúng: FREQ=WEEKLY;BYDAY=TU,TH,SA. Vui lòng truy cập https://icalendar.org/rrule-tool.html để biết thêm chi tiết.")
			}
		}
	}
	if s.ShippingServiceSelectStrategy != nil {
		for _, item := range s.ShippingServiceSelectStrategy {
			key := strings.ToLower(item.Key)
			if key == "provider" || key == "shipping_provider" {
				return cm.Errorf(cm.InvalidArgument, nil, "Vui lòng sử dụng `carrier` thay vì %v", item.Key)
			}
		}
	}
	return nil
}

type ShippingServiceSelectStrategyItem struct {
	Key   string
	Value string
}

type SurveyInfo struct {
	Key      string `json:"key"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func (s *Shop) GetShopName() string {
	if s == nil {
		return ""
	}
	return s.Name
}

func (s *Shop) GetTryOn() try_on.TryOnCode {
	if s.TryOn != 0 {
		return s.TryOn
	}
	return TryOnFromGHNNoteCode(s.GhnNoteCode)
}

var _ = sqlgenShopExtended(
	&ShopExtended{}, &Shop{}, "s",
	sq.LEFT_JOIN, &Address{}, "a", "s.address_id = a.id",
	sq.LEFT_JOIN, &User{}, "u", "s.owner_id = u.id",
)

type ShopExtended struct {
	*Shop
	Address *Address
	User    *User
}

var _ = sqlgenShopDelete(&ShopDelete{}, &Shop{})

type ShopDelete struct {
	DeletedAt time.Time
}

var _ = sqlgenPartner(&Partner{})

type Partner struct {
	ID      dot.ID
	OwnerID dot.ID
	Status  status3.Status
	IsTest  int

	Name       string
	PublicName string
	Phone      string
	Email      string
	ImageURL   string
	WebsiteURL string

	ContactPersons  []*ContactPerson
	RecognizedHosts []string
	RedirectURLs    []string
	// AvailableFromEtop: dùng để xác định `partner` này có thể xác thực shop trực tiếp từ client của Etop hay không
	// Sau khi xác thực xong sẽ trỏ trực tiếp về `redirect_url` trong field AvailableFromEtopConfig để về trang xác thực của `partner`
	AvailableFromEtop       bool
	AvailableFromEtopConfig *AvailableFromEtopConfig

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time
}

type AvailableFromEtopConfig struct {
	RedirectUrl string `json:"redirect_url"`
}

func (p *Partner) GetID() dot.ID {
	if p == nil {
		return 0
	}
	return p.ID
}

func (p *Partner) Website() string {
	return validate.DomainFromURL(p.WebsiteURL)
}

func (p *Partner) BeforeInsert() error {
	return p.validate()
}

func (p *Partner) BeforeUpdate() error {
	return p.validate()
}

func (p *Partner) validate() error {
	if p.WebsiteURL != "" {
		if !validate.Host(p.WebsiteURL) {
			return cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ website cần có dạng http(s)://example.com")
		}
		if p.Website() == "" {
			return cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ website không hợp lệ")
		}
	}
	return nil
}

var _ = sqlgenAccountAuth(&AccountAuth{})

type AccountAuth struct {
	AuthKey     string
	AccountID   dot.ID
	Status      status3.Status
	Roles       []string
	Permissions []string
	CreatedAt   time.Time `sq:"create"`
	UpdatedAt   time.Time `sq:"update"`
	DeletedAt   time.Time
}

func (m *AccountAuth) BeforeInsert() error {
	if m.AccountID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "missing Name")
	}
	if m.AuthKey == "" {
		code := gencode.GenerateCode(gencode.Alphabet54, 32)
		m.AuthKey = fmt.Sprintf("%v:%v", m.AccountID, code)
	}
	return nil
}

var _ = sqlgenAccountAuthFtPartner(
	&AccountAuthFtPartner{}, &AccountAuth{}, "aa",
	sq.JOIN, &Partner{}, "p", `aa.account_id = p.id`,
)

type AccountAuthFtPartner struct {
	*AccountAuth
	*Partner
}

var _ = sqlgenAccountAuthFtShop(
	&AccountAuthFtShop{}, &AccountAuth{}, "aa",
	sq.JOIN, &Shop{}, "s", `aa.account_id = s.id`,
)

type AccountAuthFtShop struct {
	*AccountAuth
	*Shop
}

var _ = sqlgenPartnerRelation(&PartnerRelation{})

type PartnerRelation struct {
	AuthKey           string
	PartnerID         dot.ID
	SubjectID         dot.ID
	SubjectType       SubjectType
	ExternalSubjectID string

	Nonce     dot.ID
	Status    status3.Status
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time

	Permission `sq:"inline"`
}

var _ = sqlgenPartnerRelationFtShop(
	&PartnerRelationFtShop{}, &PartnerRelation{}, "pr",
	sq.JOIN, &Shop{}, "s", "pr.subject_id = s.id",
	sq.JOIN, &User{}, "u", "s.owner_id = u.id",
)

type PartnerRelationFtShop struct {
	*PartnerRelation
	*Shop
	*User
}

type Unit struct {
	ID       string  `json:"id"`
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	FullName string  `json:"full_name"`
	Unit     string  `json:"unit"`
	UnitConv float64 `json:"unit_conv"`
	Price    int     `json:"price"`
}

type UserInner struct {
	FullName  string
	ShortName string
	Email     string
	Phone     string
}

var _ = sqlgenUser(&User{})

type User struct {
	ID dot.ID

	UserInner `sq:"inline"`

	Status status3.Status // 1: actual user, 0: stub, -1: disabled

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`

	AgreedTOSAt       time.Time
	AgreedEmailInfoAt time.Time
	EmailVerifiedAt   time.Time
	PhoneVerifiedAt   time.Time

	EmailVerificationSentAt time.Time
	PhoneVerificationSentAt time.Time

	IsTest    int
	Source    user_source.UserSource
	RefUserID dot.ID
	RefSaleID dot.ID
}

type UserExtended struct {
	User         *User
	UserInternal *UserInternal
}

// Permission ...
type Permission struct {
	Roles       []string
	Permissions []string
}

var _ = sqlgenAccountUser(&AccountUser{})

type AccountUser struct {
	AccountID dot.ID
	UserID    dot.ID

	Status         status3.Status // 1: activated, -1: rejected/disabled, 0: pending
	ResponseStatus status3.Status // 1: accepted,  -1: rejected, 0: pending

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time

	Permission `sq:"inline"`

	FullName  string
	ShortName string
	Position  string

	InvitationSentAt     time.Time
	InvitationSentBy     dot.ID
	InvitationAcceptedAt time.Time
	InvitationRejectedAt time.Time

	DisabledAt    time.Time
	DisabledBy    time.Time
	DisableReason string
}

var _ = sqlgenAccountUserExtended(
	&AccountUserExtended{}, &AccountUser{}, "au",
	sq.JOIN, &Account{}, "a", "au.account_id = a.id",
	sq.JOIN, &User{}, "u", "au.user_id = u.id",
)

type AccountUserExtended struct {
	AccountUser *AccountUser
	Account     *Account
	User        *User
}

func (m *AccountUserExtended) GetUserName() (fullName, shortName string) {
	fullName = CoalesceString2(m.AccountUser.FullName, m.User.FullName)
	shortName = CoalesceString2(m.AccountUser.ShortName, m.User.ShortName)
	return
}

var _ = sqlgenAccountUserDelete(&AccountUserDelete{}, &AccountUser{})

type AccountUserDelete struct {
	DeletedAt time.Time
}

var _ = sqlgenUserAuth(&UserAuth{})

type UserAuth struct {
	UserID   dot.ID
	AuthType string
	AuthKey  string

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

var _ = sqlgenUserInternal(&UserInternal{})

type UserInternal struct {
	ID      dot.ID
	Hashpwd string

	UpdatedAt time.Time `sq:"update"`
}

type StatusQuery struct {
	Status *status3.Status
}

type ExternalShippingLog struct {
	StateText string `json:"StateText"`
	Time      string `json:"Time"`
	Message   string `json:"Message"`
}

type FulfillmentSyncStates struct {
	SyncAt    time.Time `json:"sync_at"`
	TrySyncAt time.Time `json:"try_sync_at"`
	Error     *Error    `json:"error"`

	NextShippingState shipping.State `json:"next_shipping_state"`
}

func (f FfmAction) ToShippingState() shipping.State {
	switch f {
	case FfmActionCancel:
		return shipping.Cancelled
	case FfmActionReturn:
		return shipping.Returning
	default:
		return shipping.Unknown
	}
}

var _ = sqlgenAddress(&Address{})

type Address struct {
	ID        dot.ID `json:"id"`
	FullName  string `json:"full_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Position  string `json:"position"`
	Email     string `json:"email"`

	Country  string `json:"country"`
	City     string `json:"city"`
	Province string `json:"province"`
	District string `json:"district"`
	Ward     string `json:"ward"` // Ward may be non-empty while WardCode is empty
	Zip      string `json:"zip"`

	DistrictCode string `json:"district_code"`
	ProvinceCode string `json:"province_code"`
	WardCode     string `json:"ward_code"`

	Company     string       `json:"company"`
	Address1    string       `json:"address1"`
	Address2    string       `json:"address2"`
	Type        string       `json:"type"`
	AccountID   dot.ID       `json:"account_id"`
	Notes       *AddressNote `json:"notes"`
	CreatedAt   time.Time    `sq:"create" json:"-"`
	UpdatedAt   time.Time    `sq:"update" json:"-"`
	Coordinates *Coordinates `json:"coordinates"`
}

type Coordinates struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

func (m *Address) GetFullName() string {
	if m == nil {
		return ""
	}
	if m.FullName != "" {
		return m.FullName
	}
	return m.FirstName + " " + m.LastName
}

func (m *Address) GetShortAddress() string {
	if m == nil {
		return ""
	}
	b := strings.Builder{}
	if m.Address1 != "" {
		b.WriteString(m.Address1)
		b.WriteByte('\n')
	}
	if m.Address2 != "" {
		b.WriteString(m.Address2)
		b.WriteByte('\n')
	}
	if m.Company != "" {
		b.WriteString(m.Company)
		b.WriteByte('\n')
	}
	s := b.String()
	if s == "" {
		return ""
	}
	return s[:len(s)-1]
}

func (m *Address) GetPhone() string {
	if m == nil {
		return ""
	}
	return m.Phone
}

func (m *Address) GetProvince() string {
	if m == nil {
		return ""
	}
	return m.Province
}

func (m *Address) GetDistrict() string {
	if m == nil {
		return ""
	}
	return m.District
}

func (m *Address) GetWard() string {
	if m == nil {
		return ""
	}
	return m.Ward
}

// This function uses Ward (instead of WardCode) because WardCode may be empty
// while Ward retains raw name.
func (m *Address) GetFullAddress() string {
	b := strings.Builder{}
	if m.Address1 != "" {
		b.WriteString(m.Address1)
		b.WriteByte('\n')
	}
	if m.Address2 != "" {
		b.WriteString(m.Address2)
		b.WriteByte('\n')
	}
	if m.Company != "" {
		b.WriteString(m.Company)
		b.WriteByte('\n')
	}
	flag := false
	if m.Ward != "" {
		b.WriteString(m.Ward)
		flag = true
	}
	if m.District != "" {
		if flag {
			b.WriteString(", ")
		}
		b.WriteString(m.District)
		flag = true
	}
	if m.Province != "" {
		if flag {
			b.WriteString(", ")
		}
		b.WriteString(m.Province)
	}
	return b.String()
}

func (m *Address) UpdateAddress(phone string, fullname string) *Address {
	if phone != "" {
		m.Phone = phone
	}
	if fullname != "" {
		m.FullName = fullname
	}
	return m
}

type AddressNote struct {
	Note       string `json:"note"`
	OpenTime   string `json:"open_time"`
	LunchBreak string `json:"lunch_break"`
	Other      string `json:"other"`
}

func (m *AddressNote) GetFullNote() string {
	if m == nil {
		return ""
	}
	b := strings.Builder{}
	if m.Other != "" {
		if m.Other == "call" {
			b.WriteString("Gọi trước khi đến")
		} else if m.Other == "no-call" {
			b.WriteString("Không cần gọi trước, shop đã chuẩn bị sẵn")
		} else {
			b.WriteString(m.Other)
		}
		b.WriteString(". \n")
	}

	if m.Note != "" {
		b.WriteString(m.Note)
		if m.Note[len(m.Note)-1] != '.' {
			b.WriteByte('.')
		}
		b.WriteString(" \n")
	}
	if m.OpenTime != "" {
		b.WriteString("Giờ làm việc: ")
		b.WriteString(m.OpenTime)
		// Nếu làm việc cả buổi tối thì thêm dòng ghi chú này vào:
		// "(nếu lấy hàng không kịp vui lòng lấy buổi tối)"
		// format OpenTime: "08:00 - 21:00"
		text := strings.Split(m.OpenTime, "-")
		if len(text) > 1 {
			closedAt := strings.Split(text[1], ":")[0]
			if hour, err := strconv.Atoi(strings.TrimSpace(closedAt)); err == nil {
				if hour > 19 {
					b.WriteString(" (nếu lấy hàng không kịp vui lòng lấy buổi tối)")
				}
			}
		}
		b.WriteString(". \n")
	}
	if m.LunchBreak != "" {
		b.WriteString("giờ nghỉ trưa: ")
		b.WriteString(m.LunchBreak)
		b.WriteString(". \n")
	}
	s := b.String()
	if s == "" {
		return ""
	}
	return s[:len(s)-1]
}

type BankAccount struct {
	Name          string `json:"name"`
	Province      string `json:"province"`
	Branch        string `json:"branch"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
}

var _ = sqlgenCode(&Code{})

type Code struct {
	Code      string
	Type      CodeType
	CreatedAt time.Time `sq:"create"`
}

func (code *Code) Validate() error {
	if code.Code == "" {
		return cm.Error(cm.InvalidArgument, "Missing Code", nil)
	}
	if err := code.Type.Validate(); err != nil {
		return err
	}
	return nil
}

type CreateCodeCommand struct {
	Code *Code
}

var _ = sqlgenCredit(&Credit{})

type Credit struct {
	ID        dot.ID
	Amount    int
	ShopID    dot.ID
	Type      string
	Status    status3.Status
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	PaidAt    time.Time
}

var _ = sqlgenCreditExtended(
	&CreditExtended{}, &Credit{}, "c",
	sq.LEFT_JOIN, &Shop{}, "s", "s.id = c.shop_id",
)

type CreditExtended struct {
	*Credit
	Shop *Shop
}

type ShippingFeeLine struct {
	ShippingFeeType          shipping_fee_type.ShippingFeeType `json:"shipping_fee_type"`
	Cost                     int                               `json:"cost"`
	ExternalServiceID        string                            `json:"external_service_id"`
	ExternalServiceName      string                            `json:"external_service_name"`
	ExternalServiceType      string                            `json:"external_service_type"`
	ExternalShippingOrderID  string                            `json:"external_order_id"`
	ExternalPaymentChannelID string                            `json:"external_payment_channel_id"`
	ExternalShippingCode     string                            `json:"external_shipping_code"`
}

func GetShippingFeeShopLines(items []*ShippingFeeLine, etopPriceRule bool, mainFee dot.NullInt) []*ShippingFeeLine {
	res := make([]*ShippingFeeLine, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		line := GetShippingFeeShopLine(*item, etopPriceRule, mainFee)
		if line != nil {
			res = append(res, line)
		}
	}
	return res
}

func GetShippingFeeShopLine(item ShippingFeeLine, etopPriceRule bool, mainFee dot.NullInt) *ShippingFeeLine {
	if item.ShippingFeeType == shipping_fee_type.Main && etopPriceRule {
		item.Cost = mainFee.Apply(item.Cost)
	}
	if contains(ShippingFeeShopTypes, item.ShippingFeeType) {
		return &item
	}
	return nil
}

func GetReturnedFee(items []*ShippingFeeLine) int {
	result := 0
	for _, item := range items {
		if item.ShippingFeeType == shipping_fee_type.Return {
			result = item.Cost
			break
		}
	}
	return result
}

func GetMainFee(items []*ShippingFeeLine) int {
	result := 0
	for _, item := range items {
		if item.ShippingFeeType == shipping_fee_type.Main {
			result = item.Cost
			break
		}
	}
	return result
}

func GetTotalShippingFee(items []*ShippingFeeLine) int {
	result := 0
	for _, item := range items {
		result += item.Cost
	}
	return result
}

func UpdateShippingFees(items []*ShippingFeeLine, fee int, shippingFeeType shipping_fee_type.ShippingFeeType) []*ShippingFeeLine {
	if fee == 0 {
		return items
	}
	found := false
	for _, item := range items {
		if item.ShippingFeeType == shippingFeeType {
			item.Cost = fee
			found = true
		}
	}
	if !found {
		items = append(items, &ShippingFeeLine{
			ShippingFeeType: shippingFeeType,
			Cost:            fee,
		})
	}
	return items
}

func contains(lines []shipping_fee_type.ShippingFeeType, feeType shipping_fee_type.ShippingFeeType) bool {
	for _, line := range lines {
		if feeType == line {
			return true
		}
	}
	return false
}

var _ = sqlgenShippingSource(&ShippingSource{})

type ShippingSource struct {
	ID        dot.ID
	Name      string
	Username  string
	Type      string
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

var _ = sqlgenShippingSourceInternal(&ShippingSourceInternal{})

type ShippingSourceInternal struct {
	ID          dot.ID
	CreatedAt   time.Time `sq:"create"`
	UpdatedAt   time.Time `sq:"update"`
	LastSyncAt  time.Time
	AccessToken string
	ExpiresAt   time.Time `json:"expires_at"`
	Secret      *ShippingSourceSecret
}

type ShippingSourceSecret struct {
	CustomerID int    `json:"CustomerID"`
	Username   string `json:"Username"`
	Password   string `json:"Password"`

	// Use for VTPost
	GroupAddressID int `json:"GroupAddressID"`
}

type AvailableShippingService struct {
	Name string

	// ServiceFee: Tổng phí giao hàng (đã bao gồm phí chính + các phụ phí khác)
	ServiceFee int

	// ShippingFeeMain: Phí chính giao hàng
	ShippingFeeMain   int
	Provider          shipping_provider.ShippingProvider
	ProviderServiceID string

	ExpectedPickAt     time.Time
	ExpectedDeliveryAt time.Time
	Source             ShippingPriceSource
	ConnectionInfo     *ConnectionInfo
}

type ConnectionInfo struct {
	ID       dot.ID
	Name     string
	ImageURL string
}

func (service *AvailableShippingService) ApplyFeeMain(feeMain int) {
	service.ServiceFee = service.ServiceFee - service.ShippingFeeMain + feeMain
	service.ShippingFeeMain = feeMain
}

var _ = sqlgenWebhook(&Webhook{})

type Webhook struct {
	ID        dot.ID
	AccountID dot.ID
	Entities  []string
	Fields    []string
	URL       string
	Metadata  string
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time
}

func (m *Webhook) BeforeInsert() error {
	if m == nil {
		return cm.Errorf(cm.InvalidArgument, nil, "empty data")
	}
	if m.AccountID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "missing Name")
	}
	if !validate.URL(m.URL) {
		return cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ url không hợp lệ")
	}
	if cm.IsProd() && !strings.HasPrefix(m.URL, "https://") {
		return cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ url phải là https://")
	}
	if len(m.Entities) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "missing entity")
	}
	if len(m.Fields) > 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Thông tin fields chưa được hỗ trợ, vui lòng để trống")
	}

	mp := make(map[string]bool)
	for _, item := range m.Entities {
		if !validate.LowercaseID(item) {
			return cm.Errorf(cm.InvalidArgument, nil, `invalid entity: "%v"`, item)
		}
		switch item {
		case "order", "fulfillment":
			if mp[item] {
				return cm.Errorf(cm.InvalidArgument, nil, `duplicated entity: "%v"`, item)
			}
			mp[item] = true
		default:
			return cm.Errorf(cm.InvalidArgument, nil, `unknown entity: "%v"`, item)
		}
	}

	m.ID = cm.NewID()
	return nil
}

var _ = sqlgenChangesData(&Callback{})

type Callback struct {
	ID        dot.ID
	WebhookID dot.ID
	AccountID dot.ID
	CreatedAt time.Time `sq:"create"`
	Changes   json.RawMessage
	Result    json.RawMessage // WebhookStatesError
}

type CompanyInfo struct {
	Name                string         `json:"name"`
	TaxCode             string         `json:"tax_code"`
	Address             string         `json:"address"`
	Website             string         `json:"website"`
	LegalRepresentative *ContactPerson `json:"legal_representative"`
}

type ContactPerson struct {
	Name     string `json:"name"`
	Position string `json:"position"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}
