package accountshipnow

import (
	"time"

	"o.o/capi/dot"
)

// +gen:event:topic=event/accountshipnow

type ExternalAccountAhamove struct {
	ID                  dot.ID
	OwnerID             dot.ID
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
	ConnectionID        dot.ID
}

type ExternalAccountAhamoveCreatedEvent struct {
	ID           dot.ID
	Phone        string
	ShopID       dot.ID
	OwnerID      dot.ID
	ConnectionID dot.ID
}

type ExternalAccountShipnowVerifyRequestedEvent struct {
	ID           dot.ID
	ConnectionID dot.ID
	ShopID       dot.ID
	OwnerID      dot.ID
}

type ExternalAccountShipnowUpdateVerificationInfoEvent struct {
	ID           dot.ID
	OwnerID      dot.ID
	ConnectionID dot.ID
}
