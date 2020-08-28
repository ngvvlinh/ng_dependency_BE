package identity

import (
	"time"

	"o.o/api/main/address"
	identitytypes "o.o/api/main/identity/types"
	"o.o/api/meta"
	"o.o/api/top/types/etc/ghn_note_code"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/try_on"
	"o.o/api/top/types/etc/user_source"
	"o.o/capi/dot"
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
	WLPartnerID dot.ID
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

	BlockedAt   time.Time
	BlockedBy   dot.ID
	BlockReason string
	IsBlocked   bool
}

type UserFtRefSaff struct {
	*User
	RefAff  string
	RefSale string
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

type UserRefSaff struct {
	UserID  dot.ID
	RefAff  string
	RefSale string
}
