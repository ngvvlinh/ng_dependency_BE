package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"etop.vn/api/top/types/etc/account_type"
	"etop.vn/api/top/types/etc/ghn_note_code"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/try_on"
	"etop.vn/api/top/types/etc/user_source"
	addressmodel "etop.vn/backend/com/main/address/model"
	identitysharemodel "etop.vn/backend/com/main/identity/sharemodel"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/code/gencode"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/pkg/common/validate"
	etopmodel "etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

type (
	SubjectType string
)

const (
	SubjectTypeAccount SubjectType = "account"
	SubjectTypeUser    SubjectType = "user"
)

var _ = sqlgenExternalAccountAhamove(&ExternalAccountAhamove{})

type ExternalAccountAhamove struct {
	ID                  dot.ID
	OwnerID             dot.ID
	Phone               string
	Name                string
	ExternalID          string
	ExternalVerified    bool
	ExternalCreatedAt   time.Time
	ExternalToken       string
	CreatedAt           time.Time `sq:"create"`
	UpdatedAt           time.Time `sq:"update"`
	LastSendVerifiedAt  time.Time
	ExternalTicketID    string
	IDCardFrontImg      string
	IDCardBackImg       string
	PortraitImg         string
	WebsiteURL          string
	FanpageURL          string
	CompanyImgs         []string
	BusinessLicenseImgs []string

	ExternalDataVerified json.RawMessage

	UploadedAt time.Time
}

var _ = sqlgenSale(&Affiliate{})

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
	BankAccount       *identitysharemodel.BankAccount
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
	return etopmodel.TryOnFromGHNNoteCode(s.GhnNoteCode)
}

var _ = sqlgenShopExtended(
	&ShopExtended{}, &Shop{}, "s",
	sq.LEFT_JOIN, &addressmodel.Address{}, "a", "s.address_id = a.id",
	sq.LEFT_JOIN, &User{}, "u", "s.owner_id = u.id",
)

type ShopExtended struct {
	*Shop
	Address *addressmodel.Address
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

	ContactPersons  []*identitysharemodel.ContactPerson
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
	fullName = etopmodel.CoalesceString2(m.AccountUser.FullName, m.User.FullName)
	shortName = etopmodel.CoalesceString2(m.AccountUser.ShortName, m.User.ShortName)
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
