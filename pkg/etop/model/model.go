package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/gencode"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/validate"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

type Status3 int

const (
	S3Negative Status3 = -1 + iota
	S3Zero
	S3Positive
)

type Status4 int

const (
	S4Negative Status4 = -1 + iota
	S4Zero
	S4Positive
	S4SuperPos
)

// -1:cancelled, 0:default, 1:delivered, 2:processing
//
// 0: just created, still error, can be edited
// 2: processing, can not be edited
// 1: done
//-1: cancelled
//-2: returned
type Status5 int

const (
	S5NegSuper Status5 = -2 + iota
	S5Negative
	S5Zero
	S5Positive
	S5SuperPos
)

func (s Status3) P() *Status3 {
	return &s
}

func (s Status4) P() *Status4 {
	return &s
}

type (
	AccountType         string
	SubjectType         string
	ProcessingStatus    string
	ShippingState       string
	FulfillmentEndpoint string
	ShippingProvider    string
	ShippingPriceSource string
	OrderSourceType     string
	ShippingFeeLineType string
	UserIdentifying     string
	RoleType            string
	FfmAction           int
	CodeType            string
	TryOn               string

	ShippingRouteType    string
	ShippingDistrictType string
	DBName               string
	UserSource           string
	EtopPaymentStatus    Status4
)

func (s ShippingProvider) ToString() string {
	return string(s)
}

// Type constants
const (
	DBMain     DBName = "main"
	DBNotifier DBName = "notifier"

	TypePartner   AccountType = "partner"
	TypeShop      AccountType = "shop"
	TypeEtop      AccountType = "etop"
	TypeAffiliate AccountType = "affiliate"

	SubjectTypeAccount SubjectType = "account"
	SubjectTypeUser    SubjectType = "user"

	TypeGHN                    ShippingProvider = "ghn"
	TypeGHTK                   ShippingProvider = "ghtk"
	TypeVTPost                 ShippingProvider = "vtpost"
	TypeShippingProviderManual ShippingProvider = "manual"
	TypeShippingETOP           ShippingProvider = "etop"

	TypeShippingSourceEtop    ShippingPriceSource = "etop"
	TypeShippingSourceCarrier ShippingPriceSource = "carrier"

	TypeGHNCode  = "GHN"
	TypeGHTKCode = "GHTK"

	ShippingServiceNameStandard = "Chuẩn"
	ShippingServiceNameFaster   = "Nhanh"

	FFShop     FulfillmentEndpoint = "shop"
	FFCustomer FulfillmentEndpoint = "customer"

	OrderSourceUnknown OrderSourceType = "unknown"
	OrderSourceImport  OrderSourceType = "import"
	OrderSourceAPI     OrderSourceType = "api"
	OrderSourceSelf    OrderSourceType = "self"
	OrderSourcePOS     OrderSourceType = "etop_pos"
	OrderSourcePXS     OrderSourceType = "etop_pxs"
	OrderSourceCMX     OrderSourceType = "etop_cmx"
	OrderSourceTSApp   OrderSourceType = "ts_app"

	UserIdentifyingFull UserIdentifying = "full"
	UserIdentifyingHalf UserIdentifying = "half"
	UserIdentifyingStub UserIdentifying = "stub"

	// Don't change these values
	TagUser      = 17
	TagPartner   = 21
	TagShop      = 33
	TagAffiliate = 35
	TagEtop      = 101
	TagImport    = 111

	EtopAccountID        = TagEtop
	EtopTradingAccountID = 1000015764575267699
	StatusActive         = 1
	StatusCreated        = 0
	StatusDisabled       = -1

	StateDefault    ShippingState = "default"    //  0
	StateCreated    ShippingState = "created"    //  2
	StatePicking    ShippingState = "picking"    //  2
	StateHolding    ShippingState = "holding"    //  2
	StateDelivering ShippingState = "delivering" //  2
	StateReturning  ShippingState = "returning"  //  2
	StateDelivered  ShippingState = "delivered"  //  0
	StateReturned   ShippingState = "returned"   // -1
	StateCancelled  ShippingState = "cancelled"  // -1
	StateClosed     ShippingState = "closed"     //  1

	// Trạng thái Bồi hoàn
	StateUndeliverable ShippingState = "undeliverable" //  2
	StateUnknown       ShippingState = "unknown"       //  2

	//StateConfirmed     ShippingState = "confirmed"     //  0
	//StateProcessing    ShippingState = "processing"    //  2

	ShippingFeeTypeMain         ShippingFeeLineType = "main"           // phí dịch vụ
	ShippingFeeTypeReturn       ShippingFeeLineType = "return"         // phí trả hàng
	ShippingFeeTypeAdjustment   ShippingFeeLineType = "adjustment"     // phí chênh lệnh giá
	ShippingFeeTypeCODS         ShippingFeeLineType = "cods"           // phí thu hộ (CODS) - VTPOST tính
	ShippingFeeTypeInsurance    ShippingFeeLineType = "insurance"      // phí khai giá / phí bảo hiểm
	ShippingFeeTypeAddessChange ShippingFeeLineType = "address_change" // Phí thay đổi địa chỉ giao
	ShippingFeeTypeDiscount     ShippingFeeLineType = "discount"       // giảm giá
	ShippingFeeTypeOther        ShippingFeeLineType = "other"          // phí khác

	RoleOwner  RoleType = "owner"
	RoleAdmin  RoleType = "admin"
	RoleStaff  RoleType = "staff"
	Role3rd    RoleType = "3rd" // Delegate to third party API, which can only update/cancel itself orders
	RoleViewer RoleType = "viewer"

	PaymentMethodCOD   = "cod"
	PaymentMethodBank  = "bank"
	PaymentMethodOther = "other"
	PaymentMethodVtpay = "vtpay"

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

	// For import
	GHNNoteChoThuHang   = "cho thu hang"
	GHNNoteChoXemHang   = "cho xem hang khong thu"
	GHNNoteKhongXemHang = "khong cho xem hang"

	GHNNoteCodeChoThuHang   = "CHOTHUHANG"
	GHNNoteCodeChoXemHang   = "CHOXEMHANGKHONGTHU"
	GHNNoteCodeKhongXemHang = "KHONGCHOXEMHANG"

	TryOnNone TryOn = "none"
	TryOnOpen TryOn = "open"
	TryOnTry  TryOn = "try"

	FfmActionNothing FfmAction = iota
	FfmActionCancel
	FfmActionReturn

	ProductSourceTypeCustom = "custom"

	RouteSameProvince ShippingRouteType = "noi_tinh"
	RouteSameRegion   ShippingRouteType = "noi_mien"
	RouteNationWide   ShippingRouteType = "toan_quoc"

	ShippingDistrictTypeUrban     ShippingDistrictType = "noi_thanh"
	ShippingDistrictTypeSubUrban1 ShippingDistrictType = "ngoai_thanh_1"
	ShippingDistrictTypeSubUrban2 ShippingDistrictType = "ngoai_thanh_2"

	UserSourcePSX          UserSource = "psx"
	UserSourceEtop         UserSource = "etop"
	UserSourceTopship      UserSource = "topship"
	UserSourceTsAppAndroid UserSource = "ts_app_android"
	UserSourceTsAppIOS     UserSource = "ts_app_ios"
	UserSourceTsAppWeb     UserSource = "ts_app_web"
	UserSourcePartner      UserSource = "partner"
)

type NormalizedPhone = validate.NormalizedPhone
type NormalizedEmail = validate.NormalizedEmail

var ShippingFeeShopTypes = []ShippingFeeLineType{
	ShippingFeeTypeMain, ShippingFeeTypeReturn, ShippingFeeTypeAdjustment,
	ShippingFeeTypeAddessChange, ShippingFeeTypeCODS, ShippingFeeTypeInsurance, ShippingFeeTypeOther, ShippingFeeTypeDiscount,
}

//        SyncAt: updated when sending data to external service successfully
//     TrySyncAt: updated when sending data to external service (may unsuccessfully)
//    LastSyncAt: updated when querying external data source successfully
// LastTrySyncAt: updated when querying external data source (may unsuccessfully)

var ShippingStateMap = map[ShippingState]string{
	StateDefault:       "Mặc định",
	StateCreated:       "Mới",
	StatePicking:       "Đang lấy hàng",
	StateHolding:       "Đã lấy hàng",
	StateDelivering:    "Đang giao hàng",
	StateReturning:     "Đang trả hàng",
	StateDelivered:     "Đã giao hàng",
	StateReturned:      "Đã trả hàng",
	StateCancelled:     "Hủy",
	StateClosed:        "Hoàn tất",
	StateUndeliverable: "Bồi hoàn",
	StateUnknown:       "Không xác định",
}

func (s ShippingState) Text() string {
	res := ShippingStateMap[s]
	if res == "" {
		return "Không xác định"
	}
	return res
}

func EtopPaymentStatusLabel(s Status4) string {
	switch s {
	case S4Negative:
		return "Không cần thanh toán"
	case S4Zero:
		return "Chưa đối soát"
	case S4Positive:
		return "Đã thanh toán"
	case S4SuperPos:
		return "Cần thanh toán"
	default:
		return "Không xác định"
	}
}

func OrderStatusLabel(s Status5) string {
	switch s {
	case S5NegSuper:
		return "Trả hàng"
	case S5Negative:
		return "Huỷ"
	case S5Zero:
		return "Mới"
	case S5Positive:
		return "Thành công"
	case S5SuperPos:
		return "Đang xử lý"
	default:
		return "Không xác định"
	}
}

func PaymentStatusLabel(s Status4) string {
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

func (t AccountType) Label() string {
	switch t {
	case TypeShop:
		return "cửa hàng"
	case TypePartner:
		return "tài khoản"
	case TypeEtop:
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
		CodeTypeMoneyTransactionEtop:
		// no-op
	default:
		return cm.Error(cm.InvalidArgument, "Type is not valid", nil)
	}
	return nil
}

func VerifyPaymentMethod(s string) bool {
	switch s {
	case PaymentMethodCOD, PaymentMethodBank,
		PaymentMethodOther, PaymentMethodVtpay:
		return true
	}
	return false
}

func VerifyOrderSource(s OrderSourceType) bool {
	switch s {
	case OrderSourcePOS, OrderSourcePXS, OrderSourceCMX, OrderSourceTSApp, OrderSourceAPI, OrderSourceImport:
		return true
	}
	return false
}

func (p ShippingProvider) Label() string {
	switch p {
	case TypeGHN:
		return "Giao Hàng Nhanh"
	case TypeGHTK:
		return "Giao Hàng Tiết Kiệm"
	case TypeVTPost:
		return "Viettel Post"
	case TypeShippingProviderManual:
		return "Tự giao"
	default:
		return ""
	}
}

func VerifyShippingProvider(s ShippingProvider) bool {
	switch s {
	case TypeGHN, TypeGHTK, TypeVTPost, TypeShippingProviderManual:
		return true
	}
	return false
}

func TryOnFromGHNNoteCode(c string) TryOn {
	switch c {
	case GHNNoteCodeKhongXemHang:
		return TryOnNone
	case GHNNoteCodeChoXemHang:
		return TryOnOpen
	case GHNNoteCodeChoThuHang:
		return TryOnTry
	default:
		return ""
	}
}

func GHNNoteCodeFromTryOn(to TryOn) string {
	switch to {
	case TryOnNone:
		return GHNNoteCodeKhongXemHang
	case TryOnOpen:
		return GHNNoteCodeChoXemHang
	case TryOnTry:
		return GHNNoteCodeChoThuHang
	default:
		return ""
	}
}

var _ = sqlgenAccount(&Account{})

type Account struct {
	ID       int64
	OwnerID  int64
	Name     string
	Type     AccountType
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
		Type:     TypeShop,
		ImageURL: s.ImageURL,
		URLSlug:  "",
	}
}

func (s *Partner) GetAccount() *Account {
	return &Account{
		ID:       s.ID,
		OwnerID:  s.OwnerID,
		Name:     s.Name,
		Type:     TypePartner,
		ImageURL: s.ImageURL,
		URLSlug:  "",
	}
}

var _ = sqlgenShop(&Shop{})

type Shop struct {
	ID      int64
	Name    string
	OwnerID int64
	IsTest  int

	AddressID         int64
	ShipToAddressID   int64
	ShipFromAddressID int64
	Phone             string
	BankAccount       *BankAccount
	WebsiteURL        string
	ImageURL          string
	Email             string
	Code              string
	AutoCreateFFM     bool

	OrderSourceID int64

	Status    Status3
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time

	Address *Address `sq:"-"`

	RecognizedHosts []string

	// @deprecated use try_on instead
	GhnNoteCode string
	TryOn       TryOn
	CompanyInfo *CompanyInfo
	// MoneyTransactionRRule format:
	// FREQ=DAILY
	// FREQ=WEEKLY;BYDAY=TU,TH,SA
	// FREQ=MONTHLY;BYDAY=-1FR
	MoneyTransactionRRule         string `sq:"'money_transaction_rrule'"`
	SurveyInfo                    []*SurveyInfo
	ShippingServiceSelectStrategy []*ShippingServiceSelectStrategyItem
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

func (s *Shop) GetTryOn() TryOn {
	if s.TryOn != "" {
		return s.TryOn
	}
	return TryOnFromGHNNoteCode(s.GhnNoteCode)
}

var _ = sqlgenShopExtended(
	&ShopExtended{}, &Shop{}, sq.AS("s"),
	sq.LEFT_JOIN, &Address{}, sq.AS("a"), "s.address_id = a.id",
	sq.LEFT_JOIN, &User{}, sq.AS("u"), "s.owner_id = u.id",
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
	ID      int64
	OwnerID int64
	Status  Status3
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
	// AvailableFromEtop: dùng để xác định `partner` này có thể xác thực shop trực tiếp từ client của ETOP hay không
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

func (p *Partner) GetID() int64 {
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
	AccountID   int64
	Status      Status3
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
	&AccountAuthFtPartner{}, &AccountAuth{}, sq.AS("aa"),
	sq.JOIN, &Partner{}, sq.AS("p"), `aa.account_id = p.id`,
)

type AccountAuthFtPartner struct {
	*AccountAuth
	*Partner
}

var _ = sqlgenAccountAuthFtShop(
	&AccountAuthFtShop{}, &AccountAuth{}, sq.AS("aa"),
	sq.JOIN, &Shop{}, sq.AS("s"), `aa.account_id = s.id`,
)

type AccountAuthFtShop struct {
	*AccountAuth
	*Shop
}

var _ = sqlgenPartnerRelation(&PartnerRelation{})

type PartnerRelation struct {
	AuthKey           string
	PartnerID         int64
	SubjectID         int64
	SubjectType       SubjectType
	ExternalSubjectID string

	Nonce     int64
	Status    Status3
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time

	Permission `sq:"inline"`
}

var _ = sqlgenPartnerRelationFtShop(
	&PartnerRelationFtShop{}, &PartnerRelation{}, sq.AS("pr"),
	sq.JOIN, &Shop{}, sq.AS("s"), "pr.subject_id = s.id",
	sq.JOIN, &User{}, sq.AS("u"), "s.owner_id = u.id",
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
	ID int64

	UserInner `sq:"inline"`

	Status      Status3 // 1: actual user, 0: stub, -1: disabled
	Identifying UserIdentifying

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`

	AgreedTOSAt       time.Time
	AgreedEmailInfoAt time.Time
	EmailVerifiedAt   time.Time
	PhoneVerifiedAt   time.Time

	EmailVerificationSentAt time.Time
	PhoneVerificationSentAt time.Time

	IsTest    int
	Source    UserSource
	RefUserID int64
	RefSaleID int64
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

func (p *Permission) Validate() error {
	if len(p.Permissions) != 0 {
		return cm.Error(cm.InvalidArgument, "list of permissions is not supported yet", nil)
	}
	if len(p.Roles) != 1 {
		return cm.Error(cm.InvalidArgument, "only single role is supported", nil)
	}
	role := RoleType(p.Roles[0])
	switch role {
	case RoleOwner, RoleAdmin, RoleStaff, RoleViewer:
		return nil
	default:
		return cm.Error(cm.InvalidArgument, "invalid role", nil)
	}
}

var _ = sqlgenAccountUser(&AccountUser{})

type AccountUser struct {
	AccountID int64
	UserID    int64

	Status         Status3 // 1: activated, -1: rejected/disabled, 0: pending
	ResponseStatus Status3 // 1: accepted,  -1: rejected, 0: pending

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time

	Permission `sq:"inline"`

	FullName  string
	ShortName string
	Position  string

	InvitationSentAt     time.Time
	InvitationSentBy     int64
	InvitationAcceptedAt time.Time
	InvitationRejectedAt time.Time

	DisabledAt    time.Time
	DisabledBy    time.Time
	DisableReason string
}

var _ = sqlgenAccountUserExtended(
	&AccountUserExtended{}, &AccountUser{}, sq.AS("au"),
	sq.JOIN, &Account{}, sq.AS("a"), "au.account_id = a.id",
	sq.JOIN, &User{}, sq.AS("u"), "au.user_id = u.id",
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
	UserID   int64
	AuthType string
	AuthKey  string

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

var _ = sqlgenUserInternal(&UserInternal{})

type UserInternal struct {
	ID      int64
	Hashpwd string

	UpdatedAt time.Time `sq:"update"`
}

type StatusQuery struct {
	Status *Status3
}

type ExternalShippingLog struct {
	StateText string
	Time      string
	Message   string
}

type FulfillmentSyncStates struct {
	SyncAt    time.Time `json:"sync_at"`
	TrySyncAt time.Time `json:"try_sync_at"`
	Error     *Error    `json:"error"`

	NextShippingState ShippingState `json:"next_shipping_state"`
}

func (f FfmAction) ToShippingState() ShippingState {
	switch f {
	case FfmActionCancel:
		return StateCancelled
	case FfmActionReturn:
		return StateReturning
	default:
		return ""
	}
}

var _ = sqlgenAddress(&Address{})

type Address struct {
	ID        int64  `json:"id"`
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
	AccountID   int64        `json:"account_id"`
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
	ID        int64
	Amount    int
	ShopID    int64
	Type      string
	Status    Status3
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	PaidAt    time.Time
}

var _ = sqlgenCreditExtended(
	&CreditExtended{}, &Credit{}, sq.AS("c"),
	sq.LEFT_JOIN, &Shop{}, sq.AS("s"), "s.id = c.shop_id",
)

type CreditExtended struct {
	*Credit
	Shop *Shop
}

type ShippingFeeLine struct {
	ShippingFeeType          ShippingFeeLineType `json:"shipping_fee_type"`
	Cost                     int                 `json:"cost"`
	ExternalServiceID        string              `json:"external_service_id"`
	ExternalServiceName      string              `json:"external_service_name"`
	ExternalServiceType      string              `json:"external_service_type"`
	ExternalShippingOrderID  string              `json:"external_order_id"`
	ExternalPaymentChannelID string              `json:"external_payment_channel_id"`
	ExternalShippingCode     string              `json:"external_shipping_code"`
}

func GetShippingFeeShopLines(items []*ShippingFeeLine, etopPriceRule bool, mainFee *int) []*ShippingFeeLine {
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

func GetShippingFeeShopLine(item ShippingFeeLine, etopPriceRule bool, mainFee *int) *ShippingFeeLine {
	if item.ShippingFeeType == ShippingFeeTypeMain && etopPriceRule && mainFee != nil {
		item.Cost = *mainFee
	}
	if contains(ShippingFeeShopTypes, item.ShippingFeeType) {
		return &item
	}
	return nil
}

func GetReturnedFee(items []*ShippingFeeLine) int {
	result := 0
	for _, item := range items {
		if item.ShippingFeeType == ShippingFeeTypeReturn {
			result = item.Cost
			break
		}
	}
	return result
}

func GetMainFee(items []*ShippingFeeLine) int {
	result := 0
	for _, item := range items {
		if item.ShippingFeeType == ShippingFeeTypeMain {
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

func UpdateShippingFees(items []*ShippingFeeLine, fee int, shippingFeeType ShippingFeeLineType) []*ShippingFeeLine {
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

func contains(lines []ShippingFeeLineType, feeType ShippingFeeLineType) bool {
	for _, line := range lines {
		if feeType == line {
			return true
		}
	}
	return false
}

var _ = sqlgenShippingSource(&ShippingSource{})

type ShippingSource struct {
	ID        int64
	Name      string
	Username  string
	Type      string
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

var _ = sqlgenShippingSourceInternal(&ShippingSourceInternal{})

type ShippingSourceInternal struct {
	ID          int64
	CreatedAt   time.Time `sq:"create"`
	UpdatedAt   time.Time `sq:"update"`
	LastSyncAt  time.Time
	AccessToken string
	ExpiresAt   time.Time `json:"expires_at"`
	Secret      *ShippingSourceSecret
}

type ShippingSourceSecret struct {
	CustomerID int
	Username   string
	Password   string
	// Use for VTPost
	GroupAddressID int
}

func (s ShippingState) ToStatus4() Status4 {
	switch s {
	case StateDefault:
		return S4Zero
	case StateCancelled,
		StateReturned:
		return S4Negative
	case StateClosed:
		return S4Positive
	}
	return S4SuperPos
}

func (s ShippingState) ToShippingStatus5() Status5 {
	switch s {
	case StateDefault:
		return S5Zero
	case StateCancelled:
		return S5Negative
	case StateReturning, StateReturned:
		return S5NegSuper
	case StateDelivered:
		return S5Positive
	default:
		return S5SuperPos
	}
}

type AvailableShippingService struct {
	Name string
	// ServiceFee: Tổng phí giao hàng (đã bao gồm phí chính + các phụ phí khác)
	ServiceFee int
	// ShippingFeeMain: Phí chính giao hàng
	ShippingFeeMain   int
	Provider          ShippingProvider
	ProviderServiceID string

	ExpectedPickAt     time.Time
	ExpectedDeliveryAt time.Time
	Source             ShippingPriceSource
}

func (service *AvailableShippingService) ApplyFeeMain(feeMain int) {
	service.ServiceFee = service.ServiceFee - service.ShippingFeeMain + feeMain
	service.ShippingFeeMain = feeMain
}

var _ = sqlgenWebhook(&Webhook{})

type Webhook struct {
	ID        int64
	AccountID int64
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
	ID        int64
	WebhookID int64
	AccountID int64
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
