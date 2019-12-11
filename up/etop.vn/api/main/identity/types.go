package identity

import (
	"time"

	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/status3"
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

	AddressID         dot.ID
	ShipToAddressID   dot.ID
	ShipFromAddressID dot.ID
	Phone             string
	WebsiteURL        string
	ImageURL          string
	Email             string
	Code              string
	AutoCreateFFM     bool

	Status    status3.Status
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
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
	BankAccount *BankAccount
}

type User struct {
	ID dot.ID

	FullName  string
	ShortName string
	Email     string
	Phone     string

	Status status3.Status // 1: actual user, 0: stub, -1: disabled

	EmailVerifiedAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	RefUserID dot.ID
	RefSaleID dot.ID
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

type BankAccount struct {
	Name          string
	Province      string
	Branch        string
	AccountNumber string
	AccountName   string
}

type AccountCreatedEvent struct {
	meta.EventMeta
	ShopID dot.ID
	UserID dot.ID
}

type UserCreatedEvent struct {
	meta.EventMeta
	UserID    dot.ID
	Email     string
	FullName  string
	ShortName string

	Invitation *UserInvitation
}

type UserInvitation struct {
	Token      string
	AutoAccept bool

	FullName  string
	ShortName string
	Position  string
}

type GetCustomersByShop struct {
	meta.EventMeta
	ShopID dot.ID
}
