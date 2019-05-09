package model

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/gencode"
	sq "etop.vn/backend/pkg/common/sql"
	"etop.vn/backend/pkg/common/validate"
)

//go:generate ../../../scripts/derive.sh

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
	AccountType          string
	SubjectType          string
	ProcessingStatus     string
	ShippingState        string
	FulfillmentEndpoint  string
	ShippingProvider     string
	ShippingPriceSource  string
	OrderSourceType      string
	ShippingFeeLineType  string
	UserIdentifying      string
	RoleType             string
	FfmAction            int
	CodeType             string
	TryOn                string
	OrderFeeType         string
	ShippingRouteType    string
	ShippingDistrictType string
	DBName               string
)

// Type constants
const (
	DBMain     DBName = "main"
	DBNotifier DBName = "notifier"

	TypePartner  AccountType = "partner"
	TypeShop     AccountType = "shop"
	TypeSupplier AccountType = "supplier"
	TypeEtop     AccountType = "etop"

	SubjectTypeAccount SubjectType = "account"
	SubjectTypeUser    SubjectType = "user"

	TypeKiotviet = "kiotviet"

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
	FFSupplier FulfillmentEndpoint = "supplier"
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

	ProductSourceCustom   = "custom"
	ProductSourceKiotViet = "kiotviet"

	// Don't change these values
	TagUser     = 17
	TagPartner  = 21
	TagShop     = 33
	TagSupplier = 65
	TagEtop     = 101
	TagImport   = 111

	EtopAccountID = TagEtop

	TypeEmail  = "EM"
	TypePhone  = "PH"
	TypeAPIKey = "AK"

	StatusActive   = 1
	StatusCreated  = 0
	StatusDisabled = -1

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

	OrderFeeOther    OrderFeeType = "other"
	OrderFeeShipping OrderFeeType = "shipping"
	OrderFeeTax      OrderFeeType = "tax"

	RoleOwner  RoleType = "owner"
	RoleAdmin  RoleType = "admin"
	RoleStaff  RoleType = "staff"
	Role3rd    RoleType = "3rd" // Delegate to third party API, which can only update/cancel itself orders
	RoleViewer RoleType = "viewer"

	PaymentMethodCOD   = "cod"
	PaymentMethodBank  = "bank"
	PaymentMethodOther = "other"

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
		return "tài khoản nhà cung cấp"
	case TypeSupplier:
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
	case PaymentMethodCOD, PaymentMethodBank, PaymentMethodOther:
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
	_account()
}

func (s *Shop) _account()     {}
func (s *Supplier) _account() {}
func (s *Partner) _account()  {}

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

func (s *Supplier) GetAccount() *Account {
	return &Account{
		ID:       s.ID,
		OwnerID:  s.OwnerID,
		Name:     s.Name,
		Type:     TypeSupplier,
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

	ProductSourceID int64
	OrderSourceID   int64

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

var _ = sqlgenEtopCatgory(&EtopCategory{})

type EtopCategory struct {
	ID       int64
	Name     string
	ParentID int64

	Status    Status3
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

var _ = sqlgenSupplier(&Supplier{})

type Supplier struct {
	ID        int64
	Status    Status3
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	IsTest    int

	Name     string
	OwnerID  int64
	ImageURL string

	Rules              *SupplierPriceRules
	CompanyInfo        *CompanyInfo
	WarehouseAddressID int64
	BankAccount        *BankAccount
	ContactPersons     []*ContactPerson
	ShipFromAddressID  int64

	ProductSourceID int64
}

var _ = sqlgenSupplierInfo(&SupplierInfo{}, &Supplier{})

type SupplierInfo struct {
	Name     string
	ImageURL string

	Rules              *SupplierPriceRules
	CompanyInfo        *CompanyInfo
	WarehouseAddressID int64
	BankAccount        *BankAccount
	ContactPersons     []*ContactPerson
}

var _ = sqlgenSupplierExtended(
	&SupplierExtended{}, &Supplier{}, sq.AS("s"),
	sq.LEFT_JOIN, &Address{}, sq.AS("a"), "s.warehouse_address_id = a.id",
)

type SupplierExtended struct {
	*Supplier
	Address *Address
}

var _ = sqlgenSupplierShipFromAddress(
	&SupplierShipFromAddress{}, &Supplier{}, sq.AS("s"),
	sq.LEFT_JOIN, &Address{}, sq.AS("a"), "s.ship_from_address_id = a.id",
)

type SupplierShipFromAddress struct {
	*Supplier
	Address *Address
}

func (s *Supplier) GetPriceRules() (*SupplierPriceRules, bool) {
	priceRules := s.Rules
	if priceRules == nil {
		priceRules = &SupplierPriceRules{
			General: DefaultSupplierPriceRule(),
		}
	} else if priceRules.General == nil {
		priceRules.General = DefaultSupplierPriceRule()
	}
	return priceRules, s.Rules != nil
}

type SupplierPriceRules struct {
	General *SupplierPriceRule   `json:"general"`
	Rules   []*SupplierPriceRule `json:"rules"`
}

func (rules *SupplierPriceRules) Validate() error {
	if rules == nil {
		return cm.Error(cm.FailedPrecondition, "Empty rules", nil)
	}
	if !rules.General.IsZeroIdentifier() {
		return cm.Error(cm.FailedPrecondition, "General rule must have no identifier", nil)
	}
	if err := rules.General.Validate(); err != nil {
		return cm.Error(cm.FailedPrecondition, "General rule is not valid", err)
	}
	for i, r := range rules.Rules {
		if r.IsZeroIdentifier() {
			return cm.Error(cm.FailedPrecondition, fmt.Sprintf("Rule #%v must have identifier", i+1), nil)
		}
		if err := r.Validate(); err != nil {
			return cm.Error(cm.FailedPrecondition, fmt.Sprintf("Rule #%v is not valid", i+1), err)
		}
	}
	return nil
}

type SupplierPriceRule struct {
	SupplierCategoryID int64  `json:"supplier_category_id"`
	ExternalCategoryID string `json:"external_category_id"`
	Tag                string `json:"tag"`

	ListPriceA      float64 `json:"list_price_a"`
	ListPriceB      float64 `json:"list_price_b"`
	WholesalePriceA float64 `json:"wholesale_price_a"`
	WholesalePriceB float64 `json:"wholesale_price_b"`
	RetailPriceMinA float64 `json:"retail_price_min_a"`
	RetailPriceMinB float64 `json:"retail_price_min_b"`
	RetailPriceMaxA float64 `json:"retail_price_max_a"`
	RetailPriceMaxB float64 `json:"retail_price_max_b"`
}

func (r *SupplierPriceRule) GetListPriceA() float64 {
	if r == nil {
		return 0
	}
	return r.ListPriceA
}

func (r *SupplierPriceRule) GetListPriceB() float64 {
	if r == nil {
		return 0
	}
	return r.ListPriceB
}

func (r *SupplierPriceRule) GetWholesalePriceA() float64 {
	if r == nil {
		return 0
	}
	return r.WholesalePriceA
}

func (r *SupplierPriceRule) GetWholesalePriceB() float64 {
	if r == nil {
		return 0
	}
	return r.WholesalePriceB
}

func (r *SupplierPriceRule) GetRetailPriceMinA() float64 {
	if r == nil {
		return 0
	}
	return r.RetailPriceMinA
}

func (r *SupplierPriceRule) GetRetailPriceMinB() float64 {
	if r == nil {
		return 0
	}
	return r.RetailPriceMinB
}

func (r *SupplierPriceRule) GetRetailPriceMaxA() float64 {
	if r == nil {
		return 0
	}
	return r.RetailPriceMaxA
}

func (r *SupplierPriceRule) GetRetailPriceMaxB() float64 {
	if r == nil {
		return 0
	}
	return r.RetailPriceMaxB
}

func DefaultSupplierPriceRule() *SupplierPriceRule {
	return &SupplierPriceRule{
		ListPriceA:      1,
		WholesalePriceA: 0.5,
		RetailPriceMinA: 0.8,
		RetailPriceMaxA: 1.5,
	}
}

func (r *SupplierPriceRule) IsZeroIdentifier() bool {
	return r != nil && r.ExternalCategoryID == "" && r.SupplierCategoryID == 0 && r.Tag == ""
}

func (r *SupplierPriceRule) Validate() error {
	if r == nil {
		return cm.Error(cm.InvalidArgument, "Empty rule", nil)
	}
	if err := validatePrice("WholesalePrice", r.WholesalePriceA, r.WholesalePriceB); err != nil {
		return err
	}
	if err := validatePrice("ListPrice", r.ListPriceA, r.ListPriceB); err != nil {
		return err
	}
	if err := validatePrice("RetailPriceMin", r.RetailPriceMinA, r.RetailPriceMinB); err != nil {
		return err
	}
	if err := validatePrice("RetailPriceMax", r.RetailPriceMaxA, r.RetailPriceMaxB); err != nil {
		return err
	}

	if r.WholesalePriceA == 0 && r.ListPriceA == 0 &&
		r.RetailPriceMinA == 0 && r.RetailPriceMaxA == 0 {
		return cm.Error(cm.InvalidArgument, "All primary variables are zero", nil)
	}

	r.WholesalePriceB = math.Floor(r.WholesalePriceB)
	r.ListPriceB = math.Floor(r.ListPriceB)
	r.RetailPriceMinB = math.Floor(r.RetailPriceMinB)
	r.RetailPriceMaxB = math.Floor(r.RetailPriceMaxB)
	return nil
}

func validatePrice(name string, a, b float64) error {
	if a < 0 {
		return cm.Error(cm.InvalidArgument, "Invalid "+name, nil)
	}
	if a == 0 && b == 0 {
		return cm.Error(cm.InvalidArgument, "Invalid "+name+" (both a and b must not be zero)", nil)
	}
	if a >= 0 && a <= 5 && b >= -5e7 && b <= 5e7 {
		return nil
	}
	return cm.Error(cm.InvalidArgument, "Invalid "+name+" (a must be between 0 and 5, b must be between -5e7 and 5e7)", nil)
}

type KiotvietBranch struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Code          string    `json:"code"`
	ContactNumber string    `json:"contact_number"`
	RetailerID    string    `json:"retailer_id"`
	Address       string    `json:"address"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedAt     time.Time `json:"created_at"`
}

type SupplierKiotviet struct {
	RetailerID     string
	ClientID       string
	ClientSecret   string
	ClientToken    string
	ExpiresAt      time.Time
	KiotvietStatus int

	Branches []*KiotvietBranch
}

// var _ = sqlgenSupplierKiotviet(&SupplierKiotviet{})

// type SupplierKiotviet struct {
// 	SupplierID int64

// 	Status    int
// 	CreatedAt time.Time `sq:"create"`
// 	UpdatedAt time.Time `sq:"update"`

// 	SupplierKiotvietInner `sq:"inline"`

// 	DefaultBranchID string

// 	// We have to define  here, because the builder will use it to
// 	// generate correct query.
// 	SyncStateProducts   json.RawMessage
// 	SyncStateCategories json.RawMessage
// }

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
		return cm.Errorf(cm.InvalidArgument, nil, "missing AccountID")
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

var _ = sqlgenProductSource(&ProductSource{})

type ProductSource struct {
	ID         int64
	SupplierID int64
	Type       string
	Name       string
	Status     Status3

	ExternalStatus Status3
	ExternalKey    string
	ExternalInfo   *KiotvietExternalInfo
	ExtraInfo      *KiotvietExtraInfo

	CreatedAt  time.Time `sq:"create"`
	UpdatedAt  time.Time `sq:"update"`
	LastSyncAt time.Time

	SyncStateProducts   json.RawMessage
	SyncStateCategories json.RawMessage
}

var _ = sqlgenProductSourceSyncStates(&ProductSourceSyncStates{}, &ProductSource{})

type ProductSourceSyncStates struct {
	LastSyncAt          time.Time
	SyncStateProducts   json.RawMessage
	SyncStateCategories json.RawMessage
}

type KiotvietExternalInfo struct {
	Branches []*KiotvietBranch `json:"branches"`
}

type KiotvietExtraInfo struct {
	DefaultBranchID string `json:"default_branch_id"`
}

var _ = sqlgenProductSourceInternal(&ProductSourceInternal{})

type ProductSourceInternal struct {
	ID int64

	Secret      *KiotvietSecret
	AccessToken string
	ExpiresAt   time.Time

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

type KiotvietSecret struct {
	RetailerID   string `json:"retailer_id"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

var _ = sqlgenProductSourceInternalAccessToken(&ProductSourceInternalAccessToken{}, &ProductSourceInternal{})

type ProductSourceInternalAccessToken struct {
	AccessToken string
	ExpiresAt   time.Time
	UpdatedAt   time.Time `sq:"update"`
}

var _ = sqlgenProductSourceExtended(
	&ProductSourceExtended{}, &ProductSource{}, sq.AS("ps"),
	sq.LEFT_JOIN, &ProductSourceInternal{}, sq.AS("psi"), "ps.id = psi.id",
)

type ProductSourceExtended struct {
	*ProductSource
	*ProductSourceInternal
}

var _ = sqlgenSupplierFtProductSource(
	&SupplierFtProductSource{}, &Supplier{}, sq.AS("s"),
	sq.LEFT_JOIN, &ProductSource{}, sq.AS("ps"), "s.product_source_id = ps.id",
)

type SupplierFtProductSource struct {
	*Supplier
	*ProductSource
}

var _ = sqlgenProductSourceFtSupplier(
	&ProductSourceFtSupplier{}, &ProductSource{}, sq.AS("ps"),
	sq.LEFT_JOIN, &Supplier{}, sq.AS("s"), "s.product_source_id = ps.id",
)

type ProductSourceFtSupplier struct {
	*ProductSource
	*Supplier
}

var _ = sqlgenShopFtProductSource(
	&ShopFtProductSource{}, &Shop{}, sq.AS("s"),
	sq.LEFT_JOIN, &ProductSource{}, sq.AS("ps"), "s.product_source_id = ps.id",
	sq.LEFT_JOIN, &ProductSourceInternal{}, sq.AS("psi"), "s.product_source_id = psi.id",
)

type ShopFtProductSource struct {
	*Shop
	*ProductSource
	*ProductSourceInternal
}

var _ = sqlgenOrderSource(&OrderSource{})

type OrderSource struct {
	ID     int64
	ShopID int64
	Type   string
	Name   string
	Status int

	ExternalStatus int
	ExternalKey    string
	ExternalInfo   json.RawMessage
	ExtraInfo      json.RawMessage

	CreatedAt  time.Time `sq:"create"`
	UpdatedAt  time.Time `sq:"update"`
	LastSyncAt time.Time

	SyncStateOrders json.RawMessage
}

var _ = sqlgenOrderSourceSyncStates(&OrderSourceSyncStates{}, &OrderSource{})

type OrderSourceSyncStates struct {
	LastSyncAt      time.Time `sq:"update"`
	SyncStateOrders json.RawMessage
}

var _ = sqlgenOrderSourceInternal(&OrderSourceInternal{})

type OrderSourceInternal struct {
	ID          int64
	AccessToken string
	ExpiresAt   time.Time
	CreatedAt   time.Time `sq:"create"`
	UpdatedAt   time.Time `sq:"update"`
}

var _ = sqlgenOrderSourceExtended(
	&OrderSourceExtended{}, &OrderSource{}, sq.AS("os"),
	sq.LEFT_JOIN, &OrderSourceInternal{}, sq.AS("osi"), "os.id = osi.id",
)

type OrderSourceExtended struct {
	OrderSource         *OrderSource
	OrderSourceInternal *OrderSourceInternal
}

var _ = sqlgenShopFtOrderSource(
	&ShopFtOrderSource{}, &Shop{}, sq.AS("s"),
	sq.LEFT_JOIN, &OrderSource{}, sq.AS("os"), "s.order_source_id = os.id",
	sq.LEFT_JOIN, &OrderSourceInternal{}, sq.AS("osi"), "s.order_source_id = osi.id",
)

type ShopFtOrderSource struct {
	*Shop
	*OrderSource
	*OrderSourceInternal
}

var _ = sqlgenProductSourceCategoryExternal(&ProductSourceCategoryExternal{})

type ProductSourceCategoryExternal struct {
	ID int64

	ProductSourceID   int64
	ProductSourceType string

	ExternalID        string
	ExternalParentID  string
	ExternalCode      string
	ExternalName      string
	ExternalStatus    int
	ExternalUpdatedAt time.Time
	ExternalCreatedAt time.Time
	ExternalDeletedAt time.Time
	LastSyncAt        time.Time
}

var _ = sqlgenProductSourceCategory(&ProductSourceCategory{})

type ProductSourceCategory struct {
	ID int64

	ProductSourceID   int64
	ProductSourceType string
	SupplierID        int64
	ParentID          int64
	ShopID            int64

	Name string

	Status    int
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time
}

var _ = sqlgenProductSourceCategoryExtended(&ProductSourceCategoryExtended{},
	&ProductSourceCategory{}, sq.AS("psc"),
	sq.LEFT_JOIN, &ProductSourceCategoryExternal{}, sq.AS("psce"), "psce.id = psc.id",
)

type ProductSourceCategoryExtended struct {
	*ProductSourceCategory
	*ProductSourceCategoryExternal
}

type ProductAttributes []ProductAttribute

func (attrs ProductAttributes) Name() string {
	if len(attrs) == 0 {
		return ""
	}
	return attrs.ShortLabel()
}

func (attrs ProductAttributes) Label() string {
	if len(attrs) == 0 {
		return "Mặc định"
	}
	b := make([]byte, 0, 64)
	for _, attr := range attrs {
		if attr.Name == "" || attr.Value == "" {
			continue
		}
		if len(b) > 0 {
			b = append(b, ", "...)
		}
		b = append(b, attr.Name...)
		b = append(b, ": "...)
		b = append(b, attr.Value...)
	}
	return string(b)
}

func (attrs ProductAttributes) ShortLabel() string {
	if len(attrs) == 0 {
		return "Mặc định"
	}
	b := make([]byte, 0, 64)
	for _, attr := range attrs {
		if attr.Name == "" || attr.Value == "" {
			continue
		}
		if len(b) > 0 {
			b = append(b, ' ')
		}
		b = append(b, attr.Value...)
	}
	return string(b)
}

type ProductAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
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

var _ = sqlgenProduct(&Product{})

type Product struct {
	ID                      int64
	ProductSourceID         int64
	SupplierID              int64
	ProductSourceCategoryID string
	EtopCategoryID          int64
	ProductBrandID          int64

	Name          string
	ShortDesc     string
	Description   string
	DescHTML      string `sq:"'desc_html'"`
	EdName        string
	EdShortDesc   string
	EdDescription string
	EdDescHTML    string `sq:"'ed_desc_html'"`
	EdTags        []string
	Unit          string

	Status Status3
	Code   string
	EdCode string

	QuantityAvailable int // Available = OnHand - Reserved - 3
	QuantityOnHand    int
	QuantityReserved  int

	ImageURLs []string `sq:"'image_urls'"`

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`

	NameNorm   string // search normalization
	NameNormUa string // unaccent normalization
}

func (p *Product) BeforeInsert() error {
	p.NameNorm = validate.NormalizeSearch(p.Name)
	p.NameNormUa = validate.NormalizeUnaccent(p.Name)
	return nil
}

func (p *Product) BeforeUpdate() error {
	p.NameNorm = validate.NormalizeSearch(p.Name)
	p.NameNormUa = validate.NormalizeUnaccent(p.Name)
	return nil
}

func (p *Product) IsAvailable() bool {
	return p.QuantityAvailable > 0
}

func (p *Product) GetFullName() string {
	return coalesce(p.Name, p.EdName)
}

var _ = sqlgenProductExternal(&ProductExternal{})

type ProductExternal struct {
	ID int64

	ProductExternalCommon `sq:"inline"`

	ExternalUnits []*Unit
}

type ProductExternalCommon struct {
	ProductSourceID   int64
	ProductSourceType string

	ExternalID          string
	ExternalName        string
	ExternalCode        string
	ExternalCategoryID  string
	ExternalDescription string
	ExternalImageURLs   []string
	ExternalUnit        string

	ExternalData      json.RawMessage
	ExternalStatus    Status3
	ExternalCreatedAt time.Time
	ExternalUpdatedAt time.Time
	ExternalDeletedAt time.Time
	LastSyncAt        time.Time
}

var _ = sqlgenProductExtended(
	&ProductExtended{}, &Product{}, sq.AS("p"),
	sq.LEFT_JOIN, &ProductExternal{}, sq.AS("px"), "p.id = px.id",
	sq.LEFT_JOIN, &ProductBrand{}, sq.AS("pb"), "p.product_brand_id = pb.id",
	sq.LEFT_JOIN, &ProductSource{}, sq.AS("ps"), "p.product_source_id = ps.id",
)

type ProductExtended struct {
	*Product
	*ProductExternal
	*ProductBrand
	*ProductSource
}

type ProductFtVariant struct {
	ProductExtended
	Variants []*VariantExternalExtended
}

// func (p *Product) BeforeUpdate() {
// 	p.ExternalNameNorm = validate.NormalizeSearch(p.ExternalFullName)
// 	p.SupplierNameNorm = validate.NormalizeSearch(p.SupplierFullName)
// 	if p.SupplierNameNorm == "" {
// 		p.SupplierNameNorm = p.ExternalNameNorm
// 	}
// }

var _ = sqlgenVariant(&Variant{})

type Variant struct {
	ID              int64
	ProductID       int64
	ProductSourceID int64
	// ProductSourceType string
	SupplierID int64

	ProductSourceCategoryID int64
	EtopCategoryID          int64
	ProductBrandID          int64

	// Name          string
	ShortDesc     string
	Description   string
	DescHTML      string `sq:"'desc_html'"`
	EdName        string
	EdShortDesc   string
	EdDescription string
	EdDescHTML    string `sq:"'ed_desc_html'"`

	DescNorm string

	// key-value normalization, must be non-null. Empty attributes is '_'.
	AttrNormKv string

	Status     Status3
	EtopStatus Status3
	EdStatus   Status3
	Code       string
	EdCode     string

	WholesalePrice0 int `sq:"'wholesale_price_0'"`
	WholesalePrice  int
	ListPrice       int
	RetailPriceMin  int
	RetailPriceMax  int

	EdWholesalePrice0 int `sq:"'ed_wholesale_price_0'"`
	EdWholesalePrice  int
	EdListPrice       int
	EdRetailPriceMin  int
	EdRetailPriceMax  int

	QuantityAvailable int // Available = OnHand - Reserved - 3
	QuantityOnHand    int
	QuantityReserved  int
	ImageURLs         []string `sq:"'image_urls'"`
	SupplierMeta      json.RawMessage

	CostPrice  int
	Attributes ProductAttributes

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

func (v *Variant) GetName() string {
	if len(v.Attributes) == 0 {
		return ""
	}
	return ProductAttributes(v.Attributes).ShortLabel()
}

func (v *Variant) IsAvailable() bool {
	return v.QuantityAvailable > 0
}

func (v *Variant) BeforeInsert() error {
	v.Attributes, v.AttrNormKv = NormalizeAttributes(v.Attributes)
	return nil
}

func (v *Variant) BeforeUpdate() error {
	v.Attributes, v.AttrNormKv = NormalizeAttributes(v.Attributes)
	return nil
}

// Normalize attributes, do not sort them. Empty attributes is '_'.
func NormalizeAttributes(attrs []ProductAttribute) ([]ProductAttribute, string) {
	if len(attrs) == 0 {
		return nil, "_"
	}
	const maxAttrs = 8
	if len(attrs) > maxAttrs {
		attrs = attrs[:maxAttrs]
	}

	normAttrs := make([]ProductAttribute, 0, len(attrs))
	b := make([]byte, 0, 256)
	for _, attr := range attrs {
		attr.Name, _ = validate.NormalizeName(attr.Name)
		attr.Value, _ = validate.NormalizeName(attr.Value)
		if attr.Name == "" || attr.Value == "" {
			fmt.Println("1", attr.Name, attr.Value)
			continue
		}
		nameNorm := validate.NormalizeUnderscore(attr.Name)
		valueNorm := validate.NormalizeUnderscore(attr.Value)
		if nameNorm == "" || valueNorm == "" {
			fmt.Println("2")
			continue
		}

		normAttrs = append(normAttrs, ProductAttribute{Name: attr.Name, Value: attr.Value})
		if len(b) > 0 {
			b = append(b, ' ')
		}
		b = append(b, nameNorm...)
		b = append(b, '=')
		b = append(b, valueNorm...)
	}
	s := string(b)
	if s == "" {
		s = "_"
	}
	return normAttrs, s
}

var _ = sqlgenVariantExtended(
	&VariantExtended{}, &Variant{}, sq.AS("v"),
	sq.LEFT_JOIN, &Product{}, sq.AS("p"), "v.product_id = p.id",
	sq.LEFT_JOIN, &VariantExternal{}, sq.AS("vx"), "v.id = vx.id",
)

type VariantExtended struct {
	*Variant
	Product         *Product
	VariantExternal *VariantExternal
}

func (v *VariantExtended) GetFullName() string {
	if v.Product.Name != "" {
		return v.Product.Name + " - " + v.GetName()
	}
	return v.GetName()
}

var _ = sqlgenVariantExternal(&VariantExternal{})

type VariantExternal struct {
	ID int64

	ProductExternalCommon `sq:"inline"`

	ExternalProductID  string
	ExternalPrice      int
	ExternalBaseUnitID string
	ExternalUnitConv   float64
	ExternalAttributes []ProductAttribute
}

var _ = sqlgenVariantExternalExtended(
	&VariantExternalExtended{}, &Variant{}, sq.AS("v"),
	sq.LEFT_JOIN, &VariantExternal{}, sq.AS("vx"), "v.id = vx.id",
)

type VariantExternalExtended struct {
	*Variant
	*VariantExternal
}

var _ = sqlgenPriceDef(&PriceDef{}, &Variant{})

type PriceDef struct {
	WholesalePrice0 int `sq:"'wholesale_price_0'"`
	WholesalePrice  int
	ListPrice       int
	RetailPriceMin  int
	RetailPriceMax  int
}

func (p *PriceDef) IsValid() bool {
	return p.ListPrice > 0 &&
		p.WholesalePrice > 0 && p.WholesalePrice0 > 0 &&
		p.RetailPriceMin > 0 && p.RetailPriceMax > 0
}

func (p *PriceDef) ApplyTo(v *Variant) {
	v.ListPrice = p.ListPrice
	v.WholesalePrice0 = p.WholesalePrice0
	v.WholesalePrice = p.WholesalePrice
	v.RetailPriceMin = p.RetailPriceMin
	v.RetailPriceMax = p.RetailPriceMax
}

var _ = sqlgenVariantQuantity(&VariantQuantity{}, &Variant{})

type VariantQuantity struct {
	QuantityAvailable int // Available = OnHand - Reserved - 3
	QuantityOnHand    int
	QuantityReserved  int
}

var _ = substructEtopProduct(&EtopProduct{}, &Product{})

type EtopProduct struct {
	ID         int64
	SupplierID int64

	Name        string
	Description string
	DescHTML    string
	ShortDesc   string
	ImageURLs   []string `sq:"'image_urls'"`
	Status      Status3
	Code        string

	// Only default branch
	QuantityAvailable int // Available = OnHand - Reserved - 3
	QuantityOnHand    int
	QuantityReserved  int
}

func (p *EtopProduct) IsAvailable() bool {
	return p.QuantityAvailable > 0
}

var _ = sqlgenShopVariantExtended(
	&ShopVariantExtended{}, &ShopVariant{}, sq.AS("sv"),
	sq.LEFT_JOIN, &Variant{}, sq.AS("v"), "sv.variant_id = v.id",
	sq.LEFT_JOIN, &Product{}, sq.AS("p"), "v.product_id = p.id",
	sq.LEFT_JOIN, &ShopProduct{}, sq.AS("sp"), "sp.product_id = p.id",
)

type ShopVariantExtended struct {
	*ShopVariant
	VariantExtended
	*ShopProduct
}

func (v *ShopVariantExtended) GetFullName() string {
	var productName, variantName string
	if v.ShopProduct != nil && v.ShopProduct.Name != "" {
		productName = v.ShopProduct.Name
	} else {
		productName = v.Product.Name
	}
	if v.ShopVariant != nil && v.ShopVariant.Name != "" {
		variantName = v.ShopVariant.Name
	} else {
		variantName = v.Variant.GetName()
	}
	return productName + " - " + variantName
}

var _ = sqlgenShopVariant(&ShopVariant{})

type ShopVariant struct {
	ShopID       int64
	VariantID    int64
	CollectionID int64
	ProductID    int64

	Name        string
	Description string
	DescHTML    string
	ShortDesc   string
	ImageURLs   []string `sq:"'image_urls'"`
	Note        string
	Tags        []string

	RetailPrice int
	Status      Status3

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`

	NameNorm string
}

func (p *ShopVariant) BeforeInsert() error {
	p.NameNorm = validate.NormalizeSearch(p.Name)
	return nil
}

func (p *ShopVariant) BeforeUpdate() error {
	p.NameNorm = validate.NormalizeSearch(p.Name)
	return nil
}

var _ = sqlgenShopProduct(&ShopProduct{})

type ShopProduct struct {
	ShopID        int64
	ProductID     int64
	CollectionIDs []int64 `sq:"-"`

	Name        string
	Description string
	DescHTML    string
	ShortDesc   string
	ImageURLs   []string `sq:"'image_urls'"`
	Note        string
	Tags        []string

	RetailPrice int
	Status      Status3

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`

	NameNorm string

	ProductSourceID   int64  `sq:"-"`
	ProductSourceName string `sq:"-"`
	ProductSourceType string `sq:"-"`
}

func (p *ShopProduct) BeforeInsert() error {
	p.NameNorm = validate.NormalizeSearch(p.Name)
	return nil
}

func (p *ShopProduct) BeforeUpdate() error {
	p.NameNorm = validate.NormalizeSearch(p.Name)
	return nil
}

var _ = sqlgenShopProductFtProductFtVariantFtShopVariant(
	&ShopProductFtProductFtVariantFtShopVariant{}, &ShopProduct{}, sq.AS("sp"),
	sq.LEFT_JOIN, &Product{}, sq.AS("p"), "sp.product_id = p.id",
	sq.LEFT_JOIN, &Variant{}, sq.AS("v"), "sp.product_id = v.product_id",
	sq.LEFT_JOIN, &VariantExternal{}, sq.AS("vx"), "vx.id = v.id",
	sq.LEFT_JOIN, &ShopVariant{}, sq.AS("sv"), "sp.shop_id = sv.shop_id and sv.variant_id = v.id",
)

type ShopProductFtProductFtVariantFtShopVariant struct {
	*ShopProduct
	*Product
	*Variant
	*VariantExternal
	*ShopVariant
}

var _ = sqlgenShopProductExtended(
	&ShopProductExtended{}, &ShopProduct{}, sq.AS("sp"),
	sq.LEFT_JOIN, &Product{}, sq.AS("p"), "sp.product_id = p.id",
)

type ShopProductExtended struct {
	*ShopProduct
	*Product
	Variants []*ShopVariantExt `sq:"-"`
}

var _ = sqlgenShopVariantExt(
	&ShopVariantExt{}, &ShopVariant{}, sq.AS("sv"),
	sq.JOIN, &Variant{}, sq.AS("v"), "sp.product_id = v.product_id",
	sq.LEFT_JOIN, &VariantExternal{}, sq.AS("vx"), "vx.id = v.id",
)

type ShopVariantExt struct {
	*ShopVariant
	*Variant
	*VariantExternal
}

// select ... from variant join ... where v.product_id in (...)

var _ = sqlgenProductFtVariantFtShopProduct(
	&ProductFtVariantFtShopProduct{}, &Product{}, sq.AS("p"),
	sq.LEFT_JOIN, &Variant{}, sq.AS("v"), "v.product_id = p.id",
	sq.LEFT_JOIN, &VariantExternal{}, sq.AS("vx"), "vx.id = v.id",
	sq.LEFT_JOIN, &ShopProduct{}, sq.AS("sp"), "sp.product_id = p.id",
)

type ProductFtVariantFtShopProduct struct {
	*Product
	*Variant
	*VariantExternal
	*ShopProduct
}

var _ = sqlgenProductFtShopProduct(
	&ProductFtShopProduct{}, &Product{}, sq.AS("p"),
	sq.LEFT_JOIN, &ShopProduct{}, sq.AS("sp"), "sp.product_id = p.id",
)

type ProductFtShopProduct struct {
	*Product
	*ShopProduct
}

type ShopProductFtVariant struct {
	*ShopProduct
	*Product
	Variants []*ShopVariantExtended
}

// type ShopProductExternal struct {
// 	AccountID     int64
// 	ProductID  int64
// 	AccountID  int64
// 	ExternalID string

// 	CreatedAt time.Time `sq:"create"`
// }

var _ = sqlgenShopCollection(&ShopCollection{})

type ShopCollection struct {
	ID     int64
	ShopID int64

	Name        string
	Description string
	DescHTML    string
	ShortDesc   string

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

var _ = sqlgenProductShopCollection(&ProductShopCollection{})

type ProductShopCollection struct {
	CollectionID int64
	ProductID    int64
	ShopID       int64
	Status       int
	CreatedAt    time.Time `sq:"create"`
	UpdatedAt    time.Time `sq:"update"`
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

	IsTest int
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
	Status         *Status3
	ExternalStatus *Status3
	SupplierStatus *Status3
	EtopStatus     *Status3
}

type ProductStatusUpdate struct {
	SupplierStatus *Status3
	EtopStatus     *Status3
}

var _ = sqlgenOrder(&Order{})

type Order struct {
	ID          int64
	ShopID      int64
	Code        string
	EdCode      string
	ProductIDs  []int64
	VariantIDs  []int64
	SupplierIDs []int64
	PartnerID   int64

	Currency      string
	PaymentMethod string

	Customer        *OrderCustomer
	CustomerAddress *OrderAddress
	BillingAddress  *OrderAddress
	ShippingAddress *OrderAddress
	CustomerName    string
	CustomerPhone   string
	CustomerEmail   string

	CreatedAt    time.Time
	ProcessedAt  time.Time `sq:"create"`
	UpdatedAt    time.Time `sq:"update"`
	ClosedAt     time.Time
	ConfirmedAt  time.Time
	CancelledAt  time.Time
	CancelReason string

	CustomerConfirm Status3
	ExternalConfirm Status3
	ShopConfirm     Status3
	ConfirmStatus   Status3

	FulfillmentShippingStatus Status5
	CustomerPaymentStatus     Status3
	EtopPaymentStatus         Status4

	// -1:cancelled, 0:default, 1:delivered, 2:processing
	//
	// 0: just created, still error, can be edited
	// 2: processing, can not be edited
	// 1: done
	//-1: cancelled
	//-2: returned
	Status Status5

	FulfillmentShippingStates  []string
	FulfillmentPaymentStatuses []int

	Lines           OrderLinesList
	Discounts       []*OrderDiscount
	TotalItems      int
	BasketValue     int
	TotalWeight     int
	TotalTax        int
	OrderDiscount   int
	TotalDiscount   int
	ShopShippingFee int
	TotalFee        int
	FeeLines        OrderFeeLines
	ShopCOD         int
	TotalAmount     int

	OrderNote       string
	ShopNote        string
	ShippingNote    string
	OrderSourceType OrderSourceType
	OrderSourceID   int64
	ExternalOrderID string
	ReferenceURL    string
	ExternalURL     string
	ShopShipping    *OrderShipping
	IsOutsideEtop   bool
	Fulfillments    []*Fulfillment `sq:"-"`
	ExternalData    *OrderExternal `sq:"-"`

	// @deprecated: use try_on instead
	GhnNoteCode string
	TryOn       TryOn

	CustomerNameNorm string
	ProductNameNorm  string
}

func (m *Order) SelfURL(baseURL string, accType int) string {
	switch accType {
	case TagEtop, TagSupplier:
		return ""

	case TagShop:
		if baseURL == "" || m.ShopID == 0 || m.ID == 0 {
			return ""
		}
		return fmt.Sprintf("%v/s/%v/orders/%v", baseURL, m.ShopID, m.ID)

	default:
		panic(fmt.Sprintf("unsupported account type: %v", accType))
	}
}

func (m *Order) GetTotalFee() int {
	if m.TotalFee == 0 && m.ShopShippingFee != 0 {
		return m.ShopShippingFee
	}
	return m.TotalFee
}

func GetFeeLinesWithFallback(lines []OrderFeeLine, totalFee *int32, shopShippingFee *int32) []OrderFeeLine {
	if len(lines) == 0 &&
		(totalFee == nil || *totalFee == 0) &&
		shopShippingFee != nil && *shopShippingFee != 0 {
		return []OrderFeeLine{
			{
				Amount: int(*shopShippingFee),
				Desc:   "Phí giao hàng tính cho khách",
				Name:   "Phí giao hàng tính cho khách",
				Type:   OrderFeeShipping,
			},
		}
	}
	return lines
}

func GetShippingFeeFromFeeLines(lines []OrderFeeLine) int {
	s := 0
	for _, line := range lines {
		if line.Type == OrderFeeShipping {
			s += line.Amount
		}
	}
	return s
}

func CalcTotalFee(lines []OrderFeeLine) int {
	s := 0
	for _, line := range lines {
		s += line.Amount
	}
	return s
}

func (m *Order) GetFeeLines() []OrderFeeLine {
	if m.TotalFee == 0 && m.ShopShippingFee != 0 {
		return []OrderFeeLine{
			{
				Amount: m.ShopShippingFee,
				Desc:   "Phí giao hàng (shop tính cho khách)",
				Name:   "Phí giao hàng",
				Type:   OrderFeeShipping,
			},
		}
	}
	return m.FeeLines
}

func (m *Order) GetTryOn() TryOn {
	if m.TryOn != "" {
		return m.TryOn
	}
	return TryOnFromGHNNoteCode(m.GhnNoteCode)
}

func (m *Order) BeforeInsert() error {
	if (m.TryOn == "" || m.TryOn == "unknown") && m.GhnNoteCode != "" {
		m.TryOn = TryOnFromGHNNoteCode(m.GhnNoteCode)
	}
	if m.ShopShipping != nil {
		if err := m.ShopShipping.Validate(); err != nil {
			return err
		}
	}

	m.CustomerName = m.Customer.GetFullName()
	m.CustomerNameNorm = validate.NormalizeSearch(m.CustomerName)

	var b strings.Builder
	for _, line := range m.Lines {
		b.WriteString(line.ProductName)
		b.WriteString(" ")
	}

	m.ProductNameNorm = validate.NormalizeSearch(b.String())
	return nil
}

func (m *Order) BeforeUpdate() error {
	if (m.TryOn == "" || m.TryOn == "unknown") && m.GhnNoteCode != "" {
		m.TryOn = TryOnFromGHNNoteCode(m.GhnNoteCode)
	}
	if m.ShopShipping != nil && m.ShopShipping.ShippingProvider != "" {
		if err := m.ShopShipping.Validate(); err != nil {
			return err
		}
	}

	if m.Customer != nil {
		m.CustomerName = m.Customer.GetFullName()
		m.CustomerNameNorm = validate.NormalizeSearch(m.CustomerName)
		m.CustomerEmail = m.Customer.Email
	}
	if m.ShippingAddress != nil {
		m.CustomerPhone = m.ShippingAddress.Phone
	}
	return nil
}

func (m *Order) GetChargeableWeight() int {
	if m.ShopShipping != nil && m.ShopShipping.ChargeableWeight != 0 {
		return m.ShopShipping.ChargeableWeight
	}
	return m.TotalWeight
}

type OrderShipping struct {
	ShopAddress         *OrderAddress    `json:"shop_address"`
	ReturnAddress       *OrderAddress    `json:"return_address"`
	ExternalServiceID   string           `json:"external_service_id"`
	ExternalShippingFee int              `json:"external_shipping_fee"`
	ExternalServiceName string           `json:"external_service_name"`
	ShippingProvider    ShippingProvider `json:"shipping_provider"`
	ProviderServiceID   string           `json:"provider_service_id"`
	IncludeInsurance    bool             `json:"include_insurance"`

	Length int `json:"length"`
	Width  int `json:"width"`
	Height int `json:"height"`

	GrossWeight      int `json:"gross_weight"`
	ChargeableWeight int `json:"chargeable_weight"`
}

func (s *OrderShipping) Validate() error {
	if !VerifyShippingProvider(s.ShippingProvider) {
		return cm.Errorf(cm.InvalidArgument, nil, "Nhà vận chuyển không hợp lệ")
	}
	return nil
}

func (s *OrderShipping) GetPickupAddress() *OrderAddress {
	if s == nil {
		return nil
	}
	return s.ShopAddress
}

func (s *OrderShipping) GetShippingProvider() ShippingProvider {
	if s == nil {
		return ""
	}
	return s.ShippingProvider
}

func (s *OrderShipping) GetShippingServiceCode() string {
	if s == nil {
		return ""
	}
	return cm.Coalesce(s.ProviderServiceID, s.ExternalServiceID)
}

func (s *OrderShipping) GetPtrShippingServiceCode() *string {
	if s == nil || s.ExternalServiceID == "" {
		return nil
	}
	return &s.ExternalServiceID
}

var _ = sqlgenOrderExternal(&OrderExternal{})

type OrderExternal struct {
	ID                   int64
	OrderSourceID        int64
	ExternalOrderSource  string
	ExternalProvider     string
	ExternalOrderID      string
	ExternalOrderCode    string
	ExternalUserID       string
	ExternalCustomerID   string
	ExternalCreatedAt    time.Time
	ExternalProcessedAt  time.Time
	ExternalUpdatedAt    time.Time
	ExternalClosedAt     time.Time
	ExternalCancelledAt  time.Time
	ExternalCancelReason string
	ExternalLines        []*ExternalOrderLine
	ExternalData         json.RawMessage
}

type OrderExternalCreate struct {
	ID int64

	OrderExternalCreateOrder
	OrderExternalCreateExternal
}

var _ = sqlgenOrderExternalCreateOrder(&OrderExternalCreateOrder{}, &Order{})

type OrderExternalCreateOrder struct {
	ShopID          int64
	Code            string
	SupplierIDs     []int64 `sq:"'supplier_ids'"`
	Currency        string
	PaymentMethod   string
	Customer        *OrderCustomer
	CustomerAddress *OrderAddress
	BillingAddress  *OrderAddress
	ShippingAddress *OrderAddress
	CustomerPhone   string
	CustomerEmail   string

	CreatedAt time.Time `sq:"create"`

	ProcessedAt  time.Time
	UpdatedAt    time.Time `sq:"update"`
	ClosedAt     time.Time
	ConfirmedAt  time.Time
	CancelledAt  time.Time
	CancelReason string

	Lines         OrderLinesList
	VariantIDs    []int64 `sq:"'variant_ids'"`
	Discounts     []*OrderDiscount
	TotalItems    int
	BasketValue   int
	TotalWeight   int
	TotalTax      int
	TotalDiscount int
	TotalAmount   int
	ShopCOD       int

	ExternalConfirm Status3

	OrderSourceType OrderSourceType
	OrderSourceID   int64
	ExternalOrderID string
	IsOutsideEtop   bool
}

var _ = sqlgenOrderExternalCreateExternal(&OrderExternalCreateExternal{}, &OrderExternal{})

type OrderExternalCreateExternal struct {
	OrderSourceID        int64
	ExternalOrderSource  string
	ExternalOrderID      string
	ExternalOrderCode    string
	ExternalUserID       string
	ExternalCustomerID   string
	ExternalCreatedAt    time.Time `sq:"create"`
	ExternalProcessedAt  time.Time
	ExternalUpdatedAt    time.Time `sq:"update"`
	ExternalClosedAt     time.Time
	ExternalCancelledAt  time.Time
	ExternalCancelReason string
	ExternalLines        []*ExternalOrderLine

	ExternalData json.RawMessage
}

type OrderExternalUpdate struct {
	ID int64

	OrderExternalUpdateOrder
	OrderExternalUpdateExternal
}

var _ = sqlgenOrderExternalUpdateOrder(&OrderExternalUpdateOrder{}, &Order{})

type OrderExternalUpdateOrder struct {
	UpdatedAt    time.Time `sq:"update"`
	ClosedAt     time.Time
	CancelledAt  time.Time
	CancelReason string

	ExternalConfirm Status3
}

var _ = sqlgenOrderExternalUpdateExternal(&OrderExternalUpdateExternal{}, &OrderExternal{})

type OrderExternalUpdateExternal struct {
	ExternalUpdatedAt    time.Time `sq:"update"`
	ExternalCancelledAt  time.Time
	ExternalCancelReason string
}

type OrderLinesList []*OrderLine

func (lines OrderLinesList) GetTotalItems() int {
	s := 0
	for _, line := range lines {
		s += line.Quantity
	}
	return s
}

func (lines OrderLinesList) GetTotalWeight() int {
	s := 0
	for _, line := range lines {
		s += line.Weight
	}
	return s
}

func (lines OrderLinesList) GetTotalRetailAmount() int {
	s := 0
	for _, line := range lines {
		s += line.RetailPrice
	}
	return s
}

func (lines OrderLinesList) GetTotalPaymentAmount() int {
	s := 0
	for _, line := range lines {
		s += line.PaymentPrice
	}
	return s
}

func (lines OrderLinesList) GetSummary() string {
	var b strings.Builder
	for _, line := range lines {
		fprintf(&b, "%2d x %v", line.Quantity, line.ProductName)
		if len(line.Attributes) > 0 {
			fprintf(&b, " (")
		}
		for i, attr := range line.Attributes {
			if i > 0 {
				fprintf(&b, ", ")
			}
			fprintf(&b, "%v: %v", attr.Name, attr.Value)
		}
		if len(line.Attributes) > 0 {
			fprintf(&b, ")")
		}
		fprintf(&b, "\n")
	}
	return b.String()
}

func fprintf(w io.Writer, format string, args ...interface{}) {
	_, _ = fmt.Fprintf(w, format, args...)
}

var _ = sqlgenOrderLine(&OrderLine{})

type OrderLine struct {
	OrderID     int64  `json:"order_id"`
	VariantID   int64  `json:"variant_id"`
	ProductName string `json:"product_name"`
	ProductID   int64  `json:"product_id"`

	SupplierID int64 `json:"supplier_id"`
	ShopID     int64 `json:"shop_id"`

	ExternalVariantID       string `json:"x_variant_id"`
	ExternalSupplierOrderID string `json:"-"`

	SupplierName string `json:"supplier_name"`

	// Vendor      string `json:"vendor"` // supplier?

	// CreatedAt    time.Time `json:"-"`
	// CreatedBy    int64     `json:"-"`
	UpdatedAt time.Time `json:"-" `
	// UpdatedBy int64     `json:"-"`
	ClosedAt    time.Time `json:"-" `
	ConfirmedAt time.Time `json:"-" `
	// ConfirmedBy  int64     `json:"-"`
	CancelledAt time.Time `json:"-" `
	// CancelledBy  int64     `json:"-"`
	CancelReason string `json:"-"`

	SupplierConfirm Status3 `json:"-"`
	// EtopConfirm     int `json:"etop_confirm"`
	// ConfirmStatus int `json:"confirm_status"`

	Status Status3 `json:"-"`

	Weight          int `json:"weight"`
	Quantity        int `json:"quantity"`
	WholesalePrice0 int `json:"wholesale_price_0" sq:"'wholesale_price_0'"`
	WholesalePrice  int `json:"wholesale_price"`
	ListPrice       int `json:"list_price"`
	RetailPrice     int `json:"retail_price"`
	PaymentPrice    int `json:"payment_price"`

	LineAmount       int  `json:"line_amount"`
	TotalDiscount    int  `json:"discount"`
	TotalLineAmount  int  `json:"total_line_amount"`
	RequiresShipping bool `json:"requires_shipping"`

	ImageURL      string             `json:"image_url" `
	Attributes    []ProductAttribute `json:"attributes" sq:"-"`
	IsOutsideEtop bool               `json:"is_outside_etop"`
	Code          string             `json:"code"`
}

func (l *OrderLine) GetRetailAmount() int {
	if l == nil {
		return 0
	}
	return l.RetailPrice * l.Quantity
}

func (l *OrderLine) GetPaymentAmount() int {
	if l == nil {
		return 0
	}
	return l.PaymentPrice * l.Quantity
}

func (l *OrderLine) GetTotalDiscount() int {
	if l == nil {
		return 0
	}
	return (l.RetailPrice - l.PaymentPrice) * l.Quantity
}

var _ = sqlgenOrderLineExtended(&OrderLineExtended{},
	&OrderLine{}, sq.AS("ol"), sq.LEFT_JOIN,
	&Variant{}, sq.AS("v"), "v.id = ol.variant_id",
)

type OrderLineExtended struct {
	*OrderLine
	*Variant
}

func (olExtended *OrderLineExtended) ToOrderLine() *OrderLine {
	return &OrderLine{
		OrderID:                 olExtended.OrderID,
		VariantID:               olExtended.VariantID,
		ProductName:             olExtended.ProductName,
		ProductID:               olExtended.OrderLine.ProductID,
		SupplierID:              olExtended.OrderLine.SupplierID,
		ShopID:                  olExtended.OrderLine.ShopID,
		ExternalVariantID:       olExtended.ExternalVariantID,
		ExternalSupplierOrderID: olExtended.ExternalSupplierOrderID,
		SupplierName:            olExtended.SupplierName,
		UpdatedAt:               olExtended.OrderLine.UpdatedAt,
		ClosedAt:                olExtended.OrderLine.ClosedAt,
		ConfirmedAt:             olExtended.OrderLine.ConfirmedAt,
		CancelledAt:             olExtended.OrderLine.CancelledAt,
		CancelReason:            olExtended.OrderLine.CancelReason,
		SupplierConfirm:         olExtended.OrderLine.SupplierConfirm,
		Status:                  olExtended.OrderLine.Status,
		Weight:                  olExtended.OrderLine.Weight,
		Quantity:                olExtended.OrderLine.Quantity,
		WholesalePrice0:         olExtended.OrderLine.WholesalePrice0,
		WholesalePrice:          olExtended.OrderLine.WholesalePrice,
		ListPrice:               olExtended.OrderLine.ListPrice,
		RetailPrice:             olExtended.OrderLine.RetailPrice,
		PaymentPrice:            olExtended.OrderLine.PaymentPrice,
		LineAmount:              olExtended.OrderLine.LineAmount,
		TotalDiscount:           olExtended.OrderLine.TotalDiscount,
		TotalLineAmount:         olExtended.OrderLine.TotalLineAmount,
		RequiresShipping:        olExtended.OrderLine.RequiresShipping,
		ImageURL:                olExtended.OrderLine.ImageURL,
		Attributes:              olExtended.Variant.Attributes,
		IsOutsideEtop:           olExtended.OrderLine.IsOutsideEtop,
		Code:                    olExtended.OrderLine.Code,
	}
}

type ExternalOrderLine struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	VariantID string `json:"variant_id"`
	SKU       string `json:"sku"`

	Name        string `json:"name"`
	Title       string `json:"title"`
	VariantName string `json:"variant_name"`
	Vendor      string `json:"vendor"`
	Type        string `json:"type"`

	Weight         int `json:"weight"`
	Quantity       int `json:"quantity"`
	Price          int `json:"price"`
	PriceOriginal  int `json:"price_original"`
	PricePromotion int `json:"price_promotion"`
}

type OrderFeeLine struct {
	Amount int          `json:"amount"`
	Desc   string       `json:"desc"`
	Code   string       `json:"code"`
	Name   string       `json:"name"`
	Type   OrderFeeType `json:"type"`
}

type OrderFeeLines []OrderFeeLine

func (feeLines OrderFeeLines) GetTotalFee() int {
	s := 0
	for _, line := range feeLines {
		s += line.Amount
	}
	return s
}

func (feeLines OrderFeeLines) GetShippingFee() int {
	s := 0
	for _, line := range feeLines {
		if line.Type == OrderFeeShipping {
			s += line.Amount
		}
	}
	return s
}

func (feeLines OrderFeeLines) GetTax() int {
	s := 0
	for _, line := range feeLines {
		if line.Type == OrderFeeTax {
			s += line.Amount
		}
	}
	return s
}

func (feeLines OrderFeeLines) GetOtherFee() int {
	s := 0
	for _, line := range feeLines {
		if line.Type != OrderFeeShipping && line.Type != OrderFeeTax {
			s += line.Amount
		}
	}
	return s
}

type OrderCustomer struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	FullName      string `json:"full_name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Gender        string `json:"gender"`
	Birthday      string `json:"birthday"`
	VerifiedEmail bool   `json:"verified_email"`

	ExternalID string `json:"external_id"`
}

func (m *OrderCustomer) GetFullName() string {
	if m.FullName != "" {
		return m.FullName
	}
	return m.FirstName + " " + m.LastName
}

func (o *OrderCustomer) UpdateCustomer(fullname string) *OrderCustomer {
	if fullname != "" {
		o.FullName = fullname
	}
	return o
}

type OrderAddress struct {
	FullName  string `json:"full_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`

	Country  string `json:"country"`
	City     string `json:"city"`
	Province string `json:"province"`
	District string `json:"district"`
	Ward     string `json:"ward"`
	Zip      string `json:"zip"`

	DistrictCode string `json:"district_code"`
	ProvinceCode string `json:"province_code"`
	WardCode     string `json:"ward_code"`

	Company  string `json:"company"`
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
}

func (o *OrderAddress) UpdateAddress(phone string, fullname string) *OrderAddress {
	if phone != "" {
		o.Phone = phone
	}
	if fullname != "" {
		o.FullName = fullname
	}
	return o
}

type OrderDiscount struct {
	Code   string `json:"code"`
	Type   string `json:"type"`
	Amount int    `json:"amount"`
}

func (m *OrderAddress) GetFullName() string {
	if m.FullName != "" {
		return m.FullName
	}
	return m.FirstName + " " + m.LastName
}

var _ = sqlgenFulfillment(&Fulfillment{})

type Fulfillment struct {
	ID         int64
	OrderID    int64
	ShopID     int64
	SupplierID int64
	PartnerID  int64

	SupplierConfirm Status3
	ShopConfirm     Status3
	ConfirmStatus   Status3

	TotalItems        int
	TotalWeight       int
	BasketValue       int
	TotalDiscount     int
	TotalAmount       int
	TotalCODAmount    int
	OriginalCODAmount int

	// If a fulfillment is compensated, this field contains the actual amount
	// that carrier payback to shop.
	ActualCompensationAmount int

	ShippingFeeCustomer      int // shop charges customer/shop
	ShippingFeeShop          int // etop charges shop, actual_shipping_service_fee
	ShippingFeeShopLines     []*ShippingFeeLine
	ShippingServiceFee       int // copy from order
	ExternalShippingFee      int // provider charges eTop
	ProviderShippingFeeLines []*ShippingFeeLine
	EtopDiscount             int
	EtopFeeAdjustment        int // eTop điều chỉnh phi (phần thêm)

	ShippingFeeMain       int
	ShippingFeeReturn     int
	ShippingFeeInsurance  int
	ShippingFeeAdjustment int
	ShippingFeeCODS       int
	ShippingFeeInfoChange int
	ShippingFeeOther      int

	// EtopAdjustedShippingFeeMain: eTop điều chỉnh cước phí chính
	EtopAdjustedShippingFeeMain int
	// EtopPriceRule: true khi áp dụng bảng giá eTop, với giá `EtopAdjustedShippingFeeMain`
	EtopPriceRule bool

	VariantIDs []int64
	Lines      OrderLinesList

	TypeFrom      FulfillmentEndpoint
	TypeTo        FulfillmentEndpoint
	AddressFrom   *Address
	AddressTo     *Address
	AddressReturn *Address

	AddressToProvinceCode string
	AddressToDistrictCode string
	AddressToWardCode     string

	CreatedAt                   time.Time `sq:"create"`
	UpdatedAt                   time.Time `sq:"update"`
	ClosedAt                    time.Time
	ExpectedDeliveryAt          time.Time
	ExpectedPickAt              time.Time
	CODEtopTransferedAt         time.Time
	ShippingFeeShopTransferedAt time.Time
	ShippingCancelledAt         time.Time
	ShippingDeliveredAt         time.Time
	ShippingReturnedAt          time.Time
	ShippingCreatedAt           time.Time
	ShippingPickingAt           time.Time
	ShippingHoldingAt           time.Time
	ShippingDeliveringAt        time.Time
	ShippingReturningAt         time.Time

	MoneyTransactionID                 int64
	MoneyTransactionShippingExternalID int64

	CancelReason string

	//CreatedBy   int64
	//UpdatedBy   int64
	//CancelledBy int64

	ShippingProvider  ShippingProvider
	ProviderServiceID string
	ShippingCode      string
	ShippingNote      string
	TryOn             TryOn
	IncludeInsurance  bool

	ExternalShippingName        string
	ExternalShippingID          string // it's shipping_service_code
	ExternalShippingCode        string // it's shipping_code
	ExternalShippingCreatedAt   time.Time
	ExternalShippingUpdatedAt   time.Time
	ExternalShippingCancelledAt time.Time
	ExternalShippingDeliveredAt time.Time
	ExternalShippingReturnedAt  time.Time
	ExternalShippingClosedAt    time.Time
	ExternalShippingState       string
	ExternalShippingStateCode   string
	ExternalShippingStatus      Status5
	ExternalShippingNote        string
	ExternalShippingSubState    string

	ExternalShippingData json.RawMessage

	ShippingState     ShippingState
	ShippingStatus    Status5
	EtopPaymentStatus Status4

	// -1:cancelled, 0:default, 1:delivered, 2:processing
	//
	// 0: just created, still error, can be edited
	// 2: processing, can not be edited
	// 1: done
	//-1: cancelled
	//-2: returned
	Status Status5

	SyncStatus Status4 // -1:error, 0:new, 1:created, 2:pending
	SyncStates *FulfillmentSyncStates

	// Updated by webhook or querying GHN API
	LastSyncAt time.Time

	ExternalShippingLogs []*ExternalShippingLog
	AdminNote            string
	IsPartialDelivery    bool
}

func (f *Fulfillment) SelfURL(baseURL string, accType int) string {
	switch accType {
	case TagEtop, TagSupplier:
		return ""

	case TagShop:
		if baseURL == "" || f.ShopID == 0 || f.ID == 0 {
			return ""
		}
		return fmt.Sprintf("%v/s/%v/fulfillments/%v", baseURL, f.ShopID, f.ID)

	default:
		panic(fmt.Sprintf("unsupported account type: %v", accType))
	}
}

type ExternalShippingLog struct {
	StateText string
	Time      string
	Message   string
}

func (f *Fulfillment) BeforeInsert() error {
	return f.BeforeUpdate()
}

func (f *Fulfillment) BeforeUpdate() error {
	if f.AddressTo != nil {
		f.AddressToProvinceCode = f.AddressTo.ProvinceCode
		f.AddressToDistrictCode = f.AddressTo.DistrictCode
		f.AddressToWardCode = f.AddressTo.WardCode
	}
	return nil
}

func CalcShopShippingFee(externalFee int, ffm *Fulfillment) int {
	if ffm == nil {
		return externalFee
	}
	fee := externalFee + ffm.EtopFeeAdjustment - ffm.EtopDiscount
	if fee < 0 {
		return 0
	}
	return fee
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

	Company   string       `json:"company"`
	Address1  string       `json:"address1"`
	Address2  string       `json:"address2"`
	Type      string       `json:"type"`
	AccountID int64        `json:"account_id"`
	Notes     *AddressNote `json:"notes"`
	CreatedAt time.Time    `sq:"create" json:"-"`
	UpdatedAt time.Time    `sq:"update" json:"-"`
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
	ID         int64
	Amount     int
	ShopID     int64
	SupplierID int64
	Type       string
	Status     Status3
	CreatedAt  time.Time `sq:"create"`
	UpdatedAt  time.Time `sq:"update"`
	PaidAt     time.Time
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
	return
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
		return cm.Errorf(cm.InvalidArgument, nil, "missing AccountID")
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

func (f *Fulfillment) ApplyEtopPrice(price int) error {
	f.EtopPriceRule = true
	f.EtopAdjustedShippingFeeMain = price
	return nil
}
