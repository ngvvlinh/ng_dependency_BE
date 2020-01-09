package identity

import (
	"time"

	"etop.vn/api/main/address"
	identitytypes "etop.vn/api/main/identity/types"
	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/ghn_note_code"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/try_on"
	"etop.vn/api/top/types/etc/user_source"
	"etop.vn/capi/dot"
)

// +gen:event:topic=event/identity

type Permission struct {
	Roles       []string
	Permissions []string
}

type Shop struct {
	ID      dot.ID
	Name    string
	OwnerID dot.ID
	IsTest  int

	AddressID          dot.ID
	ShipToAddressID    dot.ID
	ShipFromAddressID  dot.ID
	Phone              string
	WebsiteURL         string
	ImageURL           string
	Email              string
	Code               string
	AutoCreateFFM      bool
	InventoryOverstock dot.NullBool
	GhnNoteCode        ghn_note_code.GHNNoteCode
	RecognizedHosts    []string
	// MoneyTransactionRRule format:
	// FREQ=DAILY
	// FREQ=WEEKLY;BYDAY=TU,TH,SA
	// FREQ=MONTHLY;BYDAY=-1FR
	MoneyTransactionRRule         string `sq:"'money_transaction_rrule'"`
	SurveyInfo                    []*SurveyInfo
	ShippingServiceSelectStrategy []*ShippingServiceSelectStrategyItem

	Status      status3.Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
	BankAccount *identitytypes.BankAccount
	TryOn       try_on.TryOnCode
	CompanyInfo *identitytypes.CompanyInfo
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

type Affiliate struct {
	ID          dot.ID
	OwnerID     dot.ID
	Name        string
	Phone       string
	Email       string
	IsTest      int
	Status      status3.Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
	BankAccount *identitytypes.BankAccount
}

type User struct {
	ID dot.ID

	FullName  string
	ShortName string
	Email     string
	Phone     string

	Status status3.Status // 1: actual user, 0: stub, -1: disabled

	EmailVerifiedAt         time.Time
	PhoneVerifiedAt         time.Time
	EmailVerificationSentAt time.Time
	PhoneVerificationSentAt time.Time

	IsTest      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	RefUserID   dot.ID
	RefSaleID   dot.ID
	WLPartnerID dot.ID
	Source      user_source.UserSource
}

type ExternalAccountAhamove struct {
	ID                  dot.ID
	Phone               string
	Name                string
	ExternalID          string
	ExternalToken       string
	ExternalVerified    bool
	CreatedAt           time.Time
	UpdatedAt           time.Time
	ExternalCreatedAt   time.Time
	LastSendVerifiedAt  time.Time
	ExternalTicketID    string
	IDCardFrontImg      string
	IDCardBackImg       string
	PortraitImg         string
	WebsiteURL          string
	FanpageURL          string
	CompanyImgs         []string
	BusinessLicenseImgs []string
	UploadedAt          time.Time
}

type AccountCreatedEvent struct {
	meta.EventMeta
	ShopID dot.ID
	UserID dot.ID
}

type UserCreatedEvent struct {
	meta.EventMeta
	UserID dot.ID

	Invitation *UserInvitation
}

type UserInvitation struct {
	Token      string
	AutoAccept bool
}

type GetCustomersByShop struct {
	meta.EventMeta
	ShopID dot.ID
}

type Partner struct {
	ID dot.ID

	OwnerID dot.ID

	Name string

	PublicName string

	ImageURL string

	WebsiteURL string

	WhiteLabelKey string
}

type ShopExtended struct {
	*Shop
	Address *address.Address
	User    *User
}
