package model

import (
	"encoding/json"
	"time"

	"o.o/capi/dot"
)

// TODO: change ExternalAccountAhamove to AccountShipnow
// +sqlgen
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

	UploadedAt   time.Time
	ConnectionID dot.ID
}
