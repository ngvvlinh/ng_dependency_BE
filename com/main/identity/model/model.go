package model

import (
	"fmt"
	"strings"
	"time"

	"o.o/api/top/types/etc/account_tag"
	"o.o/api/top/types/etc/account_type"
	"o.o/api/top/types/etc/ghn_note_code"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/try_on"
	"o.o/api/top/types/etc/user_source"
	addressmodel "o.o/backend/com/main/address/model"
	identitysharemodel "o.o/backend/com/main/identity/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/code/gencode"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/pkg/etc/typeutil"
	"o.o/capi/dot"
)

type (
	SubjectType string
)

const (
	SubjectTypeAccount SubjectType = "account"
	SubjectTypeUser    SubjectType = "user"
)

// +sqlgen
type Affiliate struct {
	ID          dot.ID
	OwnerID     dot.ID
	Name        string
	Phone       string
	Email       string
	IsTest      int
	Status      status3.Status
	CreatedAt   time.Time `sq:"create"`
	UpdatedAt   time.Time `sq:"update"`
	DeletedAt   time.Time
	BankAccount *identitysharemodel.BankAccount
}

var _ AccountInterface = &Affiliate{}

func (s *Affiliate) GetAccount() *Account {
	return &Account{
		ID:      s.ID,
		OwnerID: s.OwnerID,
		Name:    s.Name,
		Type:    account_type.Affiliate,
	}
}

// +sqlgen
type Account struct {
	ID       dot.ID
	OwnerID  dot.ID
	Name     string
	Type     account_type.AccountType
	ImageURL string
	URLSlug  string

	Rid dot.ID
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
	var _type account_type.AccountType
	switch cm.GetTag(s.ID) {
	case account_tag.TagPartner:
		_type = account_type.Partner
	case account_tag.TagCarrier:
		_type = account_type.Carrier
	default:
		panic("account type does not valid")
	}
	return &Account{
		ID:       s.ID,
		OwnerID:  s.OwnerID,
		Name:     s.Name,
		Type:     _type,
		ImageURL: s.ImageURL,
		URLSlug:  "",
	}
}

// +sqlgen
type Shop struct {
	ID      dot.ID
	OwnerID dot.ID
	IsTest  int
	Name    string

	AddressID         dot.ID
	ShipToAddressID   dot.ID
	ShipFromAddressID dot.ID
	Phone             string
	BankAccount       *identitysharemodel.BankAccount
	WebsiteURL        dot.NullString
	ImageURL          string
	Email             string
	Code              string
	AutoCreateFFM     bool

	OrderSourceID dot.ID

	Status    status3.Status
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time

	Address *addressmodel.Address `sq:"-"`

	RecognizedHosts []string

	// @deprecated use try_on instead
	GhnNoteCode ghn_note_code.GHNNoteCode
	TryOn       try_on.TryOnCode
	CompanyInfo *identitysharemodel.CompanyInfo
	// MoneyTransactionRRule format:
	// FREQ=DAILY
	// FREQ=WEEKLY;BYDAY=TU,TH,SA
	// FREQ=MONTHLY;BYDAY=-1FR
	MoneyTransactionRRule         string `sq:"'money_transaction_rrule'"`
	SurveyInfo                    []*SurveyInfo
	ShippingServiceSelectStrategy []*ShippingServiceSelectStrategyItem

	InventoryOverstock dot.NullBool

	WLPartnerID             dot.ID
	Rid                     dot.ID
	IsPriorMoneyTransaction dot.NullBool
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
	return typeutil.TryOnFromGHNNoteCode(s.GhnNoteCode)
}

// +sqlgen:           Shop    as s
// +sqlgen:left-join: Address as a on s.address_id = a.id
// +sqlgen:left-join: User    as u on s.owner_id = u.id
// +sqlgen:left-join: ShopSearch as ss on s.id = ss.id
type ShopExtended struct {
	*Shop
	ShopSearch *ShopSearch
	Address    *addressmodel.Address
	User       *User
}

// +sqlgen
type ShopSearch struct {
	ID       dot.ID
	Name     string
	NameNorm string
}

// +sqlgen=Shop
type ShopDelete struct {
	DeletedAt time.Time
}

// +sqlgen
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

	ContactPersons []*identitysharemodel.ContactPerson

	// Dùng để xác thực external_url được gửi khi tạo đơn hàng, sản phẩm, ...
	RecognizedHosts []string

	// Dùng để xác thực redirect_url được gửi khi gọi AuthorizeShop
	RedirectURLs []string

	// AvailableFromEtop: dùng để xác định `partner` này có thể xác thực shop trực tiếp từ client của Etop hay không
	// Sau khi xác thực xong sẽ trỏ trực tiếp về `redirect_url` trong field AvailableFromEtopConfig để về trang xác thực của `partner`
	AvailableFromEtop       bool
	AvailableFromEtopConfig *AvailableFromEtopConfig

	WhiteLabelKey string

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

// +sqlgen
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

// +sqlgen:      AccountAuth as aa
// +sqlgen:join: Partner     as p on aa.account_id = p.id
type AccountAuthFtPartner struct {
	*AccountAuth
	*Partner
}

// +sqlgen:      AccountAuth as aa
// +sqlgen:join: Shop        as s on aa.account_id = s.id
type AccountAuthFtShop struct {
	*AccountAuth
	*Shop
}

// +sqlgen
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

// +sqlgen:      PartnerRelation as pr
// +sqlgen:join: Shop as s on pr.subject_id = s.id
// +sqlgen:join: User as u on s.owner_id = u.id
type PartnerRelationFtShop struct {
	*PartnerRelation
	*Shop
	*User
}

// +sqlgen:      PartnerRelation as pr
// +sqlgen:join: User as u on pr.subject_id = u.id
type PartnerRelationFtUser struct {
	*PartnerRelation
	*User
}

type UserInner struct {
	FullName  string
	ShortName string
	Email     string
	Phone     string
}

// +sqlgen
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
	FullNameNorm            string
	IsTest                  int
	Source                  user_source.UserSource
	RefUserID               dot.ID
	RefSaleID               dot.ID
	WLPartnerID             dot.ID

	Rid int

	BlockedAt   time.Time
	BlockedBy   dot.ID
	BlockReason string
}

type UserExtended struct {
	User         *User
	UserInternal *UserInternal
}

// +convert:type=identity.User
// +sqlgen:      User as u
// +sqlgen:left-join: UserRefSaff as urs on u.id = urs.user_id
type UserFtRefSaff struct {
	*User
	*UserRefSaff
}

// Permission ...
type Permission struct {
	Roles       []string
	Permissions []string
}

// +convert:type=identity.AccountUser
// +sqlgen
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

	Rid dot.ID
}

// +sqlgen:      AccountUser as au
// +sqlgen:join: Account as a on au.account_id = a.id
// +sqlgen:join: User    as u on au.user_id = u.id
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

// +sqlgen=AccountUser
type AccountUserDelete struct {
	DeletedAt time.Time
}

// +sqlgen
type UserAuth struct {
	UserID   dot.ID
	AuthType string
	AuthKey  string

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

// +sqlgen
type UserInternal struct {
	ID      dot.ID
	Hashpwd string

	UpdatedAt time.Time `sq:"update"`
}

func CoalesceString2(s1, s2 string) string {
	if s1 != "" {
		return s1
	}
	return s2
}

// +sqlgen
type UserRefSaff struct {
	UserID  dot.ID
	RefAff  string
	RefSale string
}
