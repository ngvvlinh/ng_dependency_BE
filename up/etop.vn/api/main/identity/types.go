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

type ExternalAccountAhamove struct {
	ID                int64
	Phone             string
	Name              string
	ExternalToken     string
	ExternalVerified  bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ExternalCreatedAt time.Time
}
