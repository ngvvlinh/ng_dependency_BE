package identity

import (
	"time"

	etoptypes "etop.vn/api/main/etop"
)

type Shop struct {
	ID      int64
	Name    string
	OwnerID int64
	IsTest  int

	AddressID         int64
	ShipToAddressID   int64
	ShipFromAddressID int64
	Phone             string
	WebsiteURL        string
	ImageURL          string
	Email             string
	Code              string
	AutoCreateFFM     bool

	Status    etoptypes.Status3
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type User struct {
	ID int64

	FullName  string
	ShortName string
	Email     string
	Phone     string

	Status etoptypes.Status3 // 1: actual user, 0: stub, -1: disabled

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ExternalAccountAhamove struct {
	ID                  int64
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
