package model

import (
	"encoding/json"
	"time"

	"etop.vn/backend/pkg/etop/model"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenExternalAccountAhamove(&ExternalAccountAhamove{})

type ExternalAccountAhamove struct {
	ID                  int64
	OwnerID             int64
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
	ID        int64
	OwnerID   int64
	Name      string
	Phone     string
	Email     string
	IsTest    int
	Status    model.Status3
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time
}

var _ model.AccountInterface = &Affiliate{}

func (s *Affiliate) GetAccount() *model.Account {
	return &model.Account{
		ID:      s.ID,
		OwnerID: s.OwnerID,
		Name:    s.Name,
		Type:    model.TypeAffiliate,
	}
}
